# govirt

> libvirt控制器

## grpc
https://grpc.io/docs/languages/go/quickstart/
```
protoc -I=proto --go_out=proto --go_opt=paths=source_relative --go-grpc_out=proto --go-grpc_opt=paths=source_relative proto/*
```

## 配置文件示例

`.env`

```env
app_env=local
app_key=alskdjfakjsjashbdfas
app_debug=true
app_url=http://localhost:3000
app_log_level=debug
app_port=8000

db_connection=mysql
db_host=192.168.6.166
db_port=3306
db_database=govirt
db_username=govirt
db_password=aslkashfdsa
db_debug=2

log_type=daily
log_level=debug
```