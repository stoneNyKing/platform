//
package apis

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"platform/common/utils"
	"platform/ousvc/config"
	uu "platform/ousvc/dbmodels"
	"time"
)

func init() {
	orm.RegisterModel(new(uu.User))
}

func GetUserLists(siteid int64, appid int64, stypes, contents []string, order, sort string, num, start int64) (interface{}, int64, error) {
	return getUserLC(1, siteid, appid, stypes, contents, order, sort, num, start)
}
func GetUserCount(siteid int64, appid int64, stypes, contents []string, order, sort string) (int64, error) {
	_, cnt, err := getUserLC(2, siteid, appid, stypes, contents, order, sort, 0, 0)

	return cnt, err
}

func getUserLC(cate int, siteid int64, appid int64, stypes, contents []string, order, sort string, num, start int64) (interface{}, int64, error) {

	var vs []orm.Params

	o := orm.NewOrm()

	var cnt int64
	var err error

	if sort == "" {
		sort = "desc"
	}

	if num <= 0 {
		num = PAGENUM_MAX
	}

	logger.Finest("siteid=%d,appid=%d,order=%s,sort=%s,num=%d,start=%d", siteid, appid, order, sort, num, start)

	if len(stypes) != len(contents) {
		return nil, 0, errors.New("params number is not match.")
	}

	conditions := ""

	l := len(stypes)
	var v string
	for i := 0; i < l; i++ {
		if contents[i] != "" {
			v = "'%" + contents[i] + "%'"
		}

		switch stypes[i] {
		case "1":
			v = " a.name like " + v
		case "2":
			v = " a.type =" + contents[i]
		case "3":
			v = " a.created >= '" + contents[i] + "'"
		case "4":
			v = " a.created < '" + contents[i] + "'"
		case "5":
			v = " a.phone like " + v
		case "6":
			v = " a.idcard like " + v
		case "7":
			v = " a.email like " + v
		}

		if v != "" {
			if i != l-1 {
				conditions = conditions + v + " and "
			} else {
				conditions = conditions + v
			}
		}
	}

	if conditions != "" {
		conditions = " and " + conditions
	}

	logger.Finest("conditions = %s", conditions)
	var sqlstr, statement string

	if cate == 1 {
		sqlstr = SQL_USER
	} else if cate == 2 {
		sqlstr = SQL_COUNT_USER
	}

	statement = fmt.Sprintf(sqlstr, config.Config.ObjectsSchema)

	if cate == 1 {
		if siteid > 1 {
			statement = fmt.Sprintf(statement+
				"where a.site_id=%d  %s order by a.id %s limit ? offset ?", siteid, conditions, sort)
		} else {
			statement = fmt.Sprintf(statement+
				"where 1=1 %s order by a.id %s limit ? offset ?", conditions, sort)
		}
		cnt, err = o.Raw(statement, num, start).Values(&vs)

		if err != nil {
			logger.Error("不能获取用户信息列表：%v", err.Error())
			return nil, 0, err
		}
	} else if cate == 2 {
		if siteid > 1 {
			statement = fmt.Sprintf(statement+
				"where a.site_id=%d  %s order by a.id %s limit 1 ", siteid, conditions, sort)
		} else {
			statement = fmt.Sprintf(statement+
				"where 1=1 %s order by a.id %s limit 1 ", conditions, sort)
		}
		cnt, err = o.Raw(statement).Values(&vs)

		if err != nil {
			logger.Error("不能获取用户信息列表：%v", err.Error())
			return nil, 0, err
		}
		if cnt > 0 {
			cnt = utils.Convert2Int64(vs[0]["ucount"])
		}
	}

	return vs, cnt, nil
}

func GetUser(id int64, siteid int64, appid int64) (interface{}, int64, error) {

	var vs []orm.Params
	var cnt int64
	var err error

	o := orm.NewOrm()

	logger.Finest("siteid=%d,appid=%d,areaid=%d", siteid, appid, id)

	statement := fmt.Sprintf(SQL_USER, config.Config.ObjectsSchema)

	if siteid > 1 {
		cnt, err = o.Raw(statement+
			"where a.site_id=?  and a.id =?", siteid, id).Values(&vs)
	} else {
		cnt, err = o.Raw(statement+
			"where a.id =?", id).Values(&vs)
	}

	if err != nil {
		logger.Error("不能获取用户信息：%v", err.Error())
		return nil, 0, err
	}

	return vs, cnt, nil
}

func PostUser(siteid int64, action string, param map[string]interface{}) (id int64, err error) {

	if action == "add" {
		return userAdd(siteid, param)
	} else if action == "readorcreate" {
		return userReadOrCreate(siteid, param)
	}

	return 0, errors.New("action is not match")
}

