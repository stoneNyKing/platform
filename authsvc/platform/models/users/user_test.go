package users

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"union/lib/helper"
)

func init() {
	NewEngine("mysql","localhost", 3306, "root", "123456", "objects","objects")
	Truncate()
}

func Test_Users_Phone(t *testing.T) {
	assert.True(t, PhoneRegular.MatchString("13913913913"))
	assert.True(t, PhoneRegular.MatchString("18668170302"), "error")
	assert.False(t, PhoneRegular.MatchString("28668170302"))
}

func Test_Users_RegisterUserByName(t *testing.T) {
	user, err := RegisterUserByName(&User{SiteId: 1, Name: "name", Passwd: "passwd"})
	assert.Nil(t, err)
	assert.Equal(t, user.Id, 1)
	assert.Equal(t, user.Passwd, "passwd")

	user, err = GetUserByName(1, "name")
	assert.Nil(t, err)
	assert.Equal(t, user.Id, 1)
	assert.Equal(t, user.Passwd, "passwd")

}

func Test_Users_RegisterUserByPhone(t *testing.T) {
	user, err := RegisterUserByName(&User{SiteId: 1, Phone: "13913913913", Passwd: "dwssap"})
	assert.Nil(t, err)
	assert.Equal(t, user.Id, 2)
	assert.Equal(t, user.Passwd, "dwssap")

	user, err = GetUserByPhone(1, "13913913913")
	assert.Nil(t, err)
	assert.Equal(t, user.Id, 2)
	assert.Equal(t, user.Passwd, "dwssap")
}

func Test_Users_LoginUserPlain(t *testing.T) {
	//ok
	passwd := helper.Md5("passwd" + "salt")
	user, err := LoginUserPlain(1, "name", passwd, "salt")
	assert.Nil(t, err)
	assert.Equal(t, user.Id, 1)
	assert.Equal(t, user.Passwd, "passwd")

	//ok LoginUserByid
	user, err = LoginUserByid(1, 1, passwd, "salt")
	assert.Nil(t, err)
	assert.Equal(t, user.Id, 1)
	assert.Equal(t, user.Passwd, "passwd")

	//error
	passwd = helper.Md5("passwd" + "salt2")
	user, err = LoginUserPlain(1, "name", passwd, "salt")
	assert.Error(t, err)
	assert.Nil(t, user)

	//ok old
	user, err = LoginUserPlain(1, "name", "passwd", "")
	assert.Nil(t, err)
	assert.Equal(t, user.Id, 1)
	assert.Equal(t, user.Passwd, "passwd")

	//error old
	user, err = LoginUserPlain(1, "name", "passwd1x", "")
	assert.Error(t, err)
	assert.Nil(t, user)
}

func Test_Users_LoginUserPlain2(t *testing.T) {
	//ok
	passwd := helper.Md5("dwssap" + "tlas")
	user, err := LoginUserPlain(1, "13913913913", passwd, "tlas")
	assert.Nil(t, err)
	assert.Equal(t, user.Id, 2)
	assert.Equal(t, user.Passwd, "dwssap")

	//error
	passwd = helper.Md5("dwssap21" + "tlas")
	user, err = LoginUserPlain(1, "13913913913", passwd, "tlas")
	assert.Error(t, err)
	assert.Nil(t, user)

	//ok old
	user, err = LoginUserPlain(1, "13913913913", "dwssap", "")
	assert.Nil(t, err)
	assert.Equal(t, user.Id, 2)
	assert.Equal(t, user.Passwd, "dwssap")

	//error
	user, err = LoginUserPlain(1, "13913913913", "dwssap2x", "")
	assert.Error(t, err)
	assert.Nil(t, user)
}

func Test_Add_log(t *testing.T) {
	log, err := AddLog(&UserLog{SiteId: 1, UserId: 2, Level: 1})
	assert.Nil(t, err, "insert log error")
	assert.Equal(t, log.Id, 1)
	assert.Equal(t, log.UserId, 2)
}

func Test_List_Log(t *testing.T) {
	for i := 0; i < 50; i++ {
		AddLog(&UserLog{SiteId: 1, UserId: 2, Level: 1})
	}

	userlogs, err := LogList(1, 2, 20, 1, 1)
	assert.Nil(t, err, "insert log error")
	assert.Len(t, userlogs, 20)

	userlogs, err = LogList(1, 2, 20, 2, 1)
	assert.Nil(t, err, "insert log error")
	assert.Len(t, userlogs, 11)
}

func Test_Add_Feedback(t *testing.T) {
	r, err := SubmitFeedback(&Feedback{SiteId: 1, Appid: 2, Userid: 3, Content: "test"})
	assert.Nil(t, err, "insert r error")
	assert.Equal(t, r.Id, 1)
	assert.Equal(t, r.Appid, 2)
}

func Test_List_Feedback(t *testing.T) {
	for i := 0; i < 50; i++ {
		SubmitFeedback(&Feedback{SiteId: 1, Appid: 2, Userid: 3, Content: "test"})
	}

	list, err := FeedbackList(1, 20, 1)
	assert.Nil(t, err, "insert log error")
	assert.Len(t, list, 20)

	list, err = FeedbackList(1, 20, 2)
	assert.Nil(t, err, "insert log error")
	assert.Len(t, list, 11)
}
