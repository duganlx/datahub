# go micro server

在*调研v1.0*中，已经完成casbin进行鉴权的过程，并且设计了GSF进行鉴权的方式方法。但是该套方案在讨论之后发现存在延迟高和引用服务过多的问题，所以目前设计的方案为：用户中心管理并维护用户可访问资产单元列表的数据，并生成每个资产单元各自的*用户访问权限信息表*推送到nacos；GSF只需要读取并监听nacos中对应的资产单元的*用户访问权限信息表*即可。当有用户需要在gsf中下单时，首先gsf会从请求的token中获取到用户的信息，接着去*用户访问权限信息表*查看该用户是否有权限，即可完全鉴权。目前还需要调研如下内容。

- [x] casbin鉴权的单元测试改造。在*调研v1.0*中单元测试写的比较分散并且可读性不高的问题。由于*调研v1.0*中鉴权的起点是appid和appsecret，而在该版本中起点是token，所以原本的测试用例也需要调整。
- [x] 通过接口修改权限信息。casbin在启动时会读取policy.csv文件的权限配置，在运行过程存在需要更新某些权限配置。
- [x] nacos 上传数据。用户中心将casbin保存的权限配置根据各个资产单元生成一份*用户访问权限信息表*上传到nacos中。
- [] 设计上传到nacos中数据内容。gsf读取nacos对应的权限信息后自行完成请求用户的权限校验。

---

**实验**

casbin的policy.csv存放到mysql的表中 `casbin_rule{p_type, v0, v1, v2, v3, v4, v5}`，在项目*启动时*会加载该表的配置（后续如果直接对表数据进行修改，并不会生效）

测试场景搭建
```text
用户(user)
boss: 1416962189826199552
ww: 1523580757186973696
xjw: 1506439972247310336
lvx: 1559730848930992128

用户组(user_group)
admin: boss, xjw
quant: ww, boss
test: lvx, xjw, ww


资源(src)
EAM001:v1:ip:test
EAM002:v1:ip:test
EAM003:v1:ip:test
EAM011:v1:ip:prod
EAM012:v1:ip:prod
DRW001:v1:ip:test
DRW001:v1:ip:prod

资源组(src_group)
test: 最后一项为test，即EAM001:v1:ip:test, EAM002:v1:ip:test, EAM003:v1:ip:test, DRW001:v1:ip:test
prod: 最后一项为prod，即EAM011:v1:ip:prod, EAM012:v1:ip:prod
EAM: 第一项是EAM开头，即EAM001:v1:ip:test, EAM002:v1:ip:test, EAM003:v1:ip:test, EAM011:v1:ip:prod, EAM012:v1:ip:prod
DRW: 第一项是DRW开头，即DRW001:v1:ip:test, DRW001:v1:ip:prod


权限关系
lvx(user): DRW001:v1:ip:prod(src)
boss(user): prod(src_group)
admin(user_group): *(src)
quant(user_group): EAM(src_group), DRW(src_group)
test(user_group): test(src_group), EAM011:v1:ip:prod(src)
```

