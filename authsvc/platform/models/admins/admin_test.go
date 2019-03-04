package admins

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	NewEngine("mysql","localhost", 3306, "root", "123456", "objects","objects")
	Truncate()
}

func Test_Admins_CreateResource_Root(t *testing.T) {
	r, err := InitResource()

	assert.Nil(t, err, "insert r error")
	assert.Equal(t, r.Id, 1)
	assert.Equal(t, r.ParentId, 0)
}

func Test_Admins_GetResourceById(t *testing.T) {
	r, err := GetResourceById(1)
	assert.Nil(t, err, "insert r error")
	assert.Equal(t, r.Id, 1)
	assert.Equal(t, r.ParentId, 0)
}

func Test_Admins_CreateResource(t *testing.T) {
	r, err := CreateResource(&Resource{
		ParentId:    1,
		Name:        "fun1",
		Iconid:      1,
		Level:       0,
		Type:        0,
		Url:         "url",
		Description: "test",
		// EffectiveTime: time..Format("2014-09-01"),
		// ExpireTime:    "2014-12-01",
	})
	assert.Nil(t, err, "insert r error")
	assert.Equal(t, r.Id, 2)
	assert.Equal(t, r.ParentId, 1)
}

func Test_Admins_UpdateResource(t *testing.T) {
	r, err := UpdateResource(&Resource{
		Id:          2,
		ParentId:    1,
		Name:        "fun11",
		Iconid:      1,
		Level:       0,
		Type:        0,
		Url:         "url",
		Description: "test",
		// EffectiveTime: time..Format("2014-09-01"),
		// ExpireTime:    "2014-12-01",
	})
	assert.Nil(t, err, "change error")
	assert.True(t, r, "change error")
}

func Test_ListResourceById(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateResource(&Resource{
			ParentId:    1,
			Name:        "test_" + strconv.Itoa(i),
			Iconid:      i,
			Level:       0,
			Type:        0,
			Url:         "url",
			Description: "test",
			// EffectiveTime: time..Format("2014-09-01"),
			// ExpireTime:    "2014-12-01",
		})
	}
	rs, err := ListResourceById(1)
	assert.Nil(t, err, "change error")
	assert.Len(t, rs, 11, "length error")

	for i := 2; i < 10; i++ {
		CreateResource(&Resource{
			ParentId:    int64(i),
			Name:        "sub_" + strconv.Itoa(i),
			Iconid:      i,
			Level:       0,
			Type:        0,
			Url:         "url",
			Description: "test",
			// EffectiveTime: time..Format("2014-09-01"),
			// ExpireTime:    "2014-12-01",
		})
	}

	rs, err = ListResourceById(1)
	assert.Nil(t, err, "change error")
	assert.Len(t, rs, 11, "length error")
}

func Test_DeleteResourceById(t *testing.T) {
	b, err := DeleteResourceById(1)
	assert.Equal(t, err, ErrResourceIsParent, "delete error")
	assert.NotEqual(t, b, "false", "delete error2")

	b, err = DeleteResourceById(11)
	assert.Nil(t, err, "change error")
	assert.True(t, b, "delete error2")
}

/// admin
func Test_Admin_IsAdminNameExist(t *testing.T) {
	b, err := IsAdminNameExist(1, "testuser")
	assert.Nil(t, err, "error")
	assert.False(t, b, "user empty")
}

func Test_CreateAdmin(t *testing.T) {
	u, err := CreateAdmin(&Admin{
		SiteId:      1,
		RoleId:      0,
		Name:        "testuser",
		JobNumber:   "9527",
		Passwd:      "123456",
		Email:       "name@email.com",
		Phone:       "12345678901",
		Description: "string",
		// EffectiveTime: time..Format("2014-09-01"),
		// ExpireTime:    "2014-12-01",
	})
	assert.Nil(t, err, "insert r error")
	assert.Equal(t, u.Id, 1)
	assert.Equal(t, u.SiteId, 1)
	assert.Equal(t, u.RoleId, 0)
}

