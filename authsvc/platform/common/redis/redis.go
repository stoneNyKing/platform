package redis

import (
	"errors"
	"fmt"
	"strconv"
	hoisie "github.com/libra9z/hoisie-redis"
)
 
var client hoisie.Client


func Init(addr string, db int, pools int) {
	client.Addr = addr
	client.Db = db
	client.MaxPoolSize = pools
}

func GetIncr(key string) string {
	result, err := client.Incr(key)
	if checkError(err, "Get Incr Error") {
		return ""
	}
	return strconv.FormatInt(result, 10)
}

func GetIncrID(key string) int64 {
	result, err := client.Incr(key)
	if checkError(err, "Get Incr Error") {
		return 0
	}
	return result
}

func GetValue(key string) string {
	result, err := client.Get(key)
	
	if err != nil {
		return ""
	}
	
	return string(result)
}

func Del(key string) (bool,error) {
	result, err := client.Del(key)
	if checkError(err, "delete Key "+key+" Error") {
		return false,err
	}
	return result,nil
}

func SetValue(key string, value string) {
	client.Set(key, []byte(value))
}

func PushHealthData(data []byte, dq string) {
	err := client.Lpush(dq, data)
	checkError(err, "Push Health Data Error")

}
func FetchHealthData(dq string) (string, error) {
	var queues []string
	queues = append(queues, dq)
	key, value, err := client.Brpop(queues, 10)
	if checkError(err, "Fetch Health Data Error") {
		return "", err
	}
	if key == nil {
		return "", errors.New("get empty data")
	}
	return string(value), nil
}

func checkError(err error, info string) bool {
	if err != nil {
		//syssecure.AddAlarm(2, 5, "202", error.Error()+info)
		fmt.Println("ERROR:", err.Error(), info)
		return true
		//panic("ERROR: " + info + " " + error.Error()) // terminate program
	}
	return false
}

func Publish(key string, buf []byte) int64 {
	//client.Publish2(key,buf)
	data, err := client.Publish2(key, buf)
	if err != nil {
		return 0
	}

	ret := data.(int64)

	return ret
}

func Exists(key string) bool {
	b, _ := client.Exists(key)
	return b
}



func Expire(key string, t int64) bool {
	b, err := client.Expire(key, t)
	if err != nil {
		return false
	}

	return b
}

func Setnx(key string, value string) bool {
	b, _ := client.Setnx(key, []byte(value))

	return b
}

func Setex(key string, value string) bool {
	err := client.Setex(key, 604800, []byte(value))

	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func Hget(key,field string) string {
	
	res ,err := client.Hget(key,field)
	
	if err != nil {
		return ""
	}
	
	return string(res)
}

func Hset(key,field,value string) (bool,error){
	
	ok,err := client.Hset(key,field,[]byte(value))
	
	return ok,err
}

func Hsetnx(key,field,value string) (bool,error){
	
	ok,err := client.Hsetnx(key,field,[]byte(value))
	
	return ok,err
}

func Lpush(queue string,data string) error {
	return client.Lpush(queue,[]byte(data))
}

func Brpop(keys []string, timeoutSecs uint) (*string, []byte, error) {
	return client.Brpop(keys,timeoutSecs)
}
func Lrange(key string, start int, end int) ([][]byte, error) {
    return client.Lrange(key,start,end)
}

func Subscribe(subscribe <-chan string, unsubscribe <-chan string, psubscribe <-chan string, punsubscribe <-chan string, messages chan<- hoisie.Message) error {
   return client.Subscribe(subscribe,unsubscribe,psubscribe,punsubscribe,messages)
}    

func Keys(pattern string) ([]string, error) {
	return client.Keys(pattern)
}

func FetchMessage(dq string, timeout uint) (string, error) {
	var queues []string
	queues = append(queues, dq)
	key, value, err := client.Brpop(queues, timeout)
	if checkError(err, "Fetch Message Data Error") {
		return "", err
	}
	if key == nil {
		return "", errors.New("get empty data")
	}
	return string(value), nil
}

