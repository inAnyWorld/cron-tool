package apollo

import (
	"cacheflusher/config"
	"cacheflusher/config/yaml"
	"cacheflusher/tools"
	"encoding/json"
	"fmt"
	"github.com/manucorporat/try"
	"github.com/philchia/agollo/v3"
	"log"
	"strconv"
)

var ApolloConfig config.ApolloConfig

// 加载配置项
func Loading() {
	var c = NewApolloConfg()
	fmt.Printf("yaml config: [%v]", c)
	_ = agollo.StartWithConf(c)
	ApolloConfig.REDIS_ADDR  = agollo.GetStringValue(config.REDIS_ADDR, "") // Redis链接
	ApolloConfig.REDIS_DB, _ = strconv.Atoi(agollo.GetStringValue(config.REDIS_DB, strconv.Itoa(config.REDIS_DBDEFAULT)))
	ApolloConfig.REDIS_PWD   = agollo.GetStringValue(config.REDIS_PWD, "")
	ApolloConfig.CORN_DB = agollo.GetStringValue(config.CORN_DB, config.CORN_DB_DEFAULT) // 定时任务配置
	ApolloConfig.DEBUG = agollo.GetStringValue(config.DEBUG, strconv.Itoa(config.DEBUG_DEFAULT)) // DEBUG
	log.Printf("apollo获取到配置:%v", tools.StructToMapHelper(ApolloConfig))
	go monitorConfigChange()
}

func monitorConfigChange() {
	events := agollo.WatchUpdate()
	for {
		changeEvent := <-events
		try.This(func() {
			ApolloConfig.REDIS_ADDR  = changeEvent.Changes[config.REDIS_ADDR].NewValue
			ApolloConfig.REDIS_DB, _ = strconv.Atoi(changeEvent.Changes[config.REDIS_DB].NewValue)
			ApolloConfig.REDIS_PWD   = changeEvent.Changes[config.REDIS_PWD].NewValue
			ApolloConfig.CORN_DB     = changeEvent.Changes[config.CORN_DB].NewValue
			ApolloConfig.DEBUG       = changeEvent.Changes[config.DEBUG].NewValue
		}).Finally(func() {

		}).Catch(func(e try.E) {
			// Print crash
		})

		bytes, _ := json.Marshal(changeEvent.Changes)
		fmt.Println("apollo 配置发生改变 event:", string(bytes))
	}
}

func NewApolloConfg() *agollo.Conf {
	var apolloConf agollo.Conf
	apolloConf.AppID = yaml.Yaml.Apollo.Appid
	apolloConf.NameSpaceNames = append(apolloConf.NameSpaceNames, yaml.Yaml.Apollo.Namespace)
	apolloConf.MetaAddr = yaml.Yaml.Apollo.Metaaddr
	apolloConf.Cluster = yaml.Yaml.Apollo.Cluster
	apolloConf.CacheDir = yaml.Yaml.Apollo.Cachedir
	return &apolloConf
}