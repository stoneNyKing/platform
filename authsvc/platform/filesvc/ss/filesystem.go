package ss

import (
	_ "github.com/spf13/viper/remote"
	"io"
	"os"
	"path/filepath"
	"platform/filesvc/comm"
	"platform/filesvc/dbmod"
)

type FileSystem struct {
	Base
}

func (f FileSystem) GetFullPath(url string, option *dbmod.Options) (path string, err error) {
	sitepath := option.GetSitePath(f.Site, option.GetCategory(), f.FileType)

	if comm.FilePath != "" {
		sitepath = comm.FilePath
	}

	if sitepath != "" {
		path = filepath.Join(sitepath, url)
	} else {
		path = url
	}

	dir := filepath.Dir(path)

	comm.Logger.Debug("dir=%s", dir)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
	}

	comm.Logger.Debug("[GetFullPath] path =%s", path)

	return
}

func (f FileSystem) Store(url string, option *dbmod.Options, reader io.Reader) (string, error) {

	comm.Logger.Finest("store the upload file.url=%s", url)

	if fullpath, err := f.GetFullPath(url, option); err == nil {
		if dst, err := os.Create(fullpath); err == nil {
			_, err := io.Copy(dst, reader)
			return fullpath, err
		} else {
			return fullpath, err
		}
	} else {
		return fullpath, err
	}
}

func (f FileSystem) Retrieve(url string) (*os.File, error) {

	comm.Logger.Finest("call retrieve")

	if fullpath, err := f.GetFullPath(url, nil); err == nil {
		return os.Open(fullpath)
	} else {
		return nil, os.ErrNotExist
	}
}
