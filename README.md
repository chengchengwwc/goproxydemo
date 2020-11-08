### Reverse Proxy功能点
1. 更改内容支持
2. 错误信息回调
3. 支持自定义负载均衡
4. URL重写功能
5. 连接池
6. 支持websocket
7. 支持HTTPS

#### 为反向代理增加负载均衡功能
- 使用工厂方法拓展
- 使用接口统一封装

#### 熔断和降级的意义
- 熔断的意义：熔断器是当前依赖服务已经出现故障时，为了保证自生服务的正常运行不再访问依赖服务，防止
雪崩
- 降级的意义：当服务器压力巨增时，根据业务策略降级，以此释放服务资源保证业务正常。
- hystrix-go 熔断，降级，限流集成类库

#### 重新认识HTTPS和HTTP2
1. 证书创建
- CA私钥 openssl genrsa -out ca.key 2048
- CA数据证书 openssl req -x509 -new -nodes -key ca,key -subj "/CN=example1.com" -days 5000 -out ca.crt
- 服务器私钥: openssl genrsa -out server.key 2048
- 服务器证书签名请求： openssl req -new -key server.key -subj="/CN=example1.com" -out server.csr
- 上面两个生成服务器证书：openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 5000

#### TCP代理
本质上是7层反向代理，只是代理的内容是TCP协议包
- 初始化TCP服务器
- 创建上游链接
- 创建下游链接
- 上下游数据交换

TCP代理实现
- 参照 http.util.ReverseProxy实现，服务和代理逻辑分离




