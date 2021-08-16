package main

import (
	"errors"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/haolipeng/go-pg-example/conf"
)

func main() {
	var (
		err     error
		pgsqlDB *pg.DB = nil
	)

	type Item struct {
		Id    int64
		Attrs map[string]string `pg:",hstore"` // marshalled as PostgreSQL hstore
	}

	//1.连接数据库
	pgsqlDB = pg.Connect(&pg.Options{
		Addr:     conf.DbAddr,
		User:     conf.User,
		Password: conf.Password,
		Database: conf.DbName,
	})
	if pgsqlDB == nil {
		err = errors.New("pg.Connect() failed,error:")
		fmt.Println(err)
	}

	//2、忘记关闭数据库连接
	defer func(pgsqlDB *pg.DB) {
		err = pgsqlDB.Close()
		if err != nil {
			fmt.Println("close postgresql failed")
		}
	}(pgsqlDB)

	_, err = pgsqlDB.Exec(`CREATE TEMP TABLE items (id serial, attrs hstore)`)
	if err != nil {
		panic(err)
	}
	defer pgsqlDB.Exec("DROP TABLE items")

	item1 := Item{
		Id:    1,
		Attrs: map[string]string{"hello": "world"},
	}
	_, err = pgsqlDB.Model(&item1).Insert()
	if err != nil {
		panic(err)
	}

	var item Item
	err = pgsqlDB.Model(&item).Where("id = ?", 1).Select()
	if err != nil {
		panic(err)
	}
	fmt.Println(item)
}
