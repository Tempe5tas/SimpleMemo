package model

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

var DB *gorm.DB

func Init(dsn string) {
	// Connect to a database
	switch driver := viper.Get("Database.Driver").(string); driver {
	case "mysql":
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(fmt.Errorf("error reading MySQL database: %w", err))
		}
		DB = db
	case "sqlite":
		db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
		if err != nil {
			panic(fmt.Errorf("error reading sqlite database: %w", err))
		}
		DB = db
	default:
		panic(fmt.Errorf("invalid database driver, please check config file"))
	}
	// Set log mode
	sqlDB, _ := DB.DB()
	sqlDB.SetMaxIdleConns(20)  //设置连接池，空闲
	sqlDB.SetMaxOpenConns(100) //打开
	sqlDB.SetConnMaxLifetime(time.Second * 30)
	migration()
}

func migration() {
	// Auto migrate objects
	err := DB.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(&User{}, &Memo{})
	if err != nil {
		panic(fmt.Errorf("error configuring database structure: %w", err))
	}
}
