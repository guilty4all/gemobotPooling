package main

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}

type Chat struct {
	ChatId int `json:"id"`
}

type RestResponse struct {
	Result []Update `json:"result"`
}

type BotMessage struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
}

type ResearchType struct {
	RNum   string `json:"r_num"`
	RCode  string `json:"r_code"`
	RName  string `json:"r_name"`
	RTime  string `json:"r_time"`
	RPrice string `json:"r_price"`
}
