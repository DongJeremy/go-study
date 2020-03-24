package main

import (
	"net"
	"net/http"
	"time"

	"github.com/mash/go-accesslog"
)

var log = GetLogger("test main")

type logger struct {
}

func (l logger) Log(record accesslog.LogRecord) {
	log.Println(record.Method + " " + record.Uri)
}

func main() {
	l := logger{}
	rootPath := "d:/"
	http.Handle("/", accesslog.NewLoggingHandler(http.FileServer(http.Dir(rootPath)), l))

	httpServer := &http.Server{
		Addr:           "8080",          // 监听的地址和端口
		Handler:        nil,             // 所有请求需要调用的Handler（实际上这里说是ServeMux更确切）如果为空则设置为DefaultServeMux
		ReadTimeout:    0 * time.Second, // 读的最大Timeout时间
		WriteTimeout:   0 * time.Second, // 写的最大Timeout时间
		MaxHeaderBytes: 256,             // 请求头的最大长度
		TLSConfig:      nil,             // 配置TLS
	}
	h, err := net.Listen("tcp4", "0.0.0.0:8080")
	if err != nil {
		log.Errorf("start HTTP failed, %s", err)
		return
	}
	log.Info("starting pxe server...")
	if err := httpServer.Serve(h); err != nil {
		log.Errorf("HTTP server shut down: %s", err)
		return
	}
}
