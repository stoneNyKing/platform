//+build !docker
//

package main

import (
	"fmt"
	"github.com/spf13/viper"
	"platform/ousvc/config"
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
	viper.SetConfigName("ousvc")
	viper.AddConfigPath("/etc/ousvc")
	viper.AddConfigPath("$HOME/etc")
	viper.AddConfigPath(".")
	viper.SetConfigType("json")

	var addr, consuladdr, token, rpcxaddr string

	cmd.RootCmd.PersistentFlags().StringVar(&rpcxaddr, "rpcx", "", "config this rpcx service address format(ip:port).")

	cmd.RootCmd.AddCommand(versionCmd)
	cmd.RootCmd.Execute()
	addr = utils.Hostname2IPv4(cmd.Http)
	consuladdr = utils.Hostname2IPv4(cmd.Consul)
	token = cmd.Token
	isChild = cmd.IsChild
	socketOrder = cmd.SocketOrder

	once.Do(func() {
		if config.Config == nil {
			config.Config = new(config.ImConf)

			if cmd.Sdt != "" {
				config.Config.Sdt = cmd.Sdt
			} else {
				config.Config.Sdt = "consul"
			}
			if cmd.Sda != "" {
				config.Config.Sda = cmd.Sda
			} else {
				config.Config.Sda = cmd.Consul
			}
			err := viper.ReadInConfig()
			if err != nil {
				fmt.Println("不能读取本地配置！")
				ss := strings.Split(config.Config.Sda, ";")
				if len(ss) <= 0 {
					panic("无法读取配置。")
				}
				if cmd.Sdt == "etcd" {
					err = viper.AddRemoteProvider("etcd", "http://"+ss[0], "/config/ousvc.json")
				} else {
					err = viper.AddRemoteProvider("consul", ss[0], "/config/ousvc.json")
				}
				if err != nil {
					panic("设置远程provider出错")
				}
				err = viper.ReadRemoteConfig()
				if err != nil {
					panic("获取kv值出错")
				}
			}

			config.Config.HttpAddress = addr
			config.Config.ConsulToken = token
			config.Config.ConsulAddress = consuladdr
			config.Config.ReadConf()

			if cmd.DebugAddr != "" {
				config.Config.DebugAddr = utils.Hostname2IPv4(cmd.DebugAddr)
			}

			if cmd.ZipkinAddr != "" {
				config.Config.ZipkinUrl = cmd.ZipkinAddr
			}
			if cmd.KafkaAddr != "" {
				config.Config.KafkaAddress = cmd.KafkaAddr
			}
			if cmd.AppdashAddr != "" {
				config.Config.AppdashAddr = utils.Hostname2IPv4(cmd.AppdashAddr)
			}
			if cmd.LightstepToken != "" {
				config.Config.LightstepToken = cmd.LightstepToken
			}
			config.Config.Debug = cmd.Debug

			if rpcxaddr != "" {
				config.Config.RpcxAddr = utils.Hostname2IPv4(rpcxaddr)
			}
			if cmd.Http != "" {
				config.Config.HttpAddress = cmd.Http
			}
			if cmd.Https != "" {
				config.Config.HttpsAddress = cmd.Https
			}
			if cmd.ServiceConf != "" {
				config.Config.ServiceConf = cmd.ServiceConf
			}
			config.Config.RecordAddr = config.Config.HttpAddress
		}
	})

}
