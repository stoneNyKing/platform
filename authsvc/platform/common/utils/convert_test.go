
package utils

import (
	"testing"
)


func TestConvert2Int64(t *testing.T) {
	path := "2081"
	ip := Convert2Int64(path)
	t.Logf("int64: %d", ip)

	path1 := 2082
	ip1 := Convert2Int64(path1)
	t.Logf("int64: %d", ip1)
}

func TestConvert2Int(t *testing.T) {
	path := "1081"
	ip := Convert2Int(path)
	t.Logf("int: %d", ip)

	path1 := 1082
	ip1 := Convert2Int(path1)
	t.Logf("int: %d", ip1)
}

func TestConvert2Float32(t *testing.T) {
	path := "1081.001"
	ip := Convert2Float32(path)
	t.Logf("float32: %f", ip)

	path1 := 1082.913
	ip1 := Convert2Float32(path1)
	t.Logf("float32: %f", ip1)

}

func TestConvert2Float64(t *testing.T) {
	path := "641081.001"
	ip := Convert2Float64(path)
	t.Logf("float64: %f", ip)

	path1 := 641082.913
	ip1 := Convert2Float64(path1)
	t.Logf("float64: %f", ip1)
}


func TestStructToMap(t *testing.T) {

	type S1 struct {
		Siteid int64 	`form:"siteid" json:"siteid"`
		Created string `form:"created" json:"created"`
	}

	s1 := S1{Siteid:1,Created:"2009-01-01"}

	ip1,_ := StructToMap(s1)
	t.Logf("map: %+v", ip1)
}

