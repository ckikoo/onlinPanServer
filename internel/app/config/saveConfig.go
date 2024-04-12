package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/pelletier/go-toml"
	"gopkg.in/yaml.v2"
)

// SaveConfig 函数保存配置信息到指定格式的文件
func SaveConfig(fpath string) error {

	var saver configSaver
	switch {
	case strings.HasSuffix(fpath, ".json"):
		saver = &jsonSaver{}
	case strings.HasSuffix(fpath, ".toml"):
		saver = &tomlSaver{}
	case strings.HasSuffix(fpath, ".yaml") || strings.HasSuffix(fpath, ".yml"):
		saver = &yamlSaver{}
	default:
		return fmt.Errorf("unsupported config file format: %s", fpath)
	}

	return saver.Save(fpath)
}

// configSaver 接口定义了保存配置信息的方法
type configSaver interface {
	Save(fpath string) error
}

// jsonSaver 结构保存配置信息到 JSON 格式的文件
type jsonSaver struct{}

func (s *jsonSaver) Save(fpath string) error {
	file, err := os.Create(fpath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(C); err != nil {
		return fmt.Errorf("error encoding config to JSON: %w", err)
	}

	return nil
}

// tomlSaver 结构保存配置信息到 TOML 格式的文件
type tomlSaver struct{}

func (s *tomlSaver) Save(fpath string) error {
	// 创建一个新的 TOML 文件
	file, err := os.Create(fpath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	// 创建一个 TOML 树
	tree, err := toml.Marshal(C)

	// 编码 TOML 树并将其写入文件
	_, err = file.WriteString(string(tree))
	if err != nil {
		return fmt.Errorf("error writing config to file: %w", err)
	}

	return nil
}

// addConfigToTomlTree 函数将配置信息添加到 TOML 树中
func addConfigToTomlTree(tree *toml.Tree) {
	tree.Set("RunMode", C.RunMode)
	tree.Set("AppName", C.AppName)
	tree.Set("PrintConfig", C.PrintConfig)

	// 添加 HTTP 子树
	httpTree := tree.Get("HTTP").(*toml.Tree)
	httpTree.Set("Host", C.HTTP.Host)
	httpTree.Set("Port", C.HTTP.Port)
	httpTree.Set("CertFile", C.HTTP.CertFile)
	httpTree.Set("KeyFile", C.HTTP.KeyFile)
	httpTree.Set("ShutdownTimeOut", C.HTTP.ShutdownTimeout)
	httpTree.Set("MaxContentLength", C.HTTP.MaxContentLength)
	httpTree.Set("MaxReqLoggerLength", C.HTTP.MaxReqLoggerLength)
	httpTree.Set("MaxResLoggerLength", C.HTTP.MaxResLoggerLength)

	// 添加 File 子树
	fileTree := tree.Get("File").(*toml.Tree)
	fileTree.Set("InitSpaceSize", C.File.InitSpaceSize)
	fileTree.Set("FileUploadDir", C.File.FileUploadDir)

	// 添加 Email 子树
	emailTree := tree.Get("Email").(*toml.Tree)
	emailTree.Set("Port", C.Email.Port)
	emailTree.Set("Host", C.Email.Host)
	emailTree.Set("UserName", C.Email.UserName)
	emailTree.Set("Password", C.Email.Password)

	// 添加 Captcha 子树
	captchaTree := tree.Get("Captcha").(*toml.Tree)
	captchaTree.Set("Length", C.Captcha.Length)
	captchaTree.Set("Width", C.Captcha.Width)
	captchaTree.Set("Height", C.Captcha.Height)

	// 添加 RateLimiter 子树
	rateLimiterTree := tree.Get("RateLimiter").(*toml.Tree)
	rateLimiterTree.Set("Enable", C.RateLimiter.Enable)
	rateLimiterTree.Set("Count", C.RateLimiter.Count)
	rateLimiterTree.Set("RedisDB", C.RateLimiter.RedisDB)

	// 添加 Redis 子树
	redisTree := tree.Get("Redis").(*toml.Tree)
	redisTree.Set("Addr", C.Redis.Addr)
	redisTree.Set("Password", C.Redis.Password)

	// 添加 JWTAuth 子树
	jwtAuthTree := tree.Get("JWTAuth").(*toml.Tree)
	jwtAuthTree.Set("Enable", C.JWTAuth.Enable)
	jwtAuthTree.Set("SigningMethod", C.JWTAuth.SigningMethod)
	jwtAuthTree.Set("SigningKey", C.JWTAuth.SigningKey)
	jwtAuthTree.Set("Expired", C.JWTAuth.Expired)
	jwtAuthTree.Set("Store", C.JWTAuth.Store)
	jwtAuthTree.Set("FilePath", C.JWTAuth.FilePath)
	jwtAuthTree.Set("RedisDB", C.JWTAuth.RedisDB)
	jwtAuthTree.Set("RedisPrefix", C.JWTAuth.RedisPrefix)

	// 添加 Gorm 子树
	gormTree := tree.Get("Gorm").(*toml.Tree)
	gormTree.Set("Debug", C.Gorm.Debug)
	gormTree.Set("DBType", C.Gorm.DBType)
	gormTree.Set("MaxLifetime", C.Gorm.MaxLifetime)
	gormTree.Set("MaxOpenConns", C.Gorm.MaxOpenConns)
	gormTree.Set("MaxIdleConns", C.Gorm.MaxIdleConns)
	gormTree.Set("TablePrefix", C.Gorm.TablePrefix)
	gormTree.Set("EnableAutoMigrate", C.Gorm.EnableAutoMigrate)

	// 添加 MySQL 子树
	mysqlTree := tree.Get("MySQL").(*toml.Tree)
	mysqlTree.Set("Host", C.MySQL.Host)
	mysqlTree.Set("Port", C.MySQL.Port)
	mysqlTree.Set("User", C.MySQL.User)
	mysqlTree.Set("Password", C.MySQL.Password)
	mysqlTree.Set("DBName", C.MySQL.DBName)
	mysqlTree.Set("Parameters", C.MySQL.Parameters)

	// 添加 CORS 子树
	corsTree := tree.Get("CORS").(*toml.Tree)
	corsTree.Set("Enable", C.CORS.Enable)
	corsTree.Set("AllowOrigins", C.CORS.AllowOrigins)
	corsTree.Set("AllowMethods", C.CORS.AllowMethods)
	corsTree.Set("AllowHeaders", C.CORS.AllowHeaders)
	corsTree.Set("AllowCredentials", C.CORS.AllowCredentials)
	corsTree.Set("MaxAge", C.CORS.MaxAge)

	// 添加 GZIP 子树
	gzipTree := tree.Get("GZIP").(*toml.Tree)
	gzipTree.Set("Enable", C.GZIP.Enable)
	gzipTree.Set("ExcludedExtentions", C.GZIP.ExcludedExtentions)
	gzipTree.Set("ExcludedPaths", C.GZIP.ExcludedPaths)
}

// yamlSaver 结构保存配置信息到 YAML 格式的文件
type yamlSaver struct{}

func (s *yamlSaver) Save(fpath string) error {
	file, err := os.Create(fpath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	yamlData, err := yaml.Marshal(C)
	if err != nil {
		return fmt.Errorf("error marshalling config to YAML: %w", err)
	}

	_, err = file.Write(yamlData)
	if err != nil {
		return fmt.Errorf("error writing config to file: %w", err)
	}

	return nil
}
