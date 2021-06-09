package dispatch

import (
	"cacheflusher/database/redis"
	"cacheflusher/database/sqlserver"
	"cacheflusher/service"
	"cacheflusher/service/structs/table"
	"encoding/json"
	"log"
	"strings"
	"time"
)

func (this *Routers) CooperationTB(args ...interface{}) bool {
	startTime := time.Now()
	tableName := service.GetSliceStructName(&table.CooperationTB{})
	log.Printf("开始执行:[%v]", tableName)

	db := sqlserver.MSCon.Con(args[0].(string))
	rdb := redis.GetRedisInstance()
	var CooperationTBStructS []table.CooperationTB
	hmSetCooperationTBArgs := make(map[string]interface{})
	// 获取所有
	db.Table(tableName).Find(&CooperationTBStructS)
	for _,repo := range CooperationTBStructS {
		setRedis, jsonErr := json.Marshal(repo)
		if jsonErr != nil{
			log.Printf("%v Json 转换错误: [%v]\n", tableName, jsonErr)
		}
		hmSetCooperationTBArgs[strings.ToLower(repo.Cooperation)] = setRedis
	}
	_, hmErr := rdb.HMSet(rdb.Context(), tableName, hmSetCooperationTBArgs).Result()
	if hmErr != nil {
		log.Printf("%v hmset error [%v]", tableName, hmErr)
	}
	defer db.Close()
	endTime := time.Now()
	log.Printf("%v 处理耗时: [%v]", tableName, endTime.Sub(startTime))
	// 释放锁
	service.FreedLockByKey(tableName)
	return true
}
