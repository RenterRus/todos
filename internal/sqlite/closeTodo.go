package sqlite

import (
	"fmt"
	"strconv"
)

func (d *Database) CloseTodo(id int) error {
	_, err := d.DBConnection.Exec("update "+d.DBTableName+" set Status = 'close' where id = '"+strconv.Itoa(id)+"'")
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	return nil
}