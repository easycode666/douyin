// Copyright 2022 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package db

import (
	"fmt"
	"github.com/easycode666/douyin/pkg/config"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/logging/logrus"
	"gorm.io/plugin/opentelemetry/tracing"
)

var (
	DB    *gorm.DB
	viper = config.InitConfig("db-config")
)

// Init init DB
func Init() {
	var err error
	gormlogrus := logger.New(
		logrus.NewWriter(),
		logger.Config{
			SlowThreshold: time.Millisecond,
			Colorful:      true,
			LogLevel:      logger.Info,
		},
	)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		viper.GetString("MySQL.User"),
		viper.GetString("MySQL.Password"),
		viper.GetString("MySQL.Host"),
		viper.GetInt("MySQL.Port"),
		viper.GetString("MySQL.DBName"),
		viper.GetString("MySQL.CharSet"),
		viper.GetBool("MySQL.ParseTime"),
		viper.GetString("MySQL.loc"),
	)
	DB, err = gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			PrepareStmt: true,
			Logger:      gormlogrus,
		},
	)
	if err != nil {
		panic(err)
	}
	if err := DB.Use(tracing.NewPlugin()); err != nil {
		panic(err)
	}

	// 迁移数据表
	if err := DB.AutoMigrate(&User{}); err != nil {
		panic(err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		panic(err)
	}
	if err = sqlDB.Ping(); err != nil {
		panic(err)
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(viper.GetInt("MySQL.maxIdleConns"))
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(viper.GetInt("MySQL.maxOpenConns"))
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
}
