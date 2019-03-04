package models



type Admin	struct {
	Id			      	int64 `orm:"pk;auto"`
	SiteId      		int64
	RoleId      		int64

	Name      			string
	JobNumber		   	string
	Phone		    	string
	Email		    	string
	Passwd		    	string
	Description	    	string
	EffectiveTime    	string
	ExpireTime	    	string

	Created 			string
	Updated 			string
	Type 				int

	//2018-01-02
	State 				int

	ImageUrl 			string
}
type User struct {
	Id         int64 `orm:"pk;auto"`
	SiteId     int64
	Name       string
	Phone      string
	Email      string
	Idcard     string
	Type       int64
	IsActive   int
	Rfid       string
	Imei       string `json:"imei,omitempty"`
	Nickname   string
	Created    string
	Updated    string
	LastLogin  string
	LastLogout string
	Weixinid   string
	Json 	   string
	Passwd     string
	ImageUrl	string

	//2018-01-11增加
	RegisterType	int
	Sscard 			string

}


type Site struct {
	Id         int64 `orm:"pk;auto"`
	ParentId   int64
	Name       string
	Level      int64
	Key        string
	Created    string
	Updated    string
	Json       string
	Ip         string
	Port       int64
	Status     int
	RootResId  int64
	RootAreaId int64
	RootAppResid	int64

	ProxyIp         string
	ProxyPort       int64

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


type Role struct {
	Id				int64		`orm:"pk;auto"`
	SiteId			int64
	Name			string
	Iconid			int64
	Description		string
	EffectiveTime	string
	ExpireTime		string
	Created			string
	Updated			string
	Startpage		string
	WheelFlag 		int
	BlockFlag 		int
	OrgId 			int64
}


type SiteRes struct {
	ResId			int64		`orm:"pk;auto"`
	ParentId		int64
	ResourceId		int64
	Treeid			string
	Name			string
	StartTime		string
	EndTime			string
	Status			int
	Level			int
	Order			int
	Direction		int
}



type RoleRes struct {
	ResId			int64		`orm:"pk;auto"`
	ParentId		int64
	ResourceId		int64
	Treeid			string
	Name			string
	StartTime		string
	EndTime			string
	Status			int
	Level			int
	Order			int

	//操作权限
	PermSel 		int
	PermAdd 		int
	PermUpd 		int
	PermDel 		int
	PermCancel 		int
	PermAudit 		int

	//数据权限
	PermEval 		string
	PermDoc 		string

}


type SubSys struct {
	SubSysId 		int64  		`orm:"pk;auto"`
	ResId 			int64
	Appid 			int64
	SiteId 			int64
	Subtplid		int64
	Name 			string
	IconUrl 		string
	StartResourceId int64
	Created 		string
	Remark 			string
}

type RoleResource struct {
	RoleId 					int64
	ResId 					int64
	Appid 					int64
	SubSysId 				int64
	ResourceId 				int64
	StartResourceId 		int64
	RoleResourceId 			int64 		`orm:"pk;auto"`
	CreateTime 				string
}


type Resource struct {
	Id				int64		`orm:"pk;auto"`
	ParentId		int64
	Name			string
	Iconid			int64
	Level			int64
	Type			int64
	Url				string
	Description		string
	EffectiveTime	string
	ExpireTime		string
	Created			string
	Updated			string
	Proxy			string
	Icon			string

	//2018-01-20 add
	Appid 			int64
	PermDesc 		string
}

