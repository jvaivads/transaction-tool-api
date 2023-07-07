package database

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

type SQLConfig struct {
	User   string
	Pass   string
	DBName string
	Net    string
	Host   string
	Port   string
	driver string

	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
	ConnMaxIdleTime int
}

func (cfg SQLConfig) toMySQLConfig() *mysql.Config {
	config := mysql.NewConfig()
	config.User = cfg.User
	config.Passwd = cfg.Pass
	config.DBName = cfg.DBName
	config.Net = cfg.Net
	config.Addr = cfg.Host + ":" + cfg.Port

	return config
}

func NewSQLClient(cfg SQLConfig) *sql.DB {
	var (
		client *sql.DB
		err    error
	)

	if client, err = sql.Open(cfg.driver, cfg.toMySQLConfig().FormatDSN()); err != nil {
		panic(err)
	}

	client.SetMaxOpenConns(cfg.MaxOpenConns)
	client.SetMaxIdleConns(cfg.MaxIdleConns)

	if err = client.Ping(); err != nil {
		panic(fmt.Errorf("cannot connect with DSN %s due to: %w", cfg.toMySQLConfig().FormatDSN(), err))
	}

	fmt.Println("database connected")

	return client
}

func GetLocalMySQLClientConfig() SQLConfig {
	return SQLConfig{
		User:         "root",
		Pass:         "",
		DBName:       "transaction",
		Net:          "tcp",
		Host:         "transaction-db",
		Port:         "3306",
		driver:       "mysql",
		MaxOpenConns: 5,
		MaxIdleConns: 2,
	}
}
