package utils


import (
	"fmt"
	"testing"
)


func TestHostname2IPv4(t *testing.T) {
	str := "127.0.0.1:100102"

	ip := Hostname2IPv4(str)

	fmt.Printf("ip string: %s\n", ip)

	str = "localhost:100102"
	ip = Hostname2IPv4(str)
	fmt.Printf("ip2 string: %s\n", ip)

	str = "localhost"
	ip = Hostname2IPv4(str)
	fmt.Printf("ip3 string: %s\n", ip)


	str = "127.0.0.1"
	ip = Hostname2IPv4(str)
	fmt.Printf("ip4 string: %s\n", ip)

	str = ":10102"
	ip = Hostname2IPv4(str)
	fmt.Printf("ip5 string: %s\n", ip)

}

func TestParseForm(t *testing.T) {
	body := "val1=123&gender=1&age=12&name=王&address=%E5%98%89%E9%99%B5%E6%B1%9F%E4%B8%9C%E8%A1%9718%E5%8F%B7"

	param := ParseForm(body)

	fmt.Printf("param= %+v\n",param)
}