func Test_Admin_IsAdminNameExist2(t *testing.T) {
	b, err := IsAdminNameExist(1, "testuser")
	assert.Nil(t, err, "error")
	assert.True(t, b, "user empty")
}

func Test_GetAdminById(t *testing.T) {
	u, err := GetAdminById(1)
	assert.Nil(t, err, "error")
	assert.Equal(t, u.Id, 1)
	assert.Equal(t, u.SiteId, 1)
	assert.Equal(t, u.RoleId, 0)
}

func Test_GetAdminByName(t *testing.T) {
	u, err := GetAdminByName(1, "testuser")
	assert.Nil(t, err, "error")
	assert.Equal(t, u.Id, 1)
	assert.Equal(t, u.SiteId, 1)
	assert.Equal(t, u.RoleId, 0)
}

func Test_LoginAdmin(t *testing.T) {
	u, err := LoginAdmin(1, "testuser", "123456","")
	assert.Nil(t, err, "error")
	assert.Equal(t, u.Id, 1)
	assert.Equal(t, u.SiteId, 1)
	assert.Equal(t, u.RoleId, 0)

	u, err = LoginAdmin(1, "testuser1", "123456","")
	assert.Equal(t, err, ErrAdminNotExist, "error")
	assert.Nil(t, u, "empty")
}

func Test_ResetPasswd(t *testing.T) {
	b, err := ResetPasswd(2, "123456", "654321")
	assert.Equal(t, err, ErrAdminNotExist, "error")
	assert.False(t, b, "root error")

	b, err = ResetPasswd(1, "123456", "654321")
	assert.Nil(t, err, "error")
	assert.True(t, b, "root error")

	b, err = ResetPasswd(1, "123456", "654321")
	assert.Equal(t, err, ErrWarnPasswd, "error")
	assert.False(t, b, "root error")
}

func Test_UpdateAdmin(t *testing.T) {
	b, err := UpdateAdmin(&Admin{
		Id:          1,
		SiteId:      2,
		RoleId:      1,
		Name:        "usertest",
		JobNumber:   "7259",
		Passwd:      "123456",
		Email:       "name1@email.com",
		Phone:       "12345678902",
		Description: "strings",
		// EffectiveTime: time..Format("2014-09-01"),
		// ExpireTime:    "2014-12-01",
	})
	assert.Nil(t, err, "error")
	assert.True(t, b, "root error")

	u, err := GetAdminById(1)
	assert.Nil(t, err, "error")
	assert.Equal(t, u.Id, 1)
	assert.Equal(t, u.SiteId, 1)
	assert.Equal(t, u.RoleId, 1)
	assert.Equal(t, u.Name, "testuser")
	assert.Equal(t, u.JobNumber, "7259")
	assert.Equal(t, u.Passwd, "123456")
	assert.Equal(t, u.Email, "name1@email.com")
	assert.Equal(t, u.Phone, "12345678902")
	assert.Equal(t, u.Description, "strings")
}

func Test_ListAdmin(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateAdmin(&Admin{
			SiteId:      1,
			RoleId:      0,
			Type:        0,
			Name:        "testuser" + strconv.Itoa(i),
			JobNumber:   "9527" + strconv.Itoa(i),
			Passwd:      "123456" + strconv.Itoa(i),
			Email:       "name" + strconv.Itoa(i) + "@email.com",
			Phone:       "1234567890" + strconv.Itoa(i),
			Description: "string",
			// EffectiveTime: time..Format("2014-09-01"),
			// ExpireTime:    "2014-12-01",
		})
	}
	rs, err := ListAdmin(1, 0, 0)
	assert.Nil(t, err, "change error")
	assert.Len(t, rs, 10, "length error")
}

