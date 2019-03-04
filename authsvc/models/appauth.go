package models


import (
	"time"
	"errors"
)

func GetAuth(orgid,userid int64) (int64,interface{}, error) {

	// 获取所需要的license
	licenseRow, _, err := GetAppLicense(orgid,userid)

	if err != nil {
		return 1, nil,err
	}

	license := licenseRow.(map[string]interface{})["license"]
	param := license.(map[string] interface{})

	// 获取license的expire-time并判断过期时间是否过期
	if p, ok := param["produces"]; ok {
		products := p.([]interface{})
		for _, product := range products {
			productInfo := product.(map[string]interface{})
			if expireDate, ok := productInfo["expire-date"].(string); ok {
				t, _ := time.Parse("2006-01-02 15:04:05", expireDate + " 00:00:00")
				trueOrFalse := t.After(time.Now())
				if !trueOrFalse {
					return 1,nil, errors.New("the license has expired")
				}
			}

		}

		return 0, p,nil
	}

	return 1, nil,errors.New("can not find license")
}