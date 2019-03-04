package sd

import (
	"encoding/json"
	consulsd "github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"platform/common/utils"
	"platform/mskit/grace"
	mslog "platform/mskit/log"
	"reflect"
	"strconv"
	"strings"
)

var logger = mslog.Mslog

func getConsulClient(addr, schema string) consulsd.Client {
	var client consulsd.Client
	{
		consulConfig := api.DefaultConfig()
		if addr != "" {
			consulConfig.Address = addr
		}
		if schema != "" {
			consulConfig.Scheme = schema
		} else {
			consulConfig.Scheme = "http"
		}
		consulClient, err := api.NewClient(consulConfig)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		client = consulsd.NewClient(consulClient)
	}

	return client
}

func Register(app *grace.MicroService, schema, name string, prefix string, addr, consul, token string, callback ServiceCallback, params map[string]interface{}) {

	if name == "" {
		log.Fatal("name empty")
	}
	if prefix == "" {
		log.Fatal("prefix empty")
	}

	//consul address split
	cs := strings.Split(consul, ",")

	if len(cs) <= 0 {
		log.Fatal("no consul address config")
		return
	}

	consul = cs[0]

	var interval, timeout string
	if params != nil {
		if params["interval"] != nil {
			interval = utils.ConvertToString(params["interval"])
		}
		if params["timeout"] != nil {
			timeout = utils.ConvertToString(params["timeout"])
		}
	}

	if interval == "" {
		interval = "30s"
	}

	if timeout == "" {
		timeout = "2s"
	}

	prefixes := strings.Split(prefix, ",")
	host, portstr, err := net.SplitHostPort(addr)
	if err != nil {
		log.Fatal(err)
	}
	port, err := strconv.Atoi(portstr)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		log.Printf("Listening on %s serving %s", addr, prefix)
		if err := callback(app, params); err != nil {
			log.Fatal(err)
		}
	}()

	var tags []string
	for _, p := range prefixes {
		tags = append(tags, "urlprefix-"+p)
	}

	serviceID := name + "-" + addr
	service := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    name,
		Port:    port,
		Address: host,
		Tags:    tags,
		Check: &api.AgentServiceCheck{
			CheckID:  "check-" + serviceID,
			HTTP:     schema + "://" + addr + prefix + "/health",
			Interval: interval,
			Timeout:  timeout,
		},
	}

	c := getConsulClient(consul, schema)

	reg := consulsd.NewRegistrar(c, service, logger)

	reg.Register()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	reg.Deregister()

	log.Printf("Deregistered service %q in consul", name)
}

