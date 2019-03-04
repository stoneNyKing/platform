// util_test
package utils

import (
	"fmt"
	"testing"
	"math/rand"
)

func TestCaculateCRC16(t *testing.T) {

	b := []byte{0x40, 0x07, 26, 0, 0, 0, 0, 1, 0, 33, 0, 0, 0, 26, 0, 2, 0, 0, 0, 49, 7, 0, 0, 49, 8, 0, 0, 218, 137}

	sn := CaculateCRC16(b)

	fmt.Printf("The result is: %d\n", sn)
	
	fmt.Printf("Hex = %02X\n",sn)
	
	i:=0
	
	for {
		regno := rand.Int63n(100000000)
		registno := fmt.Sprintf("%08d", regno)
		
		fmt.Printf("register no: %s\n",registno)
		
		i+=1
		
		if i>10 {
			break
		}
	}
}

func TestCreateCRC32(t *testing.T) {
	str := "this a test for crc16"
	b := []byte(str)

	sn := CreateCRC32(b, uint32(len(b)))

	fmt.Printf("The result is: %d\n", sn)
}


func TestCreateContextID(t *testing.T){
	
	id := CreateContextID()
	
	fmt.Printf("Contextid = %d\n",id)
}

func TestCheckMD5Passwd(t *testing.T) {
	
	b,pwd := CheckMD5Passwd("c2d4e2562a65eee1bfe35b6e8cd86114b","InterMa140")
	
	fmt.Printf("b=%v,\tpwd=%s\n",b,pwd)
}

func TestBytesToFloat32(t *testing.T) {

	a :=[]byte{0x00, 0xc0, 0xa0 ,0x44}
	
	f:=BytesToFloat32(a)
	
	fmt.Printf("float=%.2f\n",f)

	b:= []byte{0xDA,0x0F,0x49,0x40}
	ff:=BytesToFloat32(b)
	fmt.Printf("float pi=%.5f\n",ff)
}

func TestGetTimeFormat(t *testing.T) {
	year := GetTimeFormat("2006-01-02 15:04:05")
	fmt.Printf("local time = %s\n",year)
}
func TestGetUTCTimeFormat(t *testing.T) {
	year := GetUTCTimeFormat("2006-01-02 15:04:05")
	fmt.Printf("utc time = %s\n",year)
}

func TestCompareTime(t *testing.T) {
	u := CompareTime("20160305080000","20060102150405")
	
	fmt.Printf("compare time: %d\n",u);
}


func TestGetStartEndTimeOfWeek(t *testing.T) {
	s,e := GetStartEndTimeOfWeek(0)

	fmt.Printf("start time=: %s,end time=%s \n",s,e);
}

func TestGetStartEndTimeOfWeekday(t *testing.T) {
	s,e := GetStartEndTimeOfWeekday(0,2)

	fmt.Printf("start time=: %s,end time=%s \n",s,e);
}

func TestGetStartEndTimeOfMonthday(t *testing.T) {
	s,e := GetStartEndTimeOfMonthday(0,22)

	fmt.Printf("start time=: %s,end time=%s \n",s,e);
}
