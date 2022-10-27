package models

import "time"

type Transaction struct {
	TransactionId      int
	RequestId          int
	TerminalId         int
	PartnerObjectId    int
	AmountTotal        float32
	AmountOriginal     float32
	CommissionPS       float32
	CommissionClient   float32
	CommissionProvider float32
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
