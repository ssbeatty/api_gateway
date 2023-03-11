package models

import (
	"api_gateway/internal/gateway/config"
	"api_gateway/internal/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/ssbeatty/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

const (
	defaultSort      = "id DESC"
	DBDriverMysql    = "mysql"
	DBDriverPostgres = "postgres"
	DBDriverSqlite   = "sqlite"
)

var (
	db       *DataBase
	dataPath string
	logger   *zerolog.Logger
)

type DataBase struct {
	*gorm.DB
	mu *sync.Mutex
}

func (d *DataBase) Lock() {
	if d.mu != nil {
		d.mu.Lock()
	}
}

func (d *DataBase) Unlock() {
	if d.mu != nil {
		d.mu.Unlock()
	}
}

func createDB(driver, user, pass, dsn, dbName string) error {

	switch driver {
	case DBDriverSqlite:
		return nil
	case DBDriverMysql:
		dataSource := fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8", user, pass, dsn)
		d, err := sql.Open(driver, dataSource)
		if err != nil {
			return err
		}
		defer func(d *sql.DB) {
			err := d.Close()
			if err != nil {

			}
		}(d)

		_, err = d.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8 COLLATE utf8_general_ci;", dbName))
		if err != nil {
			return err
		}
	case DBDriverPostgres:
		dsnArgs := strings.Split(dsn, ":")
		if len(dsnArgs) < 2 {
			return errors.New("dsn parse error")
		}
		dataSource := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable",
			dsnArgs[0], user, pass, dsnArgs[1],
		)
		d, err := sql.Open("pgx", dataSource)
		if err != nil {
			return err
		}
		defer func(d *sql.DB) {
			err := d.Close()
			if err != nil {

			}
		}(d)

		_, err = d.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			return err
		}
	}

	return nil
}

func InitModels(config config.DB, ctx context.Context) error {
	var d *gorm.DB
	var err error
	var dataSource string

	dataPath = config.DataPath
	logger = log.Ctx(ctx)
	logger.Info().Msg("I have to go...")
	logger.Info().Msg("Stopping server gracefully")
	newLogger := gormLogger.New(
		logger,
		gormLogger.Config{
			SlowThreshold:             time.Second,      // 慢 SQL 阈值
			LogLevel:                  gormLogger.Error, // 日志级别
			IgnoreRecordNotFoundError: true,             // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,            // 禁用彩色打印
		},
	)

	_ = createDB(config.Driver, config.User, config.Pass, config.DSN, config.DBName)
	dfConfig := &gorm.Config{
		Logger: newLogger,
	}

	if config.Driver == DBDriverMysql {
		dataSource = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", config.User, config.Pass, config.DSN, config.DBName)
		d, err = gorm.Open(mysql.Open(dataSource), dfConfig)
		db = &DataBase{d, nil}
	} else if config.Driver == DBDriverPostgres {
		dsnArgs := strings.Split(config.DSN, ":")
		if len(dsnArgs) < 2 {
			return errors.New("dsn parse error")
		}
		dataSource = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			dsnArgs[0], config.User, config.Pass, config.DBName, dsnArgs[1],
		)
		d, err = gorm.Open(postgres.Open(dataSource), dfConfig)
		db = &DataBase{d, nil}
	} else {
		if exist, _ := utils.PathExists(dataPath); !exist {
			_ = os.MkdirAll(dataPath, os.ModePerm)
		}

		// https://github.com/applikatoni/applikatoni/issues/35
		dataSource = path.Join(dataPath, "api_gateway.db?cache=shared&mode=rwc&_busy_timeout=30000")
		d, err = gorm.Open(sqlite.Open(dataSource), dfConfig)
		// 防止database locked
		sqlDB, err := d.DB()
		if err != nil {
			return err
		}
		// 设置最大开启连接数
		sqlDB.SetMaxOpenConns(1)
		sqlDB.SetMaxIdleConns(1)
		db = &DataBase{d, &sync.Mutex{}}
	}

	if err != nil {
		logger.Error().AnErr("gorm open filed", err)
		return err
	}

	return nil
}
