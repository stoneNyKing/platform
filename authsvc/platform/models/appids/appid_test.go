package appids

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	NewEngine("mysql","localhost", 3306, "root", "123456", "objects","objects")
	Truncate()
}

func Test_appid(t *testing.T) {
	appid := &Appid{
		SiteId: 1,
		Appid:  7,
		Key:    "appkey",
		ReMark: "remak",
		Status: 0,
	}
	err := appid.Insert()
	assert.Nil(t, err)

	appid1 := &Appid{Appid:appid.Appid}
	has, err := appid1.Get()
	assert.Nil(t, err)
	assert.True(t, has)
}

func Test_update(t *testing.T) {
	appid := &Appid{
		SiteId: 2,
		Appid:  12,
		Key:    "appkey2",
		ReMark: "remak",
		Status: 0,
	}
	err := appid.Insert()
	assert.Nil(t, err)

	appid1 := &Appid{Appid: appid.Appid}
	has, err := appid1.Get()
	assert.Nil(t, err)
	assert.True(t, has)

	appid1.Appid = 3
	appid1.Status = 1
	// appid1.TimeEnd =
	err = appid1.Update()
	assert.Nil(t, err)
}

func Test_list(t *testing.T) {
	var list, err = Search()
	assert.Nil(t, err)
	assert.Len(t, list, 2)
}
