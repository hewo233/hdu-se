package models

type Report struct {
	Code   int         `json:"code"`
	Result interface{} `json:"result"`
}
