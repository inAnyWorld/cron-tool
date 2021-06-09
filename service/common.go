package service

import (
	"cacheflusher/config"
	"cacheflusher/config/apollo"
	"cacheflusher/database/redis"
	"cacheflusher/service/constant"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// 获取单个结构体名称
func GetSliceStructName(i interface{}) string {
	structName := reflect.TypeOf(i).Elem().Name()
	return structName
}

// AST 获取所有的结构体名称
func StructNameWithGoFile(files string) (map[string]string, map[string]string){
	fSet := token.NewFileSet()
	// "service/structs/table/structWithTable.go"
	f, fErr := parser.ParseFile(fSet, files, nil, parser.ParseComments)
	if fErr != nil {
		log.Printf("AST 解析失败: error: [%v]", fErr)
		panic(fErr)
	}
	structMap := make(map[string]string) // 所有的表map
	structConMap := make(map[string]string)
	collectStructs := func(x ast.Node) bool {
		ts, ok := x.(*ast.TypeSpec)
		if !ok || ts.Type == nil {
			return true
		}
		// 获取结构体名称
		structName := ts.Name.Name
		structMap[structName] = structName
		s, ok := ts.Type.(*ast.StructType)
		if !ok {
			return true
		}
		for _, c := range s.Fields.List {
			if c.Doc.Text() != "" {
				splitStr := strings.Split(c.Doc.Text(), "|")
				if apollo.ApolloConfig.DEBUG == "1" {
					structConMap[structName] = strings.TrimSpace(splitStr[0])
				} else {
					structConMap[structName] = strings.TrimSpace(splitStr[1])
				}
			}
		}
		return false
	}
	ast.Inspect(f, collectStructs)
	return structMap, structConMap
}

// 上锁
func RedisLockByKey(key string) bool {
	rdb := redis.GetRedisInstance()
	execTime, _ := strconv.Atoi(CheckOfTypes(key, 2))
	log.Printf("%v 上锁锁定时间: [%v s]", key, execTime)
	setNxResult, setNxErr := rdb.SetNX(rdb.Context(), constant.LOCK_PREF + key, key, time.Second * time.Duration(execTime)).Result()
	if setNxErr != nil {
		log.Printf("tableName 上锁失败: [%v],error:[%v]", key, setNxErr)
		return false
	}
	if !setNxResult {
		log.Printf("tableName [%v]资源被占用,时间:[%v]", key, time.Now())
		return false
	}
	return true
}

// 释放锁
func FreedLockByKey(key string) bool {
	rdb := redis.GetRedisInstance()
	freed, err := rdb.Del(rdb.Context(), constant.LOCK_PREF + key).Result()
	if err != nil {
		log.Printf("key [%v] 锁释放失败", key)
	}
	if freed > 0 {
		return true
	}
	return false
}
// 返回corn表达式
// types 1,定时任务调用,2其他
func CheckOfTypes(tableName string, types int) string {
	var execTimeStr string
	tableName = strings.ToLower(tableName)
	if _, checkMap := config.DbConConfig.CornExecTime[tableName]; checkMap {
		//存在
		execTimeStr = config.DbConConfig.CornExecTime[tableName]
	}
	oldExecTime := strings.Split(execTimeStr, ":")
	if len(oldExecTime) == 2 {
		tempTime, _ := strconv.Atoi(oldExecTime[0])
		if strings.ToLower(oldExecTime[1]) == "m" {
			// 分钟
			if types == config.DISPATCH_AUTO {
				return strings.Replace("0 */corn * * * * *", "corn", strconv.Itoa(tempTime), -1)
			}
			if types == config.DISPATCH_OTHER {
				return strconv.Itoa(tempTime * 30)  // tempTime * 60 * 60 / 2
			}
		}
		if strings.ToLower(oldExecTime[1]) == "h" {
			// 小时
			if types == config.DISPATCH_AUTO {
				return strings.Replace("0 0 */corn * * * *", "corn", strconv.Itoa(tempTime), -1)
			}
			if types == config.DISPATCH_OTHER {
				return strconv.Itoa(tempTime * 60 * 30) // tempTime * 60 * 60 / 2
			}
		}

		if strings.ToLower(oldExecTime[1]) == "d" {
			// 每天固定那一刻钟执行
			if types == config.DISPATCH_AUTO {
				return strings.Replace("0 0 corn * * * *", "corn", strconv.Itoa(tempTime), -1)
			}
			if types == config.DISPATCH_OTHER {
				return strconv.Itoa(43200) // tempTime * 60 * 60 / 2
			}
		}
	}
	// 默认
	if types == config.DISPATCH_AUTO {
		return "0 0 1 * * * *"
	}
	if types == config.DISPATCH_OTHER {
		return strconv.Itoa(43200) // 86400/2
	}
	return ""
}