对应sql
```sql
-- 资源组(src_group)
INSERT INTO casbin_rule VALUES('p', 'SRCGROUP:test', 'EAM001:v1:ip:test', '*', '', '', '');
INSERT INTO casbin_rule VALUES('p', 'SRCGROUP:test', 'EAM002:v1:ip:test', '*', '', '', '');
INSERT INTO casbin_rule VALUES('p', 'SRCGROUP:test', 'EAM003:v1:ip:test', '*', '', '', '');
INSERT INTO casbin_rule VALUES('p', 'SRCGROUP:test', 'DRW001:v1:ip:test', '*', '', '', '');
INSERT INTO casbin_rule VALUES('p', 'SRCGROUP:prod', 'EAM011:v1:ip:prod', '*', '', '', '');
INSERT INTO casbin_rule VALUES('p', 'SRCGROUP:prod', 'EAM012:v1:ip:prod', '*', '', '', '');
INSERT INTO casbin_rule VALUES('p', 'SRCGROUP:EAM', 'EAM001:v1:ip:test', '*', '', '', '');
INSERT INTO casbin_rule VALUES('p', 'SRCGROUP:EAM', 'EAM002:v1:ip:test', '*', '', '', '');
INSERT INTO casbin_rule VALUES('p', 'SRCGROUP:EAM', 'EAM003:v1:ip:test', '*', '', '', '');
INSERT INTO casbin_rule VALUES('p', 'SRCGROUP:EAM', 'EAM011:v1:ip:prod', '*', '', '', '');
INSERT INTO casbin_rule VALUES('p', 'SRCGROUP:EAM', 'EAM012:v1:ip:prod', '*', '', '', '');
INSERT INTO casbin_rule VALUES('p', 'SRCGROUP:DRW', 'DRW001:v1:ip:test', '*', '', '', '');
INSERT INTO casbin_rule VALUES('p', 'SRCGROUP:DRW', 'DRW001:v1:ip:prod', '*', '', '', '');

-- 用户组(user_group)
INSERT INTO casbin_rule VALUES('p', 'USERGROUP:admin', '*', '*', '', '', '');
INSERT INTO casbin_rule VALUES('g', 'USERGROUP:quant', 'SRCGROUP:EAM', '', '', '', '');
INSERT INTO casbin_rule VALUES('g', 'USERGROUP:quant', 'SRCGROUP:DRW', '', '', '', '');
INSERT INTO casbin_rule VALUES('g', 'USERGROUP:test', 'SRCGROUP:test', '', '', '', '');
INSERT INTO casbin_rule VALUES('p', 'USERGROUP:test', 'EAM011:v1:ip:prod', '*', '', '', '');

-- 用户权限分配
INSERT INTO casbin_rule VALUES('p', 'USER:1559730848930992128', 'DRW001:v1:ip:prod', '*', '', '', '');
INSERT INTO casbin_rule VALUES('g', 'USER:1416962189826199552', 'USERGROUP:admin', '', '', '', '');
INSERT INTO casbin_rule VALUES('g', 'USER:1506439972247310336', 'USERGROUP:admin', '', '', '', '');
INSERT INTO casbin_rule VALUES('g', 'USER:1416962189826199552', 'USERGROUP:quant', '', '', '', '');
INSERT INTO casbin_rule VALUES('g', 'USER:1523580757186973696', 'USERGROUP:quant', '', '', '', '');
INSERT INTO casbin_rule VALUES('g', 'USER:1559730848930992128', 'USERGROUP:test', '', '', '', '');
INSERT INTO casbin_rule VALUES('g', 'USER:1523580757186973696', 'USERGROUP:test', '', '', '', '');
INSERT INTO casbin_rule VALUES('g', 'USER:1506439972247310336', 'USERGROUP:test', '', '', '', '');
```

## 附录

### 1. 定位问题 nacos-sdk-go 调用注册到 nacos 的服务器 http 接口时报 `code = 503 reason = NODE_NOT_FOUND message = error: code = 503 reason = no_available_node message =  metadata = map[] cause = <nil> metadata = map[] cause = <nil>` 

客户端通过 nacos 去调用服务器的 http 方式接口时，会出现问题 `code = 503 reason = NODE_NOT_FOUND message = error: code = 503 reason = no_available_node message =  metadata = map[] cause = <nil> metadata = map[] cause = <nil>`，问题定位如下。解决办法有两种，第一种是采用 grpc 去访问（推荐）；第二种是手动取服务节点转换成最终的 url(比如 `http://127.0.0.1:8000`)去访问。

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

### 2. `rpc error: code = DeadlineExceeded desc = context deadline exceeded` 问题解决

客户端每次运行前都需要将 cache/naming 中的内容删除，否则无法启动显示: rpc error: code = DeadlineExceeded desc = context deadline exceeded

### 3. 调研v1.0

目前存在两个需求点：1. 资产单元的权限管理; 2. 微服务架构下服务的权限管理;

_资产单元的权限管理_

即 right 用户的 right 模型在 right 资产单元进行下单交易。那么就需要考虑如下几个问题：

1. 如何保证 right 用户? 即资产单元允许哪些用户进行访问（资产单元的权限管理）
2. 如何保证 right 模型? 模型是用户自己创建的，是否为正确的模型是由用户进行管理，平台可以提供一套机制协助进行管理。

casbin 提供了 RBAC 的权限设计方案，可以将用户作为 sub，资产单元作为 obj，资产单元访问方式作为 act；另外还需要有`角色/组`的概念，`角色/组` 与资产单元进行绑定，表示该`角色/组`可以访问哪些资产单元，而用户可以跟`角色/组`进行绑定，表示该用户可以访问该`角色/组`中所有资产单元。这样 casbin 就帮助我们解决第一个问题了（如何保证 right 用户?）；第二个问题由访问令牌进行解决，用户自行选择对可访问的资产单元并输入 appid 和 expires 后，生成对应的 appsecret，然后在模型登录时提供 appid 和 appsecret 进行鉴权即可。

