package users

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"platform/models/common"
)

var (
	orm    *xorm.Engine
	tables []interface{}
)

func init() {
	tables = append(tables, new(User), new(UserLog), new(NameId), new(Feedback))
}


func SetEngine(driver,host string, port int, user string, passwd string, database,schema string) (err error) {

	orm,err = common.SetEngine(driver,host,port,user,passwd,database,schema)
	return err
}

func NewEngine(driver,host string, port int, user string, passwd string, database,schema string) (err error) {
	if err = SetEngine(driver,host, port, user, passwd, database,schema); err != nil {
		return err
	}
	if err = orm.Sync(tables...); err != nil {
		return fmt.Errorf("sync database struct error: %v\n", err)
	}
	return nil
}

func Truncate() {
	orm.Exec("truncate table user")
	orm.Exec("truncate table User_Log")
	orm.Exec("truncate table Name_Id")
	orm.Exec("truncate table Feedback")
}
