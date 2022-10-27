package models

import "time"

type Transaction struct {
	TransactionId      int
	RequestId          int
	TerminalId         int
	PartnerObjectId    int
	AmountTotal        float64
	AmountOriginal     float64
	CommissionPS       float64
	CommissionClient   float64
	CommissionProvider float64
	DateInput          time.Time
	DatePost           time.Time
	Status             string
	PaymentType        string
	PaymentNumber      string
	ServiceId          int
	Service            string
	PayeeId            int
	PayeeName          string
	PayeeBankMfo       int
	PayeeBankAccount   string
	PaymentNarrative   string
}
