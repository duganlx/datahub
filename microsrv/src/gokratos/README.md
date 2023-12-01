# go micro server

## nacos

注册中心
- [x] 理想情况下的服务注册与访问（grpc、http）
- [ ] 注册服务的访问权限控制
- [ ] 服务可用性监测

配置中心
- [x] 服务启动时读取 nacos 配置
- [ ] 服务运行中实时同步最新的 nacos 配置


## 实践

---

客户端每次运行前都需要将 cache/naming 中的内容删除，否则无法启动。报错信息为: rpc error: code = DeadlineExceeded desc = context deadline exceeded

---

客户端通过 nacos 去调用服务器的http方式接口时，会出现问题 `code = 503 reason = NODE_NOT_FOUND message = error: code = 503 reason = no_available_node message =  metadata = map[] cause = <nil> metadata = map[] cause = <nil>`，问题定位如下。解决办法有两种，第一种是采用grpc去访问（推荐）；第二种是手动取服务节点转换成最终的url(比如 `http://127.0.0.1:8000`)去访问。
```go
// client/main.go
conn, err := transhttp.NewClient( // NewClient returns an HTTP client.
  context.Background(),
  transhttp.WithEndpoint("discovery:///srv1.http"), // 服务名
  transhttp.WithDiscovery(r), // r: nacos registry
)

// kratos/v2@v2.7.1/transport/http/client.go
func NewClient(ctx context.Context, opts ...ClientOption) (*Client, error) {
  // ...
  // options {discovery: r (nacos registry above), block: false}
  // target {Scheme: "discovery", Endpoint: "srv1.http"}
  selector := selector.GlobalSelector().Build()
  var r *resolver
  if options.discovery != nil {
		if target.Scheme == "discovery" {
			if r, err = newResolver(ctx, options.discovery, target, selector, options.block, insecure, options.subsetSize); err != nil {
				return nil, fmt.Errorf("[http client] new resolver failed!err: %v", options.endpoint)
			}
		} 
	}
  return &Client{
		opts:     options,
		target:   target,
		insecure: insecure,
		r:        r,
		cc: &http.Client{
			Timeout:   options.timeout,
			Transport: options.transport,
		},
		selector: selector,
	}, nil
}

// kratos/v2@v2.7.1/transport/http/resolver.go
func newResolver(ctx context.Context, discovery registry.Discovery, target *Target,
	rebalancer selector.Rebalancer, block, insecure bool, subsetSize int,
) (*resolver, error) {
	watcher, err := discovery.Watch(ctx, target.Endpoint) // target.Endpoint = srv1.http
	r := &resolver{
		target:      target,
		watcher:     watcher,
		rebalancer:  rebalancer, // assign directly: selector.GlobalSelector().Build()
		insecure:    insecure,
		selecterKey: uuid.New().String(),
		subsetSize:  subsetSize,
	}
	go func() {
		for {
			// Watcher.Next returns services in the following two cases:
			// 1.the first time to watch and the service instance list is not empty.
			// 2.any service instance changes found.
			// if the above two conditions are not met, it will block until context deadline exceeded or canceled
			// 这是 watcher.Next() 的官方说明，所以会阻塞在该行
			services, err := watcher.Next()
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				log.Errorf("http client watch service %v got unexpected error:=%v", target, err)
				time.Sleep(time.Second)
				continue
			}
			r.update(services) // 如果能拿到services就更新resolver
		}
	}()
	return r, nil
}

func (r *resolver) update(services []*registry.ServiceInstance) bool {
	// ServiceInstance{ID, Name, Version, Metadata: map[string]string, Endpoints: []string}
	// Rebalancer is nodes rebalancer.
	// Rebalancer.Apply is apply all nodes when any changes happen
	r.rebalancer.Apply(nodes) 
	// 将 services 中的 Endpoints 转换成 url，接着转换成 nodes，进行应用，均衡器会选择某个node进行访问，
	// 但是在上面调用处已经阻塞，根本不会有nodes注册到 rebalancer 中。这导致在选择node访问时，其数组为空，
	// 导致报错。
	return true
}

// kratos/v2@v2.7.1/transport/http/http
// 调用栈
// - client.SayHello(context.Background(), &v1.HelloRequest{Name: "http yes!"})
// - err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
// - client.invoke(ctx, req, args, reply, c, opts...)
// - res, err := client.do(req.WithContext(ctx))
func (client *Client) do(req *http.Request) (*http.Response, error) {
	var done func(context.Context, selector.DoneInfo)
	if client.r != nil {
		var (
			err  error
			node selector.Node
		)
		// Selector is node pick balancer.
		// Selector.Select nodes. if err == nil, selected and done must not be empty.
		if node, done, err = client.selector.Select(req.Context(), selector.WithNodeFilter(client.opts.nodeFilters...)); err != nil {
			// 报错 reason = NODE_NOT_FOUND message = error: code = 503 reason = no_available_node message =  metadata = map[] cause = <nil> metadata = map[] cause = <nil>
			return nil, errors.ServiceUnavailable("NODE_NOT_FOUND", err.Error())
		}
	}
}
```

---

目前存在两个需求点：1. 资产单元的权限管理; 2. 微服务架构下服务的权限管理;

*资产单元的权限管理*，即为 right用户的right模型在right资产单元进行下单交易。那么就需要考虑如下几个问题：

