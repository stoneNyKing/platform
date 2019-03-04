package models

import (
	"crypto/aes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"platform/authsvc/imconf"
	"platform/common/redis"
	"platform/common/utils"
	"strconv"
	"time"
)

func init() {
	orm.RegisterModel(new (SecSiteInfo))
}

func GetAppLicense(orgid,userid int64) (interface{}, int64, error) {

	licenseRow := make(map[string]interface{})

	strOrgid := strconv.FormatInt(orgid, 10)
	key := "org" + strOrgid+"_"+strconv.FormatInt(userid, 10)

	// 先从redis中读取数据，如果读取不到再去数据库中进行读取
	if redis.Exists(key) == true {
		licenseRow["orgid"] = redis.Hget(key, "orgid")
		licenseRow["apikey"] = redis.Hget(key, "apikey")
		licenseRow["orgcode"] = redis.Hget(key, "orgcode")
		licenseRow["createtime"] = redis.Hget(key, "createtime")
		licenseRow["remark"]  = redis.Hget(key, "remark")
		licenseRow["license"] = redis.Hget(key, "license")
		licenseRow["userid"] = redis.Hget(key, "userid")

	} else {
		var vs []orm.Params
		o := orm.NewOrm()

		statement := fmt.Sprintf(SQL_APP_LICENSE,imconf.Config.AuthdbSchema)
		_, err := o.Raw(statement + "where a.org_id=? and a.userid=?", orgid,userid).Values(&vs)

		if err != nil {
			return nil, 0, err
		}

		if len(vs) > 0 {
			licenseRow = vs[0]
		} else {
			return 	nil, 0, errors.New("can not find")
		}

		// 读取完毕后再向redis中进行更新
		updateRedis(orgid,userid)
	}

	apikey := licenseRow["apikey"].(string)
	orgcode := licenseRow["orgcode"].(string)

	// 获取license
	license := licenseRow["license"].(string)

	// 将license进行解码获得过期时间
	// 先进行base64解码
	strEncryp, err := base64.StdEncoding.DecodeString(license)
	if err != nil {
		logger.Error("Base64解码失败： %v", err.Error())
		return nil, 0, err
	}

	// 再进行AES解密获得r.body
	strKey := []byte(fmt.Sprintf("%016X", apikey + orgcode))
	clipher, _ := aes.NewCipher(strKey[:16])

	if len(strEncryp) <= 0{
		return nil, 0, errors.New("can not find")
	}

	dst := strEncryp
	clipher.Decrypt(dst, strEncryp)

	//获取解密后的jason串
	param := make(map[string]interface{})
	err = json.Unmarshal([]byte(dst), &param)

	if err != nil {
		logger.Error("将body转为json失败： %v", err.Error())
		return nil, 0, err
	}

	licenseRow["license"] = param

	return licenseRow, 1, nil
}


func GetAppLicenseLists(siteid int64, appid int64, stypes, contents []string, order, sort string, num, start int64) (interface{}, int64, error) {
	return getAppLicenseListCount(1,siteid,appid,stypes,contents,order,sort,num,start)
}
func GetAppLicenseCount(siteid int64, appid int64, stypes, contents []string, order, sort string) (int64, error) {
	_,cnt,err := getAppLicenseListCount(2,siteid,appid,stypes,contents,order,sort,1,0)

	return cnt,err
}

func getAppLicenseListCount(cate int,siteid int64, appid int64, stypes, contents []string, order, sort string, num, start int64) (interface{}, int64, error){
	var vs []orm.Params
	var cnt int64
	var err error

	o := orm.NewOrm()
	err = SetSearchPath(o,imconf.Config.AuthdbSchema)
	if err !=nil {
		logger.Error("不能设置search_path :%v",err)
		return nil,0,err
	}

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
			v = " a.api_key like " + v
		case "2":
			v = " a.org_code like " + v
		case "3":
			v = " a.status = " +  contents[i]
		case "4":
			v = " a.create_time >= '" + contents[i] + "'"
		case "5":
			v = " a.create_time < '" + contents[i] + "'"
		case "6":
			v = " a.remark like " + v
		case "7":
			v = " a.userid = " +  contents[i]
		case "8":
			v = " a.api_key ='" + contents[i] + "'"
		case "9":
			v = " a.org_code ='" + contents[i] + "'"
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
		statement =fmt.Sprintf(SQL_APP_LICENSE,imconf.Config.AuthdbSchema)
	} else if cate == 2 {
		statement =fmt.Sprintf(SQL_COUNT_APP_LICENSE,imconf.Config.AuthdbSchema)
	}

	statement = fmt.Sprintf(statement +
		"where a.org_id=? %s order by a.org_id %s limit ? offset ?", conditions, sort)

	cnt, err = o.Raw(statement, siteid, num,start).Values(&vs)

	if err != nil {
		logger.Error("不能获取资源信息列表：%v", err.Error())
		return nil, 0, err
	}

	// 将获取的数据进行解码
	if cate ==1 {
		for index,_ := range vs{
			apikey := (vs[index])["apikey"].(string)
			orgcode := (vs[index])["orgcode"].(string)

			// 获取license
			license := (vs[index])["license"].(string)

			// 先进行base64解码
			strEncryp, err := base64.StdEncoding.DecodeString(license)
			if err != nil {
				logger.Error("Base64解码失败： %v", err.Error())
				return nil,0, err
			}

			// 再进行AES解密获得r.body
			strKey := []byte(fmt.Sprintf("%016X", apikey + orgcode))
			clipher, _ := aes.NewCipher(strKey[:16])

			if len(strEncryp) <= 0{
				return nil, 0, errors.New("can not find")
			}

			dst := strEncryp
			clipher.Decrypt(dst, strEncryp)

			//获取解密后的jason串
			param := make(map[string]interface{})
			err = json.Unmarshal([]byte(dst), &param)

			if err != nil {
				logger.Error("将body转为json失败： %v", err.Error())
				return nil, 0,err
			}

			(vs[index])["license"] = param
		}
	}else if cate == 2 {
		if cnt > 0 {
			cnt = utils.Convert2Int64(vs[0]["ucount"])
		}
	}

	return vs, cnt, nil
}