func readFile(path string) []byte {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func RegisterWithConf(app *grace.MicroService, schema string, fname string, consul, token string, callbacks ...ServiceCallback) {

	if fname == "" {
		log.Fatal("没有指定配置文件。\n")
		return
	}

	body := readFile(fname)

	var data map[string]interface{}
	var params map[string]interface{}

	err := json.Unmarshal(body, &data)

	if err != nil {
		log.Fatal("json:" + err.Error())
		return
	}

	//consul address split
	cs := strings.Split(consul, ",")

	if len(cs) <= 0 {
		log.Fatal("no consul address config")
		return
	}

	consul = cs[0]

	var p interface{}
	key := schema + "_"
	if data[key+"service"] != nil {
		p = data[key+"service"].(interface{})
	} else if data[key+"services"] != nil {
		p = data[key+"services"].(interface{})
	}

	cps := make(map[string]interface{})

	if data["TLSConfig"] != nil {
		vs := data["TLSConfig"].(map[string]interface{})
		cps["certfile"] = vs["certfile"]
		cps["keyfile"] = vs["keyfile"]
		cps["trustfile"] = vs["trustfile"]
	}

	if schema == "http" || schema == "https" {
		t := reflect.ValueOf(p)
		switch t.Kind() {
		case reflect.Slice:
			ps := p.([]interface{})
			if len(ps) != len(callbacks) {
				log.Fatal("服务数量与回调函数数量不匹配。")
				return
			}
			for i, vs := range ps {
				v := vs.(map[string]interface{})
				go registerService(app, schema, consul, token, v, callbacks[i], cps)
			}

			quit := make(chan os.Signal, 1)
			signal.Notify(quit, os.Interrupt, os.Kill)
			<-quit

			//select {}
		case reflect.Map:
			if len(callbacks) < 1 {
				log.Fatal("没有指定回调函数。")
				return
			}
			params = p.(map[string]interface{})
			registerService(app, schema, consul, token, params, callbacks[0], cps)
		}
	} else if schema == "rpcx" {
		if data["rpcx"] != nil {
			vs := data["rpcx"].([]interface{})
			for i, vv := range vs {
				v := vv.(map[string]interface{})
				m := make(map[string]interface{})
				if v["address"] != nil {
					m["host"] = utils.ConvertToString(v["address"])
				}
				if v["port"] != nil {
					m["port"] = utils.ConvertToString(v["port"])
				}
				if v["consul_address"] != nil {
					m["consul_address"] = v["consul_address"]
				}else{
					if v["sd_address"] != nil {
						m["consul_address"] = v["sd_address"]
						m["sd_address"] = v["sd_address"]
						m["sd_type"] = v["sd_type"]
					}
				}

				go callbacks[i](app, m)
			}
		} else {
			log.Fatal("没有配置参数。")
			panic("没有配置参数")
		}
	}

}

func registerService(app *grace.MicroService, schema, consul, token string, params map[string]interface{}, callback ServiceCallback, datas map[string]interface{}) {
	var name, prefix, host, addr string
	var tags []string

	//log.Print(params)

	if params["name"] != nil {
		name = utils.ConvertToString(params["name"])
	}
	if params["tags"] != nil {
		p := params["tags"].([]interface{})
		for _, v := range p {
			tags = append(tags, utils.ConvertToString(v))
		}
	}

	if params["address"] != nil {
		host = utils.Hostname2IPv4(utils.ConvertToString(params["address"]))
		datas["host"] = host
	}

	var port int
	if params["port"] != nil {
		port = utils.Convert2Int(params["port"])
		datas["port"] = port
	}

	if port == 0 {
		log.Fatal("没有指定端口号。")
		return
	}

	prefix = strings.Join(tags, ",")

	go func(po int) {
		log.Printf("Listening on %s:%d serving %s", host, po, prefix)
		if err := callback(app, datas); err != nil {
			log.Fatal(err)
		}
	}(port)

	addr = host + ":" + utils.ConvertToString(port)

	var serviceID string

	if params["id"] != nil {
		serviceID = utils.ConvertToString(params["id"])
	}

	if serviceID == "" {
		serviceID = name + "-" + addr
	}

	var checks []*api.AgentServiceCheck
	var check *api.AgentServiceCheck

	if params["checks"] != nil {
		vp := params["checks"].([]interface{})
		for _, v := range vp {
			p := v.(map[string]interface{})
			var c api.AgentServiceCheck
			if p["http"] != nil {
				c.HTTP = utils.ConvertToString(p["http"])
			}
			if p["interval"] != nil {
				c.Interval = utils.ConvertToString(p["interval"])
			}
			if p["timeout"] != nil {
				c.Timeout = utils.ConvertToString(p["timeout"])
			}
			if p["name"] != nil {
				c.Name = utils.ConvertToString(p["name"])
			}
			if p["id"] != nil {
				c.CheckID = utils.ConvertToString(p["id"])
			}
			if p["tcp"] != nil {
				c.TCP = utils.ConvertToString(p["tcp"])
			}
			if p["shell"] != nil {
				c.Shell = utils.ConvertToString(p["shell"])
			}
			if p["ttl"] != nil {
				c.TTL = utils.ConvertToString(p["ttl"])
			}
			if p["method"] != nil {
				c.Method = utils.ConvertToString(p["method"])
			}
			if p["status"] != nil {
				c.Status = utils.ConvertToString(p["status"])
			}

			if p["args"] != nil {
				vs := p["args"].([]interface{})
				for _, s := range vs {
					c.Args = append(c.Args, utils.ConvertToString(s))
				}
			}
			if p["notes"] != nil {
				c.Notes = utils.ConvertToString(p["notes"])
			}
			if p["grpc"] != nil {
				c.GRPC = utils.ConvertToString(p["grpc"])
			}

			if p["docker_container_id"] != nil {
				c.DockerContainerID = utils.ConvertToString(p["docker_container_id"])
			}

			if p["tls_skip_verify"] != nil {
				c.TLSSkipVerify = p["tls_skip_verify"].(bool)
			}
			if p["grpc_use_tls"] != nil {
				c.GRPCUseTLS = p["grpc_use_tls"].(bool)
			}

			if p["header"] != nil {
				vs := p["header"].(map[string]interface{})
				var h map[string][]string
				h = make(map[string][]string)

				for k, v := range vs {
					var ss []string
					s1 := v.([]interface{})
					for _, s := range s1 {
						ss = append(ss, utils.ConvertToString(s))
					}
					h[k] = ss
				}

				c.Header = h
			}

			checks = append(checks, &c)
		}
	}

	if checks == nil {
		check = &api.AgentServiceCheck{
			CheckID:  "check-" + serviceID,
			HTTP:     schema + "://" + addr + prefix + "/health",
			Interval: "30s",
			Timeout:  "3s",
		}
	}
	service := &api.AgentServiceRegistration{
		ID:      serviceID,
		Name:    name,
		Port:    port,
		Address: host,
		Tags:    tags,
		Check:   check,
		Checks:  checks,
	}

	if params["meta"] != nil {
		vs := params["meta"].(map[string]interface{})

		var m map[string]string
		m = make(map[string]string)
		for k, v := range vs {
			m[k] = utils.ConvertToString(v)
		}
		service.Meta = m
	}

	if params["proxy_destination"] != nil {
		service.ProxyDestination = utils.ConvertToString(params["proxy_destination"])
	}

	if params["enable_tag_override"] != nil {
		service.EnableTagOverride = params["enable_tag_override"].(bool)
	}

	if params["kind"] != nil {
		service.Kind = api.ServiceKind(utils.ConvertToString(params["kind"]))
	}

	if params["connect"] != nil {
		var proxy api.AgentServiceConnect

		p := params["connect"].(map[string]interface{})

		if p["native"] != nil {
			proxy.Native = p["native"].(bool)
		}

		if p["proxy"] != nil {
			v := p["proxy"].(map[string]interface{})
			var c api.AgentServiceConnectProxy
			if v["config"] != nil {
				c.Config = v["config"].(map[string]interface{})
			}
			if v["command"] != nil {
				vs := v["command"].([]interface{})

				for _, s := range vs {
					c.Command = append(c.Command, utils.ConvertToString(s))
				}
			}

			proxy.Proxy = &c
		}

		service.Connect = &proxy
	}

	var	caddr string
	if params["consul_address"] != nil {
		caddr = utils.ConvertToString(params["consul_address"])
	}else{
		if params["sd_address"] != nil {
			caddr = utils.ConvertToString(params["sd_address"])
		}
	}

	if caddr != "" {
		consul = caddr
	}
	c := getConsulClient(consul, schema)

	reg := consulsd.NewRegistrar(c, service, logger)

	reg.Register()

	log.Printf("Registered service %q in consul with tags %q", name, strings.Join(tags, ","))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	reg.Deregister()

	log.Printf("Deregistered service %q in consul", name)

}
