package main

import (
	"flag"
	"fmt"
	"net/http"

	"img-tools/internal/config"
	"img-tools/internal/errors"
	"img-tools/internal/handler"
	"img-tools/internal/svc"
	"img-tools/internal/types"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/img.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	httpx.SetErrorHandler(func(err error) (int, interface{}) {
		switch e := err.(type) {
		case errors.APIError:
			return http.StatusOK, types.CommonRsp{
				Code: e.Code,
				Msg:  e.Msg,
			}
		default:
			fmt.Println(err)
			return http.StatusInternalServerError, nil
		}
	})

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
