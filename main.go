package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"net/http"
)

func main() {
	initViper()
	initPrometheus()
	app := InitApp()
	err := app.GRPCServer.Serve()
	if err != nil {
		panic(err)
	}
	_ = app.GRPCServer.Close()
}

func initViper() {
	cfile := pflag.String("config", "config/config.yaml", "配置文件路径")
	pflag.Parse()
	// 直接指定文件路径
	viper.SetConfigFile(*cfile)
	// 实时监听配置变更
	viper.WatchConfig()
	// 读取配置到viper 里面
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func initPrometheus() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		// 暴露监听端口
		http.ListenAndServe(":8081", nil)
	}()
}
