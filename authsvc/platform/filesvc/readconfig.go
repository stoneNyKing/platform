//+build !docker
//

package main

import (
	"fmt"
	"platform/filesvc/comm"
	"strings"

	"sync"
	//修改为viper

	_ "github.com/spf13/viper/remote"

	"platform/common/cmd"

	"path/filepath"
	"os"
	"platform/common/utils"
	"platform/filesvc/imconf"
	"github.com/spf13/viper"
)

var once sync.Once

func GetSettings() {
	//vip设置路径
	viper.SetConfigName("filesvc")
	viper.AddConfigPath("/etc/filesvc")
	viper.AddConfigPath("$HOME/etc")
	viper.AddConfigPath(".")
	viper.SetConfigType("json")

	cmd.RootCmd.PersistentFlags().StringVar(&comm.Prefix, "prefix", "", "config the download prefix.")
	cmd.RootCmd.PersistentFlags().StringVar(&comm.FilePath, "filepath", "", "config the attach file path.")

	cmd.RootCmd.AddCommand(versionCmd)
	cmd.RootCmd.Execute()
	comm.Address = utils.Hostname2IPv4(cmd.Http)
	comm.ConsulAddress = utils.Hostname2IPv4(cmd.Consul)
	comm.ConsulToken = cmd.Token
	isChild = cmd.IsChild
	socketOrder = cmd.SocketOrder

	once.Do(func() {
		if imconf.Config == nil {
			imconf.Config = new(imconf.ImConf)
			if cmd.Sdt != "" {
				imconf.Config.Sdt = cmd.Sdt
			}else{
				imconf.Config.Sdt = "consul"
			}
			if cmd.Sda != "" {
				imconf.Config.Sda = cmd.Sda
			}else{
				imconf.Config.Sda = cmd.Consul
			}
			err := viper.ReadInConfig()
			if err != nil {
				fmt.Println("不能读取本地配置！")
				ss := strings.Split(imconf.Config.Sda,";")
				if len(ss)<=0 {
					panic("无法读取配置。")
				}
				if cmd.Sdt == "etcd" {
					err = viper.AddRemoteProvider("etcd", "http://"+ss[0], "/config/filesvc.json")
				}else{
					err = viper.AddRemoteProvider("consul", ss[0], "/config/filesvc.json")
				}
				if err != nil {
					panic("设置远程provider出错")
				}
				err = viper.ReadRemoteConfig()
				if err != nil {
					panic("获取kv值出错")
				}
			}

			imconf.Config.HttpAddress = comm.Address
			imconf.Config.ConsulToken = cmd.Token
			imconf.Config.ConsulAddress = comm.ConsulAddress
			imconf.Config.ReadConf()

			if cmd.DebugAddr != "" {
				imconf.Config.DebugAddr = utils.Hostname2IPv4(cmd.DebugAddr)
			}

			if cmd.ZipkinAddr != "" {
				imconf.Config.ZipkinUrl = cmd.ZipkinAddr
			}
			if cmd.KafkaAddr != "" {
				imconf.Config.KafkaAddress = cmd.KafkaAddr
			}
			if cmd.AppdashAddr != "" {
				imconf.Config.AppdashAddr = cmd.AppdashAddr
			}
			if cmd.LightstepToken != "" {
				imconf.Config.LightstepToken = cmd.LightstepToken
			}
			imconf.Config.Debug = cmd.Debug
			if cmd.Http != "" {
				imconf.Config.HttpAddress = cmd.Http
			}
			if cmd.Https != "" {
				imconf.Config.HttpsAddress = cmd.Https
			}

			if cmd.ServiceConf != "" {
				imconf.Config.ServiceConf = cmd.ServiceConf
			}
			imconf.Config.RecordAddr = imconf.Config.HttpAddress
		}
	})


	comm.InitApptype()

	var logf string

	logf = imconf.Config.Logfile

	fmt.Printf("log file ==%s\n", logf)

	if logf == "" {
		logf = "../logs/filesvc.log"
	} else {
		dir := filepath.Dir(logf)

		fmt.Printf("dir=%s\n", dir)

		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.MkdirAll(dir, os.ModePerm)
		}
	}

	f, err := os.Open(logf)
	if err == nil && !os.IsNotExist(err) {
		fmt.Printf("file exist!\n")
		loggff := logf + "." + utils.GetTimeFormat("[2006-01-02 15:04:05]")
		f.Close()

		os.Rename(logf, loggff)
	}

	comm.InitLogger(logf, imconf.Config.LogLevel)


}