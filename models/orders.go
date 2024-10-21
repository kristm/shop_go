package models

type Status int

const (
	Pending Status = iota
	Cancelled
	Paid
)
