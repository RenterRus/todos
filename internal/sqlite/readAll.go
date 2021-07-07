package sqlite

import (
	"encoding/json"
	"fmt"
	"strconv"
)

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