package tools

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/robfig/cron/v3"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"sort"
	"time"
)

// url:         请求地址
// response:    请求返回的内容
func GetHelper(url string) string {
	// http 超时可能会引发 error: context deadline exceeded (Client.Timeout exceeded while awaiting headers)
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}

	return result.String()
}

// url:         请求地址
// data:        POST请求提交的数据
// contentType: 请求体格式,如:application/json
// content:     请求返回的内容
func PostHelper(url string, data interface{}, contentType string) string {
	client := &http.Client{Timeout: 30 * time.Second}
	jsonStr, _ := json.Marshal(data)
	resp, err := client.Post(url, contentType, bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	return string(result)
}

// 字符串是否在字符数组中
func StrInArrHelper(target string, strArr []string) bool {
	sort.Strings(strArr)
	index := sort.SearchStrings(strArr, target)
	// index的取值：[0,len(str_array)]
	// 需要注意此处的判断，先判断 &&左侧的条件,如果不满足则结束此处判断,不会再进行右侧的判断
	if index < len(strArr) && strArr[index] == target {
		return true
	}
	return false
}

// base64编码
func Base64EncodeHelper(encodeStr string ) string {
	base64Str := encodeStr
	base64Byte := []byte(base64Str)
	return base64.StdEncoding.EncodeToString(base64Byte)
}

// 解决corn win/linux 兼容问题
func NewWithSecondHelper() *cron.Cron {
	secondParser := cron.NewParser(cron.Second | cron.Minute |
		cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	return cron.New(cron.WithParser(secondParser), cron.WithChain())
}

func StructToMapHelper(obj interface{}) interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)
	var sMap = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		sMap[obj1.Field(i).Name] = obj2.Field(i).Interface()
	}
	jsonStr, jsonErr := json.Marshal(sMap)
	if jsonErr != nil {
		log.Printf("apollo apToJson err: [%v]", jsonErr)
	}
	return string(jsonStr)
}
