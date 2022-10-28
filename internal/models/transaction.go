package models

import "time"

type Transaction struct {
	TransactionId      int       `example:"1"`
	RequestId          int       `example:"20020"`
	TerminalId         int       `example:"3506"`
	PartnerObjectId    int       `example:"1111"`
	AmountTotal        float32   `example:"1.00"`
	AmountOriginal     float32   `example:"1.00"`
	CommissionPS       float32   `example:"0.00"`
	CommissionClient   float32   `example:"0.00"`
	CommissionProvider float32   `example:"0.00"`
	DateInput          time.Time `example:"2022-08-12T11:25:27Z"`
	DatePost           time.Time `example:"2022-08-12T14:25:27Z"`
	Status             string    `example:"accepted"`
	PaymentType        string    `example:"cash"`
	PaymentNumber      string    `example:"PS16698205"`
	ServiceId          int       `example:"13980"`
	Service            string    `example:"Поповнення карток"`
	PayeeId            int       `example:"14232155"`
	PayeeName          string    `example:"pumb"`
	PayeeBankMFO       string    `example:"254751"`
	PayeeBankAccount   string    `example:"UA713451373919523"`
	PaymentNarrative   string    `example:"Перерахування коштів згідно договору про надання послуг А11/27122 від 19.11.2020 р."`
}
