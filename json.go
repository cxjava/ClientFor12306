package main

type Basic struct {
	ValidateMessagesShowId string `json:"validateMessagesShowId"`
	Status                 bool   `json:"status"`
	HttpStatus             int    `json:"httpstatus"`
}

type JsonSubmitOrderRequest struct {
	Basic
	Messages         []interface{} `json:"messages,omitempty"`
	ValidateMessages interface{}   `json:"validateMessages,omitempty"`
}

type JsonCheckOrderInfo struct {
	Basic
	Data struct {
		SubmitStatus bool `json:"submitStatus"`
	}
	Messages         []interface{} `json:"messages,omitempty"`
	ValidateMessages interface{}   `json:"validateMessages,omitempty"`
}

type JsonGetQueueCount struct {
	Basic
	Data struct {
		Count  string `json:"count"`
		Ticket string `json:"ticket"`
		OP2    string `json:"op_2"`
		CountT string `json:"countT"`
		OP1    string `json:"op_1"`
	}
	Messages         []interface{} `json:"messages,omitempty"`
	ValidateMessages interface{}   `json:"validateMessages,omitempty"`
}

type JsonConfirmSingleForQueue struct {
	Basic
	Data struct {
		SubmitStatus string `json:"submitStatus"`
	}
	Messages         []interface{} `json:"messages,omitempty"`
	ValidateMessages interface{}   `json:"validateMessages,omitempty"`
}

type JsonCheckRandCodeAnsyn struct {
	Basic
	Data             string      `json:"data"`
	Messages         []string    `json:"messages,omitempty"`
	ValidateMessages interface{} `json:"validateMessages,omitempty"`
}

/////////////////////
type QueryLeftNewDTO struct {
	Basic
	Data             []LeftTicket  `json:"data"`
	Messages         []interface{} `json:"messages,omitempty"`
	ValidateMessages interface{}   `json:"validateMessages,omitempty"`
}

type LeftTicket struct {
	Ticket         Ticket `json:"queryLeftNewDTO"`
	SecretStr      string `json:"secretStr"`
	ButtonTextInfo string `json:"buttonTextInfo"`
}

type Ticket struct {
	TrainNo              string `json:"train_no"`               //"560000K52960",
	StationTrainCode     string `json:"station_train_code"`     // "K532",
	StartStationTelecode string `json:"start_station_telecode"` //"HZH",
	StartStationName     string `json:"start_station_name"`     //"杭州",
	EndStationTelecode   string `json:"end_station_telecode"`   //"ICW",
	EndStationName       string `json:"end_station_name"`       //"成都东",
	FromStationTelecode  string `json:"from_station_telecode"`  // "WCN",
	FromStationName      string `json:"from_station_name"`      //"武昌",
	ToStationTelecode    string `json:"to_station_telecode"`    // "JCN",
	ToStationName        string `json:"to_station_name"`        //"京山",
	StartTime            string `json:"start_time"`             //"01:02",
	ArriveTime           string `json:"arrive_time"`            //"03:00",
	DayDifference        string `json:"day_difference"`         //"0",
	TrainClassName       string `json:"train_class_name"`       //"",
	Lishi                string `json:"lishi"`                  //"01:58",
	CanWebBuy            string `json:"canWebBuy"`              // "N",
	LishiValue           string `json:"lishiValue"`             //"118",
	YpInfo               string `json:"yp_info"`                //"1002353000401115000010023500003007450000",
	ControlTrianDay      string `json:"control_train_day"`      //"20991231",
	StartTrainDate       string `json:"start_train_date"`       //"20140119",
	SeatFeature          string `json:"seat_feature"`           //"W3431333",
	YpEx                 string `json:"yp_ex"`                  //"10401030",
	TrainSeatFeature     string `json:"train_seat_feature"`     //"3",
	SeatTypes            string `json:"seat_types"`             //"1413",
	LocationCode         string `json:"location_code"`          //"H3",
	FromStationNo        string `json:"from_station_no"`        //"13",
	ToStationNo          string `json:"to_station_no"`          //"15",
	ControlDay           int    `json:"control_day"`            //19,
	SaleTime             string `json:"sale_time"`              //"0800",
	IsSupportCard        string `json:"is_support_card"`        // "0",
	GGNum                string `json:"gg_num"`                 //
	GaoJiRuanWoNum       string `json:"gr_num"`                 //高级软卧
	QiTaNum              string `json:"qt_num"`                 //其他
	RuanWoNum            string `json:"rw_num"`                 //软卧
	RuanZuoNum           string `json:"rz_num"`                 //软座
	TeDengZuoNum         string `json:"tz_num"`                 //特等座
	WuZuoNum             string `json:"wz_num"`                 //无座
	YBNum                string `json:"yb_num"`                 //
	YingWoNum            string `json:"yw_num"`                 //硬卧
	YingZuoNum           string `json:"yz_num"`                 //硬座
	ErDengZuoNum         string `json:"ze_num"`                 //二等座
	YiDengZuoNum         string `json:"zy_num"`                 //一等座
	ShangWuZuoNum        string `json:"swz_num"`                //商务座
}

/////////////////////

type PassengerDTO struct {
	Basic
	Data             Data4Passenger
	Messages         []interface{} `json:"messages,omitempty"`
	ValidateMessages interface{}   `json:"validateMessages,omitempty"`
}

type Data4Passenger struct {
	IsExist          bool          `json:"isExist"`
	ExMsg            string        `json:"exMsg"`
	NoLogin          string        `json:"noLogin,omitempty"`
	NormalPassengers []Passenger   `json:"normal_passengers,omitempty"`
	DjPassengers     []interface{} `json:"dj_passengers,omitempty"`
}

type Passenger struct {
	Code                string `json:"code"`
	PassengerName       string `json:"passenger_name"`
	SexCode             string `json:"sex_code"`
	SexName             string `json:"sex_name"`
	BornDate            string `json:"born_date"`
	CountryCode         string `json:"country_code"`
	PassengerIdTypeCode string `json:"passenger_id_type_code"`
	PassengerIdTypeName string `json:"passenger_id_type_name"`
	PassengerIdNo       string `json:"passenger_id_no"`
	PassengerType       string `json:"passenger_type"`
	PassengerFlag       string `json:"passenger_flag"`
	PassengerTypeName   string `json:"passenger_type_name"`
	Mobile              string `json:"mobile_no"`
	Phone               string `json:"phone_no"`
	Email               string `json:"email"`
	Address             string `json:"address"`
	Postalcode          string `json:"postalcode"`
	FirstLetter         string `json:"first_letter"`
	RecordCount         string `json:"recordCount"`
}

//////////////