func GetApiLicenseCounts(siteid int64, appid int64, stypes, contents []string, order, sort string, num, start int64) (interface{}, int64, error){
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
			v = " a.service_id = " + contents[i]
		case "2":
			v = " a.package_id = " + contents[i]
		case "3":
			v = " a.status = " +  contents[i]
		case "4":
			v = " b.svc_code = '" +  contents[i] + "'"
		case "5":
			v = " b.route = '" +  contents[i] + "'"
		case "6":
			v = " b.svc_id = " +  contents[i]
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

	statement =fmt.Sprintf(SQL_API_LICENSE_COUNTS,imconf.Config.AuthdbSchema,imconf.Config.AuthdbSchema)

	statement = fmt.Sprintf(statement +
		"where 1=1 %s order by a.pkg_service_id %s limit ? offset ?", conditions, sort)

	cnt, err = o.Raw(statement, num,start).Values(&vs)


	if err != nil {
		logger.Error("不能获取license：%v", err.Error())
		return nil, 0, err
	}

	return vs, cnt, nil
}


func PostAppLicense(siteid int64,param map[string]interface{}) (int64, error){

	if param == nil {
		return 0, errors.New("no input")
	}

	var v SecSiteInfo

	if param["id"] != nil {
		v.LicenseId = utils.Convert2Int64(param["id"])
	}

	if param["siteid"] != nil {
		v.Siteid = utils.Convert2Int64(param["siteid"])
	}else{
		v.Siteid = siteid
	}
	if param["organizationid"] != nil {
		v.OrganizationId = utils.Convert2Int64(param["organizationid"])
	}
	if param["apikey"] != nil {
		v.ApiKey = param["apikey"].(string)
	}

	if param["orgcode"] != nil {
		v.OrgCode = param["orgcode"].(string)
	}

	if param["status"] != nil {
		v.Status = utils.Convert2Int16(param["status"])
	}

	if param["createtime"] != nil {
		v.CreateTime,_ = time.Parse("2006-01-02 15:04:05",param["createtime"].(string))
	}

	if param["userid"] != nil {
		v.Userid = utils.Convert2Int64(param["userid"])
	}

	if param["remark"] != nil {
		v.Remark = param["remark"].(string)
	}

	if param["license"] != nil {
		license := param["license"	]

		licenseJson, err := json.Marshal(license)
		if err != nil{
			return 0, logger.Error("转换失败： %v", err.Error())
		}

		// 通过appid与oracode生成key进行AES加密
		strKey := []byte(fmt.Sprintf("%016X", v.ApiKey + v.OrgCode))
		clipher, _ := aes.NewCipher(strKey[:16])

		dst := licenseJson
		clipher.Encrypt(dst, licenseJson)

		// 对加密后的字符进行BASE64编码
		v.License = base64.StdEncoding.EncodeToString(dst)
	}

	o := orm.NewOrm()

	err := SetSearchPath(o,imconf.Config.AuthdbSchema)
	if err !=nil {
		logger.Error("不能设置search_path :%v",err)
		return 0,err
	}


	id, err := o.Insert(&v)
	if err != nil {
		//o.Rollback()
		logger.Error("不能插入资源信息：%v", err.Error())
		return 0, err
	}

	// 数据库修改后，将mysql的表数据同步到redis中去
	updateRedis(v.Siteid,v.Userid)

	return id, err
}


