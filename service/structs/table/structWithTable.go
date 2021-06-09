package table

import "time"

type CooperationTB struct {
	// 测试环境db host || 生产环境db host
	Cooperation 	string 		`json:"Cooperation" gorm:"column:Cooperation"`
	DateBegin 		time.Time	`json:"DateBegin" gorm:"column:DateBegin"`
	DateEnd 		time.Time	`json:"DateEnd" gorm:"column:DateEnd"`
	IpAddress 		string		`json:"IpAddress" gorm:"column:IpAddress"`
	Status 			int			`json:"Status" gorm:"column:Status"`
	Desc 			string		`json:"Desc" gorm:"column:Desc"`
	AuthorizeKey 	string		`json:"AuthorizeKey" gorm:"column:AuthorizeKey"`
	Exclude 		string		`json:"Exclude" gorm:"column:Exclude"`
	Qps 			int64		`json:"Qps" gorm:"column:Qps"`
}

type DateInfoTBApi struct {
	// 测试环境db host || 生产环境db host
	DateInfoTBApi 	string 		`json:"DateInfoTBApi" gorm:"column:DateInfoTBApi" default:""`
}