package main

import (
	"encoding/json"
	"fmt"
	"image"
	"math/rand"
	"net/http"
	"net/url"
	"strings"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

var (
	login  = &Login{}
	ticket = &TicketQueryInfo{
		P1: &PassengerOrder{},
		P2: &PassengerOrder{},
		P3: &PassengerOrder{},
		P4: &PassengerOrder{},
		P5: &PassengerOrder{},
	}
	passenger     = &PassengerDTO{}
	mapPassengers = make(map[string]Passenger)
)

func main() {

	createUI()

}

/*
//提交订单
func submitOrderRequest(urlValues url.Values, cdn string, t Ticket) error {
	// defer func() {
	// 	<-submitChannel
	// }()
	// submitChannel <- 1

	params, _ := url.QueryUnescape(urlValues.Encode())
	Debug(params)

	body, err := DoForWardRequest(cdn, "POST", SubmitOrderRequestURL, strings.NewReader(params))
	if err != nil {
		Error("submitOrderRequest DoForWardRequest error:", err)
		return err
	}
	Debug("submitOrderRequest body:", body)

	if strings.Contains(body, `"submitStatus":true`) {
		orderResoult := &OrderResoult{}
		if err := json.Unmarshal([]byte(body), &orderResoult); err != nil {
			Error("submitOrderRequest", err)
			return err
		} else {
			dataResult := strings.Split(orderResoult.Data.Result, "#")
			//key_check_isChange=99F79C00DFB9BF8713D23EFA4A8CF06BCA8C412DAC19686DCE306476
			// leftTicketStr = 1002353600401115003110023507803007450039
			// for getQueueCount
			Info("key_check_isChange:", dataResult[1], "leftTicket:", dataResult[2])
			//获取队列
			urlValues := url.Values{}
			urlValues.Add("train_date", time.Now().String())
			urlValues.Add("train_no", t.TrainNo)
			urlValues.Add("stationTrainCode", t.StationTrainCode)
			urlValues.Add("seatType", "3")
			urlValues.Add("fromStationTelecode", t.FromStationTelecode)
			urlValues.Add("toStationTelecode", t.ToStationTelecode)
			urlValues.Add("leftTicket", dataResult[2])
			urlValues.Add("purpose_codes", Purpose_codes)
			urlValues.Add("_json_att", Json_att)

			getQueueCount(urlValues, dataResult, cdn)

			//是否同时进行提交表单
			if Config.System.GoBoth {
				//确认队列
				urlValuesForQueue := url.Values{}
				urlValuesForQueue.Add("passengerTicketStr", passengerTicketStr)
				urlValuesForQueue.Add("oldPassengerStr", oldPassengerStr)
				urlValuesForQueue.Add("randCode", dataResult[1])
				urlValuesForQueue.Add("purpose_codes", Purpose_codes)
				urlValuesForQueue.Add("key_check_isChange", dataResult[1])
				urlValuesForQueue.Add("leftTicketStr", dataResult[2])
				urlValuesForQueue.Add("train_location", dataResult[0])
				urlValuesForQueue.Add("_json_att", Json_att)

				// 需要延迟提交，提早好像要被踢！！！
				if Config.System.SubmitTime > 1000 {
					time.Sleep(time.Millisecond * time.Duration(Config.System.SubmitTime))
				}

				confirmSingleForQueue(urlValuesForQueue, cdn)
			}

		}
	} else if strings.Contains(body, `您还有未处理的订单`) {
		log.Println("订票成功！！")
		// sendMessage("订票成功！！")
	} else if strings.Contains(body, `用户未登录`) {
		log.Println("用户未登录！！")
		// sendMessage("用户未登录！！")
	} else if strings.Contains(body, `取消次数过多`) {
		log.Println("由于您取消次数过多！！")
		// sendMessage("由于您取消次数过多！！")
	} else if strings.Contains(body, `互联网售票实行实名制`) {
		log.Println("貌似你已经购买了相同的车票！！")
		// sendMessage("貌似你已经购买了相同的车票！！")
	} else {
		Warn(cdn, "订票请求警告:", body)
	}
	return nil
}*/

//查询
func (t *TicketQueryInfo) Order(cdn string) {
	//睡眠下，随机
	//time.Sleep(time.Millisecond * time.Duration(Config.System.SubmitTime))
	//time.Sleep(time.Millisecond * time.Duration(rand.Int63n(Config.System.RefreshTime)))

	// defer func() {
	// <-queryChannel
	// }()
	// queryChannel <- 1

	// queryJs(cdn)

	if tickets := t.queryLeftTicket(cdn); tickets != nil { //获取车次
		for _, trainCode := range t.Trians { //要预订的车次
			trainCode = strings.ToUpper(trainCode)
			for _, data := range tickets.Data { //每个车次
				//查询到的车次
				tkt := data.Ticket
				if tkt.StationTrainCode == trainCode { //是预订的车次
					//获取余票信息
					ticketNum := GetTicketNum(tkt.YpInfo, tkt.YpEx)
					numOfTicket := ticketNum["硬卧"]       //通过配置的seattype获取余票信息
					if numOfTicket >= t.NumOfPassenger { //想要预订席别的余票大于等于订票人的人数
						Info(cdn, "开始订票", t.TrainDate.Format("2006-01-02"), "车次", tkt.StationTrainCode, "余票", fmt.Sprintf("%v", ticketNum))
						urlValues := url.Values{}
						urlValues.Add("bed_level_order_num", Bed_level_order_num)
						urlValues.Add("cancel_flag", Cancel_flag)
						urlValues.Add("purpose_codes", Purpose_codes)
						urlValues.Add("tour_flag", Tour_flag)
						urlValues.Add("secretStr", data.SecretStr)
						urlValues.Add("train_date", t.TrainDate.Format("2006-01-02"))
						urlValues.Add("query_from_station_name", tkt.FromStationName)
						urlValues.Add("query_to_station_name", tkt.ToStationName)
						urlValues.Add("passengerTicketStr", t.PassengerTicketStr)
						urlValues.Add("oldPassengerStr", t.OldPassengerStr)

						// go submitOrderRequest(urlValues, cdn, tkt)

					} else if numOfTicket > 0 && numOfTicket < t.NumOfPassenger {
						Info("车次", tkt.StationTrainCode, "余票不足！！！", fmt.Sprintf("%v", ticketNum))
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

	go DoForWardRequest(cdn, "GET", LogQueryLeftTicketURL+leftTicketUrl, nil)
	body, err := DoForWardRequest(cdn, "GET", QueryLeftTicketURL+leftTicketUrl, nil)
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
	Info(ticket.TrainDate.Format("2006-01-02"))
	Info(ticket)
	ticket.FromStations = parseStrings(ticket.FromStationsStr)
	ticket.ToStations = parseStrings(ticket.ToStationsStr)
	ticket.Trians = parseStrings(ticket.TriansStr)

	Info(ticket)
	o, n := parseStranger(ticket)
	ticket.OldPassengerStr = o
	ticket.PassengerTicketStr = n[:len(n)-1]
	Info(ticket)

	// ticket.queryLeftTicket(Conf.CDN[0])
	ticket.Order(Conf.CDN[0])
}
func plusNum(num int, oStr, nStr string, p *PassengerOrder) (int, string, string) {
	if strings.Trim(p.Name, " ") != "" {
		nStr += p.SeatType + ",0," + p.TicketType + "," + p.Name + "," + p.PassengerIdTypeCode + "," + p.PassengerIdNo + ",,N_"
		oStr += p.Name + "," + p.PassengerIdTypeCode + "," + p.PassengerIdNo + "," + p.TicketType + "_"
		return num + 1, nStr, oStr
	}
	return num, nStr, oStr
}
func parseStranger(ticket *TicketQueryInfo) (oStr, nStr string) {
	num := 0
	num, oStr, nStr = plusNum(num, oStr, nStr, ticket.P1)
	num, oStr, nStr = plusNum(num, oStr, nStr, ticket.P2)
	num, oStr, nStr = plusNum(num, oStr, nStr, ticket.P3)
	num, oStr, nStr = plusNum(num, oStr, nStr, ticket.P4)
	num, oStr, nStr = plusNum(num, oStr, nStr, ticket.P5)

	ticket.NumOfPassenger = num
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

//获取新的验证码图片
func GetImage(cdn string) image.Image {
	req, err := http.NewRequest("GET", PassCodeNewURL, nil)
	if err != nil {
		Error("GetImage http.NewRequest error:", err)
		return nil
	}
	con, err := NewForwardClientConn(cdn, req.URL.Scheme)
	if err != nil {
		Error("GetImage newForwardClientConn error:", err)
		return nil
	}
	defer con.Close()
	resp, err := con.Do(req)
	if err != nil {
		Error("GetImage con.Do error:", err)
		return nil
	}
	defer resp.Body.Close()
	//set cookie
	login.Cookie = GetCookieFromRespHeader(resp)
	Debug("==" + login.Cookie + "==")

	img, s, err := image.Decode(resp.Body)
	Debug("image type:", s)
	if err != nil {
		Error("GetImage image.Decode:", err)
		return nil
	}
	return img
}

//从响应消息头里面获取cookie
func GetCookieFromRespHeader(resp *http.Response) (cookie string) {
	cookies := []string{}
	for k, v := range resp.Header {
		Debug("k=", k, "v=", v)
		if k == "Set-Cookie" {
			for _, b := range v {
				v := strings.Split(b, ";")[0]
				cookies = append(cookies, v)
				cookies = append(cookies, "; ")
			}
		}
	}
	d := strings.Join(cookies, "")
	if len(d) < 2 {
		return ""
	}
	cookie = d[:len(d)-2]
	return
}

//获取联系人
func getPassengerDTO() {
	for _, cdn := range Conf.CDN {
		Info("开始获取联系人！")
		body, err := DoForWardRequest(cdn, "POST", GetPassengerDTOURL, nil)
		if err != nil {
			Error("getPassengerDTO DoForWardRequest error:", err)
			continue
		}
		Debug("getPassengerDTO body:", body)

		if !strings.Contains(body, "passenger_name") {
			Error("获取联系人出错!!!!!!返回:", body)
			continue
		}

		if err := json.Unmarshal([]byte(body), &passenger); err != nil {
			Error("getPassengerDTO", cdn, err)
			continue
		} else {
			Info(cdn, "获取成功！")
			break
		}
	}
	Debug(passenger)

	go func() {
		model := []string{}
		for _, v1 := range passenger.Data.NormalPassengers {
			model = append(model, v1.PassengerName)
			mapPassengers[v1.PassengerName] = v1
		}
		myPassengers.SetModel(model)
	}()
}
