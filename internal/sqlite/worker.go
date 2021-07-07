package sqlite

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"time"
)

var DBClient *Database

type Database struct{
	DBName string
	DBTableName string
	DBConnection *sql.DB
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

func (d *Database) Write(message, status string, deadline timestamp.Timestamp) error{
	d.ExpTimeoutMech()
	_, err := d.DBConnection.Exec("insert into '"+d.DBTableName+"' (Message, Status, Deadline) values ('"+message+"', '"+status+"', '"+deadline.String()+"')")
	if err != nil{
		fmt.Println(err)
		return err
	}

	return nil
}

type Todo struct{
	Id int32 `json:"Id"`
	Message string `json:"Message"`
	Deadline int64 `json:"Deadline"`
	Status string `json:"Status"`
}

func (d *Database) ReadAll()([]byte, error, ){
	d.ExpTimeoutMech()
	result, err := d.DBConnection.Query("select * from "+d.DBTableName+" where (Status = 'open' OR Status = 'deadline')")
	if err != nil{
		return nil, err
	}

	var todos []interface{}

	for result.Next(){
		t := Todo{}
		deadline := ""
		err := result.Scan(&t.Id, &t.Message, &deadline, &t.Status)
		deadlinebuf, _ := strconv.Atoi(deadline)
		t.Deadline = int64(deadlinebuf)
		if err != nil{
			fmt.Println(err)
			continue
		}
		fmt.Printf("Id: %v, Message: %s, Deadline: %v, Status: %s\n", t.Id, t.Message, t.Deadline, t.Status)
		todos = append(todos, t)
	}

	res, err := json.Marshal(todos)
	if err != nil{
		panic(err.Error())
	}
	return res, err
}

func (d *Database) CloseTodo(id int) error {
	_, err := d.DBConnection.Exec("update "+d.DBTableName+" set Status = 'close' where id = '"+strconv.Itoa(id)+"'")
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (d *Database) ExpTimeoutMech(){
		var deadlist []int
		result, _ := d.DBConnection.Query("select * from " + d.DBTableName + " where (Status = 'deadline')")
		for result.Next() {
			t := Todo{}
			deadline := ""
			err := result.Scan(&t.Id, &t.Message, &deadline, &t.Status)
			deadlinebuf, _ := strconv.Atoi(deadline)
			t.Deadline = int64(deadlinebuf)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if t.Deadline-time.Now().Unix() <= 0 {
				deadlist = append(deadlist, int(t.Id))
			}
		}
		for _, v := range deadlist {
			d.CloseTodo(v)
		}
}

func (d *Database) SetExpTimeout(id int, deadline int64) error{
	d.ExpTimeoutMech()
	fmt.Println("DB")
	now := time.Now()
	dur := deadline - now.Unix()

	fmt.Println(dur)
	fmt.Println(time.Now().Unix() + dur)
	if dur < 0{
		return errors.New("The date is not correct or is the past")
	}
	fmt.Println(deadline)
	fmt.Printf("Duration for id: %v %v", id, dur)
	_, err := d.DBConnection.Exec("update "+d.DBTableName+" set Status = 'deadline', Deadline = '"+strconv.Itoa(int(deadline))+"' where id = '"+strconv.Itoa(id)+"'")
	if err != nil{
		return err
	}
	return nil
}

func (d *Database) UpdateTodo(id int, message string) error{
	d.ExpTimeoutMech()
	_, err := d.DBConnection.Exec("update "+d.DBTableName+" set Message = '"+message+"' where id = '"+strconv.Itoa(id)+"'")
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	return nil
}