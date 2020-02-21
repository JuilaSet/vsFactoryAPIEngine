# vsFactoryAPIEngine - lsr
## API 手册
* 注册服务
> {"servername" : "mongodb-service", "serverURL" : "mongodb://127.0.0.1:27017"}
* 检测登录：
loginCheck
* 插入数据集
> {"collectionName": "user", "data": { "name" : "tom", "dept": "CS" }}
* 查询
> {"collectionName": "user", "query": { "name" : "tom" }}
* 删除数据集
> {"collectionName": "user", "query": { "name" : "tom" }}
* 更新数据
> {"collectionName": "user", "query": { "name" : "tom" }, "update": { "name" : "tommy"}}
