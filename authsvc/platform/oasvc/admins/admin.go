//
package admins

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"platform/common/redis"
	"platform/common/utils"
	"platform/oasvc/config"
	. "platform/oasvc/dbmodels"
	"platform/pfcomm/apis"
	"strconv"
	"time"
)

func init() {
	orm.RegisterModel(new(Admin))
}

func GetAdminLists(siteid int64, appid int64, stypes, contents []string, order, sort string, num, start int64) (interface{}, int64, error) {
	return getAdminLC(1, siteid, appid, stypes, contents, order, sort, num, start)
}

func GetAdminCount(siteid int64, appid int64, stypes, contents []string, order, sort string) (int64, error) {

	_, cnt, err := getAdminLC(2, siteid, appid, stypes, contents, order, sort, 0, 0)

	return cnt, err
}

func getAdminLC(cate int, siteid int64, appid int64, stypes, contents []string, order, sort string, num, start int64) (interface{}, int64, error) {

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
			v = " a.job_number like " + v
		case "3":
			v = " a.email like " + v
		case "4":
			v = " a.phone like " + v
		case "5":
			v = " a.effective_time >= '" + contents[i] + "'"
		case "6":
			v = " a.effective_time < '" + contents[i] + "'"
		case "7":
			v = " a.expire_time >= '" + contents[i] + "'"
		case "8":
			v = " a.expire_time < '" + contents[i] + "'"
		case "9":
			v = " a.type =" + contents[i]
		case "10":
			v = " a.created >= '" + contents[i] + "'"
		case "11":
			v = " a.created < '" + contents[i] + "'"
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
		sqlstr = SQL_ADMIN
	} else if cate == 2 {
		sqlstr = SQL_COUNT_ADMIN
	}

	statement = fmt.Sprintf(sqlstr, config.Config.ObjectsSchema, config.Config.ObjectsSchema,
		config.Config.ObjectsSchema, config.Config.ObjectsSchema)

	if cate == 1 {
		if siteid > 1 {
			statement = fmt.Sprintf(statement+
				"where a.site_id=%d %s order by a.id %s limit ? offset ?", siteid, conditions, sort)
		} else {
			statement = fmt.Sprintf(statement+
				"where 1=1 %s order by a.id %s limit ? offset ?", conditions, sort)
		}
		cnt, err = o.Raw(statement, num, start).Values(&vs)

		if err != nil {
			logger.Error("不能获取用户(admin)信息列表：%v", err.Error())
			return nil, 0, err
		}

	} else if cate == 2 {
		if siteid > 1 {
			statement = fmt.Sprintf(statement+
				"where a.site_id=%d %s order by a.id %s limit 1 ", siteid, conditions, sort)
		} else {
			statement = fmt.Sprintf(statement+
				"where 1=1 %s order by a.id %s limit 1 ", conditions, sort)
		}

		cnt, err = o.Raw(statement).Values(&vs)

		if err != nil {
			logger.Error("不能获取列表数量：%v", err.Error())
			return nil, 0, err
		}

		if cnt > 0 {
			cnt = utils.Convert2Int64(vs[0]["ucount"])
		}

	}

	return vs, cnt, nil
}

func GetAdmin(id int64, siteid int64, appid int64) (interface{}, int64, error) {

	var vs []orm.Params

	o := orm.NewOrm()

	var cnt int64
	var err error

	logger.Finest("siteid=%d,appid=%d,areaid=%d", siteid, appid, id)

	statement := fmt.Sprintf(SQL_ADMIN, config.Config.ObjectsSchema, config.Config.ObjectsSchema,
		config.Config.ObjectsSchema, config.Config.ObjectsSchema)

	if siteid > 1 {
		cnt, err = o.Raw(statement+
			"where a.site_id=? and a.id =?", siteid, id).Values(&vs)
	} else {
		cnt, err = o.Raw(statement+
			"where a.id =?", id).Values(&vs)
	}

	if err != nil {
		logger.Error("不能获取用户(ADMIN)信息：%v", err.Error())
		return nil, 0, err
	}

	var maps []map[string]interface{}
	if cnt > 0 {
		for _, v := range vs {
			m := make(map[string]interface{})

			for key, value := range v {
				m[key] = value
			}
			if v["roleid"] != nil {
				m["subsys"] = getRoleSubsys(o, siteid, appid, utils.Convert2Int64(v["roleid"]))
			}
			maps = append(maps, m)
		}
	}

	return maps, cnt, nil
}

