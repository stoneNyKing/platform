//+build !docker
//

package main

import (
	"github.com/spf13/viper"
	"fmt"
	"platform/confsvc/imconf"
	"strings"

	"sync"
	//修改为viper

	_ "github.com/spf13/viper/remote"

	"platform/common/cmd"

	"platform/common/utils"
)

var once sync.Once

func GetSettings() {
	//vip设置路径
	viper.SetConfigName("confsvc")
	viper.AddConfigPath("/etc/confsvc")
	viper.AddConfigPath("$HOME/etc")
	viper.AddConfigPath(".")
	viper.SetConfigType("json")

	var addr, consuladdr, token,rpcxaddr string

	cmd.RootCmd.PersistentFlags().StringVar(&rpcxaddr, "rpcx", "", "config this rpcx service address format(ip:port).")

	cmd.RootCmd.AddCommand(versionCmd)
	cmd.RootCmd.Execute()
	addr = utils.Hostname2IPv4(cmd.Address)
	consuladdr = utils.Hostname2IPv4(cmd.Consul)
	token = cmd.Token
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
					err = viper.AddRemoteProvider("etcd", "http://"+ss[0], "/config/confsvc.json")
				}else{
					err = viper.AddRemoteProvider("consul", ss[0], "/config/confsvc.json")
				}
				if err != nil {
					panic("设置远程provider出错")
				}
				err = viper.ReadRemoteConfig()
				if err != nil {
					panic("获取kv值出错")
				}
			}

			imconf.Config.HttpAddress = addr
			imconf.Config.ConsulToken = token
			imconf.Config.ConsulAddress = consuladdr
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
			if  rpcxaddr != "" {
				imconf.Config.RpcxAddr = utils.Hostname2IPv4(rpcxaddr)
			}
			if cmd.ServiceConf != "" {
				imconf.Config.ServiceConf = cmd.ServiceConf
			}
			imconf.Config.RecordAddr = imconf.Config.HttpAddress
		}
	})

}