package models

import "time"

type SearchParams struct {
	TransactionId    int
	TerminalIds      []int
	Status           string
	PaymentType      string
	DatePostFrom     time.Time
	DatePosTo        time.Time
	PaymentNarrative string
}
