package model

type Note struct {
	ID     int     `json:"id"`
	Course string  `json:"course"`
	Value  float64 `json:"value"`
}
