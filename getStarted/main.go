package main

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/haolipeng/go-pg-example/pgdb"
)

type User struct {
	Id     int64
	Name   string
	Emails []string
}

func main() {
	var (
		err         error
		pgsqlDB     *pg.DB
		result      pg.Result
		user1       *User
		updateUser  User
		delUser     User
		userList    []User
		queryResult []User
		models      []interface{}
	)

	pgsqlDB = pgdb.Connect()

	//莫忘记关闭数据库连接
	defer pgsqlDB.Close()

	//3.创建表
	models = []interface{}{
		(*User)(nil),
	}
	err = pgdb.CreateSchema(pgsqlDB, models)
	if err != nil {
		goto ERR
	}

	//4.single 插入一条记录
	user1 = &User{
		Id:     1,
		Name:   "admin",
		Emails: []string{"admin1@admin", "admin2@admin"},
	}
	result, err = pgsqlDB.Model(user1).Insert()
	if err != nil {
		goto ERR
	}
	fmt.Printf("single insert rows affected:%d\n", result.RowsAffected())

	//5.batch 批量插入多条记录
	userList = []User{
		{
			Id:     2,
			Name:   "haolipeng",
			Emails: []string{"1078285863@qq.com"},
		},
		{
			Id:     3,
			Name:   "haolipeng",
			Emails: []string{"haolipeng12345@163.com"},
		},
	}
	result, err = pgsqlDB.Model(&userList).Insert()
	if err != nil {
		goto ERR
	}
	fmt.Printf("batch insert rows affected:%d\n", result.RowsAffected())

	//6.查询
	err = pgsqlDB.Model(&queryResult).Select()
	if err != nil {
		goto ERR
	}
	fmt.Printf("query result:%v\n", queryResult)

	//7.修改
	//修改除主键外的其他列
	updateUser = User{
		Id:     1,
		Name:   "antiy",
		Emails: []string{"haolipeng@antiy.cn"},
	}
	result, err = pgsqlDB.Model(&updateUser).WherePK().Update()
	if err != nil {
		goto ERR
	}
	fmt.Printf("update rows affected:%d\n", result.RowsAffected())

	//8.删除记录(删除id为2的记录)
	delUser = User{
		Id: 2,
	}
	result, err = pgsqlDB.Model(&delUser).WherePK().Delete()
	if err != nil {
		goto ERR
	}
	fmt.Printf("delete rows affected:%d\n", result.RowsAffected())

	//9.将当前记录查询并都打印出来
	err = pgsqlDB.Model(&queryResult).Select()
	if err != nil {
		goto ERR
	}
	fmt.Printf("query result:%v\n", queryResult)
	return
ERR:
	fmt.Println("error:", err)
	return
}
