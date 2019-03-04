package models

type (
	Resp struct {
		Ret int    `json:"Ret"`
		Msg string `json:"Msg"`
	}

	IdResp struct {
		Ret int
		Id  int64
	}

	IdSiteIdResp struct {
		Ret            int
		Id             int64
		SiteId         int64
		OrganizationId int64
	}
)
