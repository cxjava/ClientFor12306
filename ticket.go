package main

import "time"

type PassengerOrder struct {
	Name                string
	TicketType          string // 成人，学生
	PassengerIdTypeCode string
	PassengerIdNo       string
	SeatType            string //席别：硬卧，硬座
	SeatTypeName        string
}

type TicketQueryInfo struct {
	TrainDate          time.Time
	FromStations       []string
	FromStationsStr    string
	ToStations         []string
	ToStationsStr      string
	Trians             []string
	TriansStr          string
	PassengerTicketStr string
	OldPassengerStr    string
	NumOfPassenger     int
	P1                 *PassengerOrder
	P2                 *PassengerOrder
	P3                 *PassengerOrder
	P4                 *PassengerOrder
	P5                 *PassengerOrder
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

type IDTypeName struct {
	Id   string
	Name string
}

func KnownIDTypeName() []*IDTypeName {
	return []*IDTypeName{
		{"1", "二代身份证"},
		{"C", "港澳通行证"},
		{"G", "台湾通行证"},
		{"B", "护照"},
		{"H", "外国人居留证"},
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
