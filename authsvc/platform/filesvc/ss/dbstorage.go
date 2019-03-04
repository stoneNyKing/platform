package ss

import (
	"io"
	"os"
	"path/filepath"
	"platform/filesvc/comm"
	"platform/filesvc/dbmod"
)

type DbStorage struct {
	Base
}

func (f DbStorage) GetFullPath(url string, option *dbmod.Options) (path string, err error) {
	sitepath := option.GetSitePath(f.Site,option.GetCategory(), f.FileType)

	if comm.FilePath != "" {
		sitepath = comm.FilePath
	}

	if sitepath != "" {
		path = filepath.Join(sitepath, url)
	} else {
		path = url
	}

	return
}

func (f DbStorage) Store(url string, option *dbmod.Options, reader io.Reader) (string, error) {

	comm.Logger.Finest("store the upload file.url=%s", url)
	fullpath, err := f.GetFullPath(url, option)
	return fullpath, err
}

func (f DbStorage) Retrieve(url string) (*os.File, error) {

	return nil,nil
}
