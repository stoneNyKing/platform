package sites

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	NewEngine("mysql","localhost", 3306, "root", "123456", "objects","objects")
	Truncate()
}

func Test_Sites_Root(t *testing.T) {
	site, err := GetSiteById(1)
	assert.Nil(t, err, "add error")
	assert.Equal(t, site.Id, 1)
	assert.Equal(t, site.ParentId, 0)
	assert.Equal(t, site.Level, 1)
	assert.Equal(t, site.Key, "")
}

func Test_Sites_Add(t *testing.T) {
	site1, err := AddSite(&Site{Name: "第一级", ParentId: 1})
	assert.Nil(t, err, "add error")
	assert.Equal(t, site1.Id, 2)
	assert.Equal(t, site1.ParentId, 1)
	assert.Equal(t, site1.Level, 2)
	assert.Equal(t, site1.Key, "1-")

	site11, err := AddSite(&Site{Name: "第一级第一级", ParentId: site1.Id})
	assert.Nil(t, err, "add error")
	assert.Equal(t, site11.Id, 3)
	assert.Equal(t, site11.ParentId, site1.Id)
	assert.Equal(t, site11.Level, 3)
	assert.Equal(t, site11.Key, "1-2-")

	site2, err := AddSite(&Site{Name: "第二级", ParentId: 1})
	assert.Nil(t, err, "add error")
	assert.Equal(t, site2.Id, 4)
	assert.Equal(t, site2.ParentId, 1)
	assert.Equal(t, site2.Level, 2)
	assert.Equal(t, site2.Key, "1-")

	site21, err := AddSite(&Site{Name: "第二级第一级", ParentId: site2.Id})
	assert.Nil(t, err, "add error")
	assert.Equal(t, site21.Id, 5)
	assert.Equal(t, site21.ParentId, site2.Id)
	assert.Equal(t, site21.Level, 3)
	assert.Equal(t, site21.Key, "1-4-")

	site22, err := AddSite(&Site{Name: "第二级第二级", ParentId: site2.Id})
	assert.Nil(t, err, "add error")
	assert.Equal(t, site22.Id, 6)
	assert.Equal(t, site22.ParentId, site2.Id)
	assert.Equal(t, site22.Level, 3)
	assert.Equal(t, site22.Key, "1-4-")

	site221, err := AddSite(&Site{Name: "第二级第二级", ParentId: site22.Id})
	assert.Nil(t, err, "add error")
	assert.Equal(t, site221.Id, 7)
	assert.Equal(t, site221.ParentId, site22.Id)
	assert.Equal(t, site221.Level, 4)
	assert.Equal(t, site221.Key, "1-4-6-")
}

func Test_Get_Profile(t *testing.T) {
	d, err := GetProfile(1, []string{"Id", "Name", "Level", "Key", "n", "hello", "Hello"})
	assert.Nil(t, err, "get profile error")
	assert.Equal(t, d["Id"], 1)
	assert.Equal(t, d["Name"], "Union")
	assert.Equal(t, d["Level"], 1)
	assert.Equal(t, d["Key"], "")
	assert.Nil(t, d["n"], "")
	assert.Nil(t, d["hello"], "")
	assert.Nil(t, d["Hello"], "")

	d, err = GetProfile(2, []string{"Id", "Name", "Level", "Key", "n", "hello", "Hello"})
	assert.Nil(t, err, "get profile error")
	assert.Equal(t, d["Id"], 2)
	assert.Equal(t, d["Name"], "第一级")
	assert.Equal(t, d["Level"], 2)
	assert.Equal(t, d["Key"], "1-")
	assert.Nil(t, d["n"], "")
	assert.Nil(t, d["hello"], "")
	assert.Nil(t, d["Hello"], "")
}

func Test_Set_Profile(t *testing.T) {
	kv := make(map[string]interface{})
	kv["id"] = 112
	kv["Id"] = 111
	kv["level"] = 113
	kv["Level"] = 113
	kv["key"] = "a"
	kv["Key"] = "b"
	kv["name"] = "c"
	kv["Name"] = "d"
	kv["n"] = 100
	kv["hello"] = "world"
	kv["Hello"] = "World"
	kv["Ip"] = "192.168.1.1"
	kv["Port"] = float64(81)

	b, err := SetProfile(1, kv)
	assert.Nil(t, err, "set profile error")
	assert.True(t, b, "SetProfile error")

	b, err = SetProfile(2, kv)
	assert.Nil(t, err, "set profile error")
	assert.True(t, b, "SetProfile error")
}

