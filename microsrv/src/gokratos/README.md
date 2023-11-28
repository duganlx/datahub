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

客户端每次运行前都需要将 cache/naming 中的内容删除，否则无法启动。报错信息为: rpc error: code = DeadlineExceeded desc = context deadline exceeded

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


## 参考

- [SSL/TLS协议运行机制的概述](https://www.ruanyifeng.com/blog/2014/02/ssl_tls.html)
- [grpc-auth-support.md(grpc-go Documentation)](https://github.com/grpc/grpc-go/blob/master/Documentation/grpc-auth-support.md)
- [gRPC authentication guide](https://grpc.io/docs/guides/auth/)
- [jwt在线解密](https://www.box3.cn/tools/jwt.html)
- [Casbin 文档](https://casbin.org/zh/docs/overview)