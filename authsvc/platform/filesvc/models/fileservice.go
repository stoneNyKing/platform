package models

import (
	"crypto/md5"
	"errors"
	"fmt"
	borm "github.com/astaxie/beego/orm"
	"io"
	"io/ioutil"
	"platform/common/utils"
	"platform/filesvc/comm"
	mod "platform/filesvc/dbmod"
	"platform/filesvc/imconf"
	"platform/mskit/rest"
)

func init(){
	borm.RegisterModel(new(mod.FileInfo))
}

func GetFileLists(siteid int64, appid int64, stypes, contents []string, order, sort string, num, start int64) (interface{}, int64, error) {
	if num <= 0 {
		num = PAGENUM_MAX
	}

	return getFileListCount(1, siteid, appid, stypes, contents, order, sort, num, start)
}
func GetFileCount(siteid int64, appid int64, stypes, contents []string, order, sort string) (int64, error) {
	_, cnt, err := getFileListCount(2, siteid, appid, stypes, contents, order, sort, 1, 0)

	return cnt, err
}

func getFileListCount(cate int, siteid int64, appid int64, stypes, contents []string, order, sort string, num, start int64) (interface{}, int64, error) {


	var vs [] borm.Params

	var cnt int64
	var err error

	if sort == "" {
		sort = "desc"
	}

	logger.Finest("siteid=%d,appid=%d,order=%s,sort=%s", siteid, appid, order, sort)

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
			v = " a.name like " + v
		case "5":
			v = " a.file_no like " + v
		case "6":
			v = " a.file_owner like " + v
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

	o := borm.NewOrm()

	var sqlstr,statement string
	if cate == 1 {
		sqlstr = fmt.Sprintf(SQL_FileInfo)
	}else if cate == 2 {
		sqlstr = fmt.Sprintf(SQL_COUNT_FileInfo)
	}
	if siteid > 1 {
		statement = fmt.Sprintf(sqlstr +
			"where a.siteid=%d %s order by a.fileid %s limit ? offset ? ", siteid,conditions, sort)
	}else{
		statement = fmt.Sprintf(sqlstr +
			"where 1=1 %s order by a.fileid %s limit ? offset ? ", conditions, sort)
	}

	cnt, err = o.Raw(statement,num,start).Values(&vs)

	if err != nil {
		logger.Error("不能获取列表数量：%v", err.Error())
		return nil, 0, err
	}

	if cate == 2 {
		if cnt > 0 {
			cnt = utils.Convert2Int64(vs[0]["ucount"])
		}
	}

	return vs,cnt, nil}

func GetFile(id int64, siteid int64, appid int64) (interface{}, int64, error) {

	var cnt int64
	var err error

	logger.Finest("siteid=%d,appid=%d,fileid=%d", siteid, appid, id)
	file := new(mod.FileStorage)
	_, err = orm.Id(id).Get(file)

	if err != nil {
		logger.Error("不能获取信息：%v", err.Error())
		return nil, 0, err
	}
	cnt  = 1
	return file.Content, cnt, nil
}

