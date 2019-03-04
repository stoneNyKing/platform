package ss

import (
	"database/sql/driver"
	_ "github.com/spf13/viper/remote"
	"image"
	"io"
	"os"
	"platform/filesvc/comm"
	"platform/filesvc/dbmod"
)

type FileLibrary interface {
	Scan(value interface{}, oth ...string) error
	Value() (driver.Value, error)

	GetURLTemplate(*dbmod.Options) string
	GetURL(option *dbmod.Options, templater dbmod.URLTemplater) string

	GetFileHeader() fileHeader
	GetFileName() string

	GetSizes() map[string]comm.Size
	NeedCrop() bool
	Cropped(values ...bool) bool
	GetCropOption(name string) *image.Rectangle

	Store(url string, option *dbmod.Options, reader io.Reader) (string, error)
	Retrieve(url string) (*os.File, error)
	GetAppType() string

	IsImage() bool
	SetSite(string)
	SetFileType(int)
	GetSite() string
	GetFileType() int
	URL(style ...string) string
	String() string
}
