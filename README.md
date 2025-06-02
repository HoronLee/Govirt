# Govirt

> libvirt控制器

## 主要功能

支持通过 http 请求来对 libvirtd 进行各项操作，API 文档：[APIFOX 文档](https://kx5a5itjlt.apifox.cn)
当前可对虚拟机，存储池，存储卷和网络进行管理，默认有两个存储池和两个网络
两个存储池分别存储可复用的虚拟机镜像（ISO 或者 qcow2）,还有一个用于存储虚拟机的磁盘
两个网络分别用于内部通信和外部直接访问（功能暂未做完）
虚拟机可以通过 ISO 和 qcow2 来安装，可自行设置 vnc 的端口和密码

## TODO

- 内部集成 websockify 来对虚拟机的 vnc 进行转发，并且设置权限隔离
- ...

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

CON_URI=qemu:///system
```
