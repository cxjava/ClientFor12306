package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/lxn/walk"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

var (
	SubmitCaptchaStr = make(chan string)
	order            *Order
	login            = &Login{}
	ticket           = &TicketQueryInfo{
		SubmitCaptchaStr: make(chan string),
		P1:               &PassengerOrder{},
		P2:               &PassengerOrder{},
		P3:               &PassengerOrder{},
		P4:               &PassengerOrder{},
		P5:               &PassengerOrder{},
	}
	passenger     = &PassengerDTO{}
	mapPassengers = make(map[string]Passenger)
)

func main() {

	createUI()

}

func (order *Order) confirmSingleForQueue() error {

	val := url.Values{}
	val.Add("passengerTicketStr", order.PassengerTicketStr)
	val.Add("oldPassengerStr", order.OldPassengerStr)
	val.Add("randCode", order.RandCode)
	val.Add("purpose_codes", Purpose_codes)
	val.Add("key_check_isChange", order.KeyCheckIsChange)
	val.Add("leftTicket", order.Ticket.YpInfo)
	val.Add("train_location", order.TrainLocation)
	val.Add("_json_att", Json_att)
	val.Add("REPEAT_SUBMIT_TOKEN", order.RepeatSubmitToken)

	Info("confirmSingleForQueue params:", val.Encode())

	params, _ := url.QueryUnescape(val.Encode())
	h := map[string]string{"Referer": "https://kyfw.12306.cn/otn/confirmPassenger/initDc"}
	body, err := DoForWardRequestHeader(order.CDN, "POST", URLConfirmSingleForQueue, strings.NewReader(params), h)
	if err != nil {
		Error("confirmSingleForQueue DoForWardRequest error:", err)
		return err
	}
	Info("confirmSingleForQueue body:", body)
	if strings.Contains(body, `"submitStatus":true`) {
		Info("提交订单成功 body:", body)
	} else if strings.Contains(body, `订单未支付`) {
		Warn("订票成功！！")
	} else if strings.Contains(body, `用户未登录`) {
		Warn("用户未登录！！")
	} else if strings.Contains(body, `取消次数过多`) {
		Warn("由于您取消次数过多！！")
	} else if strings.Contains(body, `互联网售票实行实名制`) {
		Warn("貌似你已经购买了相同的车票！！")
	} else {
		Warn(order.CDN, "订票请求警告:", body)
	}
	return nil
}

func (order *Order) getQueueCount() error {
	// go ticket.queryLeftTicket(order.CDN)
	val := url.Values{}
	// val.Add("train_date", order.TrainDate.Local().Format(`Mon+Jan+02+2006+15:04:05+GMT-0700+(China+Standard+Time)`))
	val.Add("train_date", `Thu Jun 05 2014 00:00:00 GMT+0800 (China Standard Time)`)
	val.Add("train_no", order.Ticket.TrainNo)
	val.Add("stationTrainCode", order.Ticket.StationTrainCode)
	val.Add("seatType", order.SeatType)
	val.Add("fromStationTelecode", order.Ticket.FromStationTelecode)
	val.Add("toStationTelecode", order.Ticket.ToStationTelecode)
	val.Add("leftTicket", order.Ticket.YpInfo)
	val.Add("purpose_codes", Purpose_codes)
	val.Add("_json_att", Json_att)
	val.Add("REPEAT_SUBMIT_TOKEN", order.RepeatSubmitToken)

	Info("getQueueCount params:", val)

	params, _ := url.QueryUnescape(val.Encode())
	body, err := DoForWardRequestHeader(order.CDN, "POST", URLQueueCount, strings.NewReader(params), nil)
	if err != nil {
		Error("getQueueCount DoForWardRequest error:", err)
		return err
	}
	Info("getQueueCount body:", body)
	// gqc := &JsonGetQueueCount{}
	// if err := json.Unmarshal([]byte(body), &gqc); err != nil {
	// 	Error("checkOrderInfo json.Unmarshal:", err)
	// 	return err
	// }

	// if !coi.Data.SubmitStatus {
	// 	return fmt.Errorf("checkOrderInfo 出错！body:%v", body)
	// }
	time.Sleep(time.Second * 3)
	return order.confirmSingleForQueue()
}

