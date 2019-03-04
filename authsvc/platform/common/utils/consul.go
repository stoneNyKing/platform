// consul.go
package utils

import (
	"flag"
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
)

type ConsulClient struct {
	Address string
	client  *consulapi.Client
}

func (consul *ConsulClient) Init() {
	var addr string
	flag.StringVar(&addr, "consul", "127.0.0.1:8500", "host:port of the consul agent")
	flag.Parse()

	consul.Address = addr

	fmt.Printf("consul addr=%s\n", addr)

	config := consulapi.DefaultConfig()
	config.Address = consul.Address
	cc, err := consulapi.NewClient(config)

	consul.client = cc

	if err != nil {
		fmt.Printf("不能获取consul句柄：%v\n", err.Error())
	}
}

func (consul *ConsulClient) Get(key string, q *consulapi.QueryOptions) string {
	kv := consul.client.KV()

	kvp, _, err := kv.Get(key, nil)
	if err != nil {
		return ""
	} else {
		return string(kvp.Value)
	}
}

func (consul *ConsulClient) Put(key string, value string, w *consulapi.WriteOptions) (err error) {
	kv := consul.client.KV()

	d := &consulapi.KVPair{Key: key, Value: []byte(value)}

	_, err = kv.Put(d, w)

	return err
}

func (consul *ConsulClient) Delete(key string, w *consulapi.WriteOptions) (err error) {
	kv := consul.client.KV()

	_, err = kv.Delete(key, w)
	return err
}
