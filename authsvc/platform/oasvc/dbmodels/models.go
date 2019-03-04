package dbmodels

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"platform/models/common"
)

var (
	orm    *xorm.Engine
	tables []interface{}
)

func init() {
	tables = append(tables, new(AdminLog), new(Admin),new(Appid))
}

func SetEngine(driver, host string, port int, user string, passwd string, database, schema string) (err error) {

	orm, err = common.SetEngine(driver, host, port, user, passwd, database, schema)
	return err
}

func NewEngine(driver, host string, port int, user string, passwd string, database, schema string) (err error) {
	if err = SetEngine(driver, host, port, user, passwd, database, schema); err != nil {
		return err
	}
	if err = orm.Sync(tables...); err != nil {
		return fmt.Errorf("sync database struct error: %v\n", err)
	}
	return nil
}

func Truncate() {
	orm.Exec("truncate table Admin")
	orm.Exec("truncate table Admin_Log")
}
