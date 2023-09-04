package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func executeMultipleStatements(db *sql.DB, sqlScript string) error {
	statements := strings.Split(sqlScript, ";")

	for _, statement := range statements {
		trimmedStatement := strings.TrimSpace(statement)
		if trimmedStatement != "" {
			_, err := db.Exec(trimmedStatement)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func importSQLFile(filename, dsn string) {
	// 连接数据库
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 读取.sql文件内容
	sqlScript, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	if err := executeMultipleStatements(db, string(sqlScript)); err != nil {
		log.Fatal(err)
	}

	fmt.Println("SQL文件导入成功")
}

func main() {
	user := flag.String("u", "root", "username of mysql")
	passwd := flag.String("p", "root", "password of mysql")
	file := flag.String("f", "./integration/blank.sql", "sql file")
	flag.Parse()

	dsn := fmt.Sprintf("%s:%s@tcp(localhost)/", *user, *passwd)

	importSQLFile(*file, dsn)
}
