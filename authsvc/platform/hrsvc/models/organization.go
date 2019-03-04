package models

import (
	"context"
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"platform/common/utils"
	"platform/hrsvc/imconf"
	"strings"
	"time"
)

func init() {

	orm.RegisterModel(new(NcOrganization))
}

func GetOrganizationTree(ctx context.Context,siteid int64, appid int64, stypes, contents []string, order, sort string, num, start int64) (interface{}, int64, error) {

	var maps []map[string]interface{}
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
			v = " a.status= " + contents[i]
		case "2":
			v = " a.create_time >  '" + contents[i] + "'"
		case "3":
			v = " a.create_time <=  '" + contents[i] + "'"
		case "4":
			v = " a.update_time >  '" + contents[i] + "'"
		case "5":
			v = " a.update_time <=  '" + contents[i] + "'"
		case "6":
			v = " a.creatorid =  " + contents[i]
		case "7":
			v = " a.updaterid =  " + contents[i]
		case "8":
			v = " a.reg_div_id =  " + contents[i]
		case "9":
			v = " a.org_code like  " + v
		case "10":
			v = " a.reg_code like  " + v
		case "11":
			v = " a.parent_id =  " + contents[i]
		case "12":
			v = " a.org_type =  '" + contents[i] + "'"
		case "13":
			v = " a.org_ln like  " + v
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
	}else{
		conditions = " and a.parent_id=0 "
	}

	logger.Finest("conditions = %s", conditions)
	var statement string
	statement = fmt.Sprintf(SQL_NcOrganization)

	if siteid > 1 {
		statement = fmt.Sprintf(statement+
			"where (1=%d) %s order by a.organization_id %s limit ? offset ?", 1, conditions, sort)
	} else {
		statement = fmt.Sprintf(statement+
			"where 1=1 %s order by a.organization_id %s limit ? offset ?", conditions, sort)

	}

	cnt, err = o.Raw(statement, num, start).Values(&vs)

	if err != nil {
		logger.Error("不能获取列表信息：%v", err.Error())
		return nil, 0, err
	}

	if cnt > 0 {
		for _, v := range vs {
			rid := utils.ConvertToString(v["id"])
			child, _, _ := GetOrganizationTreeChild(o, siteid, rid)
			v["child"] = child
			maps = append(maps, v)
		}
	}

	return maps, 0, nil
}

