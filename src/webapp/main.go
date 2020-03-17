package main

import (
	"flag"
	"fmt"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"io/ioutil"
	"os"
	"runtime"
	"alert"
)

func parseArgs() (string, uint, string) {
	host := flag.String("h", "0.0.0.0", "监听地址")
	port := flag.Uint("p", 8080, "监听端口")
	pidFile := flag.String("pid-file", "", "进程pid文件")

	flag.Parse()

	return *host, *port, *pidFile
}

func recovery() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if err := recover(); err != nil {
					stack := make([]byte, 4<<10) // 4 KB
					length := runtime.Stack(stack, false)
					key := fmt.Sprintf("%v", err)
					errMsg := fmt.Sprintf("[PANIC RECOVER] %v\n %s\n", err, stack[:length])
					alert.Alert(key, errMsg)
				}
			}()

			return next(c)
		}
	}
}

func main() {
	Init()

	e := echo.New()

	// recover and alert middle ware
	//e.Use(middleware.Recover())
	e.Use(recovery())//在捕获到panic后，按时刻发送出去
	e.Use(middleware.Logger())


	// 注册路由
	for uri, handler := range Routes() {
		e.POST(uri, handler)
		e.GET(uri, handler)
	}

	// 获取启动参数
	host, port, pidFile := parseArgs()
	//echopprof.Wrap(e)

	e.Server.Addr = fmt.Sprintf("%s:%d", host, port)

	// 默认pidFile为空，不保存pid到文件。pidFile 不为空，则保存pid到文件
	if pidFile != "" {
		pid := fmt.Sprintf("%d\n", os.Getpid())
		err := ioutil.WriteFile(pidFile, []byte(pid), 0644)
		if err != nil {
			panic(err)
		}
	}


	err := gracehttp.Serve(e.Server)
	if err != nil {
		fmt.Println("webapp start err:", err)
	}
}
