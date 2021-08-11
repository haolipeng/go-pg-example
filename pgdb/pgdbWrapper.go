package pgdb

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/haolipeng/go-pg-example/conf"
	"sync"
)

type PgWrapper struct {
	Db *pg.DB
}

var g_dbWrapper *PgWrapper
var once sync.Once

//GetInstance 单例模式实现
func GetInstance() *PgWrapper {
	once.Do(func() {
		g_dbWrapper = new(PgWrapper)

		//1.连接数据库
		pgsqlDB := pg.Connect(&pg.Options{
			Addr:     conf.DbAddr,
			User:     conf.User,
			Password: conf.Password,
			Database: conf.DbName,
		})
		if pgsqlDB == nil {
			fmt.Println("pg.Connect() failed,error:")
		}

		g_dbWrapper.Db = pgsqlDB
	})

	return g_dbWrapper
}

//通过定义的models来创建数据库表
func (w *PgWrapper) createSchema(models []interface{}) error {
	for _, model := range models {
		err := w.Db.Model(model).CreateTable(&orm.CreateTableOptions{
			//Temp: true,//建表是临时的
			IfNotExists: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

//通过结构体来删除表
func (w *PgWrapper) deleteSchema(models []interface{}) error {
	err := w.Db.Model(&models).DropTable(&orm.DropTableOptions{
		IfExists: true,
		Cascade:  true,
	})
	return err
}