func GetOrganizationTreeChild(o orm.Ormer, siteid int64, rid string) (interface{}, int64, error) {
	var maps []map[string]interface{}
	var vs []orm.Params

	if o == nil {
		o = orm.NewOrm()
	}

	var cnt int64
	var err error

	var statement string
	statement = fmt.Sprintf(SQL_NcOrganization)

	if siteid > 1 {
		statement = fmt.Sprintf(statement+
			"where (1=%d and a.parent_id=?) order by a.organization_id asc ", 1)
	} else {
		statement = fmt.Sprintf(statement +
			"where a.parent_id=? order by a.organization_id asc ")

	}

	cnt, err = o.Raw(statement, rid).Values(&vs)

	if err != nil {
		logger.Error("不能获取信息列表：%v", err.Error())
		return nil, 0, err
	}

	var ids string = ""
	if cnt > 0 {
		for _, v := range vs {
			m := make(map[string]interface{})

			for key, value := range v {
				m[key] = value
			}
			if v["id"] != nil {
				ids = v["id"].(string)
				if ids != "" {
					pas, count, err := GetOrganizationTreeChild(o, siteid, ids)

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



func GetOrganizationLists(ctx context.Context,siteid int64, appid int64, stypes, contents []string, order, sort string, num, start int64) (interface{}, int64, error) {
	if num <= 0 {
		num=PAGENUM_MAX
	}

	return getOrganizationListCount(ctx,1,siteid,appid,stypes,contents,order,sort,num,start)
}
func GetOrganizationCount(ctx context.Context,siteid int64, appid int64, stypes, contents []string, order, sort string) (int64, error) {
	_,cnt,err := getOrganizationListCount(ctx,2,siteid,appid,stypes,contents,order,sort,1,0)

	return cnt,err
}


/*
	搜索条件：
	stype = 1： 状态
	stype = 2： 开始时间大于某个时间
	stype = 3： 开始时间小于某个时间
	stype = 4： 更新时间大于某个时间
	stype = 5： 更新时间小于某个时间
	stype = 6： 创建者标识
	stype = 7： 更新者标识
	stype = 8： 行政区标识
	stype = 9： 机构代码
	stype = 10： 行政区代码
	stype = 11： 上一级标识
	stype = 12： 机构类型
	stype = 13： 机构执业许可证
	stype = 14： 机构名称
*/
func getOrganizationListCount(ctx context.Context,cate int,siteid int64, appid int64, stypes, contents []string, order, sort string, num, start int64) (interface{}, int64, error) {

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
			v = " a.status= " + contents[i]
		case "2":
			v = " a.create_time >  '" + contents[i] + "'"
		case "3":
			v = " a.create_time <=  '" + contents[i] + "'"
		case "4":
			v = " a.update_time >  '" + contents[i] + "'"
		case "5":
			v = " a.update_time <=  '" + contents[i] + "'"
		case "6":
			v = " a.creatorid =  " + contents[i]
		case "7":
			v = " a.updaterid =  " + contents[i]
		case "8":
			v = " a.reg_div_id =  " + contents[i]
		case "9":
			v = " a.org_code like  " + v
		case "10":
			v = " a.reg_code like  " + v
		case "11":
			v = " a.parent_id =  " + contents[i]
		case "12":
			v = " a.org_type =  '" + contents[i] + "'"
		case "13":
			v = " a.org_ln like  " + v
		case "14":
			v = " a.reg_name like  " + v
		case "15":
			v = " a.org_name like  " + v
		case "16":
			v = " a.manage_cate =  " + contents[i]
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
	var statement string

	if cate == 1 {
		statement = fmt.Sprintf(SQL_NcOrganization)
	}else if cate ==2 {
		statement = fmt.Sprintf(SQL_COUNT_NcOrganization)
	}

	if siteid > 1 {
		statement = fmt.Sprintf(statement+
			"where (a.siteid=%d) %s order by a.organization_id %s limit ? offset ?", siteid, conditions, sort)
	} else {
		statement = fmt.Sprintf(statement+
			"where 1=1 %s order by a.organization_id %s limit ? offset ?", conditions, sort)

	}

	cnt, err = o.Raw(statement, num, start).Values(&vs)

	if err != nil {
		logger.Error("不能获取列表信息：%v", err.Error())
		return nil, 0, err
	}

	if cate == 2 {
		if cnt > 0 {
			cnt = utils.Convert2Int64(vs[0]["ucount"])
		}
	}

	return vs, cnt, nil
}

func GetOrganization(ctx context.Context,id int64, siteid int64, appid int64) (interface{}, int64, error) {

	var vs []orm.Params

	o := orm.NewOrm()

	var cnt int64
	var err error

	logger.Finest("siteid=%d,appid=%d,areaid=%d", siteid, appid, id)

	statement := fmt.Sprintf(SQL_NcOrganization)

	if siteid > 1 {
		cnt, err = o.Raw(statement+
			"where (a.siteid=%d) and a.organization_id =?", siteid, id).Values(&vs)
	} else {
		cnt, err = o.Raw(statement+
			"where a.organization_id =?", id).Values(&vs)
	}

	if err != nil {
		logger.Error("不能获取信息：%v", err.Error())
		return nil, 0, err
	}

	return vs, cnt, nil
}


func PostOrganization(ctx context.Context,appid, siteid int64, token string, param map[string]interface{}) (id int64, err error) {

	if param == nil {
		return 0, errors.New("no input")
	}

	var v NcOrganization

	if param["status"] != nil {
		v.Status = utils.Convert2Int16(param["status"])
	}
	if param["nature"] != nil {
		v.Nature = utils.Convert2Int8(param["nature"])
	}

	if param["updaterid"] != nil {
		v.Updaterid = utils.Convert2Int64(param["updaterid"])
	}

	if param["updatetime"] != nil {
		v.UpdateTime,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["updatetime"]))
	}else{
		v.UpdateTime,_ = time.Parse("2006-01-02 15:04:05",utils.GetTimeFormat("2006-01-02 15:04:05"))
	}

	if param["approvaldepartment"] != nil {
		v.ApprovalDepartment = utils.ConvertToString(param["approvaldepartment"])
	}

	if param["address"] != nil {
		v.Address = utils.ConvertToString(param["address"])
	}

	if param["membershiplevel"] != nil {
		v.MembershipLevel = utils.ConvertToString(param["membershiplevel"])
	}

	if param["remark"] != nil {
		v.Remark = utils.ConvertToString(param["remark"])
	}

	if param["regcode"] != nil {
		v.RegCode = utils.ConvertToString(param["regcode"])
	}
	if param["regname"] != nil {
		v.RegName = utils.ConvertToString(param["regname"])
	}

	if param["orgname"] != nil {
		v.OrgName = utils.ConvertToString(param["orgname"])
	}
	if param["managecate"] != nil {
		v.ManageCate = utils.Convert2Int8(param["managecate"])
	}

	if param["parentid"] != nil {
		v.ParentId = utils.Convert2Int64(param["parentid"])
	}

	if param["approvalbednum"] != nil {
		v.ApprovalBedNum = utils.Convert2Int(param["approvalbednum"])
	}

	if param["legalperson"] != nil {
		v.LegalPerson = utils.ConvertToString(param["legalperson"])
	}

	if param["approvaltime"] != nil {
		v.ApprovalTime,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["approvaltime"]))
	}else{
		v.ApprovalTime = time.Now()
	}

	if param["branchesnum"] != nil {
		v.BranchesNum = utils.Convert2Int(param["branchesnum"])
	}

	if param["managementarea"] != nil {
		v.ManagementArea = utils.ConvertToString(param["managementarea"])
	}

	if param["orglevel"] != nil {
		v.OrgLevel = utils.ConvertToString(param["orglevel"])
	}

	if param["phone"] != nil {
		v.Phone = utils.ConvertToString(param["phone"])
	}

	if param["orgtype"] != nil {
		v.OrgType = utils.ConvertToString(param["orgtype"])
	}

	if param["id"] != nil {
		v.OrganizationId = utils.Convert2Int64(param["id"])
	}

	if param["creatorid"] != nil {
		v.Creatorid = utils.Convert2Int64(param["creatorid"])
	}

	if param["createtime"] != nil {
		v.CreateTime,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["createtime"]))
	}else{
		v.CreateTime,_ = time.Parse("2006-01-02 15:04:05",utils.GetTimeFormat("2006-01-02 15:04:05"))
	}

	if param["regdivid"] != nil {
		v.RegDivId = utils.Convert2Int64(param["regdivid"])
	}

	if param["orgcode"] != nil {
		v.OrgCode = utils.ConvertToString(param["orgcode"])
	}

	if param["acreage"] != nil {
		v.Acreage = utils.Convert2Float32(param["acreage"])
	}

	if param["orgln"] != nil {
		v.OrgLn = utils.ConvertToString(param["orgln"])
	}

	if param["effdate"] != nil {
		v.EffectiveDate,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["effdate"]))
	}
	if param["expdate"] != nil {
		v.ExpireDate,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["expdate"]))
	}

	o := orm.NewOrm()

	err = SetSearchPath(o, imconf.Config.DefaultSchema)
	if err != nil {
		logger.Error("不能设置search_path :%v", err)
		return 0, err
	}

	id, err = o.Insert(&v)

	if err != nil {
		logger.Error("不能插入信息：%v", err.Error())
		return 0, err
	}

	return id, err
}