//上传文件
func UploadFile(r *rest.Request,category int)(vs interface{},err error) {

	request,err := decodeUploadRequest(nil,r.OriginRequest)
	if request == nil {
		return "",err
	}
	req := request.(uploadRequest)
	svc := fileService{}
	f, err := svc.Upload(&req)

	if err != nil {
		return nil,err
	}

	fid :=""
	logger.Finest("上传的文件名为：%s",f)

	o := borm.NewOrm()

	prefix := "/oss"
	h := md5.New()
	io.WriteString(h, f)
	cipherStr := h.Sum(nil)
	fid = fmt.Sprintf("%02x", cipherStr)
	logger.Fine("file id = %s", fid)

	var cf mod.FileInfo
	cf.Filekey = fid
	cf.StorageType = req.Storage
	cf.OrigUrl = f
	cf.Status = 1
	cf.Siteid = utils.Convert2Int64(req.Site)
	cf.Hash = req.Hash
	cf.CreateTime = utils.GetTimeFormat("2006-01-02 15:04:05")
	cf.Name = req.FileName

	err = BSetSearchPath(o, imconf.Config.FiledbSchema)
	if err != nil {
		logger.Error("不能设置search_path :%v", err)
		return "", err
	}

	id, err := o.Insert(&cf)

	if err != nil {
		logger.Error(" 不能插入文件名到数据库：%v", err)
	} else {
		if cf.StorageType == comm.STORAGE_TYPE_ALIOSS {
			f = prefix + "/" + fid
		}else if cf.StorageType == comm.STORAGE_TYPE_FILESYSTEM {
			f = cf.OrigUrl
		}else if cf.StorageType == comm.STORAGE_TYPE_DATABASE {
			f = imconf.Config.Prefix + "/file/" + utils.ConvertToString(id)

			if len(req.FileHeader)>0 {
				r,err := req.FileHeader[0].Open()
				if err != nil {
					logger.Error("打开文件失败: %v",err)
					return "",err
				}

				b,err:=ioutil.ReadAll(r)
				hb := md5.Sum(b)
				cf.Hash = fmt.Sprintf("%x",hb)
				var fs mod.FileStorage
				fs.Fileid = id
				fs.Status = 1
				fs.Content = []uint8(b)
				_,err = orm.Insert(&fs)
				if err != nil {
					logger.Error("插入存储数据失败: %v",err)
					return "",err
				}
			}
		}
		cf.Location = imconf.Config.DomainName + f
		_,err = orm.Id(id).Cols("hash","location").Update(&cf)
		if err != nil {
			logger.Error("不能更新hash值: %v",err)
		}
	}

	f = imconf.Config.DomainName + f

	version := utils.Convert2Float32(r.Version)
	ver := int(version * 10)
	if ver <20 {
		vs = f
	}else {
		m := make(map[string]interface{})
		m["id"] = id
		m["url"] = f
		vs = m
	}

	logger.Finest("filename=%+v",vs)

	return vs,err
}



func PutFileInfo(id int64,appid,siteid int64,token string, param map[string]interface{}) (cnt int64, err error) {

	if param == nil {
		return 0, errors.New("no input")
	}

	var v mod.FileInfo

	if param["id"] != nil {
		v.Fileid = utils.Convert2Int64(param["id"])
	}

	if id != v.Fileid {
		return 0, errors.New("id is not match.")
	}

	o := borm.NewOrm()
	err = BSetSearchPath(o, imconf.Config.FiledbSchema)
	if err !=nil {
		logger.Error("不能设置search_path :%v",err)
		return 0,err
	}

	err = o.Read(&v)
	if err != nil {
		logger.Error("cannot read record.")
		return 0, err
	}

	if param["fileno"] != nil {
		v.FileNo = utils.ConvertToString(param["fileno"])
	}

	if param["status"] != nil {
		v.Status = utils.Convert2Int(param["status"])
	}

	if param["filesize"] != nil {
		v.FileSize = utils.Convert2Int64(param["filesize"])
	}

	if param["siteid"] != nil {
		v.Siteid = utils.Convert2Int64(param["siteid"])
	}

	if param["fileowner"] != nil {
		v.FileOwner = utils.ConvertToString(param["fileowner"])
	}

	if param["name"] != nil {
		v.Name = utils.ConvertToString(param["name"])
	}

	if param["location"] != nil {
		v.Location = utils.ConvertToString(param["location"])
	}

	if param["redirect"] != nil {
		v.Redirect = utils.ConvertToString(param["redirect"])
	}

	if param["remark"] != nil {
		v.Remark = utils.ConvertToString(param["remark"])
	}


	cnt, err = o.Update(&v)
	if err != nil {
		logger.Error("不能更新记录：%v", err.Error())
		return 0, err
	}

	return cnt, nil
}

