package redis

import (
    "testing"
    "fmt"
)

func TestGetDeviceAllProperties(t *testing.T) {
	
	var val map[string]string
	
	val = make(map[string]string)
	
	GetDeviceAllProperties("8657210202182310",val)
	
	fmt.Printf("\nredis content: %+v\n",val)
}

func TestGetDeviceProperty(t *testing.T) {
	
	
	val:=GetDeviceProperty("8657210202182310","status")
	
	fmt.Printf("\nredis content: %+v\n",val)
}