func DeleteAppLicense(param map[string]interface{}) (num int64, err error){

	if param == nil {
		return 0, errors.New("no input")
	}


	var v SecSiteInfo

	var ids []interface{}

	if param["id"] != nil {
		ids = param["id"].([]interface{})
	}

	o := orm.NewOrm()
	err = SetSearchPath(o,imconf.Config.AuthdbSchema)
	if err !=nil {
		logger.Error("不能设置search_path :%v",err)
		return 0,err
	}

	var cnt int64 = 0
	for _, rid := range ids {
		v.LicenseId = utils.Convert2Int64(rid)
		err = o.Read(&v)
		// 数据库修改后，将mysql的表数据同步到redis中去
		strOrgid := strconv.FormatInt(v.Siteid, 10)
		key := "org" + strOrgid+"_"+strconv.FormatInt(v.Userid, 10)

		redis.Del(key)

		num, err = o.Delete(&v)
		cnt += num
	}


	return cnt, err
}

func PutAppLicense(id int64, param map[string]interface{}) (int64, error){

	if param == nil {
		return 0, errors.New("no input")
	}

	var v SecSiteInfo

	if param["id"] != nil {
		v.LicenseId = utils.Convert2Int64(param["id"])
	}

	if id != v.LicenseId {
		return 0, errors.New("id is not match.")
	}

	o := orm.NewOrm()
	err := SetSearchPath(o,imconf.Config.AuthdbSchema)
	if err !=nil {
		logger.Error("不能设置search_path :%v",err)
		return 0,err
	}


	err = o.Read(&v)
	if err != nil {
		logger.Error("cannot read record.")
		return 0, err
	}

	if param["siteid"] != nil {
		v.Siteid = utils.Convert2Int64(param["siteid"])
	}
	if param["organizationid"] != nil {
		v.OrganizationId = utils.Convert2Int64(param["organizationid"])
	}

	if param["apikey"] != nil {
		v.ApiKey = param["apikey"].(string)
	}

	if param["orgcode"] != nil {
		v.OrgCode = param["orgcode"].(string)
	}
	if param["userid"] != nil {
		v.Userid = utils.Convert2Int64(param["userid"])
	}

	if param["status"] != nil {
		v.Status = utils.Convert2Int16(param["status"])
	}

	if param["createtime"] != nil {
		v.CreateTime,_ = time.Parse("2006-01-02 15:04:05",param["createtime"].(string))
	}

	if param["remark"] != nil {
		v.Remark = param["remark"].(string)
	}

	if param["license"] != nil {
		license := param["license"]

		licenseJson, err := json.Marshal(license)
		if err != nil{
			return 0, logger.Error("转换失败： %v", err.Error())
		}

		// 通过appid与oracode生成key进行AES加密
		strKey := []byte(fmt.Sprintf("%016X", v.ApiKey + v.OrgCode))
		clipher, _ := aes.NewCipher(strKey[:16])

		dst := licenseJson
		clipher.Encrypt(dst, licenseJson)

		// 对加密后的字符进行BASE64编码
		v.License = base64.StdEncoding.EncodeToString(dst)
	}

	//o := orm.NewOrm()

	cnt ,err := o.Update(&v)
	if err != nil {
		//o.Rollback()
		logger.Error("不能更新资源信息：%v", err.Error())
		return 0, err
	}

	// 数据库修改后，将mysql的表数据同步到redis中去
	updateRedis(v.Siteid,v.Userid)

	return cnt, err
}


func updateRedis (siteid,userid int64){
	// 获取sec_site_info的数据库所有记录信息
	var vs []orm.Params
	o := orm.NewOrm()

	statement := SQL_APP_LICENSE
	_, err := o.Raw(statement+" where a.org_id=? and a.userid=?", siteid,userid).Values(&vs)

	if err != nil{
		logger.Error("获取License信息失败：%v", err.Error())
		return
	}

	strOrgid := strconv.FormatInt(siteid, 10)
	key := "org" + strOrgid+"_"+strconv.FormatInt(userid, 10)

	// 将数据写到redis中去， 存放两个
	 for _, userlicense := range vs{

	 	var id,organizationid,siteid, apikey, orgcode, status, createtime, remark, license,userid string
	 	id = utils.ConvertToString(userlicense["id"])
	 	siteid = utils.ConvertToString(userlicense["siteid"])
	 	organizationid = utils.ConvertToString(userlicense["organizationid"])
	 	apikey = utils.ConvertToString(userlicense["apikey"])
	 	orgcode = utils.ConvertToString(userlicense["orgcode"])
	 	status = utils.ConvertToString(userlicense["status"])
		createtime = utils.ConvertToString(userlicense["createtime"])
		remark = utils.ConvertToString(userlicense["remark"])
		license = utils.ConvertToString(userlicense["license"])
		userid = utils.ConvertToString(userlicense["userid"])

	 	redis.Hset(key, "id", id)
	 	redis.Hset(key, "siteid", siteid)
	 	redis.Hset(key, "organizationid", organizationid)
		redis.Hset(key, "apikey", apikey)
		redis.Hset(key, "orgcode", orgcode)
		redis.Hset(key, "status", status)
		redis.Hset(key, "createtime", createtime)
		redis.Hset(key, "remark", remark)
		redis.Hset(key, "license", license)
		redis.Hset(key, "userid", userid)
	 }
	// 释放redis
	return
}