package mapreduce

import (
	"fmt"
	"sync"
	"log"
)

//
// schedule() starts and waits for all tasks in the given phase (Map
// or Reduce). the mapFiles argument holds the names of the files that
// are the inputs to the map phase, one per map task. nReduce is the
// number of reduce tasks. the registerChan argument yields a stream
// of registered workers; each item is the worker's RPC address,
// suitable for passing to call(). registerChan will yield all
// existing registered workers (if any) and new ones as they register.
//
func schedule(jobName string, mapFiles []string, nReduce int, phase jobPhase, registerChan chan string) {
	var ntasks int
	var n_other int // number of inputs (for reduce) or outputs (for map)
	switch phase {
	// map阶段
	case mapPhase:
		ntasks = len(mapFiles)
		n_other = nReduce
	// reduce阶段
	case reducePhase:
		ntasks = nReduce
		n_other = len(mapFiles)
	}

	fmt.Printf("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, n_other)

	// All ntasks tasks have to be scheduled on workers, and only once all of
	// them have been completed successfully should the function return.
	// Remember that workers may fail, and that any given worker may finish
	// multiple tasks.
	//
	// TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO
	/**
	 * @author tantexian<my.oschina.net/tantexian>
	 * @since 2017/3/4
	 * @params
	 */
	var idleWorker string
	var wg sync.WaitGroup

	for i, fileName := range mapFiles {
		wg.Add(1)
		log.Printf("1----####idleWorker==%v,i==%v, fileName==%v", idleWorker, i, fileName)

		go func() {
			StartCall:
			// 获取空闲worker
			idleWorker = <-registerChan
			// 函数退出时候执行
			defer wg.Add(-1)
			// 根据net/rpc[Go官方库RPC开发指南:https://my.oschina.net/tantexian/blog/851914]
			// 可知 此处rpc调用Worker结构体的DoTask方法
			log.Printf("####idleWorker==%v,i==%v, fileName==%v", idleWorker, i, fileName)
			success := call(idleWorker, "Worker.DoTask", DoTaskArgs{jobName, fileName, phase, i, n_other}, nil)
			if success == true {
				// 执行成功则继续放入到空闲worker池中
				registerChan <- idleWorker
			} else {
				fmt.Printf("Master: RPC %s DoTask error\n", idleWorker)
				// 执行失败则重新执行
				goto StartCall
			}

		}()
	}

	// 一直等到wg为0
	wg.Wait()
	fmt.Printf("Schedule: %v phase done\n", phase)
}
