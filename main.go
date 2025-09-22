package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/colin-404/logx"
	"github.com/spf13/viper"
	"github.com/xid-protocol/SIEM/aws"
)

var sig = make(chan os.Signal, 1)

func initConfig() string {
	confPath := flag.String("c", "/opt/xidp/conf/config.yml", "config file path")
	flag.Parse()
	// confPath_str := common.NormalizePath(*confPath)
	//如果配置文件不存在，则报错并关闭程序
	if _, err := os.Stat(*confPath); os.IsNotExist(err) {
		logx.Errorf("config file not found: %s", *confPath)
		os.Exit(1)
	}
	//加载配置
	return *confPath
}

func initLog() {

	//初始化日志
	logOpts := logx.Options{
		LogFile:    "/var/log/siem/siem.log",
		MaxSize:    10,
		MaxAge:     100,
		MaxBackups: 100,
		TimeFormat: logx.TimeFormats.RFC3339,
	}
	loger := logx.NewLoger(&logOpts)
	logx.InitLogger(loger)
}

func init() {
	initLog()
	//优雅关闭
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	//获取配置路径
	confPath := initConfig()

	//使用viper加载配置
	viper.SetConfigFile(confPath)
	viper.ReadInConfig()

	// initLog()
	// db.InitMongoDatabase(viper.GetString("mongodb.uri"), viper.GetString("mongodb.database"))

}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 监听信号并取消 context
	go func() {
		<-sig
		logx.Infof("Received signal, cancelling context...")
		cancel()
	}()
	awsCloud := aws.NewAWSCloud()
	go awsCloud.Handler(ctx)
	//3分钟检测一次
	ticker := time.NewTicker(175 * time.Second)
	defer ticker.Stop()

	awsCloud.GuardDuty(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			awsCloud.GuardDuty(ctx)
		}
	}

}
