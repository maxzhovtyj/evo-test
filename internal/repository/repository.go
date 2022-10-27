package repository

import (
	"evo-test/internal/models"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx"
	"log"
)

const (
	transactionTable = "transaction"
)

type Repository interface {
	InsertData([]models.Transaction) error
	GetTransactions(params models.SearchParams) ([]models.Transaction, error)
}

type repository struct {
	db           *pgx.Conn
	queryBuilder sq.StatementBuilderType
}

func New(conn *pgx.Conn, qb sq.StatementBuilderType) Repository {
	return &repository{
		db:           conn,
		queryBuilder: qb,
	}
}

func (r *repository) InsertData(transactions []models.Transaction) error {
	queryInsertTransaction := fmt.Sprintf(
		`
		INSERT INTO %s
			(transactionId, 
			requestId,
			terminalId,
			partnerObjectId,
			amountTotal,
			amountOriginal,
			commissionPS,
			commissionClient,
			commissionProvider,
			dateInput,
			datePost,
			status,
			paymentType,
			paymentNumber,
			serviceId,
			service,
			payeeId,
			payeeName,
			payeeBankAccount,
			paymentNarrative)
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
		`,
		transactionTable,
	)
	for _, tr := range transactions {
		_, err := r.db.Exec(queryInsertTransaction,
			tr.TransactionId,
			tr.RequestId,
			tr.TerminalId,
			tr.PartnerObjectId,
			tr.AmountTotal,
			tr.AmountOriginal,
			tr.CommissionPS,
			tr.CommissionClient,
			tr.CommissionProvider,
			tr.DateInput,
			tr.DatePost,
			tr.Status,
			tr.PaymentType,
			tr.PaymentNumber,
			tr.ServiceId,
			tr.Service,
			tr.PayeeId,
			tr.PayeeName,
			tr.PayeeBankAccount,
			tr.PaymentNarrative,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) GetTransactions(params models.SearchParams) ([]models.Transaction, error) {
	log.Println(params)
	return nil, nil
}
