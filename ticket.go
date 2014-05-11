package main

import "time"

type TicketQueryInfo struct {
	TrainDate          time.Time
	FromStations       []string
	FromStationsStr    string
	ToStations         []string
	ToStationsStr      string
	TicketType         string //车票类型：成人票，学生票
	SeatType           string //席别：硬卧，硬座
	Trians             []string
	TriansStr          string
	Passengers         []string
	PassengerTicketStr string
	OldPassengerStr    string
}

type TicketTypeName struct {
	Id   string
	Name string
}

func KnownTicketTypeName() []*TicketTypeName {
	return []*TicketTypeName{
		{"1", "成人票"},
		{"2", "小孩票"},
		{"3", "学生票"},
		{"4", "伤残军人票"},
	}
}

type SeatTypeName struct {
	Name  string
	Value string
	Code  string
}

func KnownSeatTypeName() []*SeatTypeName {
	return []*SeatTypeName{
		{"一等座", "M", "ZY"},
		{"二等座", "O", "ZE"},
		{"商务座", "9", "SWZ"},
		{"特等座", "P", "TZ"},
		{"硬座", "1", "YZ"},
		{"软座", "2", "RZ"},
		{"硬卧", "3", "YW"},
		{"软卧", "4", "RW"},
		{"高级软卧", "6", "GR"},
		{"无座", "WZ", "WZ"},
	}
}
