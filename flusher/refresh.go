package flusher

import (
	"cacheflusher/config"
	"cacheflusher/service"
	"cacheflusher/service/dispatch"
	"github.com/gorhill/cronexpr"
	"log"
	"reflect"
	"time"
)

type CronJob struct {
	expr *cronexpr.Expression
	nextTime time.Time
}

// 定义控制器函数Map类型，便于后续快捷使用
type ControllerMapsType map[string]reflect.Value

var (
	ControllerMaps ControllerMapsType // 声明控制器函数Map类型变量
	expr *cronexpr.Expression
	now time.Time
	cronJob *CronJob
	cronSchedule map[string]*CronJob
)

func InitRefresh() {
	for t := range config.DbConConfig.Table {
		run(config.DbConConfig.Table[t], config.DISPATCH_OTHER)
	}
}

// 计划任务配置规则
// 初步分为三种
// 1.间隔多少分钟执行一次
//		最低每分钟执行一次
// 2.间隔多少小时执行一次
//		最长 24小时
// 3.不选择则默认凌晨1点执行一次
func AutoRefresh() {
	log.Printf("定时任务开始执行,时间:[%v] \n", time.Now())
	cronSchedule = make(map[string]*CronJob)
	now = time.Now()
	for k := range config.DbConConfig.Table {
		execCorn := service.CheckOfTypes(config.DbConConfig.Table[k], 1)
		log.Printf("[%v] corn表达式:[%v]", config.DbConConfig.Table[k], execCorn)
		expr = cronexpr.MustParse(execCorn)
		cronJob = &CronJob{
			expr: expr,
			nextTime:expr.Next(now),
		}
		cronSchedule[config.DbConConfig.Table[k]] = cronJob
	}

	// 启动一个调度协程
	go func() {
		var (
			jobName string
			cronJob *CronJob
			_now time.Time
		)
		// 定时检查一下任务调度表
		for {
			_now = time.Now()
			for jobName, cronJob = range cronSchedule {
				// 判断是否过期
				if cronJob.nextTime.Before(_now) || cronJob.nextTime.Equal(_now) {
					// 启动一个协程, 执行这个任务
					go func(jobName string) {
						log.Printf("[%v]本次开始执行时间:[%s]", jobName, time.Now())
						run(jobName, config.DISPATCH_AUTO)
					}(jobName)
					// 计算下一次调度时间
					cronJob.nextTime = cronJob.expr.Next(_now)
					log.Printf("[%v]下次执行时间:[%v]", jobName, cronJob.nextTime)
				}
			}
			select {
				case <-time.NewTimer(100 * time.Millisecond).C: // 睡眠
			}
		}
	}()
}

func ManualRefresh(tableName string) bool {
	return run(tableName, config.DISPATCH_OTHER)
}


func run(tableName string, types int) bool {
	if types == config.DISPATCH_AUTO {
		if !service.RedisLockByKey(tableName) {
			return false
		}
	}
	if _, err := CallFunc(tableName, config.DbConConfig.Db[tableName]); err != nil {
		log.Printf("call err: [%v]", err)
		return false
	}
	return true
}

func CallFunc(tableName string, args ... interface{}) (result []reflect.Value, err error) {
	var router dispatch.Routers
	ControllerMap := make(ControllerMapsType, 0)
	rf := reflect.ValueOf(&router)
	rft := rf.Type()
	funcNum := rf.NumMethod()
	for i := 0; i < funcNum; i ++ {
		mName := rft.Method(i).Name
		ControllerMap[mName] = rf.Method(i)
	}
	parameter := make([]reflect.Value, len(args))
	for k, arg := range args {
		parameter[k] = reflect.ValueOf(arg)
	}
	result = ControllerMap[tableName].Call(parameter)
	return 
}