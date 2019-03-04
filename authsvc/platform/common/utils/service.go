package utils

import (
	"strings"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"errors"
	"net"
	"time"
	"context"
)

func ServiceGet(uri string)(interface{},error) {

	resp, err := http.Get(uri)

	if resp == nil || err != nil  {
		fmt.Printf("no response body return or has error(%v).\n",err)
		return 0,errors.New("no response body return or has error")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		fmt.Println(string(body))
		return 0,err
	}


	var vs interface{}
	err = json.Unmarshal(body,&vs)

	if err != nil {
		fmt.Printf("cannot convert json.\n")
		return 0,err
	}

	return vs,err
}


func ServiceGetWithTimeout(uri string,timeout int)(interface{},error) {

	if timeout < 0 {
		timeout = 5
	}

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context,netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*time.Duration(timeout)) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Second * time.Duration(timeout))) //设置发送接受数据超时
				return conn, nil
			},
		},
	}

	res, err := client.Get(uri)
	if err != nil {
		return 0,err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		// handle error
		fmt.Println(string(body))
		return 0,err
	}


	var vs interface{}
	err = json.Unmarshal(body,&vs)

	if err != nil {
		fmt.Printf("cannot convert json.\n")
		return 0,err
	}

	return vs,err
}


func ServicePost(uri string,params string)(interface{},error) {

	return servPost("application/json",uri,params)
}

func ServicePostFormData(uri string,params string)(interface{},error) {

	return servPost("application/x-www-form-urlencoded",uri,params)
}


func servPost(contenttype string,uri string,params string)(interface{},error) {
	//"application/json"

	resp, err := http.Post(uri,contenttype,strings.NewReader(params))

	if resp == nil || err != nil  {
		fmt.Printf("no response body return or has error(%v).\n",err)
		return nil,errors.New("no response body return or has error")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		fmt.Println(string(body))
		return nil,err
	}else {
		fmt.Println(string(body))
	}


	var vs map[string]interface{}
	err = json.Unmarshal(body,&vs)

	if err != nil {
		fmt.Printf("cannot convert json.\n")
		return nil,err
	}

	return vs,err
}

func ServicePostWithTimeout(uri string,params string,timeout int)(interface{},error) {

 	return servPostWithTimeout("application/json",uri,params,timeout)
}

func ServicePostFDWithTimeout(uri string,params string,timeout int)(interface{},error) {

	return servPostWithTimeout("application/x-www-form-urlencoded",uri,params,timeout)
}


func servPostWithTimeout(contenttype string,uri string,params string,timeout int)(interface{},error) {
	if timeout < 0 {
		timeout = 5
	}

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context,netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*time.Duration(timeout)) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Second * time.Duration(timeout))) //设置发送接受数据超时
				return conn, nil
			},
		},
	}

	resp, err := client.Post(uri,contenttype,strings.NewReader(params))

	if resp == nil || err != nil  {
		fmt.Printf("no response body return or has error(%v).\n",err)
		return nil,errors.New("no response body return or has error")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		fmt.Println(string(body))
		return nil,err
	}else {
		fmt.Println(string(body))
	}


	var vs map[string]interface{}
	err = json.Unmarshal(body,&vs)

	if err != nil {
		fmt.Printf("cannot convert json.\n")
		return nil,err
	}

	return vs,err
}


func ServicePut(uri string,params string)(interface{},error) {
	return servPut("application/json",uri,params)
}

func ServicePutFormData(uri string,params string)(interface{},error) {
	return servPut("application/x-www-form-urlencoded",uri,params)
}


func servPut(contenttype string,uri string,params string)(interface{},error) {
	client := &http.Client{}
	req,err := http.NewRequest(http.MethodPut,uri,strings.NewReader(params))

	req.Header.Add("Content-Type",contenttype)

	resp, err := client.Do(req)

	if resp == nil || err != nil  {
		fmt.Printf("no response body return or has error(%v).\n",err)
		return 0,errors.New("no response body return or has error")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		fmt.Println(string(body))
		return 0,err
	}


	var vs map[string]interface{}
	err = json.Unmarshal(body,&vs)

	if err != nil {
		fmt.Printf("cannot convert json.\n")
		return 0,err
	}

	return vs,err
}

func ServicePutWithTimeout(uri string,params string,timeout int)(interface{},error) {
	return servPutFDWithTimeout("application/json",uri,params,timeout)
}

func ServicePutFDWithTimeout(uri string,params string,timeout int)(interface{},error) {
	return servPutFDWithTimeout("application/x-www-form-urlencoded",uri,params,timeout)
}


func servPutFDWithTimeout(contenttype, uri string,params string,timeout int)(interface{},error) {
	if timeout < 0 {
		timeout = 5
	}

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context,netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*time.Duration(timeout)) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Second * time.Duration(timeout))) //设置发送接受数据超时
				return conn, nil
			},
		},

	}
	req,err := http.NewRequest(http.MethodPut,uri,strings.NewReader(params))

	req.Header.Add("Content-Type",contenttype)

	resp, err := client.Do(req)

	if resp == nil || err != nil  {
		fmt.Printf("no response body return or has error(%v).\n",err)
		return 0,errors.New("no response body return or has error")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		fmt.Println(string(body))
		return 0,err
	}


	var vs map[string]interface{}
	err = json.Unmarshal(body,&vs)

	if err != nil {
		fmt.Printf("cannot convert json.\n")
		return 0,err
	}

	return vs,err
}

func ServiceDelete(uri string,params string)(interface{},error) {

	return servDel("application/json",uri,params)
}

func ServiceDeleteFormData(uri string,params string)(interface{},error) {

	return servDel("application/x-www-form-urlencoded",uri,params)
}

func servDel(contenttype string,uri string,params string)(interface{},error) {
	client := &http.Client{}
	req,err := http.NewRequest(http.MethodDelete,uri,strings.NewReader(params))

	req.Header.Add("Content-Type",contenttype)

	resp, err := client.Do(req)

	if resp == nil || err != nil  {
		fmt.Printf("no response body return or has error(%v).\n",err)
		return 0,errors.New("no response body return or has error")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		fmt.Println(string(body))
		return 0,err
	}


	var vs map[string]interface{}
	err = json.Unmarshal(body,&vs)

	if err != nil {
		fmt.Printf("cannot convert json.\n")
		return 0,err
	}

	return vs,err
}

func ServiceDeleteWithTimeout(uri string,params string,timeout int)(interface{},error) {
	return servDelWithtimeout("application/json",uri,params,timeout)
}

func ServiceDeleteFDWithTimeout(uri string,params string,timeout int)(interface{},error) {
	return servDelWithtimeout("application/x-www-form-urlencoded",uri,params,timeout)
}

func servDelWithtimeout(contenttype string,uri string,params string,timeout int)(interface{},error) {
	if timeout < 0 {
		timeout = 5
	}

	client := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context,netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*time.Duration(timeout)) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Second * time.Duration(timeout))) //设置发送接受数据超时
				return conn, nil
			},
		},

	}
	req,err := http.NewRequest(http.MethodDelete,uri,strings.NewReader(params))

	req.Header.Add("Content-Type",contenttype)

	resp, err := client.Do(req)

	if resp == nil || err != nil  {
		fmt.Printf("no response body return or has error(%v).\n",err)
		return 0,errors.New("no response body return or has error")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		fmt.Println(string(body))
		return 0,err
	}


	var vs map[string]interface{}
	err = json.Unmarshal(body,&vs)

	if err != nil {
		fmt.Printf("cannot convert json.\n")
		return 0,err
	}

	return vs,err

}