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

func Init(ctx context.Context, opts ...Option) (func(), error) {
	var o options
	for _, opt := range opts {
		opt(&o)
	}

	config.MustLoad(o.ConfigFile)

	config.PrintWithJSON()

	injector, injectorCleanFunc, err := BuildInjector()
	if err != nil {
		return nil, err
	}

	httpServerCleanFunc := InitHttpServer(ctx, injector.Engine)

	return func() {
		httpServerCleanFunc()
		injectorCleanFunc()
	}, nil
}

func InitHttpServer(ctx context.Context, handler http.Handler) func() {
	cfg := config.C.HTTP
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	fmt.Printf("addr: %v\n", addr)
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
		// ReadTimeout: 10 * time.Second,
		// WriteTimeout: 30 * time.Second,
		// IdleTimeout: 30 * time.Second,
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
