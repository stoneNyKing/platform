package models

import (
	"errors"
	"golang.org/x/net/context"
	"mime/multipart"
	"net/http"

	"platform/filesvc/comm"
	"platform/filesvc/ss"
	"platform/filesvc/ss/aliyun"

	_ "github.com/spf13/viper/remote"
	"io"
	"os"
	"platform/common/utils"
	mod "platform/filesvc/dbmod"
	"platform/filesvc/imconf"
	"strings"
)

type FileService interface {
	Upload(req *uploadRequest) (string, error)
	Download(root,url string) (string, error)
	Health(url string) (string, error)
}

type fileService struct{}


type uploadRequest struct {
	Content    string `json:"content"`
	Url        string `json:"-"`
	Appid      string `json:"appid"`
	Site       string `json:"site"`
	FileType   int     `json:"filetype"`
	Storage    int    `json:"storage"`
	Category    int    `json:"category"`
	FileName	string  `json:"-"`
	Hash		string  `json:"-"`
	FileHeader []*multipart.FileHeader
}



var meta_upload_name string

func InitFilename() {
	meta_upload_name = comm.GetKeyString("meta.upload_filename")
	if meta_upload_name == "" {
		meta_upload_name = "uploadfile"
	}
}

func (svc fileService) Upload(req *uploadRequest) (string, error) {

	var flibrary ss.FileLibrary

	if req.Storage == comm.STORAGE_TYPE_FILESYSTEM {
		flibrary = new(ss.FileSystem)
	} else if req.Storage == comm.STORAGE_TYPE_ALIOSS {
		flibrary = new(aliyun.OSS)
	}else if req.Storage == comm.STORAGE_TYPE_DATABASE {
		flibrary = new(ss.DbStorage)
	}

	if flibrary == nil {
		return "", errors.New("no storage type set.")
	}

	var file multipart.File
	var err error

	err = flibrary.Scan(req.FileHeader, req.Appid, req.Site, utils.ConvertToString(req.FileType))
	if err != nil {
		logger.Error("上传失败：%v",err)
		return "",err
	}

	handler := getHandler(req.FileType)

	if fileHeader := flibrary.GetFileHeader(); fileHeader != nil {
		file, err = flibrary.GetFileHeader().Open()
	} else {
		comm.Logger.Debug("msg = %s", req.Site+"|"+utils.ConvertToString(req.FileType))
		file, err = flibrary.Retrieve(flibrary.URL("original"))
	}

	logger.Finest("开始处理文件...")

	url := ""
	hash := ""
	op := &mod.Options{}
	op.SetCategory(req.Category)

	if file != nil {
		defer file.Close()
		var handled = false
		if handler.CouldHandle(flibrary) {
			if hash,url, req.Url, err = handler.Handle(flibrary, file, op); err == nil {
				logger.Debug("save file success.url = %s", url)
				handled = true
				req.Hash = hash
			}else{
				return "",err
			}
		} else {
			logger.Error("not a recognized file type.")
		}

		// Save File
		if !handled {

			hash,rd,err, fileurl :=getHash(file)
			if err != nil {
				return "",err
			}

			if fileurl != ""{
				return  fileurl, nil
			}

			req.Hash = hash
			req.Url, err = flibrary.Store(flibrary.GetURL(op, flibrary), op, rd)
		}

		if req.Storage == comm.STORAGE_TYPE_FILESYSTEM {
			comm.Logger.Fine("return url=%s", url)
			prefix := op.GetSitePrefix(req.Site,req.Category, req.FileType)
			if comm.Prefix != "" {
				prefix = "/"+comm.Prefix
			}
			req.Url = prefix + url
		}


	} else {
		logger.Error("error:", err)
	}

	req.FileName = flibrary.GetFileName()

	return req.Url, nil
}

func (svc fileService) Download(root,url string) (string, string, error) {

	fns := strings.Split(url,"/")
	fn := ""
	if len(fns) > 0 {
		fn = fns[len(fns)-1]
	}

	res, err := http.Get(url)
	if err != nil {
		return "","",err
	}
	f, err := os.Create(comm.FilePath+"/"+fn)
	if err != nil {
		return "","",err
	}
	io.Copy(f, res.Body)

	f.Close()


	return imconf.Config.RedirectHost+"/"+comm.Prefix+"/" + fn,comm.FilePath+"/"+fn,nil

}
func (svc fileService) Health(url string) (string, error) {
	return "", nil

}


func decodeUploadRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request uploadRequest
	var appid string

	//fmt.Printf("run here\n")

	if r == nil {
		return nil, errors.New("request is nil.")
	}

	s := r.URL.Query().Get("site")
	if s == "" {
		logger.Error("(decodeUploadRequest) site is empty.")
		return nil, errors.New("request site is empty.")
	}

	site := "site" + s
	request.Site = site

	logger.Finest("(decodeUploadRequest) site is %s.", site)

	appid = r.URL.Query().Get("appid")
	if appid == "" {
		return nil, errors.New("request appid is empty.")
	}

	request.Appid = appid
	if appid == "" {
		return nil, errors.New("request appid is empty.")
	}

	request.Category = utils.Convert2Int(r.URL.Query().Get("category"))

	if r.Method != "POST" && r.Method != "PUT" && r.Method != "OPTIONS"  {
		return nil, errors.New("request method error.")
	}

	if r.Method == "OPTIONS" {
		return request,nil
	}

	r.ParseMultipartForm(32 << 22)

	if r.MultipartForm == nil {
		return nil,errors.New("MultipartForm wrong(uploadfile).")
	}

	if r.MultipartForm.File[meta_upload_name] == nil {
		return nil,errors.New("file upload type wrong(uploadfile).")
	}

	fh := r.MultipartForm.File[meta_upload_name]

	if fh == nil {
		return nil, errors.New("no file upload.")
	}
	ftype := comm.GetFileType(fh[0].Filename)

	request.FileHeader = fh

	request.FileType = ftype

	b,ft := mod.Option.GetSiteStorage(site,request.Category, ftype)

	if !b {
		return nil, errors.New("storage location is not set.")
	}

	logger.Finest("storage type is : %v", ft)

	request.Storage = ft
	return request, nil
}

