## 自动缓存

### 使用

定义好业务相关struct,在dispatch实现对应逻辑即可


### 目录结构

|-- config

|-- -- apollo apollo配置

|-- -- dbConnect 数据库链接配置

|-- -- yaml yaml配置

|-- -- config.go 缓存基础配置

|-- database

|-- -- mysql

|-- -- redis

|-- -- sqlserver

|-- docs 说明文档

|-- -- README.md readme

|-- flusher

|-- -- refresh.go 缓存系统核心实现

|-- service

|-- -- constant 业务常量配置

|-- -- dispathch 业务逻辑实现

|-- -- structs 业务struct配置

|-- -- common.go 业务助手函数

|-- tools

|-- -- common.go 系统助手函数

|-- Dockerfile

|-- main.go 启动文件




### 系统配置


#### 定时任务执行数据库

db host

#### 定时任务调度表

|  字段   | 类型  |空|默认|注释|
|  ----  | ----  |----|----|----|
| Id  | int |NOT NULL|0|自增主键
| TableName  | string |NOT NULL | ''|表名|
| Interval  | int |NOT NULL | 0|执行时间间隔|
| IntervalType  | string |NOT NULL | ''|m分钟,h小时,d默认凌晨1点执行|
| IsDel  | int |NOT NULL | 0|是否删除,0不删除,1删除|
| CreateTime  | int |NOT NULL | 0|创建时间|
| CreateTime  | int |NOT NULL | 0|更新时间|
| BusinessType  | string |NOT NULL | 0|业务类型说明|


#### apollo

|  key   |注释|
|  ----  | ----|
| REDIS_ADDR  | redis 连接字符串 |
| REDIS_DB  | redis 连接库 |
| REDIS_PWD  | redis 密码 |
| CORN_DB  | 定时任务配置表 |