- authcode{id, appid, appsecret, expires, aucodes, allow, userid}

整个流程如下

1. 管理员先配置好用户可访问的资产单元列表
1. 用户生成某个资产单元（au1）的访问令牌 appsecret，另外也可以生成能访问所有可访问资产单元（\*）的令牌 appsecret
1. 用户编写的模型要访问资产单元进行下单前，需要提供 appid, appsecret, aucode 进行鉴权，鉴权通过之后会生成 jwt 信息
1. 后续访问携带 jwt 信息访问

**实验测试**

管理员设置了资产单元的访问规则如下。为了更好的进行资产单元的管理，将资产单元按组为单位进行划分后，再将整个组分配给特定用户。

```text
==== policy.csv 内容 ====
# 产品 & 资产单元对应关系
p, PRODUCT_EAMLS1, AU_300016, *
p, PRODUCT_EAMLS1, AU_88853899_ww, r
p, PRODUCT_EAMLS1, AU_EAMLS1ZT_00, w
p, PRODUCT_EAMLS1, AU_EAMLS1ZTX_00, *
p, PRODUCT_DRW004, AU_121000, *

# 投资经理 & 资产单元对应关系
p, MANAGER_WW, AU_0148P1016_ww, *
p, MANAGER_WW, AU_88853899_ww, r
p, MANAGER_WW, AU_DRWZQ1ZT_03, w
# p, MANAGER_WSY, AU_DRW001ZTX_04, *

# 临时配置某个用户对某个资产单元的配置
p, USER_wsy, AU_DRW001ZTX_04, *

# 定义关联: 用户 - 可访问的资产单元(组)
g, USER_ww, MANAGER_WW
g, USER_xjw, PRODUCT_EAMLS1
g, USER_yrl, MANAGER_WW
# g, USER_wsy, MANAGER_WSY


==== policy.csv 解读 ====
用户ww可以访问 MANAGER_WW组中的资产单元（MANAGER_WW组 = {0148P1016_ww, 88853899_ww, DRWZQ1ZT_03}）
用户xjw可以访问 PRODUCT_EAMLS1组中的资产单元（PRODUCT_EAMLS1组 = {300016, 88853899_ww, EAMLS1ZT_00, EAMLS1ZTX_00}）
用户wsy可以访问资产单元DRW001ZTX_04
用户yrl可以访问 MANAGER_WW组中的资产单元
```

场景设计如下，样例中所说的成功/失败表示预期的鉴权结果（成功：鉴权通过；失败：鉴权不通过）

1. 用户 ww 生成*只能*访问资产单元`[0148P1016_ww]`的访问令牌，并访问资产单元`0148P1016_ww` —— 成功
2. 用户 ww 生成*只能*访问资产单元`[0148P1016_ww, 88853899_ww]`的访问令牌，并访问资产单元`0148P1016_ww` —— 成功
3. 用户 ww 生成*只能*访问资产单元`[0148P1016_ww]`的访问令牌，并访问资产单元`88853899_ww` —— 失败，该令牌没有访问`88853899_ww`的权限
4. 用户 ww 生成可访问*所有*资产单元（`MANAGER_WW组`）的访问令牌，并访问资产单元`88853899_ww` —— 成功
5. 用户 ww 生成可访问*所有*资产单元（`MANAGER_WW组`）的访问令牌，并访问资产单元`EAMLS1ZT_00` —— 失败，用户 ww 没有访问`EAMLS1ZT_00`的权限
6. 用户 ww 生成*不能*访问资产单元`[0148P1016_ww]`的访问令牌，并访问资产单元`0148P1016_ww` —— 失败
7. 用户 ww 生成*不能*访问资产单元`[0148P1016_ww]`的访问令牌，并访问资产单元`88853899_ww` —— 成功，`88853899_ww` 在 MANAGER_WW 组中，但不在不能访问的列表中
8. 用户 ww 生成*不能*访问资产单元`[0148P1016_ww]`的访问令牌，并访问资产单元`EAMLS1ZT_00` —— 失败，用户 ww 没有访问`EAMLS1ZT_00`的权限
9. 用户 wsy 生成*只能*访问资产单元`[DRW001ZTX_04]`的访问令牌，并访问资产单元`DRW001ZTX_04` —— 成功
10. 用户 xjw 生成*只能*访问资产单元`[EAMLS1ZT_00]`的访问令牌，并访问资产单元`EAMLS1ZT_00` —— 成功

