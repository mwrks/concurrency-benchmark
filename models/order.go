package models

import "time"

type Order struct {
	OrderID   int       `json:"order_id"`
	TicketID  int       `json:"ticket_id"`
	OrderedBy string    `json:"ordered_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
