# emqx 配置
侧边栏的选项有：
监控
    集群概览-能够管理集群
    客户端-能够看到所有的ClientID
    订阅管理-能够看到所有的Topic信息
访问控制
    在这里设置认证的方式，比如使用HTTP认证
    设置禁用（黑名单）等
管理
    MQTT基础配置
    监听器
    日志
    延迟发布
    网关
数据集成
    设置在线/离线事件的WebHook，首先设置规则，接着设置规则桥接的目标服务，再把规则和桥接关联起来
    在Flows里面能够看到规则和桥接的关联图
问题分析
    通常不会遇到问题的
系统设置
    创建非管理员新用户
    设置API密钥

### emqx 在线离线WebHook配置
在左侧边栏找到，【数据集成】，首先添加规则名称
device_online
和对应的SQL语句
SELECT
*
FROM
"$events/client_connected"
离线规则：
device_offline
和对应的SQL语句
SELECT
*
FROM
"$events/client_disconnected"

接着添加数据桥接：
set_device_online
POST http://host.docker.internal:30070/v1/mqtt/emqx/setDeviceOnline
content-type application/json 以及其它鉴权的头信息
和消息体（我写了个测试用例能自动生成消息体）：
{

}
离线的数据桥接
set_device_offline
POST http://host.docker.internal:30070/v1/mqtt/emqx/setDeviceOffline
content-type application/json 其它鉴权信息
请求体
{

}

在规则中添加动作 使用数据桥接转发 数据桥接 webhook:set_device_online
测试，这时将得到连接的请求：
[GIN] 2023/04/26 - 14:33:38 | 200 |    4.871166ms |       127.0.0.1 | POST     "/v1/mqtt/emqx/setDeviceOffline"
以及连接断开的请求：
[GIN] 2023/04/26 - 14:33:38 | 200 |    6.209042ms |       127.0.0.1 | POST     "/v1/mqtt/emqx/setDeviceOnline"