代码实现上，设计了两张表来维护数据，分别为*访问令牌表*和*用户表*，表结构和上述实验的数据如下所示。

```text
用户表 User {Id, UserName, Mobile}
访问令牌表 AccessToken {Id, AppId, AppSecret, UserId, AuCodes, Allow, Expires}

== 实验数据 ==
用户表 User 数据:
[
	{Id: 15739, UserName: "ww", Mobile: "15308681364"},
	{Id: 15743, UserName: "xjw", Mobile: "13608681364"},
	{Id: 15747, UserName: "wsy", Mobile: "13708681364"},
]

访问令牌表 AccessToken 数据:
[
	// ww(id:15739) has generated authToken to access 0148P1016_ww
	{Id: 1, AppId: "asdj", AppSecret: "d54sdfejbd561sa", UserId: 15739, AuCodes: ["0148P1016_ww"], Allow: true, Expires: "Seven days later"},

	// xjw(id:15743) has generated authToken to access EAMLS1ZT_00
	{Id: 2, AppId: "kfuks", AppSecret: "4fd1ufklnksbry9", UserId: 15743, AuCodes: ["EAMLS1ZT_00"], Allow: true, Expires: "Seven days later"},

	// ww(id:15739) has generated authToken to access all au which ww can access.
	{Id: 3, AppId: "jkwsx", AppSecret: "luwxtuf5twprw5l", UserId: 15739, AuCodes: ["*"], Allow: true, Expires: "Seven days later"},

	// ww(id:15739) has generated authToken to access all au which ww can access except 0148P1016_ww.
	{Id: 4, AppId: "ggTks", AppSecret: "psuhl055bwaeTIjk", UserId: 15739, AuCodes: ["0148P1016_ww"], Allow: false, Expires: "Seven days later"},

	// ww(id:15739) has generated authToken to access [0148P1016_ww, 88853899_ww]
	{Id: 5, AppId: "xstt", AppSecret: "abeo5tgrt754arh57", UserId: 15739, AuCodes: ["0148P1016_ww", "88853899_ww"], Allow: true, Expires: "Seven days later"},

	// wsy(id:15747) has generated authToken to access DRW001ZTX_04
	{Id: 5, AppId: "ko8w", AppSecret: "8hw416ery9ah4foig", UserId: 15747, AuCodes: ["DRW001ZTX_04"], Allow: true, Expires: "Seven days later"},

	// xjw(id:15743) has generated authToken to access EAMLS1ZT_00
	{Id: 5, AppId: "eut2", AppSecret: "tyt1ra48is13awer6", UserId: 15743, AuCodes: ["EAMLS1ZT_00"], Allow: true, Expires: "Seven days later"},
]
```

_微服务架构下服务的权限管理_

用 nacos 作为服务注册&发现中心，各个`资产单元`和`用户中心`都会在 nacos 进行注册。当用户的某个模型需要在资产单元`Au1`中下单时，带上 appid 和 appsecret，`Au1`会去访问`用户中心`的接口进行鉴权，当鉴权通过之后，则将该对 appid 和 appsecret 保存在内存中，下次如果再遇到该对时就不用再去用户中心鉴权而直接放行。这样的设计会导致

### 4. nacos 调研进度

注册中心

- [x] 理想情况下的服务注册与访问（grpc、http）
- [ ] 注册服务的访问权限控制
- [ ] 服务可用性监测

配置中心

- [x] 服务启动时读取 nacos 配置
- [ ] 服务运行中实时同步最新的 nacos 配置

## 参考

- [SSL/TLS 协议运行机制的概述](https://www.ruanyifeng.com/blog/2014/02/ssl_tls.html)
- [grpc-auth-support.md(grpc-go Documentation)](https://github.com/grpc/grpc-go/blob/master/Documentation/grpc-auth-support.md)
- [gRPC authentication guide](https://grpc.io/docs/guides/auth/)
- [jwt 在线解密](https://www.box3.cn/tools/jwt.html)
- [Casbin 文档](https://casbin.org/zh/docs/overview)
- [Basic Role-Based HTTP Authorization in Go with Casbin](https://zupzup.org/casbin-http-role-auth/)
