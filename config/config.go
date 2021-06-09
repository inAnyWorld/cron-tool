package config

const (
	// 执行方式
	DISPATCH_AUTO  = 1 // 定时任务
	DISPATCH_OTHER = 2 // 初始化或其他

	// Redis
	REDIS_ADDR       = "REDIS_ADDR"
	REDIS_DB         = "REDIS_DB"
	REDIS_DBDEFAULT  = 14
	REDIS_PWD        = "REDIS_PWD"
	CORN_DB          = "CORN_DB" // 定时任务sqlserver链接库
	CORN_DB_DEFAULT  = "db host"
	DEBUG			 = "DEBUG"
	DEBUG_DEFAULT	 = 1
)
// db
var DbConConfig DbConfig

type DbConfig struct {
	Db map[string]string // 所有缓存的表对应的数据库链接
	Table map[string]string // 所有缓存的表
	CornExecTime map[string]string // 定时任务执行时间
}

// 手动刷新缓存
type TbConfig struct {
	Tb string `json:"tb"`
}

// corn config
type CronPlain struct {
	Id 				int 		`json:"Id" gorm:"column:Id"`
	TableName 		string 		`json:"TableName" gorm:"column:TableName"`
	Interval 		int 		`json:"Interval" gorm:"column:Interval"`
	IntervalType 	string 		`json:"IntervalType" gorm:"column:IntervalType"`
	IsDel 			int 		`json:"IsDel" gorm:"column:IsDel"`
	CreateTime 		int 		`json:"CreateTime" gorm:"column:CreateTime"`
	UpdateTime 		int 		`json:"UpdateTime" gorm:"column:UpdateTime"`
	BusinessType 	string 		`json:"BusinessType" gorm:"column:BusinessType"`
}

// apollo
type ApolloConfig struct {
	REDIS_ADDR     string // redis 链接地址
	REDIS_PWD      string // redis 密码
	REDIS_DB       int    // redis 库
	CORN_DB		   string // 定时任务保存的数据库
	DEBUG		   string // 定时任务保存的数据库
}

// yaml
type YamlConfig struct {
	Apollo YamlApollo
}

type YamlApollo struct {
	Metaaddr  string `yaml:"metaaddr"`
	Cluster   string `yaml:"cluster"`
	Appid     string `yaml:"appid"`
	Namespace string `yaml:"namespace"`
	Cachedir  string `yaml:"cachedir"`
}