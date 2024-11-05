package models

type Socials struct {
	Id         int    `json:"id"`
	CustomerId int    `json:"customer_id"`
	Subscribe  bool   `json:"subscribe"`
	Socials    string `json:"socials"`
}
