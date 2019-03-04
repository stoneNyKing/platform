package models

import (
	"os"
	mod "platform/filesvc/dbmod"
	"platform/filesvc/imconf"
	"platform/mskit/rest"
	"strings"
)

func DownloadFile(r *rest.Request)(f string,err error) {

	svc := fileService{}
	or := r.OriginRequest

	logger.Finest("request=%+v,path=%s,url=%+v",or,or.URL.Path,or.URL)

	fileid := ""
	url := ""

	err = SetSearchPath(imconf.Config.FiledbSchema)
	if err !=nil {
		logger.Error("不能设置search_path :%v",err)
		return "",err
	}


	var fi  mod.FileInfo

	if or.URL.Path != "" {
		ss := strings.Split(or.URL.Path,"/")
		if len(ss)>0 {
			fileid =  ss[len(ss)-1]
			fi.Filekey = fileid

			_,err = orm.Get(&fi)
			if err != nil {
				logger.Error( "不能获取文件id对应的url：%v" ,err)
				return "",err
			}

			_,err =os.Stat(fi.LocalFile)
			if err == nil {
				logger.Finest("文件已经存在不需要重新获取。")
				return fi.Redirect,nil
			}

			url = fi.OrigUrl
		}
	}

	var root string
	lf := ""
	f,lf, err = svc.Download(root,url)

	if err == nil && lf != "" {
		fi.LocalFile = lf
		fi.Redirect = f
		_,err = orm.Update(&fi)
		if err != nil {
			logger.Error("不能更新本地文件名: %v",err)
		}
	}

	logger.Finest("返回结果：url=%v, lf=%s,err = %v ",f,lf,err)

	return f,nil
}