func userAdd(siteid int64, param map[string]interface{}) (id int64, err error) {
	if param == nil {
		return 0, errors.New("no input")
	}

	var v uu.User

	if param["id"] != nil {
		v.Id = utils.Convert2Int64(param["id"])
	}

	if param["siteid"] != nil {
		v.SiteId = utils.Convert2Int64(param["siteid"])
	} else {
		v.SiteId = siteid
	}

	if param["name"] != nil {
		v.Name = utils.ConvertToString(param["name"])
	}
	if param["type"] != nil {
		v.Type = utils.Convert2Int(param["type"])
	}
	if param["isactive"] != nil {
		v.IsActive = utils.Convert2Int(param["isactive"])
	}

	if param["phone"] != nil {
		v.Phone = utils.ConvertToString(param["phone"])
	}
	if param["imageurl"] != nil {
		v.ImageUrl = utils.ConvertToString(param["imageurl"])
	}

	if param["email"] != nil {
		v.Email = utils.ConvertToString(param["email"])
	}
	if param["idcard"] != nil {
		v.Idcard = utils.ConvertToString(param["idcard"])
	}
	if param["rfid"] != nil {
		v.Rfid = utils.ConvertToString(param["rfid"])
	}

	if param["imei"] != nil {
		v.Imei = utils.ConvertToString(param["imei"])
	}
	if param["nickname"] != nil {
		v.Nickname = utils.ConvertToString(param["nickname"])
	}

	if param["passwd"] != nil {
		v.Passwd = utils.ConvertToString(param["passwd"])
	} else {
		pwd := "123456"
		m5 := md5.Sum([]byte(pwd))
		str := fmt.Sprintf("%02x", m5)
		v.Passwd = str
	}
	if param["json"] != nil {
		v.Json = utils.ConvertToString(param["json"])
	}

	if param["lastlogin"] != nil {
		v.LastLogin, _ = time.Parse("2006-01-02 15:04:05", utils.ConvertToString(param["lastlogin"]))
	}
	if param["lastlogout"] != nil {
		v.LastLogout, _ = time.Parse("2006-01-02 15:04:05", utils.ConvertToString(param["lastlogout"]))
	}

	if param["sscard"] != nil {
		v.Sscard = utils.ConvertToString(param["sscard"])
	}

	if param["weixinid"] != nil {
		v.Weixinid = utils.ConvertToString(param["weixinid"])
	}

	if param["accounttype"] != nil {
		v.RegisterType = utils.Convert2Int(param["accounttype"])
	}

	v.Created = time.Now()
	v.Updated = v.Created

	o := orm.NewOrm()
	err = SetSearchPath(o, config.Config.ObjectsSchema)
	if err != nil {
		logger.Error("不能设置search_path :%v", err)
		return 0, err
	}

	id, err = o.Insert(&v)
	if err != nil {
		logger.Error("不能插入用户信息：%v", err.Error())
		return 0, err
	}

	return id, err
}

func userReadOrCreate(siteid int64, param map[string]interface{}) (id int64, err error) {
	if param == nil {
		return 0, errors.New("no input")
	}

	var v uu.User

	if param["id"] != nil {
		v.Id = utils.Convert2Int64(param["id"])
	}

	if param["siteid"] != nil {
		v.SiteId = utils.Convert2Int64(param["siteid"])
	} else {
		v.SiteId = siteid
	}

	if param["name"] != nil {
		v.Name = utils.ConvertToString(param["name"])
	}
	if param["type"] != nil {
		v.Type = utils.Convert2Int(param["type"])
	}
	if param["isactive"] != nil {
		v.IsActive = utils.Convert2Int(param["isactive"])
	}

	if param["phone"] != nil {
		v.Phone = utils.ConvertToString(param["phone"])
	}

	if param["email"] != nil {
		v.Email = utils.ConvertToString(param["email"])
	}
	if param["idcard"] != nil {
		v.Idcard = utils.ConvertToString(param["idcard"])
	}
	if param["rfid"] != nil {
		v.Rfid = utils.ConvertToString(param["rfid"])
	}

	if param["imei"] != nil {
		v.Imei = utils.ConvertToString(param["imei"])
	}
	if param["nickname"] != nil {
		v.Nickname = utils.ConvertToString(param["nickname"])
	}
	if param["imageurl"] != nil {
		v.ImageUrl = utils.ConvertToString(param["imageurl"])
	}

	if param["passwd"] != nil {
		v.Passwd = utils.ConvertToString(param["passwd"])
	} else {
		pwd := "123456"
		m5 := md5.Sum([]byte(pwd))
		str := fmt.Sprintf("%02x", m5)
		v.Passwd = str
	}
	if param["json"] != nil {
		v.Json = utils.ConvertToString(param["json"])
	}

	if param["lastlogin"] != nil {
		v.LastLogin, _ = time.Parse("2006-01-02 15:04:05", utils.ConvertToString(param["lastlogin"]))
	}
	if param["lastlogout"] != nil {
		v.LastLogout, _ = time.Parse("2006-01-02 15:04:05", utils.ConvertToString(param["lastlogout"]))
	}

	if param["weixinid"] != nil {
		v.Weixinid = utils.ConvertToString(param["weixinid"])
	}
	if param["sscard"] != nil {
		v.Sscard = utils.ConvertToString(param["sscard"])
	}

	v.Created = time.Now()
	v.Updated = v.Created

	var accounttype int

	if param["accounttype"] != nil {
		accounttype = utils.Convert2Int(param["accounttype"])
	}

	if accounttype == REGISTER_TYPE_NONE {
		accounttype = REGISTER_TYPE_PHONE
	}

	v.RegisterType = accounttype

	o := orm.NewOrm()
	err = SetSearchPath(o, config.Config.ObjectsSchema)
	if err != nil {
		logger.Error("不能设置search_path :%v", err)
		return 0, err
	}

	if v.Phone != "" && accounttype == REGISTER_TYPE_PHONE {
		_, id, err = o.ReadOrCreate(&v, "phone", "site_id")
	} else if v.Idcard != "" && accounttype == REGISTER_TYPE_IDCARD {
		_, id, err = o.ReadOrCreate(&v, "idcard", "site_id")
	} else if v.Name != "" && accounttype == REGISTER_TYPE_NAME {
		_, id, err = o.ReadOrCreate(&v, "name", "site_id")
	} else if v.Imei != "" && accounttype == REGISTER_TYPE_IMEI {
		_, id, err = o.ReadOrCreate(&v, "imei", "site_id")
	} else if v.Email != "" && accounttype == REGISTER_TYPE_EMAIL {
		_, id, err = o.ReadOrCreate(&v, "email", "site_id")
	} else if v.Sscard != "" && accounttype == REGISTER_TYPE_SSCARD {
		_, id, err = o.ReadOrCreate(&v, "sscard", "site_id")
	} else {
		logger.Error("不能插入用户信息：%v", err.Error())
		return 0, err
	}

	//o.Commit()
	return id, err
}

