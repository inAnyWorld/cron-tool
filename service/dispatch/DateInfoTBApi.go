package dispatch

import (
	"cacheflusher/database/redis"
	"cacheflusher/service"
	"cacheflusher/service/constant/dateInfoTBApi"
	"cacheflusher/service/structs/table"
	"cacheflusher/tools"
	"encoding/json"
	"log"
	"net/url"
	"strings"
	"time"
)

func (this *Routers) DateInfoTBApi(args ...interface{}) bool {
	tableName := service.GetSliceStructName(&table.DateInfoTBApi{})
	log.Printf("开始执行:[%v]", tableName)
	startTime := time.Now()
	var apiConfigKey [8]string
	apiConfigKey[0] = dateInfoTBApi.VOCATION_ZH_CN
	apiConfigKey[1] = dateInfoTBApi.VOCATION_ZH_HK
	apiConfigKey[2] = dateInfoTBApi.VOCATION_ZH_TW
	apiConfigKey[3] = dateInfoTBApi.VOCATION_ZH_MAC
	apiConfigKey[4] = dateInfoTBApi.FESTIVAL_ZH_CN
	apiConfigKey[5] = dateInfoTBApi.FESTIVAL_ZH_HK
	apiConfigKey[6] = dateInfoTBApi.FESTIVAL_ZH_TW
	apiConfigKey[7] = dateInfoTBApi.FESTIVAL_ZH_MAC
	rdb := redis.GetRedisInstance()
	apiAddr, _ := url.QueryUnescape(args[0].(string))
	hmSetApiArgs := make(map[string]interface{})
	for i := 0; i < len(apiConfigKey); i++ {
		uri := strings.Replace(apiAddr, "{0}", apiConfigKey[i], -1)
		responseJson := tools.GetHelper(uri)
		var mapResponse map[string]interface{}
		jsonErr := json.Unmarshal([]byte(responseJson), &mapResponse)
		if jsonErr != nil {
			log.Printf("%v json err: [%v]",tableName, jsonErr)
		}
		hmSetApiArgs[apiConfigKey[i]] = mapResponse["msg"]
	}
	_, hmErr := rdb.HMSet(rdb.Context(), tableName, hmSetApiArgs).Result()
	if hmErr != nil {
		log.Printf("DateInfoApi hmset error [%v]", hmErr)
	}
	endTime := time.Now()
	log.Printf("%v 处理耗时: [%v]", tableName, endTime.Sub(startTime))
	service.FreedLockByKey(tableName)
	return true
}
