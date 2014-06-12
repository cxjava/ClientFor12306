package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

var (
	SubmitCaptchaStr = make(chan string)
	login            = &Login{}
	ticket           = &TicketQueryInfo{
		SubmitCaptchaStr: make(chan string),
		P1:               &PassengerOrder{},
		P2:               &PassengerOrder{},
		P3:               &PassengerOrder{},
		P4:               &PassengerOrder{},
		P5:               &PassengerOrder{},
	}

	mapPassengers = make(map[string]Passenger)
)

//查询
func (t *TicketQuery) Order() (order *Order) {
	order = nil
	if err, tickets := t.queryLeftTicket(); err == nil { //获取车次
		for _, trainCode := range t.Trians { //要预订的车次
			trainCode = strings.ToUpper(trainCode)
			for _, data := range tickets.Data { //每个车次
				//查询到的车次
				tkt := data.Ticket
				if tkt.StationTrainCode == trainCode { //是预订的车次
					//获取余票信息
					ticketNum := GetTicketNum(tkt.YpInfo, tkt.YpEx)
					if validateNum(ticketNum, t.NumOfSeatType) { //想要预订席别的余票大于等于订票人的人数
						Info(t.CDN, "开始订票", t.TrainDate, "车次", tkt.StationTrainCode, "余票", fmt.Sprintf("%v", ticketNum))
						order = &Order{
							CDN:                t.CDN,
							PassengerTicketStr: t.PassengerTicketStr,
							OldPassengerStr:    t.OldPassengerStr,
							Ticket:             tkt,
							SecretStr:          data.SecretStr,
							SubmitCaptchaStr:   make(chan string),
							TrainDate:          t.TrainDate,
							SeatType:           getSeatType(t.NumOfSeatType),
						}
						break
					} else {
						Info("车次", tkt.StationTrainCode, "余票不足！！！剩余票：", fmt.Sprintf("%v", ticketNum), "订购的票:", fmt.Sprintf("%v", t.NumOfSeatType))
						break
					}
				}
			}
		}
	} else {
		Error(t.CDN, "余票查询错误", tickets, err)
	}
	return
}
func getSeatType(seatTypeNum map[string]int) (t string) {
	t = "3"
	max, name := 0, "硬卧"
	for k, v := range seatTypeNum {
		if v > max {
			max, name = v, k
		}
	}
	t = SeatTypeNameToV[name]
	return
}
func validateNum(ticketNum, seatTypeNum map[string]int) (b bool) {
	b = true
	for k, v := range seatTypeNum {
		if ticketNum[k] < v {
			b = false
			break
		}
	}
	return
}

//查询余票
func (t *TicketQuery) queryLeftTicket() (error, *QueryLeftNewDTO) {
	fr := t.FromStations
	to := t.ToStations
	leftTicketUrl := ""
	leftTicketUrl += "leftTicketDTO.train_date=" + t.TrainDate + "&"
	leftTicketUrl += "leftTicketDTO.from_station=" + StationMap[fr[rand.Intn(len(fr))]] + "&"
	leftTicketUrl += "leftTicketDTO.to_station=" + StationMap[to[rand.Intn(len(to))]] + "&"
	leftTicketUrl += "purpose_codes=" + Purpose_codes

	Info("queryLeftTicket url:", leftTicketUrl)

	h := map[string]string{
		"Cache-Control":     "no-cache",
		"x-requested-with":  "XMLHttpRequest",
		"Referer":           "https://kyfw.12306.cn/otn/leftTicket/init",
		"If-Modified-Since": time.Now().Local().Format(time.RFC1123Z),
		"If-None-Match":     strconv.FormatInt(time.Now().UnixNano(), 10),
	}

	go DoForWardRequestHeader(t.CDN, "GET", URLQueryLog+leftTicketUrl, nil, h)

	body, err := DoForWardRequestHeader(t.CDN, "GET", URLQuery+leftTicketUrl, nil, h)
	if err != nil {
		Error("queryLeftTicket DoForWardRequest error:", err)
		return err, nil
	}
	Debug("queryLeftTicket body:", body)

	if !strings.Contains(body, "queryLeftNewDTO") {
		Error("查询余票出错，返回:", body, "查询链接:", leftTicketUrl)
		return err, nil
	}
	leftTicket := &QueryLeftNewDTO{}

	if err := json.Unmarshal([]byte(body), &leftTicket); err != nil {
		Error("queryLeftTicket", t.CDN, err)
		return err, nil
	} else {
		Info(t.CDN, "获取成功！")
	}

	return nil, leftTicket
}

func (t *TicketQuery) parseTicket() {
	t.FromStations = parseStrings(t.Start)
	t.ToStations = parseStrings(t.End)
	t.Trians = parseStrings(t.Train)
}
func parseStrings(str string) (s []string) {
	str = strings.Replace(str, "，", ",", -1)
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

func (t *TicketQuery) parseStranger() {
	oStr, nStr := "", ""
	m := make(map[string]int)
	oStr, nStr, m = plusNum(oStr, nStr, t.P1.SeatType1, t.P1.TicketType1, t.P1.PassengerName1, t.P1.IDType1, t.P1.IDNumber1, m)
	oStr, nStr, m = plusNum(oStr, nStr, t.P2.SeatType2, t.P2.TicketType2, t.P2.PassengerName2, t.P2.IDType2, t.P2.IDNumber2, m)
	oStr, nStr, m = plusNum(oStr, nStr, t.P3.SeatType3, t.P3.TicketType3, t.P3.PassengerName3, t.P3.IDType3, t.P3.IDNumber3, m)
	oStr, nStr, m = plusNum(oStr, nStr, t.P4.SeatType4, t.P4.TicketType4, t.P4.PassengerName4, t.P4.IDType4, t.P4.IDNumber4, m)
	oStr, nStr, m = plusNum(oStr, nStr, t.P5.SeatType5, t.P5.TicketType5, t.P5.PassengerName5, t.P5.IDType5, t.P5.IDNumber5, m)

	t.NumOfSeatType = m
	t.OldPassengerStr = oStr
	t.PassengerTicketStr = nStr[:len(nStr)-1]
	return
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
