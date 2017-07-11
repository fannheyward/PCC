[PCC][1] 高可用架构 PCC 性能挑战赛

[1]:https://github.com/archnotes/PCC

* Golang + Redis 实现
* 点赞关系用 Sorted Set 存储，时间戳作为 score
* 好友关系用 Sorted Set 存储，时间戳作为 score
* 用户／对象信息 Key-Value 存储

测试：

```
清空／创建 Redis 实例
./bin/boot.sh

生成 100用户+100对象，其中 1% 的用户有 50 个好友，5% 20个好友，其他用户只有 10个好友
curl -X POST "http://127.0.0.1:9009/test/gen?user=100&object=100"

curl "http://127.0.0.1:9009/pcc?action=like&oid=1&uid=1"
curl "http://127.0.0.1:9009/pcc?action=is_like&oid=1&uid=1"

MBP 本机压测 is_like 能达到 100K QPS，距离 300K 还有不少差距
```

改进：

* SSDB／LevelDB 做持久化存储
* 加上负载均衡的话，理论值是能支撑 300K，需要进一步验证