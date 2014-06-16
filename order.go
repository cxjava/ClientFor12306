package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type Order struct {
	CDN                string
	SeatType           string
	PassengerTicketStr string
	OldPassengerStr    string
	RepeatSubmitToken  string
	SecretStr          string
	TrainDate          string
	Ticket             Ticket
	SubmitCaptchaStr   chan string
	RandCode           string
	KeyCheckIsChange   string
	TrainLocation      string
}

func (order *Order) confirmSingleForQueue() error {
	val := url.Values{}
	val.Add("REPEAT_SUBMIT_TOKEN", order.RepeatSubmitToken)
	val.Add("_json_att", Json_att)
	val.Add("train_location", order.TrainLocation)
	val.Add("leftTicketStr", order.Ticket.YpInfo)
	val.Add("key_check_isChange", order.KeyCheckIsChange)
	val.Add("purpose_codes", Purpose_codes2)
	val.Add("randCode", order.RandCode)
	val.Add("oldPassengerStr", order.OldPassengerStr)
	val.Add("passengerTicketStr", order.PassengerTicketStr)

	Info("confirmSingleForQueue params:", val.Encode())
	h := map[string]string{
		"Accept":           "application/json, text/javascript, */*; q=0.01",
		"X-Requested-With": "XMLHttpRequest",
		"Referer":          "https://kyfw.12306.cn/otn/confirmPassenger/initDc",
	}
	body, err := DoForWardRequestHeader(order.CDN, "POST", URLConfirmSingleForQueue, strings.NewReader(val.Encode()), h)
	if err != nil {
		Error("confirmSingleForQueue DoForWardRequestHeader error:", err)
		return err
	}
	Info("confirmSingleForQueue body:", body)
	if strings.Contains(body, `"submitStatus":true`) {
		Info("提交订单成功 body:", body)
		fmt.Println("提交订单成功 body:", body)
	} else if strings.Contains(body, `订单未支付`) {
		Warn("订票成功！！")
		fmt.Println("订票成功！！")
	} else if strings.Contains(body, `用户未登录`) {
		Warn("用户未登录！！")
		fmt.Println("用户未登录！！")
	} else if strings.Contains(body, `取消次数过多`) {
		Warn("由于您取消次数过多！！")
		fmt.Println("由于您取消次数过多！！")
	} else if strings.Contains(body, `互联网售票实行实名制`) {
		Warn("貌似你已经购买了相同的车票！！")
		fmt.Println("貌似你已经购买了相同的车票！！")
	} else {
		Warn(order.CDN, "订票请求警告:", body)
		fmt.Println(order.CDN, "订票请求警告:", body)
	}
	return nil
}

func (order *Order) getQueueCount() error {
	val := url.Values{}
	val.Add("REPEAT_SUBMIT_TOKEN", order.RepeatSubmitToken)
	val.Add("_json_att", Json_att)
	val.Add("purpose_codes", Purpose_codes2)
	val.Add("leftTicket", order.Ticket.YpInfo)
	val.Add("toStationTelecode", order.Ticket.ToStationTelecode)
	val.Add("fromStationTelecode", order.Ticket.FromStationTelecode)
	val.Add("seatType", order.SeatType)
	val.Add("stationTrainCode", order.Ticket.StationTrainCode)
	val.Add("train_no", order.Ticket.TrainNo)
	val.Add("train_date", time.Now().Local().Format(`Mon Jun 16 2014 19:52:23 GMT+0800 (China Standard Time)`))
	// val.Add("train_date", `Wed Jun 25 00:00:00 UTC+0800 2014`)

	Info("getQueueCount params:", val.Encode())
	h := map[string]string{
		"Accept":           "application/json, text/javascript, */*; q=0.01",
		"X-Requested-With": "XMLHttpRequest",
		"Referer":          "https://kyfw.12306.cn/otn/confirmPassenger/initDc",
	}
	body, err := DoForWardRequestHeader(order.CDN, "POST", URLQueueCount, strings.NewReader(val.Encode()), h)
	if err != nil {
		Error("getQueueCount DoForWardRequestHeader error:", err)
		return err
	}
	Info("getQueueCount body:", body)
	return nil
}

