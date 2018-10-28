package model

import (
	"github.com/jinzhu/gorm"
	"fmt"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Database struct {
	Local *gorm.DB
}

var DB *Database

func (db *Database) Init() {
	DB = &Database{
		Local: GetLocalDB(),
	}
	if err := DB.Local.AutoMigrate(&UserModel{}).Error; err != nil {
		log.Debug(err.Error())
	}
	log.Info("Database connect successfully.")
}
func (db *Database) Close() {
	DB.Local.Close()
}

func openDB(username, password, addr, name string) *gorm.DB {
	config := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s",
		username,
		password,
		addr,
		name,
		true,
		//"Asia/Shanghai"),
		"Local")

	db, err := gorm.Open("mysql", config)
	if err != nil {
		log.Errorf(err, "Database: %s connection failed.", name)
	}
	setupDB(db)
	return db
}

func setupDB(db *gorm.DB) {
	db.LogMode(viper.GetBool("gormlog"))
	db.DB().SetMaxIdleConns(0)
}

func InitLocalDB() *gorm.DB {
	res := openDB(viper.GetString("db.username"),
		viper.GetString("db.password"),
		viper.GetString("db.addr"),
		viper.GetString("db.name"))
	return res
}
func GetLocalDB() *gorm.DB {
	return InitLocalDB()
}
