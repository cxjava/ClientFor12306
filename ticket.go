package main

import "strings"

var (
	SeatTypeNameToV = map[string]string{
		"一等座":  "M",
		"二等座":  "O",
		"商务座":  "9",
		"特等座":  "P",
		"硬座":   "1",
		"软座":   "2",
		"硬卧":   "3",
		"软卧":   "4",
		"高级软卧": "6",
		"无座":   "WZ",
	}
	SeatTypeValueToN = map[string]string{
		"M":  "一等座",
		"O":  "二等座",
		"9":  "商务座",
		"P":  "特等座",
		"1":  "硬座",
		"2":  "软座",
		"3":  "硬卧",
		"4":  "软卧",
		"6":  "高级软卧",
		"WZ": "无座",
	}
)

type TicketQuery struct {
	Start     string `form:"start" binding:"required"`
	End       string `form:"end" binding:"required"`
	Train     string `form:"train" binding:"required"`
	TrainDate string `form:"date" binding:"required"`
	P1        struct {
		PassengerName1 string `form:"passengerName1" binding:"required"`
		TicketType1    string `form:"ticketType1" binding:"required"`
		SeatType1      string `form:"seatType1" binding:"required"`
		IDType1        string `form:"IDType1" binding:"required"`
		IDNumber1      string `form:"IDNumber1" binding:"required"`
	}
	P2 struct {
		PassengerName2 string `form:"passengerName2" `
		TicketType2    string `form:"ticketType2" `
		SeatType2      string `form:"seatType2" `
		IDType2        string `form:"IDType2" `
		IDNumber2      string `form:"IDNumber2" `
	}
	P3 struct {
		PassengerName3 string `form:"passengerName3" `
		TicketType3    string `form:"ticketType3" `
		SeatType3      string `form:"seatType3" `
		IDType3        string `form:"IDType3" `
		IDNumber3      string `form:"IDNumber3" `
	}
	P4 struct {
		PassengerName4 string `form:"passengerName4" `
		TicketType4    string `form:"ticketType4" `
		SeatType4      string `form:"seatType4" `
		IDType4        string `form:"IDType4" `
		IDNumber4      string `form:"IDNumber4" `
	}
	P5 struct {
		PassengerName5 string `form:"passengerName5"`
		TicketType5    string `form:"ticketType5"`
		SeatType5      string `form:"seatType5"`
		IDType5        string `form:"IDType5"`
		IDNumber5      string `form:"IDNumber5"`
	}
}
type TicketTypeName struct {
	Id   string
	Name string
}

type IDTypeName struct {
	Id   string
	Name string
}

type SeatTypeName struct {
	Name  string
	Value string
	Code  string
}

func KnownTicketTypeName() []*TicketTypeName {
	return []*TicketTypeName{
		{"1", "成人票"},
		{"2", "小孩票"},
		{"3", "学生票"},
		{"4", "伤残军人票"},
	}
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

func (t *TicketQuery) parseTicket(q *Query) {
	q.FromStations = parseStrings(t.Start)
	q.ToStations = parseStrings(t.End)
	q.TrianCodes = parseStrings(t.Train)
	q.TrainDate = t.TrainDate
}

func parseStrings(str string) (s []string) {
	str = strings.Replace(str, "，", ",", -1)
	str = strings.Replace(str, "；", ",", -1)
	str = strings.Replace(str, ";", ",", -1)
	if !strings.Contains(str, ",") {
		s = append(s, str)
	} else {
		for _, v := range strings.Split(str, ",") {
			if ts := strings.TrimSpace(v); ts != "" {
				s = append(s, ts)
			}
		}
	}
	return
}

func (t *TicketQuery) parseStranger(q *Query) {
	oStr, nStr := "", ""
	m := make(map[string]int)
	oStr, nStr, m = plusNum(oStr, nStr, t.P1.SeatType1, t.P1.TicketType1, t.P1.PassengerName1, t.P1.IDType1, t.P1.IDNumber1, m)
	oStr, nStr, m = plusNum(oStr, nStr, t.P2.SeatType2, t.P2.TicketType2, t.P2.PassengerName2, t.P2.IDType2, t.P2.IDNumber2, m)
	oStr, nStr, m = plusNum(oStr, nStr, t.P3.SeatType3, t.P3.TicketType3, t.P3.PassengerName3, t.P3.IDType3, t.P3.IDNumber3, m)
	oStr, nStr, m = plusNum(oStr, nStr, t.P4.SeatType4, t.P4.TicketType4, t.P4.PassengerName4, t.P4.IDType4, t.P4.IDNumber4, m)
	oStr, nStr, m = plusNum(oStr, nStr, t.P5.SeatType5, t.P5.TicketType5, t.P5.PassengerName5, t.P5.IDType5, t.P5.IDNumber5, m)

	q.NumOfSeatType = m
	q.OldPassengerStr = oStr
	q.PassengerTicketStr = nStr[:len(nStr)-1]
}

func plusNum(oStr, nStr, SeatType, TicketType, Name, PassengerIdTypeCode, PassengerIdNo string, m map[string]int) (string, string, map[string]int) {
	if strings.TrimSpace(Name) != "" {
		name := SeatTypeValueToN[SeatType]
		m[name] = m[name] + 1
		nStr += SeatType + ",0," + TicketType + "," + Name + "," + PassengerIdTypeCode + "," + PassengerIdNo + ",,N_"
		oStr += Name + "," + PassengerIdTypeCode + "," + PassengerIdNo + "," + TicketType + "_"
		return oStr, nStr, m
	}
	return oStr, nStr, m
}