func (order *Order) checkOrderInfo() error {
	//等待输入验证码
	submitCode := <-order.SubmitCaptchaStr
	order.RandCode = submitCode
	Info("验证码:", order.RandCode)

	val := url.Values{}
	val.Add("cancel_flag", Cancel_flag)
	val.Add("bed_level_order_num", Bed_level_order_num)
	val.Add("passengerTicketStr", order.PassengerTicketStr)
	val.Add("oldPassengerStr", order.OldPassengerStr)
	val.Add("tour_flag", Tour_flag)
	val.Add("randCode", order.RandCode)
	val.Add("_json_att", Json_att)
	val.Add("REPEAT_SUBMIT_TOKEN", order.RepeatSubmitToken)

	Info("checkOrderInfo params:", val.Encode())
	h := map[string]string{
		"X-Requested-With": "XMLHttpRequest",
		"Referer":          "https://kyfw.12306.cn/otn/confirmPassenger/initDc",
	}
	body, err := DoForWardRequestHeader(order.CDN, "POST", URLCheckOrderInfo, strings.NewReader(val.Encode()), h)
	if err != nil {
		Error("checkOrderInfo DoForWardRequestHeader error:", err)
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

	return nil
}

func (order *Order) checkRandCodeAnsyn(randCode string) (r bool, msg []string) {
	val := url.Values{}
	val.Add("REPEAT_SUBMIT_TOKEN", order.RepeatSubmitToken)
	val.Add("_json_att", Json_att)
	val.Add("rand", Randp)
	val.Add("randCode", randCode)

	Info("checkRandCodeAnsyn params:", val.Encode())
	h := map[string]string{
		"X-Requested-With": "XMLHttpRequest",
		"Referer":          "https://kyfw.12306.cn/otn/confirmPassenger/initDc",
	}
	body, err := DoForWardRequestHeader(order.CDN, "POST", URLCheckRandCodeAnsyn, strings.NewReader(val.Encode()), h)
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

//获取联系人
func (order *Order) GetPassengerDTO() {
	val := url.Values{}
	val.Add("_json_att", Json_att)
	val.Add("REPEAT_SUBMIT_TOKEN", order.RepeatSubmitToken)
	params := val.Encode()
	Info("getPassengerDTO params:", params)
	h := map[string]string{
		"X-Requested-With": "XMLHttpRequest",
		"Referer":          "https://kyfw.12306.cn/otn/confirmPassenger/initDc",
	}
	body, err := DoForWardRequestHeader(order.CDN, "POST", URLGetPassengerDTOs, strings.NewReader(params), h)
	if err != nil {
		Error("getPassengerDTO DoForWardRequest error:", err)
	}
	Debug("getPassengerDTO body:", body)
}

func (order *Order) initDc() error {
	h := map[string]string{
		"Accept":  "text/html, application/xhtml+xml, */*",
		"Referer": "https://kyfw.12306.cn/otn/leftTicket/init",
	}
	body, err := DoForWardRequestHeader(order.CDN, "POST", URLInitDc, strings.NewReader("_json_att="), h)
	if err != nil {
		Error("initDc DoForWardRequest error:", err)
		return err
	}
	fmt.Println("initDc body:", body)

	if !strings.Contains(body, `'key_check_isChange':'`) {
		return fmt.Errorf("获取网页出错！body:%v", body)
	}

	str := strings.Split(body, `var globalRepeatSubmitToken = '`)
	token := str[1][:32]
	fmt.Println("token:", token)
	Info("token:", token)
	order.RepeatSubmitToken = token

	str2 := strings.Split(str[1], `'key_check_isChange':'`)
	key_check_isChange := str2[1][:56]
	fmt.Println("key_check_isChange:", key_check_isChange)
	Info("key_check_isChange:", key_check_isChange)
	order.KeyCheckIsChange = key_check_isChange

	str3 := strings.Split(str2[1], `'tour_flag':'dc','train_location':'`)
	train_location := str3[1][:2]
	fmt.Println("train_location:", train_location)
	Info("train_location:", train_location)
	order.TrainLocation = train_location

	return nil
}
func (order *Order) submitOrderRequest() error {
	val := url.Values{}
	val.Add("secretStr", order.SecretStr)
	val.Add("train_date", order.TrainDate)
	val.Add("back_train_date", order.TrainDate)
	val.Add("tour_flag", Tour_flag)
	val.Add("purpose_codes", Purpose_codes)
	val.Add("query_from_station_name", order.Ticket.FromStationName)
	val.Add("query_to_station_name", order.Ticket.ToStationName)

	params, _ := url.QueryUnescape(val.Encode())
	params = params + "&undefined"
	Info("submitOrderRequest params:", params)
	h := map[string]string{
		"X-Requested-With": "XMLHttpRequest",
		"Referer":          "https://kyfw.12306.cn/otn/leftTicket/init",
	}
	body, err := DoForWardRequestHeader(order.CDN, "POST", URLSubmitOrderRequest, strings.NewReader(params), h)
	if err != nil {
		Error("submitOrderRequest DoForWardRequestHeader error:", err)
		return err
	}
	Info("submitOrderRequest body:", body)

	if !strings.Contains(body, `"status":true`) {
		return errors.New("提交订票请求出错！返回信息:" + body)
	}
	sor := &JsonSubmitOrderRequest{}
	if err := json.Unmarshal([]byte(body), &sor); err != nil {
		Error("submitOrderRequest json.Unmarshal:", err)
		return err
	}
	if sor.HttpStatus == http.StatusOK {
		return nil
	}
	return fmt.Errorf("提交订票请求出错！返回信息:%v", body)
}
