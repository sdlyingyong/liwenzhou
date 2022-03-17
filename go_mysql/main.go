package go_mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql" //用来驱动注册到全局map
)

func main() {
	// DSN:Data Source Name
	dsn := "user:password@tcp(127.0.0.1:3306)/dbname"
	//用来实际操作mysql的操作是三方库实现,所有库sql操作都是实现同一套方法,可以替换
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close() // 注意这行代码要写在上面err判断的下面
}
