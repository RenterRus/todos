package sqlite

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

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