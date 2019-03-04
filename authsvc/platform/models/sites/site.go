package sites

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	ErrSiteAlreadyExist   = errors.New("Site already exist")
	ErrSiteNotExist       = errors.New("Site does not exist")
	ErrParentSiteNotExist = errors.New("Parent Site does not exist")
)

type Site struct {
	Id       int64
	ParentId int64  `xorm:"not null"`
	Name     string `xorm:"not null"`
	Level    int    `xorm:"not null"`
	Key      string `xorm:"not null"`
	Ip       string
	Port     int
	Created  time.Time              `xorm:"created"`
	Updated  time.Time              `xorm:"updated"`
	Json     map[string]interface{} `xorm:text`
	Status 	int
	RootResId		int64
	RootAreaId		int64
	RootAppResid	int64
	
	ProxyIp         string
	ProxyPort       int

	//2018-12-03
	OrgCode 		string
	Address 		string
	MembershipLevel	string
	Nature			string
	OrgLevel		string
	LegalPerson		string
	Phone 			string
	OrgType 		string
	OrgLn 			string
	Fax 			string

}

func IsSiteIdExist(id int64) (bool, error) {
	if id == 0 {
		return false, nil
	}
	return orm.Get(&Site{Id: id})
}

func GetSiteById(id int64) (*Site, error) {
	site := new(Site)
	has, err := orm.Id(id).Get(site)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrSiteNotExist
	}
	return site, nil
}

func GetSiteIdByHost(host string, port int) (int64, error) {
	site := &Site{Ip: strings.ToLower(host), Port: port}
	has, err := orm.Get(site)
	if err != nil {
		return 0, err
	}
	if !has {

		site = &Site{ProxyIp: strings.ToLower(host), ProxyPort: port}
		has, err = orm.Get(site)
		if err != nil {
			return 0, err
		}

		if !has {
			return 0, ErrSiteNotExist
		}
	}
	return site.Id, nil
}

func GetChildSites(id int64, keys []string) ([]interface{}, error) {
	psite, err := GetSiteById(id)
	if err != nil {
		return nil, err
	} else if psite.Id == 0 {
		return nil, ErrParentSiteNotExist
	}

	var d = make([]interface{}, 0)
	rows, err := orm.Where("`level` = ? and `key` like ?", psite.Level+1, psite.Key+strconv.FormatInt(psite.Id, 10)+"-%").Asc("key").Rows(new(Site))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	site := new(Site)
	for rows.Next() {
		err = rows.Scan(site)
		if err != nil {
			return nil, err
		}
		d1, err := _GetProfile(site, keys)
		if err != nil {
			return nil, err
		}
		d = append(d, d1)
	}
	return d, nil
}

func GetPosteritySites(id int64, keys []string) ([]interface{}, error) {
	psite, err := GetSiteById(id)
	if err != nil {
		return nil, err
	} else if psite.Id == 0 {
		return nil, ErrParentSiteNotExist
	}
	var d = make([]interface{}, 0)
	rows, err := orm.Where("`key` like ?", psite.Key+strconv.FormatInt(psite.Id, 10)+"-%").Rows(new(Site))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	site := new(Site)
	for rows.Next() {
		err = rows.Scan(site)
		if err != nil {
			return nil, err
		}
		d1, err := _GetProfile(site, keys)
		if err != nil {
			return nil, err
		}
		d = append(d, d1)
	}
	return d, nil
}

func AddSite(site *Site) (*Site, error) {
	psite, err := GetSiteById(site.ParentId)
	if err != nil {
		return nil, err
	} else if psite.Id == 0 {
		return nil, ErrParentSiteNotExist
	}
	site.Level = psite.Level + 1
	site.Key = psite.Key + strconv.FormatInt(psite.Id, 10) + "-"

	if _, err = orm.Insert(site); err != nil {
		return nil, err
	}
	return site, err
}

func _GetProfile(site *Site, keys []string) (map[string]interface{}, error) {
	var d = make(map[string]interface{})
	if site.Json == nil {
		site.Json = make(map[string]interface{})
	}

	r := reflect.ValueOf(site)
	for _, k := range keys {
		f := reflect.Indirect(r).FieldByName(k)
		if f.IsValid() == true {
			switch f.Kind() {
			case reflect.String:
				d[k] = f.String()
			case reflect.Int:
				d[k] = f.Int()
			case reflect.Int64:
				d[k] = f.Int()
			case reflect.Map:
				d[k] = f.Interface()
			default:
				if k == "Create" || k == "Updated" {
					// d[k] = f.
				}
			}
		} else {
			d[k] = site.Json[k]
		}
	}
	return d, nil
}

func GetProfile(id int64, keys []string) (map[string]interface{}, error) {
	site, err := GetSiteById(id)
	if err != nil {
		return nil, err
	}
	return _GetProfile(site, keys)
}

func SetProfile(id int64, d map[string]interface{}) (bool, error) {
	site, err := GetSiteById(id)
	if err != nil {
		return false, err
	}

	if site.Json == nil {
		site.Json = make(map[string]interface{})
	}

	bjson := false
	r := reflect.ValueOf(site)
	for k, v := range d {
		f := reflect.Indirect(r).FieldByName(k)
		if f.IsValid() == true {
			if k == "Name" {
				site.Name = v.(string)
			} else if k == "Ip" {
				site.Ip = v.(string)
			} else if k == "Port" {
				site.Port = int(v.(float64))
			} else if k == "Json" {
				site.Json = v.(map[string]interface{})
				bjson = true
			}
			// mutable := reflect.ValueOf(site).Elem()
			// mutable.FieldByName(k).Set(v)
		} else {
			if bjson != true {
				site.Json[k] = v
			}
		}
	}

	affected, err := orm.Id(id).Update(site)
	return affected == 1, err
}
