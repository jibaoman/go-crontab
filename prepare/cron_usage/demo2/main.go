package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)


//代表一个任务
type ConrJob struct {
	expr *cronexpr.Expression
	nextTime time.Time
}

func main() {
	var (
		cronJob *ConrJob
		expr *cronexpr.Expression
		now time.Time
		scheduleTable map[string]*ConrJob //调度表
	)

	scheduleTable = make(map[string]*ConrJob)

	//当前时间
	now = time.Now()

	//定义两个cronjob
	expr = cronexpr.MustParse("*/5 * * * * * *")
	cronJob = &ConrJob{
		expr:     expr,
		nextTime: expr.Next(now),
	}
	//任务注册到调度表
	scheduleTable["job1"] = cronJob

	expr = cronexpr.MustParse("*/5 * * * * * *")
	cronJob = &ConrJob{
		expr:     expr,
		nextTime: expr.Next(now),
	}

	scheduleTable["job2"] = cronJob


	//启动一个调度协程
	go func() {
		var (
			jobName string
			cronJob *ConrJob
			now time.Time
		)

		//定时检查下任务调度表
		for {
			now = time.Now()

			for jobName,cronJob = range scheduleTable {

				//判断是否过期
				if cronJob.nextTime.Before(now) || cronJob.nextTime.Equal(now) {
					//启动一个协程执行任务
					go func(jobName string) {
						fmt.Println("执行：",jobName)
					}(jobName)

					//计算下一次调度
					cronJob.nextTime = cronJob.expr.Next(now)
					fmt.Println(jobName,"下次执行时间",cronJob.nextTime)
				}
			}

			select {
			case <-time.NewTimer(100 * time.Millisecond).C: //
			}
		}
	}()

	time.Sleep(10 * time.Second)
}
