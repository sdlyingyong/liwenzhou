package main

import (
	"database/sql/driver"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	client = initDb()
)

type User struct {
	ID   int    `db:"id"`
	Age  int    `db:"age"`
	Name string `db:"name"`
}

func main() {
	queryRowDemo()
	queryMultiRow()
	updateRow()
	showNamedExec()
	showNamedQuery()
	showBatchInsert()
}

func (u User) Value() (driver.Value, error) {
	return []interface{}{u.Name, u.Age}, nil
}

func BatchInsert(users []User) error {
	query, args, _ := sqlx.In("INSERT INTO users (name,age) VALUES (?),(?),(?)",
		users)
	fmt.Println("query ", query)
	fmt.Println("args ", args)
	_, err := client.Exec(query, args...)
	return err
}

func showBatchInsert() {
	u1 := User{nil, 28, "七米"}
	u2 := User{nil, 18, "qimi"}
	u3 := User{nil, 38, "小王子"}
	users := []User{u1, u2, u3}
	BatchInsert(users)
}

func showNamedQuery() {
	//用map格式查询条件
	sqlStr := "SELECT * FROM users WHERE name=:name"
	queryData := map[string]interface{}{
		"name": "七米",
	}
	rows, err := client.NamedQuery(sqlStr, queryData)
	if err != nil {
		fmt.Println("NamedQuery fail, err :", err)
		return
	}
	defer rows.Close()
	//用结构体查询条件
	u := User{Name: "七米"}
	rows2, err := client.NamedQuery(sqlStr, u)
	if err != nil {
		fmt.Println("NamedQuery fail, err :", err)
		return
	}
	defer rows2.Close()
	//遍历取出
	for rows.Next() {
		var getU User
		err := rows.StructScan(&getU)
		if err != nil {
			fmt.Println("StructScan fail, err : ", err)
			return
		}
		fmt.Printf("user: %#v \n", getU)
	}

}

func showNamedExec() {
	insertSrt := `insert into users (name,age) values (:name,:age)`
	insertData := map[string]interface{}{
		"name": "七米",
		"age":  32,
	}
	ret, err := client.NamedExec(insertSrt, insertData)
	if err != nil {
		fmt.Println("insert fail, err :", err)
		return
	}
	num, err := ret.RowsAffected()
	if err != nil {
		fmt.Println("get rows affected fail, err :", err)
		return
	}
	fmt.Println("insert row num is ", num)
}

func updateRow() {
	sqlStr := "update users set age = ? where id = ?"
	//delStr := "delete from user where id = ? "
	//insertStr := "insert into user (name,age) values(?)"
	ret, err := client.Exec(sqlStr, 23, 2)
	if err != nil {
		fmt.Println("update fail,err :", err)
		return
	}
	nums, err := ret.RowsAffected()
	if err != nil {
		fmt.Println("get rows affected fail, err :", err)
		return
	}
	fmt.Println("exec success,ret :", nums)
}

func queryRowDemo() {
	sqlStr := "select id,name ,age from users where id = ?"
	var u1 User
	//1.需要首字母大写 2.需要结构体db:"字段"
	//传指针
	err := client.Get(&u1, sqlStr, 1)
	if err != nil {
		fmt.Println("query fail,err :", err)
		return
	}
	fmt.Println("u1 is", u1)
}

func queryMultiRow() {
	sqlStr := "select id ,name, age from users where id > ?"
	var ulist []User
	err := client.Select(&ulist, sqlStr, 0)
	if err != nil {
		fmt.Println("query fail,err :", err)
		return
	}
	fmt.Println("ulist are ", ulist)
}

func initDb() (db *sqlx.DB) {
	dsn := "root:root@tcp(127.0.0.1:3306)/sqlx?charset=utf8mb4&parseTime=True"
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		panic(err)
	}
	//必须连接 出问题就直接panic
	db = sqlx.MustConnect("mysql", dsn)

	//最大连接数
	db.SetMaxOpenConns(20)
	//最大处理数
	db.SetMaxIdleConns(10)
	fmt.Println("conn mysql success")
	return
}
