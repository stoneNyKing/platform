package models

import ()

type AccountItem struct {
	Item    string
	Max     int64
	Min     int64
	Average int64 //ms为单位
	Times   int64
}
type AccountResult struct {
	Module    int
	AppID     int
	Begintime string
	Endtime   string
	Account   []AccountItem
}

type TraceMsg struct {
	Mtype   string //”1”:Device,”2”:一体机
	ID      string //用于区分跟踪什么内容的ID，比如mType为1时，ID为DeviceID
	AppID   string //目前都填写1
	Module  string
	Event   string
	Time    string
	Success string //”0”失败，”1”成功
	Input   string
	Output  string
	Notice  string
}