//role
func Test_CreateRole(t *testing.T) {
	r, err := CreateRole(&Role{
		SiteId:      1,
		Name:        "admins",
		Iconid:      1,
		Description: "string",
		// EffectiveTime: time..Format("2014-09-01"),
		// ExpireTime:    "2014-12-01",
	})
	assert.Nil(t, err, "insert r error")
	assert.Equal(t, r.Id, 1)
	assert.Equal(t, r.SiteId, 1)
}

func Test_GetRoleById(t *testing.T) {
	r, err := GetRoleById(1)
	assert.Nil(t, err, "insert r error")
	assert.Equal(t, r.Id, 1)
	assert.Equal(t, r.SiteId, 1)
}

func Test_UpdateRole(t *testing.T) {
	b, err := UpdateRole(&Role{
		Id:          1,
		SiteId:      2,
		Name:        "admins2",
		Iconid:      2,
		Description: "string",
		// EffectiveTime: time..Format("2014-09-01"),
		// ExpireTime:    "2014-12-01",
	})

	assert.Nil(t, err, "error")
	assert.True(t, b, "root error")

	u, err := GetRoleById(1)
	assert.Nil(t, err, "error")
	assert.Equal(t, u.Id, 1)
	assert.Equal(t, u.SiteId, 1)
	assert.Equal(t, u.Name, "admins2")
}

func Test_ListRole(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRole(&Role{
			SiteId:      1,
			Name:        "admins" + strconv.Itoa(i),
			Iconid:      2,
			Description: "string",
			// EffectiveTime: time..Format("2014-09-01"),
			// ExpireTime:    "2014-12-01",
		})
	}
	rs, err := ListRole(1)
	assert.Nil(t, err, "change error")
	assert.Len(t, rs, 10, "length error")
}

//RoleResource
func Test_IsRoleResourceExist(t *testing.T) {
	b, err := IsRoleResourceExist(1, 1)
	assert.Nil(t, err, "error")
	assert.False(t, b, "root error")
}

func Test_CreateRoleResource(t *testing.T) {
	res1, err := CreateResource(&Resource{
		ParentId:    1,
		Name:        "bind_1",
		Iconid:      1,
		Level:       0,
		Type:        1,
		Url:         "url1",
		Description: "test",
		// EffectiveTime: time..Format("2014-09-01"),
		// ExpireTime:    "2014-12-01",
	})
	assert.Nil(t, err, "insert r error")

	res2, err := CreateResource(&Resource{
		ParentId:    res1.Id,
		Name:        "bind_2",
		Iconid:      2,
		Level:       0,
		Type:        2,
		Url:         "url1",
		Description: "test",
		// EffectiveTime: time..Format("2014-09-01"),
		// ExpireTime:    "2014-12-01",
	})
	assert.Nil(t, err, "insert r error")
	assert.Equal(t, res1.Id+1, res2.Id)

	role, err := GetRoleById(1)
	resource, err := GetResourceById(res2.Id)
	r, err := CreateRoleResource(&RoleResource{
		Role:     *role,
		Resource: *resource,
	})
	assert.Nil(t, err, "insert r error")
	assert.Equal(t, r.Role.Id, 1)
	assert.Equal(t, r.Resource.Id, res2.Id)
}

func Test_IsRoleResourceExist2(t *testing.T) {
	b, err := IsRoleResourceExist(1, 22)
	assert.Nil(t, err, "error")
	assert.True(t, b, "root error")
}

func Test_ListRoleResourceByRoleId(t *testing.T) {
	role, _ := GetRoleById(int64(1))
	for i := 1; i < 10; i++ {
		resource, _ := GetResourceById(int64(i))
		CreateRoleResource(&RoleResource{
			Role:     *role,
			Resource: *resource,
		})

	}
	rs, err := ListRoleResourceByRoleId(1)
	assert.Nil(t, err, "change error")
	assert.Len(t, rs, 1, "length error")
}

func Test_DeleteRoleResource(t *testing.T) {

	role, err := GetRoleById(1)
	b, err := DeleteRoleResource(&RoleResource{
		Role: *role,
	})

	assert.Nil(t, err, "change error")
	assert.True(t, b, "length error")
}