func PostBulkOrganization(ctx context.Context,appid, siteid int64, token string, data map[string]interface{}) (id int64, err error) {

	if data == nil {
		return 0, errors.New("no input")
	}

	var vs []interface{}
	if data["items"] != nil {
		vs = data["items"].([]interface{})
	}

	var iis []NcOrganization
	var c int = 0
	for _, vv := range vs {
		param := vv.(map[string]interface{})
		var v NcOrganization

		if param["orglevel"] != nil {
			v.OrgLevel = utils.ConvertToString(param["orglevel"])
		}

		if param["phone"] != nil {
			v.Phone = utils.ConvertToString(param["phone"])
		}

		if param["orgtype"] != nil {
			v.OrgType = utils.ConvertToString(param["orgtype"])
		}

		if param["approvaltime"] != nil {
			v.ApprovalTime,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["approvaltime"]))
		}

		if param["branchesnum"] != nil {
			v.BranchesNum = utils.Convert2Int(param["branchesnum"])
		}

		if param["managementarea"] != nil {
			v.ManagementArea = utils.ConvertToString(param["managementarea"])
		}

		if param["regdivid"] != nil {
			v.RegDivId = utils.Convert2Int64(param["regdivid"])
		}

		if param["orgcode"] != nil {
			v.OrgCode = utils.ConvertToString(param["orgcode"])
		}
		if param["orgname"] != nil {
			v.OrgName = utils.ConvertToString(param["orgname"])
		}
		if param["managecate"] != nil {
			v.ManageCate = utils.Convert2Int8(param["managecate"])
		}

		if param["acreage"] != nil {
			v.Acreage = utils.Convert2Float32(param["acreage"])
		}

		if param["orgln"] != nil {
			v.OrgLn = utils.ConvertToString(param["orgln"])
		}

		if param["id"] != nil {
			v.OrganizationId = utils.Convert2Int64(param["id"])
		}

		if param["creatorid"] != nil {
			v.Creatorid = utils.Convert2Int64(param["creatorid"])
		}

		if param["createtime"] != nil {
			v.CreateTime,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["createtime"]))
		}else{
			v.CreateTime,_ =time.Parse("2006-01-02 15:04:05",utils.GetTimeFormat("2006-01-02 15:04:05"))
		}

		if param["approvaldepartment"] != nil {
			v.ApprovalDepartment = utils.ConvertToString(param["approvaldepartment"])
		}

		if param["address"] != nil {
			v.Address = utils.ConvertToString(param["address"])
		}

		if param["membershiplevel"] != nil {
			v.MembershipLevel = utils.ConvertToString(param["membershiplevel"])
		}

		if param["status"] != nil {
			v.Status = utils.Convert2Int16(param["status"])
		}

		if param["updaterid"] != nil {
			v.Updaterid = utils.Convert2Int64(param["updaterid"])
		}

		if param["updatetime"] != nil {
			v.UpdateTime,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["updatetime"]))
		}else{
			v.UpdateTime,_ = time.Parse("2006-01-02 15:04:05",utils.GetTimeFormat("2006-01-02 15:04:05"))
		}

		if param["approvalbednum"] != nil {
			v.ApprovalBedNum = utils.Convert2Int(param["approvalbednum"])
		}

		if param["legalperson"] != nil {
			v.LegalPerson = utils.ConvertToString(param["legalperson"])
		}

		if param["remark"] != nil {
			v.Remark = utils.ConvertToString(param["remark"])
		}

		if param["regcode"] != nil {
			v.RegCode = utils.ConvertToString(param["regcode"])
		}
		if param["regname"] != nil {
			v.RegName = utils.ConvertToString(param["regname"])
		}
		if param["nature"] != nil {
			v.Nature = utils.Convert2Int8(param["nature"])
		}

		if param["parentid"] != nil {
			v.ParentId = utils.Convert2Int64(param["parentid"])
		}
		if param["effdate"] != nil {
			v.EffectiveDate,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["effdate"]))
		}
		if param["expdate"] != nil {
			v.ExpireDate,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["expdate"]))
		}

		iis = append(iis, v)
		c++
	}

	o := orm.NewOrm()
	err = SetSearchPath(o, imconf.Config.DefaultSchema)
	if err != nil {
		logger.Error("不能设置search_path :%v", err)
		return 0, err
	}

	cnt, err := o.InsertMulti(c, &iis)

	if err != nil {
		logger.Error("不能插入信息：%v", err.Error())
		return 0, err
	}

	return cnt, nil
}

