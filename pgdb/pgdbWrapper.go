package pgdb

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/haolipeng/go-pg-example/conf"
	"sync"
)

var pgdbObj *pg.DB
var once sync.Once

//GetInstance 单例模式实现数据连接的初始化
func Connect() *pg.DB {
	once.Do(func() {
		//1.连接数据库
		pgdbObj = pg.Connect(&pg.Options{
			Addr:     conf.DbAddr,
			User:     conf.User,
			Password: conf.Password,
			Database: conf.DbName,
		})
		if pgdbObj == nil {
			fmt.Println("pg.Connect() failed,error:")
			panic(0)
		}
	})

	return pgdbObj
}

//CreateSchema 通过定义的models来创建数据库表
func CreateSchema(db *pg.DB, models []interface{}) error {
	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			//Temp: true,//建表是临时的
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

//DeleteSchema 通过结构体来删除表
func DeleteSchema(db *pg.DB, models []interface{}) error {
	err := db.Model(&models).DropTable(&orm.DropTableOptions{
		IfExists: true,
		Cascade:  true,
	})
	return err
}
