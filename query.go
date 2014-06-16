package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Query struct {
	CDN                string
	FromStations       []string
	ToStations         []string
	TrianCodes         []string // 车次,T260,T258 等等
	PassengerTicketStr string
	OldPassengerStr    string
	NumOfSeatType      map[string]int
	TrainDate          string
}

//查询
func (q *Query) Order() (or *Order) {

	or = nil
	if err, ticketInfo := q.queryLeftTicket(); err == nil { //获取车次
		for _, trainCode := range q.TrianCodes { //要预订的车次
			trainCode = strings.ToUpper(trainCode)
			for _, data := range ticketInfo.Data { //查询结果的每个车次
				//查询到的车次
				if tkt := data.Ticket; tkt.StationTrainCode == trainCode { //是预订的车次
					//获取余票信息
					ticketNum := GetTicketNum(tkt.YpInfo, tkt.YpEx)
					if validateNum(ticketNum, q.NumOfSeatType) { //想要预订席别的余票大于等于订票人的人数
						Info(q.CDN, "开始订票", q.TrainDate, "车次", tkt.StationTrainCode, "余票", fmt.Sprintf("%v", ticketNum))
						or = &Order{
							CDN:                q.CDN,
							PassengerTicketStr: q.PassengerTicketStr,
							OldPassengerStr:    q.OldPassengerStr,
							Ticket:             tkt,
							SecretStr:          data.SecretStr,
							SubmitCaptchaStr:   make(chan string),
							TrainDate:          q.TrainDate,
							SeatType:           getSeatType(q.NumOfSeatType),
						}
						break
					} else if data.ButtonTextInfo != "预订" {
						Warn("车次", tkt.StationTrainCode, data.ButtonTextInfo)
					} else {
						Warn("车次", tkt.StationTrainCode, "余票不足！！！")
						Warn("剩余的票:", fmt.Sprintf("%v", ticketNum))
						Warn("订购的票:", fmt.Sprintf("%v", q.NumOfSeatType))
						break
					}
				}
			}
		}
	} else {
		Error(q.CDN, "余票查询错误", ticketInfo, err)
	}
	return
}

//查询余票
func (q *Query) queryLeftTicket() (error, *QueryLeftNewDTO) {
	leftTicketUrl := ""
	leftTicketUrl += "leftTicketDTO.train_date=" + q.TrainDate + "&"
	//随机获取出发站的code去查询，防止缓存
	leftTicketUrl += "leftTicketDTO.from_station=" + StationMap[q.FromStations[rand.Intn(len(q.FromStations))]] + "&"
	//随机获取终点站的code去查询，防止缓存
	leftTicketUrl += "leftTicketDTO.to_station=" + StationMap[q.ToStations[rand.Intn(len(q.ToStations))]] + "&"
	leftTicketUrl += "purpose_codes=" + Purpose_codes

	Info("queryLeftTicket url:", URLQueryLog+leftTicketUrl)

	h := map[string]string{
		"Cache-Control":     "no-cache",
		"X-Requested-With":  "XMLHttpRequest",
		"Referer":           "https://kyfw.12306.cn/otn/leftTicket/init",
		"If-Modified-Since": time.Now().Local().Format(time.RFC1123Z),
		"If-None-Match":     strconv.FormatInt(time.Now().UnixNano(), 10),
	}

	go DoForWardRequestHeader(q.CDN, "GET", URLQueryLog+leftTicketUrl, nil, h)

	body, err := DoForWardRequestHeader(q.CDN, "GET", URLQuery+leftTicketUrl, nil, h)
	if err != nil {
		Error("queryLeftTicket DoForWardRequestHeader error:", err)
		return err, nil
	}
	Debug("queryLeftTicket body:", body)

	if !strings.Contains(body, "queryLeftNewDTO") {
		Error("CDN:"+q.CDN+"查询余票出错,返回:", body, "查询链接:", leftTicketUrl)
		return err, nil
	}
	leftTicket := &QueryLeftNewDTO{}

	if err := json.Unmarshal([]byte(body), &leftTicket); err != nil {
		Error(q.CDN, "queryLeftTicket error:", err)
		return err, nil
	} else {
		Info(q.CDN, "获取余票成功！")
	}
	return nil, leftTicket
}

// 获取 类似:[硬卧:1,二等座:2]==>二等座
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

//验证余票是否充足
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