func PutOrganization(ctx context.Context,id int64, appid, siteid int64, token string, param map[string]interface{}) (cnt int64, err error) {

	if param == nil {
		return 0, errors.New("no input")
	}

	var v NcOrganization

	if param["id"] != nil {
		v.OrganizationId = utils.Convert2Int64(param["id"])
	}

	if id != v.OrganizationId {
		return 0, errors.New("id is not match.")
	}

	o := orm.NewOrm()
	err = SetSearchPath(o, imconf.Config.DefaultSchema)
	if err != nil {
		logger.Error("不能设置search_path :%v", err)
		return 0, err
	}

	err = o.Read(&v)
	if err != nil {
		logger.Error("cannot read record.")
		return 0, err
	}

	if param["address"] != nil {
		v.Address = utils.ConvertToString(param["address"])
	}

	if param["membershiplevel"] != nil {
		v.MembershipLevel = utils.ConvertToString(param["membershiplevel"])
	}
	if param["orgname"] != nil {
		v.OrgName = utils.ConvertToString(param["orgname"])
	}
	if param["managecate"] != nil {
		v.ManageCate = utils.Convert2Int8(param["managecate"])
	}

	if param["status"] != nil {
		v.Status = utils.Convert2Int16(param["status"])
	}

	if param["updaterid"] != nil {
		v.Updaterid = utils.Convert2Int64(param["updaterid"])
	}

	if param["updatetime"] != nil {
		v.UpdateTime,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["updatetime"]))
	}else{
		v.UpdateTime,_ = time.Parse("2006-01-02 15:04:05",utils.GetTimeFormat("2006-01-02 15:04:05"))
	}

	if param["approvaldepartment"] != nil {
		v.ApprovalDepartment = utils.ConvertToString(param["approvaldepartment"])
	}

	if param["legalperson"] != nil {
		v.LegalPerson = utils.ConvertToString(param["legalperson"])
	}

	if param["remark"] != nil {
		v.Remark = utils.ConvertToString(param["remark"])
	}

	if param["regcode"] != nil {
		v.RegCode = utils.ConvertToString(param["regcode"])
	}
	if param["regname"] != nil {
		v.RegName = utils.ConvertToString(param["regname"])
	}
	if param["nature"] != nil {
		v.Nature = utils.Convert2Int8(param["nature"])
	}

	if param["parentid"] != nil {
		v.ParentId = utils.Convert2Int64(param["parentid"])
	}

	if param["approvalbednum"] != nil {
		v.ApprovalBedNum = utils.Convert2Int(param["approvalbednum"])
	}

	if param["phone"] != nil {
		v.Phone = utils.ConvertToString(param["phone"])
	}

	if param["orgtype"] != nil {
		v.OrgType = utils.ConvertToString(param["orgtype"])
	}

	if param["approvaltime"] != nil {
		v.ApprovalTime,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["approvaltime"]))
	}

	if param["branchesnum"] != nil {
		v.BranchesNum = utils.Convert2Int(param["branchesnum"])
	}

	if param["managementarea"] != nil {
		v.ManagementArea = utils.ConvertToString(param["managementarea"])
	}

	if param["orglevel"] != nil {
		v.OrgLevel = utils.ConvertToString(param["orglevel"])
	}

	if param["orgcode"] != nil {
		v.OrgCode = utils.ConvertToString(param["orgcode"])
	}

	if param["acreage"] != nil {
		v.Acreage = utils.Convert2Float32(param["acreage"])
	}

	if param["orgln"] != nil {
		v.OrgLn = utils.ConvertToString(param["orgln"])
	}

	if param["id"] != nil {
		v.OrganizationId = utils.Convert2Int64(param["id"])
	}

	if param["creatorid"] != nil {
		v.Creatorid = utils.Convert2Int64(param["creatorid"])
	}

	if param["createtime"] != nil {
		v.CreateTime,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["createtime"]))
	}else{
		v.CreateTime,_ = time.Parse("2006-01-02 15:04:05",v.CreateTime.Format("2006-01-02 15:04:05"))
	}

	if param["regdivid"] != nil {
		v.RegDivId = utils.Convert2Int64(param["regdivid"])
	}
	if param["nature"] != nil {
		v.Nature = utils.Convert2Int8(param["nature"])
	}
	if param["effdate"] != nil {
		v.EffectiveDate,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["effdate"]))
	}
	if param["expdate"] != nil {
		v.ExpireDate,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["expdate"]))
	}

	cnt, err = o.Update(&v)
	if err != nil {
		logger.Error("不能更新记录：%v", err.Error())
		return 0, err
	}

	return cnt, nil
}

