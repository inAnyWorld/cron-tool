package mysql

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"net/url"
	"sync"
	"time"
)

var MYCon MYDB

func init()  {
	MYCon = MYDB{dbs:sync.Map{}}
}

type MYDB struct {
	dbs sync.Map
}

func (mydb MYDB)Con(connString string) *gorm.DB {
	db, ok := mydb.dbs.Load(connString)
	if !ok {
		db = GetMYInstance(connString)
		mydb.dbs.Store(connString, db)
	}
	dbTemp := db.(*gorm.DB)
	return dbTemp
}

// Database 在中间件中初始化mysql链接
func GetMYInstance(connString string) *gorm.DB {
	connString, _ = url.QueryUnescape(connString)
	mysqlDb, dbErr := gorm.Open("mysql", connString)
	mysqlDb.LogMode(true)
	if dbErr != nil {
		log.Printf("mysql数据库链接失败,链接字符串: [%v],error:[%v]", connString, dbErr)
		panic("数据库类型 [mysql] 链接失败 \n")
	}
	//空闲
	mysqlDb.DB().SetMaxIdleConns(50)
	//打开
	mysqlDb.DB().SetMaxOpenConns(100)
	//超时
	mysqlDb.DB().SetConnMaxLifetime(time.Second * 60)
	return mysqlDb
}