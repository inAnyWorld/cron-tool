package sqlserver

import (
	"github.com/jinzhu/gorm"
	"log"
	"net/url"
	"sync"
	"time"
)

var MSCon MSDB

func init()  {
	MSCon = MSDB{dbs:sync.Map{}}
}

type MSDB struct {
	dbs sync.Map
}

func (msdb MSDB)Con(connString string) *gorm.DB {
	db, ok := msdb.dbs.Load(connString)
	if !ok {
		db = GetMSInstance(connString)
		msdb.dbs.Store(connString, db)
	}
	dbTemp := db.(*gorm.DB)
	return dbTemp
}


func GetMSInstance(connString string)  *gorm.DB {
	connString, _ = url.QueryUnescape(connString)
	sqlserverDb, dbErr := gorm.Open("mssql", connString)
	//sqlserverDb.LogMode(true)
	if dbErr != nil {
		log.Printf("sqlserver数据库链接失败,链接字符串: [%v],error [%v]", connString, dbErr)
		panic("数据库类型 [sqlserver] 链接失败\n")
	}
	//设置连接池
	//空闲
	sqlserverDb.DB().SetMaxIdleConns(50)
	//打开
	sqlserverDb.DB().SetMaxOpenConns(100)
	//超时
	sqlserverDb.DB().SetConnMaxLifetime(time.Second * 60)
	return sqlserverDb
}