func getRoleSubsys(o orm.Ormer, siteid, appid, roleid int64) interface{} {

	var vs []orm.Params

	statement := fmt.Sprintf("select a.role_resource_id as id,a.appid,a.sub_sys_id as subsysid,a.res_id as resid,ifnull(ss.subtplid,0) as subtplid,ss.icon_url as iconurl " +
		"from role_resource a " +
		"left join sub_sys ss on a.sub_sys_id=ss.sub_sys_id " +
		"where a.role_id=? ")

	_, err := o.Raw(statement, roleid).Values(&vs)
	if err != nil {
		logger.Error("不能获取角色子系统：%v", err)
	}

	for _, v := range vs {
		if v["resid"] != nil {
			var stypes, contents []string
			stypes = append(stypes, "4")
			resid := utils.Convert2Int64(v["resid"])
			contents = append(contents, strconv.FormatInt(resid, 10))
			v["resources"], _ = RpcxGetRoleResTree(siteid, appid, "", "", stypes, contents)
		}
	}

	return vs
}

func RpcxGetRoleResTree(siteid int64, appid int64, order, sort string, stypes []string, contents []string) (interface{}, error) {
	param := make(map[string]interface{})
	param["siteid"] = siteid
	param["stypes"] = stypes
	param["contents"] = contents

	_, r, err := apis.RpcxGetService(appid, siteid, "", "SysmgrJSONRpc", "GetRoleResTree", config.Config.RpcxSysmgrBasepath,
		config.Config.ConsulAddress, param)

	return r, err
}

func GetRoleResTree(siteid int64, appid int64, order, sort string, stypes []string, contents []string) ([]map[string]interface{}, error) {
	var maps []map[string]interface{}

	var err error
	var cnt int64

	if sort == "" {
		sort = "asc"
	}

	if len(stypes) != len(contents) {
		return nil, errors.New("params number is not match.")
	}

	logger.Finest("siteid=%d,appid=%d,order=%s,sort=%s", siteid, appid, order, sort)

	conditions := ""

	res_id := ""

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
			v = " a.parent_id =" + contents[i]
		case "3":
			v = " a.level =" + contents[i]
		case "4":
			v = ""
			res_id = contents[i]
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

	if res_id != "" {

		statement := fmt.Sprintf("select a.res_id as id,a.parent_id as parentid,a.name,a.resource_id as resid,a.treeid,a.status, "+
			"a.start_time as starttime,a.end_time as endtime,a.level,a.order, "+
			"a.perm_sel,a.perm_add,a.perm_upd,a.perm_del,a.perm_cancel,a.perm_audit, "+
			"a.perm_eval,a.perm_doc, "+
			"b.name as resname,b.url,b.proxy,b.type,b.icon,b.perm_desc as permdesc "+
			"from role_res a "+
			"left join resource b on a.resource_id=b.id "+
			"where a.res_id=%s order by a.order asc; ", res_id)

		//logger.Finest("statement=%s", statement)

		o := orm.NewOrm()
		var vs []orm.Params

		cnt, err = o.Raw(statement).Values(&vs)

		if err != nil {
			logger.Error("不能获取角色资源tree信息列表：%v", err.Error())
			return nil, err
		}

		if cnt > 0 {
			for _, v := range vs {
				m := make(map[string]interface{})

				for key, value := range v {
					m[key] = value
				}

				child, _, _ := GetRoleResTreeChild(siteid, res_id)
				m["child"] = child
				maps = append(maps, m)
			}
		}
	}

	return maps, err
}

