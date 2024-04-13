package gormx

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// config GORM config
type Config struct {
	Debug        bool
	DBType       string
	DSN          string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
	TablePrefix  string
}

func New(c *Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch strings.ToLower(c.DBType) {
	case "mysql":
		cfg, err := mysqlDriver.ParseDSN(c.DSN)
		if err != nil {
			return nil, err
		}

		err = createDatabaseWithMysql(cfg)

		if err != nil {
			return nil, err
		}

		dialector = mysql.Open(c.DSN)
	}

	gconfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   c.TablePrefix,
			SingularTable: true,
		},
	}

	db, err := gorm.Open(dialector, gconfig)
	if err != nil {
		return nil, err
	}
	if c.Debug {
		db = db.Debug()
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(c.MaxIdleConns)
	sqlDB.SetMaxOpenConns(c.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(c.MaxLifetime) * time.Second)
	return db, nil
}

func createDatabaseWithMysql(c *mysqlDriver.Config) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/", c.User, c.Passwd, c.Addr)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	defer db.Close()

	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET = `utf8mb4`;", c.DBName)
	_, err = db.Exec(query)
	return err
}
