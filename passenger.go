package main

type Basic struct {
	ValidateMessagesShowId string `json:"validateMessagesShowId"`
	Status                 bool   `json:"status"`
	HttpStatus             int    `json:"httpstatus"`
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
