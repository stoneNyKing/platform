package storages

import (
	"fmt"
	// "os"
	// "path"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var (
	orm    *xorm.Engine
	tables []interface{}
)

func init() {
	tables = append(tables, new(Storage))
}

func SetEngine(host string, port int, user string, passwd string, database string) (err error) {
	orm, err = xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		user, passwd, host+":"+strconv.Itoa(port), database))
	if err != nil {
		return fmt.Errorf("models.init(fail to conntect database): %v", err)
	}

	// logPath := "./log/xorm.log"
	// os.MkdirAll(path.Dir(logPath), os.ModePerm)

	/*
		f, err := os.Create(logPath)
		if err != nil {
			return fmt.Errorf("models.init(fail to create xorm.log): %v", err)
		}
		orm.Logger = f
	*/

	// orm.ShowSQL = true
	// orm.ShowDebug = true
	// orm.ShowErr = true
	return nil
}

func NewEngine(host string, port int, user string, passwd string, database string) (err error) {
	if err = SetEngine(host, port, user, passwd, database); err != nil {
		return err
	}
	if err = orm.Sync(tables...); err != nil {
		return fmt.Errorf("sync database struct error: %v\n", err)
	}
	return nil
}

func Truncate() {
	orm.Exec("truncate table user")
}
