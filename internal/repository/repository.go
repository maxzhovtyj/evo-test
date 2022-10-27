package repository

import (
	"evo-test/internal/models"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx"
	"strings"
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
			payeeBankMFO,
			payeeBankAccount,
			paymentNarrative)
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
		`, transactionTable,
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
			tr.PayeeBankMFO,
			tr.PayeeBankAccount,
			tr.PaymentNarrative,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

var transactionTableColumns = []string{
	"transactionId",
	"requestId",
	"terminalId",
	"partnerObjectId",
	"amountTotal",
	"amountOriginal",
	"commissionPS",
	"commissionClient",
	"commissionProvider",
	"dateInput",
	"datePost",
	"status",
	"paymentType",
	"paymentNumber",
	"serviceId",
	"service",
	"payeeId",
	"payeeName",
	"payeeBankMFO",
	"payeeBankAccount",
	"paymentNarrative",
}

func (r *repository) GetTransactions(params models.SearchParams) ([]models.Transaction, error) {
	var transactions []models.Transaction

	query := r.queryBuilder.
		Select(transactionTableColumns...).
		From(transactionTable)

	if params.TransactionId != 0 {
		query = query.Where(sq.Eq{"transactionId": params.TransactionId})
	}

	if len(params.TerminalIds) != 0 && params.TerminalIds[0] != "" {
		ids := strings.Join(params.TerminalIds, ", ")
		query = query.Where(fmt.Sprintf("terminalId IN (%s)", ids))
	}

	if params.Status != "" {
		query = query.Where(sq.Eq{"status": params.Status})
	}

	if params.PaymentNarrative != "" {
		narrativeLower := strings.ToLower(params.PaymentNarrative)
		query = query.Where(sq.Like{"LOWER(paymentNarrative)": "%" + narrativeLower + "%"})
	}

	if params.PaymentType != "" {
		query = query.Where(sq.Eq{"paymentType": params.PaymentType})
	}

	if params.DatePostFrom != "" {
		query = query.Where(sq.GtOrEq{"datePost": params.DatePostFrom})
	}

	if params.DatePostTo != "" {
		query = query.Where(sq.LtOrEq{"datePost": params.DatePostTo})
	}

	querySelectTransactionsSql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(querySelectTransactionsSql, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var tr models.Transaction

		err = rows.Scan(
			&tr.TransactionId,
			&tr.RequestId,
			&tr.TerminalId,
			&tr.PartnerObjectId,
			&tr.AmountTotal,
			&tr.AmountOriginal,
			&tr.CommissionPS,
			&tr.CommissionClient,
			&tr.CommissionProvider,
			&tr.DateInput,
			&tr.DatePost,
			&tr.Status,
			&tr.PaymentType,
			&tr.PaymentNumber,
			&tr.ServiceId,
			&tr.Service,
			&tr.PayeeId,
			&tr.PayeeName,
			&tr.PayeeBankMFO,
			&tr.PayeeBankAccount,
			&tr.PaymentNarrative,
		)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, tr)
	}

	return transactions, nil
}