func PutUser(id int64, param map[string]interface{}) (cnt int64, err error) {

	if param == nil {
		return 0, errors.New("no input")
	}

	var v uu.User

	if param["id"] != nil {
		v.Id = utils.Convert2Int64(param["id"])
	}

	if id != v.Id {
		return 0, errors.New("id is not match.")
	}

	o := orm.NewOrm()
	err = SetSearchPath(o, config.Config.ObjectsSchema)
	if err != nil {
		logger.Error("不能设置search_path :%v", err)
		return 0, err
	}

	err = o.Read(&v)
	if err != nil {
		logger.Error("cannot read record.")
		return 0, err
	}

	if param["name"] != nil {
		v.Name = utils.ConvertToString(param["name"])
	}
	if param["type"] != nil {
		v.Type = utils.Convert2Int(param["type"])
	}
	if param["isactive"] != nil {
		v.IsActive = utils.Convert2Int(param["isactive"])
	}

	if param["phone"] != nil {
		v.Phone = utils.ConvertToString(param["phone"])
	}

	if param["email"] != nil {
		v.Email = utils.ConvertToString(param["email"])
	}
	if param["idcard"] != nil {
		v.Idcard = utils.ConvertToString(param["idcard"])
	}
	if param["rfid"] != nil {
		v.Rfid = utils.ConvertToString(param["rfid"])
	}

	if param["imei"] != nil {
		v.Imei = utils.ConvertToString(param["imei"])
	}
	if param["nickname"] != nil {
		v.Nickname = utils.ConvertToString(param["nickname"])
	}

	if param["passwd"] != nil {
		v.Passwd = utils.ConvertToString(param["passwd"])
	}

	if param["json"] != nil {
		v.Json = utils.ConvertToString(param["json"])
	}

	if param["lastlogin"] != nil {
		v.LastLogin, _ = time.Parse("2006-01-02 15:04:05", utils.ConvertToString(param["lastlogin"]))
	}
	if param["lastlogout"] != nil {
		v.LastLogout, _ = time.Parse("2006-01-02 15:04:05", utils.ConvertToString(param["lastlogout"]))
	}

	if param["weixinid"] != nil {
		v.Weixinid = utils.ConvertToString(param["weixinid"])
	}
	if param["imageurl"] != nil {
		v.ImageUrl = utils.ConvertToString(param["imageurl"])
	}
	if param["sscard"] != nil {
		v.Sscard = utils.ConvertToString(param["sscard"])
	}

	v.Updated = time.Now()

	cnt, err = o.Update(&v)
	if err != nil {
		logger.Error("不能更新用户信息：%v", err.Error())
		return 0, err
	}

	return cnt, nil
}

func DeleteUser(id int64, param map[string]interface{}) (num int64, err error) {

	o := orm.NewOrm()
	err = SetSearchPath(o, config.Config.ObjectsSchema)
	if err != nil {
		logger.Error("不能设置search_path :%v", err)
		return 0, err
	}

	var v uu.User

	var ids []interface{}

	if param["id"] != nil {
		ids = param["id"].([]interface{})
	}

	var cnt int64 = 0

	for _, rid := range ids {
		v.Id = utils.Convert2Int64(rid)

		num, err = o.Delete(&v)
		cnt += num
	}

	return cnt, err
}
