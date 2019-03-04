package aliyun

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	//neturl "net/url"

	"strings"
	//"fmt"
	"platform/filesvc/comm"
	"platform/filesvc/dbmod"
	"platform/filesvc/ss"

	l4g "github.com/libra9z/log4go"
	//"platform/common/utils"

	//"github.com/libra9z/ali-oss"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"platform/filesvc/imconf"
)

type OSS struct {
	ss.Base
}

var logger l4g.Logger = comm.Logger

var (
	AliOSSBucket       string
	AliOSSAccessKey    string
	AliOSSAccessSecret string
	AliOSSEndpoint     string
	AliOSSRegion       string
	AliOSSPrivRegion   string
)

const (
	PRIV_NET   = 1
	PUBLIC_NET = 0
)

var Network int

func InitOs() {

	AliOSSBucket = viper.GetString("oss.bucket")
	AliOSSAccessKey = viper.GetString("oss.accessKey")
	AliOSSAccessSecret = viper.GetString("oss.accessSecret")
	AliOSSRegion = viper.GetString("oss.region")

	AliOSSPrivRegion = viper.GetString("oss.priv_region")
	Network = viper.GetInt("oss.network")
	//aliossClient = alioss.NewClient( AliOSSAccessKey,AliOSSAccessSecret )
}

func getBucket(apptype string, option *dbmod.Options) string {

	key := "oss." + apptype + ".bucket"

	if bucket := option.GetString(key); bucket != "" {
		return bucket
	}

	return AliOSSBucket
}

func getEndpoint(apptype string, option *dbmod.Options, net int) string {

	key := "oss." + apptype + ".endpoint"

	if endpoint := option.GetString(key); endpoint != "" {
		return endpoint
	}

	var ret string

	if net == PUBLIC_NET {
		//ret = getBucket(apptype,option) + "." + AliOSSRegion
		ret = AliOSSRegion
	} else if net == PRIV_NET {
		ret = AliOSSPrivRegion
	}

	return ret
}

func (s OSS) GetURLTemplate(option *dbmod.Options) (path string) {

	if path = option.GetString("oss.endpoint"); path == "" {
		path = "/{{appid}}/{{site}}/{{filetype}}/{{filename_with_hash}}"
	}

	apptype := s.GetAppType()

	//fmt.Printf("apptype=%s\n",apptype)

	path = "//" + getEndpoint(apptype, option, PRIV_NET) + path

	return
}


func (s OSS) Store(url string, option *dbmod.Options, reader io.Reader) (u string, err error) {

	apptype := s.GetAppType()

	logger.Finest("(store) oss url= %s,apptype=%s ",url,apptype)

	priv_endpoint := getEndpoint(apptype, option, PRIV_NET)
	endpoint := getEndpoint(apptype, option, PUBLIC_NET)

	path := strings.Replace(url, "//"+priv_endpoint, "", -1)

	logger.Finest("oss priv_endpoint=%s", priv_endpoint)
	logger.Finest("oss endpoint=%s", endpoint)
	logger.Finest("oss path=%s", path)

	client, err := oss.New("http://"+getEndpoint(apptype, option, Network), AliOSSAccessKey, AliOSSAccessSecret)
	if err != nil {
		logger.Error("cannot create a new client")
		return
	}

	lsRes, err := client.ListBuckets()
	if err != nil {
		// HandleError(err)
		logger.Error("cannot list buckets: %s", err.Error())
	}

	for _, bucket := range lsRes.Buckets {
		logger.Finest("Buckets: %s", bucket.Name)
	}

	logger.Finest("oss bucket=%s", getBucket(apptype, option))

	bucket, err := client.Bucket(getBucket(apptype, option))
	if err != nil {
		logger.Error("cannot create a new bucket")
		return
	}

	var vpath string = path
	if strings.HasPrefix(path, "/") {
		vpath = strings.TrimPrefix(path, "/")
	}

	err = bucket.PutObject(vpath, reader)
	if err != nil {
		logger.Error("cannot put object: err= %s", err.Error())
		return
	}

	logger.Finest("oss AliOSSRegion=%s", AliOSSRegion)

	if imconf.Config.OSSNetworkType == PRIV_NET {
		u = "http://" + getBucket(apptype, option) + "." + priv_endpoint + path

	}else if imconf.Config.OSSNetworkType == PUBLIC_NET {
		u = "http://" + getBucket(apptype, option) + "." + endpoint + path
	}


	logger.Finest("(store)u=%s", u)

	return
}

func (s OSS) Retrieve(url string) (*os.File, error) {
	response, err := http.Get("http:" + url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if file, err := ioutil.TempFile("/tmp", "OSS"); err == nil {
		_, err := io.Copy(file, response.Body)
		return file, err
	} else {
		return nil, err
	}
}
