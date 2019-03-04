package common


import (

	"github.com/go-xorm/xorm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"fmt"
	"strconv"
)


func SetEngine(driver,host string, port int, user string, passwd string, database,schema string) (orm *xorm.Engine,err error) {

	if driver == "mysql" {
		orm, err = xorm.NewEngine("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
			user, passwd, host+":"+strconv.Itoa(port), database))
		if err != nil {
			return nil,fmt.Errorf("models.init(fail to conntect database): %v", err)
		}

	}else if driver == "pgsql" {
		//consgtr: "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full"
		orm, err = xorm.NewEngine("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?connect_timeout=10&sslmode=disable",
			user, passwd, host+":"+strconv.Itoa(port), database))
		if err != nil {
			return nil,fmt.Errorf("models.init(fail to conntect database): %v", err)
		}
		smt := fmt.Sprintf("SET SEARCH_PATH TO \"%s\"",schema)
		orm.Exec(smt)
	}

	return orm,nil
}

func NewEngine(driver,host string, port int, user string, passwd string, database,schema string) (orm *xorm.Engine,err error) {

	if orm,err = SetEngine(driver,host, port, user, passwd, database,schema); err != nil {
		return nil,err
	}

	return orm,nil
}
