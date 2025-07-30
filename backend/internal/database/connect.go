package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectToPostgres 連線到 Postgres Database
//
// param host string 主機地址
// param user string 使用者名稱
// param password string 使用者密碼
// param name string 資料庫名稱
// param port string 服務端口
// returns *gorm.DB 資料庫連接
// returns error 錯誤
func ConnectToPostgres(
	host,
	user,
	password,
	name,
	port string,
	doLog bool,
) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Taipei",
		host, user, password, name, port,
	)
	if doLog {
		return gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	} else {
		return gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}

}

// ConnectToSQLLine
//
//	param filename string 資料檔案名稱
//	return *gorm.DB database connect
//	return error
func ConnectToSQLLine(filename string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(filename), &gorm.Config{})
}
