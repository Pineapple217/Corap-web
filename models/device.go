package models

type Device struct {
	Deveui     string `json:"deveui"`
	Name       string `json:"name"`
	Hashedname string `json:"hashedname"`
	IsDefect   bool   `json:"is_defect"`
}
