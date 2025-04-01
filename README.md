# GoHub

> 这是一个开发模板项目，基于[Summer大佬的Gohub教学项目](https://github.com/summerblue/gohub)

## 更改部分

- 没有做手机号注册的完整逻辑
- 没有添加限流功能
- 优化短信发送逻辑，支持配置文件选择TLS方式连接SMTP服务器

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
db_database=gohub
db_username=gohub
db_password=aslkashfdsa
db_debug=2

log_type=daily
log_level=debug

redis_host=192.168.6.166
redis_password=alsdkjfhkasd

MAIL_HOST=smtp.qq.com
MAIL_PORT=465
MAIL_FROM_ADDRESS=horonlee@foxmail.com
MAIL_FROM_NAME=GoHub
MAIL_USERNAME=horonlee@foxmail.com
MAIL_PASSWORD=paosdfalsdkjf
MAIL_TLS=true
```