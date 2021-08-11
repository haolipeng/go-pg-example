package main

import (
	"errors"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/haolipeng/go-pg-example/conf"
	"github.com/haolipeng/go-pg-example/pgdb"
)

type Profile struct {
	tableName struct{} `pg:"profile"`
	ID        int      `pg:"id"`
	Lang      string
}

type User struct {
	tableName struct{} `pg:"users"`
	ID        int      `pg:"user_id"`
	Name      string
	ProfileID int
	Profile   *Profile `pg:"rel:has-one"` //note the has-one relation
}

func main() {
	var (
		err     error
		pgsqlDB *pg.DB = nil
		models  []interface{}
		result  pg.Result
	)

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

	//3.先删除以前的表,再创建新表
	result, err = pgsqlDB.Exec("DROP TABLE users;")
	if err != nil {
		fmt.Println("drop table users failed:", err)
	}

	result, err = pgsqlDB.Exec("DROP TABLE profile;")
	if err != nil {
		fmt.Println("drop table profile failed:", err)
	}

	models = []interface{}{
		(*Profile)(nil),
		(*User)(nil),
	}
	err = pgdb.CreateSchema(pgsqlDB, models)
	if err != nil {
		fmt.Println("createSchema failed:", err)
	}

	//4.插入一条profile记录
	profiles := []Profile{
		{
			ID:   1,
			Lang: "chinese",
		},
		{
			ID:   2,
			Lang: "english",
		},
		{
			ID:   3,
			Lang: "japan",
		},
	}
	result, err = pgsqlDB.Model(&profiles).Insert()
	if err != nil {
		fmt.Println("Profile Insert failed:", err)
		return
	}
	fmt.Printf("insert Profile record affected:%d\n", result.RowsAffected())

	//5.插入一条user记录
	user1 := &User{
		ID:        1,
		Name:      "haolipeng",
		ProfileID: 1,
		Profile:   &Profile{ID: 1, Lang: "chinese"},
	}
	result, err = pgsqlDB.Model(user1).Insert()
	if err != nil {
		fmt.Println("User Insert failed:", err)
		return
	}
	fmt.Printf("insert User record affected:%d\n", result.RowsAffected())

	//6.利用Relation关系查找记录
	var user User
	err = pgsqlDB.Model(&user).Relation("Profile").Where("user_id = ?", 1).Select()
	//err = pgsqlDB.Model(&user).Relation("Profile").Select()
	if err != nil {
		fmt.Println("Select with Relation failed:", err)
		return
	}

	fmt.Println("user:", user)
	fmt.Println("program exit normal!")
}
