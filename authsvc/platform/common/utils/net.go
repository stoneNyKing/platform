package utils

import (
	"strings"
	"net"
	"net/url"
)

func Hostname2IPv4( hostn string ) ( ip string) {
	if hostn == "" {
		return ip
	}

	ss := strings.Split(hostn,":")
	
	s1,_ := net.LookupHost(ss[0])

	var p string

	for _,v := range s1 {
		if v != "::1" {
			p = v
			break
		}
	}

	if len(ss)>1 {
		ip = p +":"+ss[1]
	}else{
		ip = p
	}

	return
}

func ParseForm(body string)(ret map[string]interface{}) {

	if body == "" {
		return
	}

	vs,err := url.ParseQuery(body)
	if err != nil {
		return nil
	}
	ret = make(map[string]interface{})

	for k, v1 := range vs {
		if len(v1) > 0 {
			ret[k] = v1[0]
		}
	}

	return ret
}
