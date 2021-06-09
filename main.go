package main

import (
	"cacheflusher/config"
	"cacheflusher/config/apollo"
	"cacheflusher/config/dbConnect"
	"cacheflusher/config/yaml"
	"cacheflusher/flusher"
	"encoding/json"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"io"
	"net/http"
	_ "net/http/pprof"
)

func main()  {
	// 加载配置
	yaml.Loading()
	apollo.Loading()
	dbConnect.Loading()
	go flusher.InitRefresh() // 部署时执行一次
	flusher.AutoRefresh()
	http.HandleFunc("/cacheFlusher", cacheFlusher)
	http.HandleFunc("/helloWorld", helloWorld)
	_ = http.ListenAndServe("0.0.0.0:8080", nil)
}

// 手动刷新缓存
func cacheFlusher(writer http.ResponseWriter,  request *http.Request)  {
	var flushJson config.TbConfig
	if flushErr := json.NewDecoder(request.Body).Decode(&flushJson); flushErr != nil {
		_ = request.Body.Close()
		fmt.Printf("flushErr Error:[%v]\n", flushErr)
	}
	fmt.Printf("接口更新缓存 cacheFlusherApi: [%v] result is %v\n", flushJson.Tb, flusher.ManualRefresh(flushJson.Tb))
	writer.Header().Set("Content-type", "application/text")
	writer.WriteHeader(200)
	_, _ = io.WriteString(writer, `{"code": "200","message": "success"}`)
}

// k8s 部署检查回调地址
func helloWorld(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("content-type", "application/json")
	writer.WriteHeader(200)
	_, _ = io.WriteString(writer, `{"code": "200","message": "hello world"}`)
}