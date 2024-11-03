package models

type Path string

type Photo struct {
	ProductId string `json:"product_id"`
	Paths     []Path `json:"images"`
}
