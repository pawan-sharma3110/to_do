package models

type ToDo struct {
	ID       int    `json:"id"`
	Complete bool   `json:"complete"`
	Body     string `json:"body"`
}
