package dbConnect

import (
	"cacheflusher/config"
	"cacheflusher/config/apollo"
	"cacheflusher/database/sqlserver"
	"cacheflusher/service"
	"log"
	"strconv"
)

func Loading() {
	// 数据库等链接配置
	config.DbConConfig.Table, config.DbConConfig.Db = service.StructNameWithGoFile("service/structs/table/structWithTable.go")
	// 获取定时任务执行时间
	tableName := service.GetSliceStructName(&config.CronPlain{})
	db := sqlserver.MSCon.Con(apollo.ApolloConfig.CORN_DB)
	var CronPlainStruct []config.CronPlain
	db.Table(tableName).Where("IsDel = ?", 0).Find(&CronPlainStruct)
	cornMap := make(map[string]string)
	if len(CronPlainStruct) == 0 {
		log.Printf("%v 获取到的执行时间配置为空", tableName)
		for _, table := range config.DbConConfig.Table {
			cornMap[table] = "1:d" // 每天凌晨1点执行
		}
	} else {
		log.Printf("corn执行获取到配置:%v", CronPlainStruct)
		for _, cornExpr := range CronPlainStruct {
			cornMap[cornExpr.TableName] = strconv.Itoa(cornExpr.Interval) + ":" + cornExpr.IntervalType
		}
	}
	config.DbConConfig.CornExecTime = cornMap
}
