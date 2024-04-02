package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/koding/multiconfig"
)

var (
	C    = new(Config)
	once sync.Once
)

func PrintWithJSON() {
	if C.PrintConfig {
		b, err := json.MarshalIndent(C, "", "")
		if err != nil {
			os.Stdout.WriteString("[CONFIG] JSON marshal error: " + err.Error())
			return
		}
		os.Stdout.WriteString(string(b) + "\n")
	}
}

func MustLoad(fpaths ...string) {
	once.Do(func() {
		loaders := []multiconfig.Loader{
			&multiconfig.TagLoader{},
			&multiconfig.EnvironmentLoader{},
		}

		for _, path := range fpaths {
			if strings.HasSuffix(path, "toml") {
				loaders = append(loaders, &multiconfig.TOMLLoader{Path: path})
			}
			if strings.HasSuffix(path, "json") {
				loaders = append(loaders, &multiconfig.JSONLoader{Path: path})
			}
			if strings.HasSuffix(path, "yaml") {
				loaders = append(loaders, &multiconfig.YAMLLoader{Path: path})
			}
		}

		m := multiconfig.DefaultLoader{
			Loader:    multiconfig.MultiLoader(loaders...),
			Validator: multiconfig.MultiValidator(&multiconfig.RequiredValidator{}),
		}
		m.MustLoad(C)
	})
}

type Config struct {
	AppName     string
	PrintConfig bool
	RunMode     string
	HTTP        HTTP
	JWTAuth     JWTAuth
	Captcha     Captcha
	RateLimiter RateLimiter
	CORS        CORS
	Email       Email
	Redis       Redis
	Gorm        Gorm
	MySQL       MySQL
	GZIP        GZIP
	File        File
}

type File struct {
	DefaultSpace uint64
}
type Email struct {
	Host     string
	Port     string
	UserName string
	Password string
}

type GZIP struct {
	Enable             bool
	ExcludedExtentions []string
	ExcludedPaths      []string
}

type Captcha struct {
	Length int
	Width  int
	Height int
}

type RateLimiter struct {
	Enable  bool
	Count   int64
	RedisDB int
}

type HTTP struct {
	Host               string
	Port               int
	CertFile           string
	KeyFile            string
	ShutdownTimeout    int
	MaxContentLength   int64
	MaxReqLoggerLength int `default:"1024"`
	MaxResLoggerLength int `default:"1024"`
}

type JWTAuth struct {
	Enable        bool
	SigningMethod string
	SigningKey    string
	Expired       int
	Store         string
	FilePath      string
	RedisDB       int
	RedisPrefix   string
}
type CORS struct {
	Enable           bool
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	MaxAge           int
}
type Redis struct {
	Addr     string
	Password string
}
type Gorm struct {
	Debug             bool
	DBType            string
	MaxLifetime       int
	MaxOpenConns      int
	MaxIdleConns      int
	TablePrefix       string
	EnableAutoMigrate bool
}
type MySQL struct {
	Host       string
	Port       int
	User       string
	Password   string
	DBName     string
	Parameters string
}

func (a MySQL) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		a.User, a.Password, a.Host, a.Port, a.DBName, a.Parameters)
}
