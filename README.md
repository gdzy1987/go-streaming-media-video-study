# go-streaming-media-video-study
Go语言实战流媒体视频网站 学习笔记
## 目录
```
├── api             api
├── common
├── config
├── doc
├── logger
├── prepare 
├── scheduler       scheduler
├── streamserver    streamserver
├── templates       templates html页面
└── web             web
```

## 前置
- mysql创建数据库，导入表
```
create database video_server;
use video_server;
```
```
doc/initdb.sql
```
- 配置阿里云OSS
```
endpoint
key
secret
bucket
```
## 部署，执行
- 执行脚本
```
git clone xx.git
build.sh mac|linux
deploy.sh
```
- 访问
```
api             localhost:8000
streamserver    localhost:9000
scheduler       localhost:9001
web             localhost:8080
```
- 测试
```
localhost:8080
```
### 知识点
- mysql "github.com/go-sql-driver/mysql"
- session
- UUID 
- proxy转发
- 阿里云OSS使用 github.com/aliyun/aliyun-oss-go-sdk/oss