func PutBulkOrganization(ctx context.Context,id int64, appid, siteid int64, token string, data map[string]interface{}) (cnt int64, err error) {

	if data == nil {
		return 0, errors.New("no input")
	}

	o := orm.NewOrm()
	err = SetSearchPath(o, imconf.Config.DefaultSchema)
	if err != nil {
		logger.Error("不能设置search_path :%v", err)
		return 0, err
	}

	var vs []interface{}
	if data["items"] != nil {
		vs = data["items"].([]interface{})
	}

	var counts int64 = 0
	for _, vv := range vs {
		param := vv.(map[string]interface{})
		var v NcOrganization

		if param["id"] != nil {
			v.OrganizationId = utils.Convert2Int64(param["id"])
		}

		if id != v.OrganizationId {
			return 0, errors.New("id is not match.")
		}

		err = o.Read(&v)
		if err != nil {
			logger.Error("cannot read record.")
			return 0, err
		}

		if param["id"] != nil {
			v.OrganizationId = utils.Convert2Int64(param["id"])
		}
		if param["orgname"] != nil {
			v.OrgName = utils.ConvertToString(param["orgname"])
		}
		if param["managecate"] != nil {
			v.ManageCate = utils.Convert2Int8(param["managecate"])
		}

		if param["creatorid"] != nil {
			v.Creatorid = utils.Convert2Int64(param["creatorid"])
		}

		if param["createtime"] != nil {
			v.CreateTime,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["createtime"]))
		}else{
			v.CreateTime,_ = time.Parse("2006-01-02 15:04:05",v.CreateTime.Format("2006-01-02 15:04:05"))
		}

		if param["regdivid"] != nil {
			v.RegDivId = utils.Convert2Int64(param["regdivid"])
		}

		if param["orgcode"] != nil {
			v.OrgCode = utils.ConvertToString(param["orgcode"])
		}

		if param["acreage"] != nil {
			v.Acreage = utils.Convert2Float32(param["acreage"])
		}

		if param["orgln"] != nil {
			v.OrgLn = utils.ConvertToString(param["orgln"])
		}

		if param["status"] != nil {
			v.Status = utils.Convert2Int16(param["status"])
		}
		if param["nature"] != nil {
			v.Nature = utils.Convert2Int8(param["nature"])
		}

		if param["updaterid"] != nil {
			v.Updaterid = utils.Convert2Int64(param["updaterid"])
		}

		if param["updatetime"] != nil {
			v.UpdateTime,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["updatetime"]))
		}else{
			v.UpdateTime,_ = time.Parse("2006-01-02 15:04:05",utils.GetTimeFormat("2006-01-02 15:04:05"))
		}

		if param["approvaldepartment"] != nil {
			v.ApprovalDepartment = utils.ConvertToString(param["approvaldepartment"])
		}

		if param["address"] != nil {
			v.Address = utils.ConvertToString(param["address"])
		}

		if param["membershiplevel"] != nil {
			v.MembershipLevel = utils.ConvertToString(param["membershiplevel"])
		}

		if param["remark"] != nil {
			v.Remark = utils.ConvertToString(param["remark"])
		}

		if param["regcode"] != nil {
			v.RegCode = utils.ConvertToString(param["regcode"])
		}
		if param["regname"] != nil {
			v.RegName = utils.ConvertToString(param["regname"])
		}

		if param["parentid"] != nil {
			v.ParentId = utils.Convert2Int64(param["parentid"])
		}

		if param["approvalbednum"] != nil {
			v.ApprovalBedNum = utils.Convert2Int(param["approvalbednum"])
		}

		if param["legalperson"] != nil {
			v.LegalPerson = utils.ConvertToString(param["legalperson"])
		}

		if param["approvaltime"] != nil {
			v.ApprovalTime,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["approvaltime"]))
		}

		if param["branchesnum"] != nil {
			v.BranchesNum = utils.Convert2Int(param["branchesnum"])
		}

		if param["managementarea"] != nil {
			v.ManagementArea = utils.ConvertToString(param["managementarea"])
		}

		if param["orglevel"] != nil {
			v.OrgLevel = utils.ConvertToString(param["orglevel"])
		}

		if param["phone"] != nil {
			v.Phone = utils.ConvertToString(param["phone"])
		}

		if param["orgtype"] != nil {
			v.OrgType = utils.ConvertToString(param["orgtype"])
		}
		if param["effdate"] != nil {
			v.EffectiveDate,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["effdate"]))
		}
		if param["expdate"] != nil {
			v.ExpireDate,_ = time.Parse("2006-01-02 15:04:05",utils.ConvertToString(param["expdate"]))
		}

		cnt, err = o.Update(&v)
		if err != nil {
			logger.Error("不能更新记录：%v", err.Error())
			return counts, err
		}
		counts += cnt
	}

	return counts, nil
}

func DeleteOrganization(ctx context.Context,id int64, appid, siteid int64, token string, param map[string]interface{}) (num int64, err error) {

	o := orm.NewOrm()
	err = SetSearchPath(o, imconf.Config.DefaultSchema)
	if err != nil {
		logger.Error("不能设置search_path :%v", err)
		return 0, err
	}

	var v NcOrganization

	var ids []interface{}

	if param["id"] != nil {
		ids = param["id"].([]interface{})
	}

	var cnt int64 = 0

	for _, rid := range ids {
		v.OrganizationId = utils.Convert2Int64(rid)
		num, err = o.Delete(&v)
		if err != nil {
			if strings.Contains(err.Error(),"1451") {
				err = errors.New("医院已关联其他数据,不能删除。")
			}
			return 0, err
		}

		cnt += num
	}

	return cnt, err
}
