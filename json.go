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
