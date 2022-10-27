package handler

import (
	"encoding/csv"
	_ "evo-test/docs"
	"evo-test/internal/apperror"
	"evo-test/internal/models"
	"evo-test/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"       // swagger embed files
	"github.com/swaggo/gin-swagger" // gin-swagger middleware
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func newErrorResponse(ctx *gin.Context, statusCode int, message string) {
	ctx.AbortWithStatusJSON(statusCode, apperror.Error{Message: message})
}

type Handler interface {
	Register() *gin.Engine
}

type handler struct {
	service service.Service
}

func New(service service.Service) Handler {
	return &handler{service: service}
}

func (h *handler) Register() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		api.POST("/load-csv", h.ParseCSVFile)
		api.GET("/transaction", h.GetTransaction)
	}

	return router
}

func (h *handler) ParseCSVFile(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "form/json")
	err := ctx.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	files, ok := ctx.Request.MultipartForm.File["file"]
	if len(files) == 0 {
		if !ok {
			newErrorResponse(ctx, http.StatusBadRequest, "something wrong with file you provided")
			return
		} else {
			newErrorResponse(ctx, http.StatusBadRequest, "file not provided")
			return
		}
	}

	fileInfo := files[0]
	fileReader, err := fileInfo.Open()
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var parsedData []models.Transaction
	r := csv.NewReader(fileReader)

	skipFirst := 0
	for {
		if skipFirst == 0 {
			skipFirst++
			continue
		}

		read, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			continue
		}

		transactionId, err := strconv.Atoi(read[0])
		requestId, err := strconv.Atoi(read[1])
		terminalId, err := strconv.Atoi(read[2])
		partnerObjectId, err := strconv.Atoi(read[3])
		amountTotal, err := strconv.ParseFloat(read[4], 32)
		amountOriginal, err := strconv.ParseFloat(read[5], 32)
		commissionPS, err := strconv.ParseFloat(read[6], 32)
		commissionClient, err := strconv.ParseFloat(read[7], 32)
		commissionProvider, err := strconv.ParseFloat(read[8], 32)
		dateInput, err := time.Parse("2006-01-02 15:04:05", read[9])
		datePost, err := time.Parse("2006-01-02 15:04:05", read[10])
		serviceId, err := strconv.Atoi(read[14])
		payeeId, err := strconv.Atoi(read[16])
		payeeBankMfo, err := strconv.Atoi(read[18])

		if err != nil {
			newErrorResponse(ctx, http.StatusBadRequest, err.Error())
			return
		}

		parsedData = append(parsedData, models.Transaction{
			TransactionId:      transactionId,
			RequestId:          requestId,
			TerminalId:         terminalId,
			PartnerObjectId:    partnerObjectId,
			AmountTotal:        float32(amountTotal),
			AmountOriginal:     float32(amountOriginal),
			CommissionPS:       float32(commissionPS),
			CommissionClient:   float32(commissionClient),
			CommissionProvider: float32(commissionProvider),
			DateInput:          dateInput,
			DatePost:           datePost,
			Status:             read[11],
			PaymentType:        read[12],
			PaymentNumber:      read[13],
			ServiceId:          serviceId,
			Service:            read[15],
			PayeeId:            payeeId,
			PayeeName:          read[17],
			PayeeBankMfo:       payeeBankMfo,
			PayeeBankAccount:   read[19],
			PaymentNarrative:   read[20],
		})
	}

	err = h.service.LoadData(parsedData)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, "data successfully loaded")
}

// GetTransaction godoc
// @Summary      Gets a transaction
// @Description  get transaction
// @Tags         api
// @Accept       json
// @Produce      json
// @Param		 transactionId query int false "Transaction id"
// @Param 		 terminalIds query string false "Array with terminal ids, example: '1,2,3,4,5'"
// @Param		 status query string false "Transaction status"
// @Param		 paymentType query string false "Transaction payment type"
// @Param		 datePost query string false "Transaction date post interval, example: '<fromTime>:<toTime>'"
// @Param		 paymentNarrative query string false "Transaction payment narrative"
// @Success      200  {array}  	models.Transaction
// @Failure      400  {object}  object
// @Failure      500  {object}  object
// @Router       /api/transaction [get]
func (h *handler) GetTransaction(ctx *gin.Context) {
	var err error

	var transactionIdInt int
	transactionId := ctx.Query("transactionId")
	if transactionId != "" {
		transactionIdInt, err = strconv.Atoi(transactionId)
		if err != nil {
			newErrorResponse(ctx, http.StatusBadRequest, "invalid transaction id")
			return
		}
	}

	terminalIds := strings.Split(ctx.Query("terminalIds"), ",")
	if terminalIds[0] != "" {
		for _, id := range terminalIds {
			_, err = strconv.Atoi(id)
			if err != nil {
				newErrorResponse(ctx, http.StatusBadRequest, "invalid terminal ids")
				return
			}
		}
	}

	transactions, err := h.service.GetTransactions(models.SearchParams{
		TransactionId:    transactionIdInt,
		TerminalIds:      terminalIds,
		Status:           ctx.Query("status"),
		PaymentType:      ctx.Query("paymentType"),
		DatePostFrom:     ctx.Query("datePostFrom"),
		DatePostTo:       ctx.Query("datePostTo"),
		PaymentNarrative: ctx.Query("paymentNarrative"),
	})

	ctx.JSON(http.StatusOK, transactions)
}
