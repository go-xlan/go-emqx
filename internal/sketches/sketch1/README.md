# emqx 把上线和下线消息发给 http 服务端

use mqtt [emqx](https://github.com/emqx/emqx)

emqx 旧版的 8081 端口已被合并至 新版的 18083 端口

```bash
docker run -d --name emqx -p 1883:1883 -p 8083:8083 -p 8084:8084 -p 8883:8883 -p 18083:18083 emqx/emqx
```

访问 emqx 控制台
http://localhost:18083

创建连接器
http://127.0.0.1:18083/#/connector

创建规则-在命中规则时走连接器
http://127.0.0.1:18083/#/rule/rules

接着启动服务
```bash
go run main.go
```
即可监听到鉴权请求和上下线请求