func GetRoleResTreeChild(siteid int64, resids string) ([]map[string]interface{}, int64, error) {
	var maps []map[string]interface{}

	var vs []orm.Params

	o := orm.NewOrm()

	var cnt int64
	var err error

	statement := fmt.Sprintf("select a.res_id as id,a.parent_id as parentid,a.name,a.resource_id as resid,a.treeid,a.status, "+
		"a.start_time as starttime,a.end_time as endtime,a.level,a.order, "+
		"a.perm_sel,a.perm_add,a.perm_upd,a.perm_del,a.perm_cancel,a.perm_audit, "+
		"a.perm_eval,a.perm_doc, "+
		"b.name as resname,b.url,b.proxy,b.type,b.icon,b.perm_desc as permdesc "+
		"from role_res a "+
		"left join resource b on a.resource_id=b.id "+
		"where a.parent_id=%s order by a.order asc; ", resids)

	//logger.Finest("statement=%s", statement)

	cnt, err = o.Raw(statement).Values(&vs)

	if err != nil {
		logger.Error("不能获取角色资源tree信息列表：%v", err.Error())
		return nil, 0, err
	}

	var res_ids string = ""
	if cnt > 0 {
		for _, v := range vs {
			m := make(map[string]interface{})

			for key, value := range v {
				m[key] = value
			}
			if v["id"] != nil {
				res_ids = v["id"].(string)
				if res_ids != "" {
					pas, count, err := GetRoleResTreeChild(siteid, res_ids)

					if err == nil && count > 0 {
						m["child"] = pas
					}
				} else {
					m["child"] = "[]"
				}
			}
			maps = append(maps, m)
		}

	}

	return maps, cnt, nil
}

func PostAdmin(siteid int64, action string, param map[string]interface{}) (id int64, err error) {

	if action == "add" || action == "registry" {
		return adminAdd(siteid, param)
	} else if action == "readorcreate" {
		return adminReadOrCreate(siteid, param)
	}

	return 0, errors.New("action is not match")
}

func adminAdd(siteid int64, param map[string]interface{}) (id int64, err error) {
	if param == nil {
		return 0, errors.New("no input")
	}

	var v Admin

	if param["id"] != nil {
		v.Id = utils.Convert2Int64(param["id"])
	}

	if param["siteid"] != nil {
		v.SiteId = utils.Convert2Int64(param["siteid"])
	} else {
		v.SiteId = siteid
	}

	if param["roleid"] != nil {
		v.RoleId = utils.Convert2Int64(param["roleid"])
	} else {
		v.RoleId = 0
	}

	if param["organizationid"] != nil {
		v.OrganizationId = utils.Convert2Int64(param["organizationid"])
	}

	if param["name"] != nil {
		v.Name = utils.ConvertToString(param["name"])
	}
	if param["type"] != nil {
		v.Type = utils.Convert2Int(param["type"])
	}

	if param["realname"] != nil {
		v.Realname = utils.ConvertToString(param["realname"])
	}
	if param["gender"] != nil {
		v.Gender = utils.Convert2Int(param["gender"])
	}

	if param["phone"] != nil {
		v.Phone = utils.ConvertToString(param["phone"])
	}

	if param["email"] != nil {
		v.Email = utils.ConvertToString(param["email"])
	}
	if param["jobnumber"] != nil {
		v.JobNumber = utils.ConvertToString(param["jobnumber"])
	} else {
		uk := JOBNUMBER_PREFIX + strconv.FormatInt(siteid, 10)
		v.JobNumber = redis.GetIncr(uk)

		logger.Finest("uk=%s,Jobnumber = %s", uk, v.JobNumber)
	}

	if param["imageurl"] != nil {
		v.ImageUrl = utils.ConvertToString(param["imageurl"])
	}

	if param["effectivetime"] != nil {
		v.EffectiveTime, _ = time.Parse("2006-01-02 15:04:05", utils.ConvertToString(param["effectivetime"]))
	}

	if param["expiretime"] != nil {
		v.ExpireTime, _ = time.Parse("2006-01-02 15:04:05", utils.ConvertToString(param["expiretime"]))
	}

	if param["passwd"] != nil {
		v.Passwd = utils.ConvertToString(param["passwd"])
	} else {
		pwd := "123456"
		m5 := md5.Sum([]byte(pwd))
		str := fmt.Sprintf("%02x", m5)
		v.Passwd = str
	}
	if param["description"] != nil {
		v.Description = utils.ConvertToString(param["description"])
	}

	if param["state"] != nil {
		v.State = utils.Convert2Int(param["state"])
	}

	v.Created ,_ = time.Parse("2006-01-02 15:04:05",utils.GetTimeFormat("2006-01-02 15:04:05"))
	v.Updated = v.Created

	o := orm.NewOrm()
	err = SetSearchPath(o, config.Config.ObjectsSchema)
	if err != nil {
		logger.Error("不能设置search_path :%v", err)
		return 0, err
	}

	//err = duplicateJobnumber(o,v.SiteId,v.JobNumber)
	//if err != nil {
	//	return 0,err
	//}
	var vs []orm.Params
	cnt, err := o.Raw("select id from admin where phone=? ", v.Phone).Values(&vs)
	if err == nil && cnt > 0 {
		return 0, errors.New("电话号码重复。")
	}

	/* modified by yangr 用户可以重名
	cnt, err = o.Raw("select id from admin where name=? and site_id=? ", v.Name, v.SiteId).Values(&vs)
	if err == nil && cnt > 0 {
		return 0, errors.New("用户姓名重复。")
	}*/

	id, err = o.Insert(&v)
	if err != nil {
		logger.Error("不能插入用户(admin)信息：%v", err.Error())
		return 0, err
	}

	return id, err
}