func (order *Order) checkOrderInfo() error {
	val := url.Values{}
	val.Add("cancel_flag", Cancel_flag)
	val.Add("bed_level_order_num", Bed_level_order_num)
	val.Add("passengerTicketStr", order.PassengerTicketStr)
	val.Add("oldPassengerStr", order.OldPassengerStr)
	val.Add("tour_flag", Tour_flag)
	Info("验证码:", order.RandCode)
	val.Add("randCode", order.RandCode)
	val.Add("_json_att", Json_att)
	val.Add("REPEAT_SUBMIT_TOKEN", order.RepeatSubmitToken)
	Info("checkOrderInfo params:", val)

	params, _ := url.QueryUnescape(val.Encode())
	body, err := DoForWardRequestHeader(order.CDN, "POST", URLCheckOrderInfo, strings.NewReader(params), nil)
	if err != nil {
		Error("checkOrderInfo DoForWardRequest error:", err)
		return err
	}
	Info("checkOrderInfo body:", body)

	coi := &JsonCheckOrderInfo{}
	if err := json.Unmarshal([]byte(body), &coi); err != nil {
		Error("checkOrderInfo json.Unmarshal:", err)
		return err
	}

	if !coi.Data.SubmitStatus {
		return fmt.Errorf("checkOrderInfo 出错！body:%v", body)
	}
	time.Sleep(time.Second * 2)
	return order.getQueueCount()
}
func (order *Order) checkRandCodeAnsyn(randCode string) (r bool, msg []string) {
	val := url.Values{}
	val.Add("randCode", randCode)
	val.Add("rand", Randp)
	val.Add("_json_att", Json_att)
	val.Add("REPEAT_SUBMIT_TOKEN", order.RepeatSubmitToken)
	Info("checkRandCodeAnsyn params:", val)

	params, _ := url.QueryUnescape(val.Encode())
	body, err := DoForWardRequestHeader(order.CDN, "POST", URLCheckRandCodeAnsyn, strings.NewReader(params), nil)
	if err != nil {
		Error("checkRandCodeAnsyn DoForWardRequest error:", err)
		return false, []string{err.Error()}
	}
	Info("checkRandCodeAnsyn body:", body)

	crca := &JsonCheckRandCodeAnsyn{}
	if err := json.Unmarshal([]byte(body), &crca); err != nil {
		Error("checkRandCodeAnsyn json.Unmarshal:", err)
		return false, []string{err.Error()}
	}
	return crca.Data == "Y", crca.Messages
}

