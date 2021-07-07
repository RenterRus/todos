package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var DBClient *Database

type Database struct{
	DBName string
	DBTableName string
	DBConnection *sql.DB
}

type Todo struct{
	Id int32 `json:"Id"`
	Message string `json:"Message"`
	Deadline int64 `json:"Deadline"`
	Status string `json:"Status"`
}

func Initial(dbname, tablename string) *Database{
	d := new(Database)
	d.DBName = dbname
	d.DBTableName = tablename
	var err error
	d.DBConnection, err = sql.Open("sqlite3", d.DBName)
	if err != nil {
		panic(err.Error())
	}
	return d
}

func (d *Database) DisableConnect(){
	d.DBConnection.Close()
}