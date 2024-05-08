package app

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"onlineCLoud/internel/app/config"
	"os"
	"os/signal"

	"syscall"
	"time"
)

type options struct {
	ConfigFile string
	Version    string
}

type Option func(*options)

func SetConfigFile(s string) Option {
	return func(o *options) {
		fmt.Printf("%s\n", s)
		o.ConfigFile = s
	}
}

func SetVersion(s string) Option {
	return func(o *options) {
		o.Version = s
	}
}

// Init 初始化函数
// ctx: 上下文，用于控制函数的生命周期。
// opts: 一个或多个选项函数，用于配置初始化过程。
// 返回值: 初始化完成后需要调用的清理函数和可能发生的错误。
func Init(ctx context.Context, opts ...Option) (func(), error) {
	var o options // 定义选项结构体变量用于存储配置。
	// 遍历并应用所有提供的选项函数来配置初始化过程。
	for _, opt := range opts {
		opt(&o)
	}

	// 加载配置文件。
	config.MustLoad(o.ConfigFile)

	// 打印配置信息。
	config.PrintWithJSON()

	// 构建依赖注入器，并获取其清理函数。
	injector, injectorCleanFunc, err := BuildInjector()
	if err != nil {
		return nil, err // 如果构建依赖注入器失败，则返回错误。
	}

	// 初始化HTTP服务器，并获取其清理函数。
	httpServerCleanFunc := InitHttpServer(ctx, injector.Engine)

	// 返回一个组合的清理函数，用于清理HTTP服务器和依赖注入器。
	return func() {
		httpServerCleanFunc() // 清理HTTP服务器。
		injectorCleanFunc()   // 清理依赖注入器。
	}, nil
}

func InitHttpServer(ctx context.Context, handler http.Handler) func() {
	cfg := config.C.HTTP
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	fmt.Printf("addr: %v\n", addr)
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		var err error
		if cfg.CertFile != "" && cfg.KeyFile != "" {
			srv.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
			err = srv.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(cfg.ShutdownTimeout))
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			os.Stdout.WriteString(err.Error())
		}
	}
}

func Run(ctx context.Context, opts ...Option) error {
	state := 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	cleanFunc, err := Init(ctx, opts...)
	if err != nil {
		return err
	}

EXIT:
	for {
		sig := <-sc
		os.Stdout.WriteString(fmt.Sprintf("Receive signal[%s]", sig.String()))
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			state = 0
			break EXIT
		case syscall.SIGHUP:
		default:
			break EXIT
		}
	}
	cleanFunc()
	os.Stdout.WriteString("Server exit\n")
	time.Sleep(time.Second)
	os.Exit(state)
	return nil
}