func adminReadOrCreate(siteid int64, param map[string]interface{}) (id int64, err error) {
	if param == nil {
		return 0, errors.New("no input")
	}

	var v Admin

	if param["id"] != nil {
		v.Id = utils.Convert2Int64(param["id"])
	}

	if param["siteid"] != nil {
		v.SiteId = utils.Convert2Int64(param["siteid"])
	} else {
		v.SiteId = siteid
	}

	if param["roleid"] != nil {
		v.RoleId = utils.Convert2Int64(param["roleid"])
	} else {
		v.RoleId = 0
	}
	if param["organizationid"] != nil {
		v.OrganizationId = utils.Convert2Int64(param["organizationid"])
	}

	if param["name"] != nil {
		v.Name = param["name"].(string)
	}
	if param["type"] != nil {
		v.Type = utils.Convert2Int(param["type"])
	}
	if param["realname"] != nil {
		v.Realname = utils.ConvertToString(param["realname"])
	}
	if param["gender"] != nil {
		v.Gender = utils.Convert2Int(param["gender"])
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
	if param["jobnumber"] != nil {
		v.JobNumber = utils.ConvertToString(param["jobnumber"])
	} else {
		uk := JOBNUMBER_PREFIX + strconv.FormatInt(siteid, 10)
		v.JobNumber = redis.GetIncr(uk)
	}

	if param["effectivetime"] != nil {
		v.EffectiveTime, _ = time.Parse("2006-01-02 15:04:05", utils.ConvertToString(param["effectivetime"]))
	}

	if param["expiretime"] != nil {
		v.ExpireTime, _ = time.Parse("2006-01-02 15:04:05", utils.ConvertToString(param["expiretime"]))
	}

	if param["passwd"] != nil {
		v.Passwd = utils.ConvertToString(param["passwd"])
	} else {
		pwd := "123456"
		m5 := md5.Sum([]byte(pwd))
		str := fmt.Sprintf("%02x", m5)
		v.Passwd = str
	}
	if param["description"] != nil {
		v.Description = utils.ConvertToString(param["description"])
	}

	if param["state"] != nil {
		v.State = utils.Convert2Int(param["state"])
	}

	v.Created,_ = time.Parse("2006-01-02 15:04:05",utils.GetTimeFormat("2006-01-02 15:04:05"))
	v.Updated = v.Created

	o := orm.NewOrm()
	err = SetSearchPath(o, config.Config.ObjectsSchema)
	if err != nil {
		logger.Error("不能设置search_path :%v", err)
		return 0, err
	}

	if v.Phone != "" {
		_, id, err = o.ReadOrCreate(&v, "phone")
	} else {
		logger.Error("不能插入用户(admin)信息：%v", err.Error())
		return 0, err
	}

	return id, err
}

func duplicateJobnumber(o orm.Ormer, siteid int64, jobno string) error {
	var vs []orm.Params
	cnt, err := o.Raw("select job_number from admin where site_id=? and job_number=? ", siteid, jobno).Values(&vs)

	if err != nil {
		if err != orm.ErrNoRows {
			return errors.New("操作数据库出错。")
		}
	}
	if cnt > 0 {
		return errors.New("员工工号重复，请重新输入。")
	}

	return nil
}

func PutAdmin(id int64, param map[string]interface{}) (cnt int64, err error) {

	if param == nil {
		return 0, errors.New("no input")
	}

	var v Admin

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

	if param["roleid"] != nil {
		v.RoleId = utils.Convert2Int64(param["roleid"])
	}
	if param["organizationid"] != nil {
		v.OrganizationId = utils.Convert2Int64(param["organizationid"])
	}

	if param["name"] != nil {
		v.Name = param["name"].(string)
	}
	if param["type"] != nil {
		v.Type = utils.Convert2Int(param["type"])
	}

	if param["realname"] != nil {
		v.Realname = utils.ConvertToString(param["realname"])
	}
	if param["gender"] != nil {
		v.Gender = utils.Convert2Int(param["gender"])
	}

	if param["phone"] != nil {
		v.Phone = utils.ConvertToString(param["phone"])
	}

	if param["email"] != nil {
		v.Email = utils.ConvertToString(param["email"])
	}
	if param["jobnumber"] != nil {
		v.JobNumber = utils.ConvertToString(param["jobnumber"])
	}
	if param["effectivetime"] != nil {
		t1, _ := time.Parse("2006-01-02 15:04:05", utils.ConvertToString(param["effectivetime"]))
		if t1.IsZero() {
			v.EffectiveTime, _ = time.Parse("2006-01-02 15:04:05", v.EffectiveTime.Format("2006-01-02 15:04:05"))
		}else{
			v.EffectiveTime = t1
		}
	}else{
		v.EffectiveTime, _ = time.Parse("2006-01-02 15:04:05", v.EffectiveTime.Format("2006-01-02 15:04:05"))
	}

	if param["expiretime"] != nil {
		t1, _ := time.Parse("2006-01-02 15:04:05", utils.ConvertToString(param["expiretime"]))
		if t1.IsZero() {
			v.ExpireTime, _ = time.Parse("2006-01-02 15:04:05", v.ExpireTime.Format("2006-01-02 15:04:05"))
		}else{
			v.ExpireTime = t1
		}
	}else{
		v.ExpireTime, _ = time.Parse("2006-01-02 15:04:05", v.ExpireTime.Format("2006-01-02 15:04:05"))
	}

	if param["passwd"] != nil {
		v.Passwd = utils.ConvertToString(param["passwd"])
	}

	if param["description"] != nil {
		v.Description = utils.ConvertToString(param["description"])
	}

	if param["imageurl"] != nil {
		v.ImageUrl = utils.ConvertToString(param["imageurl"])
	}

	v.Created,_ = time.Parse("2006-01-02 15:04:05",v.Created.Format("2006-01-02 15:04:05"))
	v.Updated,_ = time.Parse("2006-01-02 15:04:05",utils.GetTimeFormat("2006-01-02 15:04:05"))

	if param["state"] != nil {
		v.State = utils.Convert2Int(param["state"])
	}

	cnt, err = o.Update(&v)
	if err != nil {
		logger.Error("不能更新用户信息：%v", err.Error())
		return 0, err
	}

	return cnt, nil
}

func DeleteAdmin(id int64, param map[string]interface{}) (num int64, err error) {

	o := orm.NewOrm()
	err = SetSearchPath(o, config.Config.ObjectsSchema)
	if err != nil {
		logger.Error("不能设置search_path :%v", err)
		return 0, err
	}

	var v Admin

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

func InitUKJobNumber() {

	var vs []orm.Params
	var sites []orm.Params
	var siteid int64
	var cnt int64
	var err error

	o := orm.NewOrm()
	err = SetSearchPath(o, config.Config.ObjectsSchema)
	if err != nil {
		logger.Error("不能设置search_path :%v", err)
		return
	}

	count, err := o.Raw("select id from site  ").Values(&sites)

	if err != nil {
		logger.Error("不能获取site信息：%v", err.Error())
		return

	}

	if count <= 0 {
		logger.Info("没有对应的租户记录。")
		return
	}
	for _, site := range sites {
		if site["id"] != nil {
			siteid = utils.Convert2Int64(site["id"])

			cnt, err = o.Raw("select a.job_number as jobnumber,a.site_id as siteid from admin a where a.site_id=? order by a.job_number desc limit 1 ", siteid).Values(&vs)

			if err != nil {
				logger.Error("不能获取jobnumber最大id值：%v", err.Error())
				return
			}

			uk := JOBNUMBER_PREFIX + strconv.FormatInt(siteid, 10)
			var maxid string
			if cnt > 0 {
				for _, v := range vs {
					if v["jobnumber"] != nil {
						maxid = utils.ConvertToString(v["jobnumber"])
						break
					}
				}
				redis.SetValue(uk, maxid)
			} else {
				redis.SetValue(uk, "0")
			}
		}
	}

	return
}
