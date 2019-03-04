package ss

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"image"
	"io"
	"mime/multipart"
	"os"
	"path"
	"platform/filesvc/comm"
	"strings"
	"platform/common/utils"
	"platform/filesvc/dbmod"

	_ "github.com/spf13/viper/remote"
)

var ErrNotImplemented = errors.New("not implemented")

type CropOption struct {
	X, Y, Width, Height int
}

type fileHeader interface {
	Open() (multipart.File, error)
}

type fileWrapper struct {
	*os.File
}

func (fileWrapper *fileWrapper) Open() (multipart.File, error) {
	return fileWrapper.File, nil
}

type Base struct {
	Appid       string `json:"-"` //	客户端的appid
	FileName    string
	Url         string
	CropOptions map[string]*CropOption `json:",omitempty"`
	Crop        bool                   `json:"-"`
	Valid       bool                   `json:"-"`
	FileHeader  fileHeader             `json:"-"`
	Reader      io.Reader              `json:"-"`
	cropped     bool                   `json:"-"`
	Site        string                 `json:"-"`
	FileType    int                 `json:"-"`
}

func (b *Base) Scan(data interface{}, stype ...string) (err error) {
	switch values := data.(type) {
	case *os.File:
		b.FileHeader = &fileWrapper{values}
		b.FileName = path.Base(values.Name())
		b.Valid = true
	case []*multipart.FileHeader:
		if len(values) > 0 {
			file := values[0]
			b.FileHeader, b.FileName, b.Valid = file, file.Filename, true
		}
	case []byte:
		if err = json.Unmarshal(values, b); err == nil {
			b.Valid = true
		}
		var doCrop struct{ Crop bool }
		if err = json.Unmarshal(values, &doCrop); err == nil && doCrop.Crop {
			b.Crop = true
		}
	case string:
		b.Scan([]byte(values))
	case []string:
		for _, str := range values {
			b.Scan(str)
		}
	default:
		err = errors.New("unsupported driver -> Scan pair for MediaLibrary")
	}

	if len(stype) > 2 {
		comm.Logger.Finest("msg:(site=%s)(filetype=%s)", stype[1], stype[2])
		b.Appid = stype[0]
		b.SetSite(stype[1])
		b.SetFileType(utils.Convert2Int(stype[2]))
	}

	return
}

func (b Base) Value() (driver.Value, error) {
	if b.Valid {
		result, err := json.Marshal(b)
		return string(result), err
	}
	return nil, nil
}

func (b Base) URL(styles ...string) string {
	if b.Url != "" && len(styles) > 0 {
		ext := path.Ext(b.Url)
		return fmt.Sprintf("%v.%v%v", strings.TrimSuffix(b.Url, ext), styles[0], ext)
	}
	return b.Url
}

func (b Base) String() string {
	return b.URL()
}

func (b Base) GetFileName() string {
	return b.FileName
}

func (b Base) GetFileHeader() fileHeader {
	return b.FileHeader
}

func (b Base) GetURLTemplate(option *dbmod.Options) (path string) {
	if path = option.GetSiteURLTemplate(b.GetSite(),option.GetCategory(), b.GetFileType()); path == "" {
		path = "/{{appid}}/{{site}}/{{filetype}}/{{filename_with_hash}}"
	}
	return
}

func (b Base) Store(url string, option *dbmod.Options, reader io.Reader) (string, error) {
	return "", nil
}

func getFuncMap(appid string, site string, filetype int, filename string) template.FuncMap {
	hash := func() string {
		// return strings.Replace(time.Now().Format("20060102150506.000000000"), ".", "", -1)
		return ""
	}
	return template.FuncMap{
		"appid":    func() string { return appid },
		"site":     func() string { return site },
		"filetype": func() int  { return filetype },
		"filename": func() string {
			var fn string
			_, fn = path.Split(filename)
			return fn
		},
		"basename": func() string { return strings.TrimSuffix(path.Base(filename), path.Ext(filename)) },
		"hash":     hash,
		"filename_with_hash": func() string {
			var fn string
			_, fn = path.Split(filename)
			return fmt.Sprintf("%v.%v%v", strings.TrimSuffix(fn, path.Ext(fn)), hash(), path.Ext(fn))
		},
		"extension": func() string { return strings.TrimPrefix(path.Ext(filename), ".") },
	}
}

func (b Base) GetURL(option *dbmod.Options, templater dbmod.URLTemplater) string {
	if path := templater.GetURLTemplate(option); path != "" {
		tmpl := template.New("").Funcs(getFuncMap(b.Appid, b.GetSite(), b.GetFileType(), b.GetFileName()))
		if tmpl, err := tmpl.Parse(path); err == nil {
			var result = bytes.NewBufferString("")
			data := make(map[string]interface{})
			data["appid"] = b.Appid
			data["site"] = b.GetSite()
			data["filetype"] = b.GetFileType()
			data["filename"] = b.GetFileName()

			if err := tmpl.Execute(result, data); err == nil {
				comm.Logger.Finest("(GetUrl)result=%s", result.String())
				return result.String()
			} else {
				comm.Logger.Error("(GetURL) error = %v", err.Error())
			}
		}
	}
	return ""
}

func (b *Base) Cropped(values ...bool) (result bool) {
	result = b.cropped
	for _, value := range values {
		b.cropped = value
	}
	return result
}

func (b *Base) NeedCrop() bool {
	return b.Crop
}

func (b *Base) GetCropOption(name string) *image.Rectangle {
	if cropOption := b.CropOptions[strings.Split(name, "@")[0]]; cropOption != nil {
		return &image.Rectangle{
			Min: image.Point{X: cropOption.X, Y: cropOption.Y},
			Max: image.Point{X: cropOption.X + cropOption.Width, Y: cropOption.Y + cropOption.Height},
		}
	} else {
		return nil
	}
}

func (b Base) Retrieve(url string) (*os.File, error) {
	return nil, ErrNotImplemented
}

func (b Base) GetSizes() map[string]comm.Size {
	return map[string]comm.Size{}
}

func (b Base) IsImage() bool {
	_, err := utils.GetImageFormat(b.GetFileName())
	return err == nil
}

func (b *Base) SetSite(site string) {
	b.Site = site
}
func (b *Base) SetFileType(ft int) {
	b.FileType = ft
}
func (b *Base) GetSite() string {
	return b.Site
}
func (b *Base) GetFileType() int {
	return b.FileType
}

func (b *Base) GetAppType() string {
	return comm.GetAppType(b.Appid)
}