//提交订单
func (order *Order) getPassCodeNew() error {
	Info("getPassCodeNew")
	req, err := http.NewRequest("GET", URLPassCodeNewPassenger+strconv.Itoa(rand.Intn(99999)), nil)
	if err != nil {
		Error("getPassCodeNew http.NewRequest error:", err)
		return err
	}
	AddReqestHeader(req, "GET", map[string]string{"Accept": "image/webp,*/*;q=0.8"})
	con, err := NewForwardClientConn(order.CDN, req.URL.Scheme)
	if err != nil {
		Error("getPassCodeNew newForwardClientConn error:", err)
		return err
	}
	defer con.Close()
	resp, err := con.Do(req)
	if err != nil {
		Error("getPassCodeNew con.Do error:", err)
		return err
	}
	defer resp.Body.Close()

	img, s, err := image.Decode(resp.Body)
	Debug("image type:", s)
	if err != nil {
		Error("getPassCodeNew image.Decode:", err)
		return err
	}
	Im, err := walk.NewBitmapFromImage(img)
	if err != nil {
		Error("getPassCodeNew walk.NewBitmapFromImage:", err)
		return err
	}
	submitCaptchaImage.SetImage(Im)
	submitCaptchaEdit.SetText("")
	submitCaptchaEdit.SetFocus()
	go getPassengerDTO()
	return nil
}
func (order *Order) initDc() error {
	val := url.Values{}
	val.Add("_json_att", Json_att)
	Info("initDc params:", val)

	params, _ := url.QueryUnescape(val.Encode())

	body, err := DoForWardRequest(order.CDN, "POST", URLInitDc, strings.NewReader(params))
	if err != nil {
		Error("initDc DoForWardRequest error:", err)
		return err
	}
	Debug("initDc body:", body)

	if !strings.Contains(body, `var globalRepeatSubmitToken = '`) {
		return fmt.Errorf("获取网页出错！body:%v", body)
	}

	str := strings.Split(body, `var globalRepeatSubmitToken = '`)
	token := str[1][:32]
	Info(token)
	order.RepeatSubmitToken = token

	str2 := strings.Split(str[1], `'key_check_isChange':'`)
	key_check_isChange := str2[1][:56]
	Info(key_check_isChange)
	order.KeyCheckIsChange = key_check_isChange

	str3 := strings.Split(str2[1], `'tour_flag':'dc','train_location':'`)
	train_location := str3[1][:2]
	Info(train_location)
	order.TrainLocation = train_location

	return nil
}
func loginCheckBefore(cdn string) error {
	body, err := DoForWardRequest(cdn, "GET", URLLoginInit, nil)
	if err != nil {
		Error("loginCheckBefore DoForWardRequest error:", err)
		return err
	}
	Debug("loginCheckBefore body:", body)

	body, err = DoForWardRequest(cdn, "GET", URLLoginJs, nil)
	if err != nil {
		Error("loginCheckBefore DoForWardRequest error:", err)
		return err
	}
	Debug("loginCheckBefore body:", body)

	return nil
}
func loginCheck(cdn string) error {
	body, err := DoForWardRequest(cdn, "GET", URLInitQueryUserInfo, nil)
	if err != nil {
		Error("loginCheck DoForWardRequest error:", err)
		return err
	}
	Debug("loginCheck body:", body)

	body, err = DoForWardRequest(cdn, "POST", URLCheckUser, nil)
	if err != nil {
		Error("loginCheck DoForWardRequest error:", err)
		return err
	}
	Debug("loginCheck body:", body)

	body, err = DoForWardRequest(cdn, "GET", URLInit, nil)
	if err != nil {
		Error("loginCheck DoForWardRequest error:", err)
		return err
	}
	Debug("loginCheck body:", body)

	body, err = DoForWardRequest(cdn, "GET", URLQueryJs, nil)
	if err != nil {
		Error("loginCheck DoForWardRequest error:", err)
		return err
	}
	Debug("loginCheck body:", body)

	GetImage(cdn, false)

	return nil
}
func (order *Order) queryJs() error {
	body, err := DoForWardRequest(order.CDN, "GET", URLQueryJs, nil)
	if err != nil {
		Error("submitOrderRequest DoForWardRequest error:", err)
		return err
	}
	Debug("submitOrderRequest body:", body)

	return nil
}
func (order *Order) inits() error {
	body, err := DoForWardRequest(order.CDN, "GET", URLInit, nil)
	if err != nil {
		Error("submitOrderRequest DoForWardRequest error:", err)
		return err
	}
	Debug("submitOrderRequest body:", body)
	return nil
}
func (order *Order) submitOrderRequest() error {
	val := url.Values{}
	val.Add("secretStr", order.SecretStr)
	val.Add("train_date", order.TrainDate.Format("2006-01-02"))
	val.Add("back_train_date", order.TrainDate.Format("2006-01-02"))
	val.Add("tour_flag", Tour_flag)
	val.Add("purpose_codes", Purpose_codes)
	val.Add("query_from_station_name", order.Ticket.FromStationName)
	val.Add("query_to_station_name", order.Ticket.ToStationName)
	Info("submitOrderRequest params:", val)

	params, _ := url.QueryUnescape(val.Encode())

	body, err := DoForWardRequest(order.CDN, "POST", URLSubmitOrderRequest, strings.NewReader(params+"&undefined"))
	if err != nil {
		Error("submitOrderRequest DoForWardRequest error:", err)
		return err
	}
	Info("submitOrderRequest body:", body)

	if !strings.Contains(body, `"status":true`) {
		return errors.New("提交订票请求出错！")
	}
	sor := &JsonSubmitOrderRequest{}
	if err := json.Unmarshal([]byte(body), &sor); err != nil {
		Error("submitOrderRequest json.Unmarshal:", err)
		return err
	}
	if sor.HttpStatus == http.StatusOK {
		return nil
	}
	return fmt.Errorf("提交订票请求出错！body:%v", body)
}

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
						DoForWardRequest(order.CDN, "POST", URLCheckUser, nil)

						err = order.submitOrderRequest()
						if err != nil {
							Error(err)
							return
						}
						// order.inits()
						err = order.initDc()
						if err != nil {
							Error(err)
							return
						}
						order.queryJs()
						// go ticket.queryLeftTicket(order.CDN)
						go order.getPassCodeNew()
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

	go DoForWardRequest(cdn, "GET", LogQueryLeftTicketURL+leftTicketUrl, nil)

	h := HeaderMap[QueryLeftTicketURL]
	h["If-Modified-Since"] = time.Now().Local().Format(time.RFC1123Z)
	h["If-None-Match"] = strconv.FormatInt(time.Now().UnixNano(), 10)

	body, err := DoForWardRequestHeader(cdn, "GET", QueryLeftTicketURL+leftTicketUrl, nil, h)
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

//获取新的验证码图片
func GetImage(cdn string, setCookie bool) image.Image {
	req, err := http.NewRequest("GET", PassCodeNewURL, nil)
	if err != nil {
		Error("GetImage http.NewRequest error:", err)
		return nil
	}
	if !setCookie {
		AddReqestHeader(req, "GET", map[string]string{"Accept": "image/webp,*/*;q=0.8"})
		req.Header.Del("X-Requested-With")
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
	if setCookie {
		login.Cookie = GetCookieFromRespHeader(resp)
	}
	Info("==" + login.Cookie + "==")

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

}
