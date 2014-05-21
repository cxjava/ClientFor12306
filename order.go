package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/lxn/walk"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

var (
	order *Order
)

type Order struct {
	CDN                string
	SeatType           string
	PassengerTicketStr string
	OldPassengerStr    string
	RepeatSubmitToken  string
	SecretStr          string
	TrainDate          time.Time
	Ticket             Ticket
	SubmitCaptchaStr   chan string
	RandCode           string
	KeyCheckIsChange   string
	TrainLocation      string
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

	params, _ := url.QueryUnescape(val.Encode())
	Info("confirmSingleForQueue params:", params)
	h := map[string]string{"Accept": "application/json, text/javascript, */*; q=0.01"}
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
	val := url.Values{}
	val.Add("train_date", order.TrainDate.Local().Format(`Mon+Jan+02+2006+15%3A04%3A05+GMT%2B0700+(China+Standard+Time)`))
	// val.Add("train_date", `Thu Jun 05 2014 00:00:00 GMT 0800 (China Standard Time)`)
	val.Add("train_no", order.Ticket.TrainNo)
	val.Add("stationTrainCode", order.Ticket.StationTrainCode)
	val.Add("seatType", order.SeatType)
	val.Add("fromStationTelecode", order.Ticket.FromStationTelecode)
	val.Add("toStationTelecode", order.Ticket.ToStationTelecode)
	val.Add("leftTicket", order.Ticket.YpInfo)
	val.Add("purpose_codes", Purpose_codes)
	val.Add("_json_att", Json_att)
	val.Add("REPEAT_SUBMIT_TOKEN", order.RepeatSubmitToken)
	params, _ := url.QueryUnescape(val.Encode())
	Info("getQueueCount params:", params)
	h := map[string]string{"Accept": "application/json, text/javascript, */*; q=0.01"}
	body, err := DoForWardRequestHeader(order.CDN, "POST", URLQueueCount, strings.NewReader(params), h)
	if err != nil {
		Error("getQueueCount DoForWardRequest error:", err)
		return err
	}
	Info("getQueueCount body:", body)

	// time.Sleep(time.Second * 8)
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

	params, _ := url.QueryUnescape(val.Encode())
	Info("checkOrderInfo params:", params)
	h := map[string]string{"Accept": "application/json, text/javascript, */*; q=0.01"}
	body, err := DoForWardRequestHeader(order.CDN, "POST", URLCheckOrderInfo, strings.NewReader(params), h)
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
	// time.Sleep(time.Second * 8)
	return order.getQueueCount()
}

func (order *Order) checkRandCodeAnsyn(randCode string) (r bool, msg []string) {
	val := url.Values{}
	val.Add("randCode", randCode)
	val.Add("rand", Randp)
	val.Add("_json_att", Json_att)
	val.Add("REPEAT_SUBMIT_TOKEN", order.RepeatSubmitToken)

	params, _ := url.QueryUnescape(val.Encode())
	Info("checkRandCodeAnsyn params:", params)
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

//获取验证码
func (order *Order) getPassCodeNew() error {
	Info("getPassCodeNew")
	req, err := http.NewRequest("GET", URLPassCodeNewPassenger, nil)
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
	return nil
}

func (order *Order) initDc() error {
	val := url.Values{}
	val.Add("_json_att", Json_att)
	Info("initDc params:", val)

	params, _ := url.QueryUnescape(val.Encode())
	h := map[string]string{"Referer": "https://kyfw.12306.cn/otn/leftTicket/init"}
	body, err := DoForWardRequestHeader(order.CDN, "POST", URLInitDc, strings.NewReader(params), h)
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

func (order *Order) submitOrderRequest() error {
	val := url.Values{}
	val.Add("secretStr", order.SecretStr)
	val.Add("train_date", order.TrainDate.Format("2006-01-02"))
	val.Add("back_train_date", order.TrainDate.Format("2006-01-02"))
	val.Add("tour_flag", Tour_flag)
	val.Add("purpose_codes", Purpose_codes)
	val.Add("query_from_station_name", order.Ticket.FromStationName)
	val.Add("query_to_station_name", order.Ticket.ToStationName)

	params, _ := url.QueryUnescape(val.Encode())
	Info("submitOrderRequest params:", params)
	h := map[string]string{"Referer": "https://kyfw.12306.cn/otn/leftTicket/init"}
	body, err := DoForWardRequestHeader(order.CDN, "POST", URLSubmitOrderRequest, strings.NewReader(params+"&undefined"), h)
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
