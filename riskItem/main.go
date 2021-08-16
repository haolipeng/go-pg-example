package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/haolipeng/go-pg-example/pgdb"
	"io/ioutil"
)

type SettingRiskItem struct {
	Table    string `pg:"table"`          //一级大项
	Item     string `pg:"item"`           //二类小项
	Err      string `pg:"err"`            //检测错误
	Expected string `pg:"expected"`       //检测期望值
	Desc     string `pg:"desc"`           //检测描述
	ResClass string `pg:"resource_class"` //检测类型 文件、文件夹、内核参数、etcd、系统参数等
	ResParam string `pg:"resource_param"` //检测对象
	FuncName string `pg:"func_name"`      //检测项对应的函数名称
}

func InsertRecordFromJsonfile(pgsqlDB *pg.DB, filePath string) error {
	var (
		bytes  []byte
		items  []SettingRiskItem
		err    error
		result pg.Result
	)

	//从json文件中读取数据
	//"account_risk_item.json" etc
	bytes, err = ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &items)
	if err != nil {
		return err
	}

	result, err = pgsqlDB.Model(&items).Insert()
	if err != nil {
		return err
	}
	fmt.Printf("single insert rows affected:%d\n", result.RowsAffected())

	return nil
}

func main() {
	var (
		err     error
		pgsqlDB *pg.DB
		models  []interface{}

		queryResult []SettingRiskItem
	)

	pgsqlDB = pgdb.Connect()

	//莫忘记关闭数据库连接
	defer pgsqlDB.Close()

	//3.创建表
	models = []interface{}{
		(*SettingRiskItem)(nil),
	}
	err = pgdb.CreateSchema(pgsqlDB, models)
	if err != nil {
		goto ERR
	}

	{
		accountFilePath := "account_risk_item.json"
		appFilePath := "app_risk_item.json"
		dockerFilePath := "docker_risk_item.json"
		kubernetesFilePath := "kubernetes_risk_item.json"
		systemFilePath := "system_risk_item.json"

		//account
		err = InsertRecordFromJsonfile(pgsqlDB, accountFilePath)
		//application
		err = InsertRecordFromJsonfile(pgsqlDB, appFilePath)
		//docker
		err = InsertRecordFromJsonfile(pgsqlDB, dockerFilePath)
		//kubernetes
		err = InsertRecordFromJsonfile(pgsqlDB, kubernetesFilePath)
		//system
		err = InsertRecordFromJsonfile(pgsqlDB, systemFilePath)
	}

	//6.查询
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