func Test_Get_Profile2(t *testing.T) {
	d, err := GetProfile(1, []string{"Id", "Name", "Level", "Key", "n", "hello", "Hello", "Json"})
	assert.Nil(t, err, "get profile error")
	assert.Equal(t, d["Id"], 1)
	assert.Equal(t, d["Name"], "d")
	assert.Equal(t, d["Level"], 1)
	assert.Equal(t, d["Key"], "")
	assert.Equal(t, d["n"], 100)
	assert.Equal(t, d["hello"], "world")
	assert.Equal(t, d["Hello"], "World")

	jsonvalue := d["Json"].(map[string]interface{})
	assert.Equal(t, jsonvalue["n"], 100)
	assert.Equal(t, jsonvalue["hello"], "world")
	assert.Equal(t, jsonvalue["Hello"], "World")

	d, err = GetProfile(2, []string{"Id", "Name", "Level", "Key", "n", "hello", "Hello"})
	assert.Nil(t, err, "get profile error")
	assert.Equal(t, d["Id"], 2)
	assert.Equal(t, d["Name"], "d")
	assert.Equal(t, d["Level"], 2)
	assert.Equal(t, d["Key"], "1-")
	assert.Equal(t, d["n"], 100)
	assert.Equal(t, d["hello"], "world")
	assert.Equal(t, d["Hello"], "World")
}

func Test_GetChildSites(t *testing.T) {
	d, err := GetChildSites(1, []string{"Id", "Name", "Level", "Key", "n", "hello", "Hello"})
	assert.Nil(t, err, "GetChildSites error")
	assert.Equal(t, len(d), 2)

	d, err = GetChildSites(2, []string{"Id", "Name", "Level", "Key", "n", "hello", "Hello"})
	assert.Nil(t, err, "GetChildSites error")
	assert.Equal(t, len(d), 1)

	d, err = GetChildSites(3, []string{"Id", "Name", "Level", "Key", "n", "hello", "Hello"})
	assert.Nil(t, err, "GetChildSites error")
	assert.Equal(t, len(d), 0)
}

func Test_GetPosteritySites(t *testing.T) {
	d, err := GetPosteritySites(1, []string{"Id", "Name", "Level", "Key", "n", "hello", "Hello"})
	assert.Nil(t, err, "GetPosteritySites error")
	assert.Equal(t, len(d), 6)

	d, err = GetPosteritySites(2, []string{"Id", "Name", "Level", "Key", "n", "hello", "Hello"})
	assert.Nil(t, err, "GetPosteritySites error")
	assert.Equal(t, len(d), 1)

	d, err = GetPosteritySites(3, []string{"Id", "Name", "Level", "Key", "n", "hello", "Hello"})
	assert.Nil(t, err, "GetPosteritySites error")
	assert.Equal(t, len(d), 0)
}

func Test_Set_Profile2(t *testing.T) {
	kv := make(map[string]interface{})
	kv["id"] = 112
	kv["Id"] = 111
	kv["level"] = 113
	kv["Level"] = 113
	kv["key"] = "a"
	kv["Key"] = "b"
	kv["name"] = "c"
	kv["Name"] = "d"
	kv["n"] = 100
	kv["hello"] = "world"
	kv["Hello"] = "World"
	kv["Ip"] = "192.168.1.1"
	kv["Port"] = float64(81)

	kv2 := make(map[string]interface{})
	kv2["name1"] = "cd"
	kv2["Name1"] = "dd"
	kv["Json"] = kv2

	b, err := SetProfile(1, kv)
	assert.Nil(t, err, "set profile error")
	assert.True(t, b, "SetProfile error")
}

func Test_Get_Profile3(t *testing.T) {
	d, err := GetProfile(1, []string{"Id", "Name", "Level", "Key", "n", "hello", "Hello", "name1", "Name1", "Json"})
	assert.Nil(t, err, "get profile error")
	assert.Equal(t, d["Id"], 1)
	assert.Equal(t, d["Name"], "d")
	assert.Equal(t, d["Level"], 1)
	assert.Equal(t, d["Key"], "")
	assert.Nil(t, d["n"])
	assert.Nil(t, d["hello"])
	assert.Nil(t, d["Hello"])
	assert.Equal(t, d["name1"], "cd")
	assert.Equal(t, d["Name1"], "dd")

	jsonvalue := d["Json"].(map[string]interface{})
	assert.Nil(t, jsonvalue["n"])
	assert.Equal(t, jsonvalue["name1"], "cd")
	assert.Equal(t, jsonvalue["Name1"], "dd")
}
