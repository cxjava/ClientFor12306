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
func (t *TicketQueryInfo) Order(cdn string) {

	if tickets := t.queryLeftTicket(cdn); tickets != nil { //获取车次
		for _, trainCode := range t.Trians { //要预订的车次
			trainCode = strings.ToUpper(trainCode)
			for _, data := range tickets.Data { //每个车次
				//查询到的车次
				tkt := data.Ticket
				if tkt.StationTrainCode == strings.ToUpper(trainCode) { //是预订的车次
					//获取余票信息
					ticketNum := GetTicketNum(tkt.YpInfo, tkt.YpEx)
					if validateNum(ticketNum, t.NumOfSeatType) { //想要预订席别的余票大于等于订票人的人数
						Info(cdn, "开始订票", t.TrainDate.Format("2006-01-02"), "车次", tkt.StationTrainCode, "余票", fmt.Sprintf("%v", ticketNum))
						order = &Order{
							CDN:                cdn,
							PassengerTicketStr: t.PassengerTicketStr,
							OldPassengerStr:    t.OldPassengerStr,
							Ticket:             tkt,
							SecretStr:          data.SecretStr,
							SubmitCaptchaStr:   make(chan string),
							TrainDate:          t.TrainDate,
							SeatType:           getSeatType(t.NumOfSeatType),
						}
						var err error
						go checkUser(Conf.CDN[0])
						err = order.submitOrderRequest()
						if err != nil {
							Error(err)
							return
						}
						// leftTicketInit(Conf.CDN[0])
						// t.queryLeftTicket(cdn)
						// order.inits()
						err = order.initDc()
						if err != nil {
							Error(err)
							return
						}
						dyQueryJs(Conf.CDN[0])
						// go ticket.queryLeftTicket(order.CDN)
						go order.getPassCodeNew()
						break
					} else {
						Info("车次", tkt.StationTrainCode, "余票不足！！！剩余票：", fmt.Sprintf("%v", ticketNum), "订购的票:", fmt.Sprintf("%v", t.NumOfSeatType))
					}
				} else { //不是预订的车次
					//Debug(tkt.StationTrainCode, "余票", fmt.Sprintf("%v", getTicketNum(tkt.YpInfo, tkt.YpEx)))
				}
			}
		}
	} else {
		Error(cdn, "余票查询错误", tickets)
	}
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
		Info(k, v, ticketNum[k])
		if ticketNum[k] < v {
			b = false
			break
		}
	}
	return
}

//查询余票
func (t *TicketQueryInfo) queryLeftTicket(cdn string) *QueryLeftNewDTO {
	fr := t.FromStations
	to := t.ToStations
	leftTicketUrl := ""
	leftTicketUrl += "leftTicketDTO.train_date=" + t.TrainDate.Format("2006-01-02") + "&"
	leftTicketUrl += "leftTicketDTO.from_station=" + StationMap[fr[rand.Intn(len(fr))]] + "&"
	leftTicketUrl += "leftTicketDTO.to_station=" + StationMap[to[rand.Intn(len(to))]] + "&"
	leftTicketUrl += "purpose_codes=ADULT"

	Info("queryLeftTicket url:", leftTicketUrl)

	go DoForWardRequest(cdn, "GET", URLQueryLog+leftTicketUrl, nil)

	h := map[string]string{"If-Modified-Since": time.Now().Local().Format(time.RFC1123Z),
		"If-None-Match": strconv.FormatInt(time.Now().UnixNano(), 10)}

	body, err := DoForWardRequestHeader(cdn, "GET", URLQuery+leftTicketUrl, nil, h)
	if err != nil {
		Error("queryLeftTicket DoForWardRequest error:", err)
		return nil
	}
	Debug("queryLeftTicket body:", body)

	if !strings.Contains(body, "queryLeftNewDTO") {
		Error("查询余票出错，返回:", body, "查询链接:", leftTicketUrl)
		//删除废弃的CDN
		// if len(availableCDN) > 5 {
		// delete(availableCDN, cdn)
		// }
		return nil
	}
	leftTicket := &QueryLeftNewDTO{}

	if err := json.Unmarshal([]byte(body), &leftTicket); err != nil {
		Error("queryLeftTicket", cdn, err)
		return nil
	} else {
		Info(cdn, "获取成功！")
	}

	return leftTicket
}

func parseTicket() {
	ticket.FromStations = parseStrings(ticket.FromStationsStr)
	ticket.ToStations = parseStrings(ticket.ToStationsStr)
	ticket.Trians = parseStrings(ticket.TriansStr)

	o, n := parseStranger(ticket)
	ticket.OldPassengerStr = o
	ticket.PassengerTicketStr = n[:len(n)-1]

	ticket.Order(Conf.CDN[0])
}
func plusNum(num int, oStr, nStr string, p *PassengerOrder, m map[string]int) (int, string, string, map[string]int) {
	if strings.Trim(p.Name, " ") != "" {
		name := SeatTypeValueToN[p.SeatType]
		m[name] = m[name] + 1
		nStr += p.SeatType + ",0," + p.TicketType + "," + p.Name + "," + p.PassengerIdTypeCode + "," + p.PassengerIdNo + ",,N_"
		oStr += p.Name + "," + p.PassengerIdTypeCode + "," + p.PassengerIdNo + "," + p.TicketType + "_"
		return num + 1, oStr, nStr, m
	}
	return num, oStr, nStr, m
}
func parseStranger(ticket *TicketQueryInfo) (oStr, nStr string) {
	num := 0
	m := make(map[string]int)
	num, oStr, nStr, m = plusNum(num, oStr, nStr, ticket.P1, m)
	num, oStr, nStr, m = plusNum(num, oStr, nStr, ticket.P2, m)
	num, oStr, nStr, m = plusNum(num, oStr, nStr, ticket.P3, m)
	num, oStr, nStr, m = plusNum(num, oStr, nStr, ticket.P4, m)
	num, oStr, nStr, m = plusNum(num, oStr, nStr, ticket.P5, m)

	ticket.NumOfPassenger = num
	ticket.NumOfSeatType = m
	return
}
func parseStrings(str string) (s []string) {
	if strings.ContainsRune(str, rune('，')) {
		for _, v := range strings.Split(str, "，") {
			if v != "" {
				s = append(s, v)
			}
		}
	}
	if strings.ContainsRune(str, rune(',')) {
		for _, v := range strings.Split(str, ",") {
			if v != "" {
				s = append(s, v)
			}
		}
	}
	if len(s) == 0 {
		s = append(s, str)
	}
	return
}
