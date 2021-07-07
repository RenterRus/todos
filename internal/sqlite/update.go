package sqlite

import (
	"fmt"
	"strconv"
)

func (d *Database) UpdateTodo(id int, message string) error{
	d.ExpTimeoutMech()
	_, err := d.DBConnection.Exec("update "+d.DBTableName+" set Message = '"+message+"' where id = '"+strconv.Itoa(id)+"'")
	if err != nil{
		fmt.Println(err.Error())
		return err
	}
	return nil
}