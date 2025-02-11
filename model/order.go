package model

import "time"

type Order struct {
	ID        int
	UserID    int
	Amount    float64
	OrderDate time.Time
}
