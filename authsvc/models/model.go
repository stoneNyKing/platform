package models

import "time"

type ApiPackage struct{
	PackageId           int64			`orm:"pk"`
	Name                string
	Price 			    int
	ChargeModel        	int
	SubSysId          	int64
	Status              int16			`orm:"default(1)"`
	CreateTime         	time.Time		`orm:"auto_now_add;type(datetime)"`
	Remark              string			`orm:"type(text)"`
	PackageCode         string
}

type ApiPackageService struct{
	PkgServiceId        int64			`orm:"pk;auto"`
	ServiceId           int
	PackageId 			int
	TotalCounts        	int
	DailyCounts        	int
	Status              int16			`orm:"default(1)"`
	CreateTime         	time.Time		`orm:"auto_now_add;type(datetime)"`
}

type ApiService struct{
	ServiceId        	int64			`orm:"pk"`
	SvcCode           	string
	SvcId 				int
	Route        		string
	WebUrl          	string
	ApiVer          	string
	Status				int16			`orm:"default(1)"`
	Remark              string			`orm:"type(text)"`
}

type SecSiteInfo struct{
	LicenseId           int64			`orm:"pk;auto"`
	Siteid              int64
	ApiKey              string
	OrgCode 			string
	Status              int16			`orm:"default(1)"`
	CreateTime         	time.Time		`orm:"auto_now_add;type(datetime)"`
	Remark              string			`orm:"type(text)"`
	License             string			`orm:"type(text)"`
	Userid				int64
	OrganizationId		int64
}
