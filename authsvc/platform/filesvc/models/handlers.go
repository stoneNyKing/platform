package models

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"

	"github.com/disintegration/imaging"

	"platform/common/utils"
	"platform/filesvc/comm"
	"platform/filesvc/ss"
	mod "platform/filesvc/dbmod"

	_ "github.com/spf13/viper/remote"
)

var fileHandlers = make(map[string]FileHandler)



func init() {
	RegisterFileHandler("image_file_handler", imageHandler{})
	RegisterFileHandler("regular_file_handler", regularHandler{})
}


type FileHandler interface {
	CouldHandle(media ss.FileLibrary) bool
	Handle(media ss.FileLibrary, file multipart.File, option *mod.Options) (string,string, string, error)
}

func RegisterFileHandler(name string, handler FileHandler) {
	fileHandlers[name] = handler
}

// Register default image handler
type imageHandler struct{}

type regularHandler struct{}

func (imageHandler) CouldHandle(media ss.FileLibrary) bool {
	return media.IsImage()
}

func (imageHandler) Handle(media ss.FileLibrary, file multipart.File, option *mod.Options) (hash,url string, fullpath string, err error) {

	url = media.GetURL(option, media)

	if url == "" {
		return "","", "", errors.New("no url found.")
	}

	hash,rd,err,fileurl := getHash(file)
	if err != nil {
		return "","","",err
	}
	if fileurl != ""{
		return  hash, url, fileurl, nil
	}

	comm.Logger.Finest("(imageHandle)url=%s", url)

	if fullpath, err = media.Store(url, option, rd); err == nil {

		logger.Finest(" 文件已经存储。。。")
		if media.NeedCrop() {
			file.Seek(0, 0)

			if img, err := imaging.Decode(file); err == nil {
				if format, err := utils.GetImageFormat(media.URL()); err == nil {
					if cropOption := media.GetCropOption("original"); cropOption != nil {
						img = imaging.Crop(img, *cropOption)
					}

					// Save default image
					var buffer bytes.Buffer
					imaging.Encode(&buffer, img, *format)
					fullpath, err = media.Store(media.URL(), option, &buffer)

					for key, size := range media.GetSizes() {
						newImage := img
						if cropOption := media.GetCropOption(key); cropOption != nil {
							newImage = imaging.Crop(newImage, *cropOption)
						}

						dst := imaging.Thumbnail(newImage, size.Width, size.Height, imaging.Lanczos)
						var buffer bytes.Buffer
						imaging.Encode(&buffer, dst, *format)
						media.Store(media.URL(key), option, &buffer)
					}
					return hash,url, fullpath, nil
				} else {
					return hash,url, fullpath, err
				}
			} else {
				return hash,url, fullpath, err
			}
		}
	} else {
		logger.Error("不能存储文件: %v",err)
		return hash,url, fullpath, err
	}

	return hash,url, fullpath, err
}

/**
* 针对一般文件的处理
**/
func (regularHandler) CouldHandle(media ss.FileLibrary) bool {
	return true
}

func (regularHandler) Handle(media ss.FileLibrary, file multipart.File, option *mod.Options) (hash,url string, fullpath string, err error) {

	url = media.GetURL(option, media)
	if url == "" {
		return "","", "", errors.New("no url found.")
	}

	hash,rd,err, fileurl :=getHash(file)
	if err != nil {
		return "","","",err
	}

	if fileurl != ""{
		return hash, url, fileurl, nil
	}

	if fullpath, err = media.Store(url, option, rd); err != nil {
		comm.Logger.Error("(Handle)", err)
	}

	return hash,url, fullpath, err
}

func GetFileHandler(key string) FileHandler {
	if handler, ok := fileHandlers[key]; ok {
		return handler
	} else {
		return nil
	}
}



func getHandler(ftype int) FileHandler {
	var fh FileHandler
	switch ftype {
	case comm.FILE_TYPE_AUDIO:
		fh = GetFileHandler("regular_file_handler")
	case comm.FILE_TYPE_IMAGE:
		fh = GetFileHandler("image_file_handler")
	case comm.FILE_TYPE_REGULAR:
		fh = GetFileHandler("regular_file_handler")
	}

	return fh
}


func getHash(r io.Reader)(string,*bytes.Reader,error, string) {
	b,err:=ioutil.ReadAll(r)
	if err != nil {
		logger.Error("不能读取文件内容: %v",err)
	}
	hb :=md5.Sum(b)
	hash := fmt.Sprintf("%x",hb)

	fi := &mod.FileInfo{Hash: hash}
	has,err := orm.Get(fi)

	logger.Finest("len=%d,hash=%s,has=%v,err=%v",len(b),hash,has,err)

	// 已经存在的数据直接返回url
	/*
	if has {
		return "",nil,errors.New("已存在相同的文件")
	}*/

	rd := bytes.NewReader(b)

	if has{
		return hash,rd,nil, fi.Location
	}else{
		return hash,rd,nil, ""
	}

}