package utils

import (
    "testing"
    "fmt"
)

func TestServicePost(t *testing.T) {
	uri := "http://dev.laoyou99.cn:80/service/331/interview?appid=7&siteid=1&token="
	params :=`
	{
	"orgid":"1",
	"content":"准备下个月入院。",
	"interviewtime":"2017-02-25 12:00:00",
	"customerid":"1",
	"communicationmode":"1",
	"visittype":"1",
	"operatorid":"2",
	"remark":""
	}
	`

	result,err := ServicePost(uri,params)
	
	fmt.Printf("POST:result=%v,err=%v\n",result,err)
	
}


func TestServicePut(t *testing.T) {
	uri := "http://dev.laoyou99.cn:80/service/331/interview/4?appid=7&siteid=1&token="
	params :=`
	{
	"id":"4",
	"orgid":"1",
	"content":"3月5日入院。",
	"interviewtime":"2017-02-25 13:00:00",
	"customerid":"1",
	"communicationmode":"1",
	"visittype":"1",
	"operatorid":"2",
	"remark":""
	}
	`
	result,err := ServicePut(uri,params)

	fmt.Printf("PUT:result=%v,err=%v\n",result,err)
}

func TestServiceDelete(t *testing.T) {

	uri := "http://dev.laoyou99.cn:80/service/331/interview/4?appid=7&siteid=1&token="
	params :=`
	{
	"id":[4]
	}
	`

	result,err := ServiceDelete(uri,params)

	fmt.Printf("DELETE:result=%v,err=%v\n",result,err)
}