package models

type SearchParams struct {
	TransactionId    int
	TerminalIds      []string
	Status           string
	PaymentType      string
	DatePostFrom     string
	DatePostTo       string
	PaymentNarrative string
}
