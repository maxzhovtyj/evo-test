package repository

import (
	"context"
	"evo-test/internal/models"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"strings"
)

const (
	transactionTable   = "transaction"
	transactionId      = "transaction_id"
	requestId          = "request_id"
	terminalId         = "terminal_id"
	partnerObjectId    = "partner_object_id"
	amountTotal        = "amount_total"
	amountOriginal     = "amount_original"
	commissionPs       = "commission_ps"
	commissionClient   = "commission_client"
	commissionProvider = "commission_provider"
	dateInput          = "date_input"
	datePost           = "date_post"
	status             = "status"
	paymentType        = "payment_type"
	paymentNumber      = "payment_number"
	serviceId          = "service_id"
	service            = "service"
	payeeId            = "payee_id"
	payeeName          = "payee_name"
	payeeBankMfo       = "payee_bank_mfo"
	payeeBankAccount   = "payee_bank_account"
	paymentNarrative   = "payment_narrative"
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

var transactionInsertColumns = []string{
	transactionTable,
	transactionId,
	requestId,
	terminalId,
	partnerObjectId,
	amountTotal,
	amountOriginal,
	commissionPs,
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
	payeeBankMfo,
	payeeBankAccount,
	paymentNarrative,
}

func (r *repository) InsertData(transactions []models.Transaction) error {
	for _, tr := range transactions {
		queryInsertTransaction, args, err := r.queryBuilder.
			Insert(strings.Join(transactionInsertColumns, ", ")).
			Into(transactionTable).
			Values(
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
			).ToSql()
		if err != nil {
			return err
		}

		_, err = r.db.Exec(context.Background(), queryInsertTransaction, args...)
		if err != nil {
			return err
		}
	}

	return nil
}

var transactionTableColumns = []string{
	transactionId,
	requestId,
	terminalId,
	partnerObjectId,
	amountTotal,
	amountOriginal,
	commissionPs,
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
	payeeBankMfo,
	payeeBankAccount,
	paymentNarrative,
}

func (r *repository) GetTransactions(params models.SearchParams) ([]models.Transaction, error) {
	var transactions []models.Transaction

	query := r.queryBuilder.
		Select(transactionTableColumns...).
		From(transactionTable)

	if params.TransactionId != 0 {
		query = query.Where(sq.Eq{transactionId: params.TransactionId})
	}

	if len(params.TerminalIds) != 0 && params.TerminalIds[0] != "" {
		ids := strings.Join(params.TerminalIds, ", ")
		IdsInRangeStr := fmt.Sprintf("%s IN (%s)", terminalId, ids)
		query = query.Where(IdsInRangeStr)
	}

	if params.Status != "" {
		query = query.Where(sq.Eq{status: params.Status})
	}

	if params.PaymentNarrative != "" {
		narrativeLower := strings.ToLower(params.PaymentNarrative)

		narrativeKey := fmt.Sprintf("LOWER(%s)", paymentNarrative)
		narrativeValue := fmt.Sprintf("%%%s%%", narrativeLower)

		query = query.Where(sq.Like{narrativeKey: narrativeValue})
	}

	if params.PaymentType != "" {
		query = query.Where(sq.Eq{paymentType: params.PaymentType})
	}

	if params.DatePostFrom != "" {
		query = query.Where(sq.GtOrEq{datePost: params.DatePostFrom})
	}

	if params.DatePostTo != "" {
		query = query.Where(sq.LtOrEq{datePost: params.DatePostTo})
	}

	querySelectTransactionsSql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(context.Background(), querySelectTransactionsSql, args...)
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
