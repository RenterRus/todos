package sqlite

import (
	"fmt"
	"strconv"
	"time"
)

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