1. 如何保证 right用户? 即资产单元允许哪些用户进行访问（资产单元的权限管理）
2. 如何保证 right模型? 模型是用户自己创建的，是否为正确的模型是由用户进行管理，平台可以提供一套机制协助进行管理。

casbin提供了RBAC的权限设计方案，可以将用户作为sub，资产单元作为obj，资产单元访问方式作为act；另外还需要有`角色/组`的概念，`角色/组` 与资产单元进行绑定，表示该`角色/组`可以访问哪些资产单元，而用户可以跟`角色/组`进行绑定，表示该用户可以访问该`角色/组`中所有资产单元。这样casbin就帮助我们解决第一个问题了（如何保证 right用户?）；第二个问题由访问令牌进行解决，用户自行选择对可访问的资产单元并输入appid和expires后，生成对应的appsecret，然后在模型登录时提供appid 和 appsecret 进行鉴权即可。

- authcode{id, appid, appsecret, expires, aucodes, allow, userid}

整个流程如下
1. 管理员先配置好用户可访问的资产单元列表
1. 用户生成某个资产单元（au1）的访问令牌appsecret，另外也可以生成能访问所有可访问资产单元（*）的令牌appsecret
1. 用户编写的模型要访问资产单元进行下单前，需要提供appid, appsecret, aucode进行鉴权，鉴权通过之后会生成jwt信息
1. 后续访问携带jwt信息访问


**实验测试**

管理员设置了资产单元的访问规则如下。为了更好的进行资产单元的管理，将资产单元按组为单位进行划分后，再将整个组分配给特定用户。

```text
# 解读：用户ww可以访问 MANAGER_WW组中的资产单元
用户ww: MANAGER_WW组 = {0148P1016_ww, 88853899_ww, DRWZQ1ZT_03}
# 解读：用户xjw可以访问 PRODUCT_EAMLS1组中的资产单元
用户xjw: PRODUCT_EAMLS1组 = {300016, 88853899_ww, EAMLS1ZT_00, EAMLS1ZTX_00}
# 解读：用户wsy可以访问资产单元DRW001ZTX_04
用户wsy: DRW001ZTX_04
# 解读：用户yrl可以访问 MANAGER_WW组中的资产单元
用户yrl: MANAGER_WW组 
```

场景设计如下，样例中所说的成功/失败表示预期的鉴权结果（成功：鉴权通过；失败：鉴权不通过）

1. 用户ww生成*只能*访问资产单元`[0148P1016_ww]`的访问令牌，并访问资产单元`0148P1016_ww` —— 成功
2. 用户ww生成*只能*访问资产单元`[0148P1016_ww, 88853899_ww]`的访问令牌，并访问资产单元`0148P1016_ww` —— 成功
3. 用户ww生成*只能*访问资产单元`[0148P1016_ww]`的访问令牌，并访问资产单元`88853899_ww` —— 失败，该令牌没有访问`88853899_ww`的权限
4. 用户ww生成可访问*所有*资产单元（`MANAGER_WW组`）的访问令牌，并访问资产单元`88853899_ww` —— 成功
5. 用户ww生成可访问*所有*资产单元（`MANAGER_WW组`）的访问令牌，并访问资产单元`EAMLS1ZT_00` —— 失败，用户ww没有访问`EAMLS1ZT_00`的权限
6. 用户ww生成*不能*访问资产单元`[0148P1016_ww]`的访问令牌，并访问资产单元`0148P1016_ww` —— 失败
7. 用户ww生成*不能*访问资产单元`[0148P1016_ww]`的访问令牌，并访问资产单元`88853899_ww` —— 成功，`88853899_ww` 在 MANAGER_WW组中，但不在不能访问的列表中
8. 用户ww生成*不能*访问资产单元`[0148P1016_ww]`的访问令牌，并访问资产单元`EAMLS1ZT_00` —— 失败，用户ww没有访问`EAMLS1ZT_00`的权限
9. 用户wsy生成*只能*访问资产单元`[DRW001ZTX_04]`的访问令牌，并访问资产单元`DRW001ZTX_04` —— 成功
10. 用户xjw生成*只能*访问资产单元`[EAMLS1ZT_00]`的访问令牌，并访问资产单元`EAMLS1ZT_00` —— 成功


*微服务架构下服务的权限管理*

用nacos作为服务注册&发现中心，各个`资产单元`和`用户中心`都会在nacos进行注册。当用户的某个模型需要在资产单元`Au1`中下单时，带上 appid 和 appsecret，`Au1`会去访问`用户中心`的接口进行鉴权，当鉴权通过之后，则将该对appid和appsecret保存在内存中，下次如果再遇到该对时就不用再去用户中心鉴权而直接放行。

## 参考

- [SSL/TLS协议运行机制的概述](https://www.ruanyifeng.com/blog/2014/02/ssl_tls.html)
- [grpc-auth-support.md(grpc-go Documentation)](https://github.com/grpc/grpc-go/blob/master/Documentation/grpc-auth-support.md)
- [gRPC authentication guide](https://grpc.io/docs/guides/auth/)
- [jwt在线解密](https://www.box3.cn/tools/jwt.html)
- [Casbin 文档](https://casbin.org/zh/docs/overview)
- [Basic Role-Based HTTP Authorization in Go with Casbin](https://zupzup.org/casbin-http-role-auth/)