# GoFrame Etcd Registry

Use `etcd` as service registration and discovery management.

Fork by https://github.com/gogf/gf/tree/master/contrib/registry/etcd

To adjust the connection parameters to accommodate the need for account credentials in an etcd service, modify the connection parameters to use etcd3.Config.

## Installation

## Example

### Reference example
server
```go
package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/gsvc"
	"mkluo123.xyz/gofg/gf/contrib/registry/etcd/v2"
)

func main() {
	grpcx.Resolver.Register(etcd.New(etcd3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
		Username:  "username",
		Password:  "pwd",
	}))
	s := g.Server(`hello.svc`)
	s.BindHandler("/", func(r *ghttp.Request) {
		g.Log().Info(r.Context(), `request received`)
		r.Response.Write(`Hello world`)
	})
	s.Run()
}
```
client

```go
package main

import (
	"fmt"
	"time"

	"github.com/gogf/gf/contrib/registry/etcd/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsel"
	"github.com/gogf/gf/v2/net/gsvc"
	"mkluo123.xyz/gofg/gf/contrib/registry/etcd/v2"
)
func main() {
	gsvc.SetRegistry(etcd.New(etcd3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
		Username:  "username",
		Password:  "pwd",
	}))
	gsel.SetBuilder(gsel.NewBuilderRoundRobin())

	client := g.Client()
	for i := 0; i < 100; i++ {
		res, err := client.Get(gctx.New(), `http://hello.svc/`)
		if err != nil {
			panic(err)
		}
		fmt.Println(res.ReadAllString())
		res.Close()
		time.Sleep(time.Second)
	}
}
```