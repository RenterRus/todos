package sqlite

import (
	"fmt"
	"github.com/golang/protobuf/ptypes/timestamp"
)

func (d *Database) Write(message, status string, deadline timestamp.Timestamp) error{
	d.ExpTimeoutMech()
	_, err := d.DBConnection.Exec("insert into '"+d.DBTableName+"' (Message, Status, Deadline) values ('"+message+"', '"+status+"', '"+deadline.String()+"')")
	if err != nil{
		fmt.Println(err)
		return err
	}

	return nil
}