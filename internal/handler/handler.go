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
	//csvFile, err := os.Open("data.csv")
	//if err != nil {
	//	newErrorResponse(ctx, http.StatusInternalServerError, "failed to open csv file")
	//	return
	//}

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
		if err != nil {
			continue
		}

		requestId, err := strconv.Atoi(read[1])
		if err != nil {
			continue
		}

		terminalId, err := strconv.Atoi(read[2])
		if err != nil {
			continue
		}

		partnerObjectId, err := strconv.Atoi(read[3])
		if err != nil {
			continue
		}

		amountTotal, err := strconv.ParseFloat(read[4], 64)
		if err != nil {
			continue
		}

		amountOriginal, err := strconv.ParseFloat(read[5], 64)
		if err != nil {
			return
		}

		commissionPS, err := strconv.ParseFloat(read[6], 64)
		if err != nil {
			return
		}

		commissionClient, err := strconv.ParseFloat(read[7], 64)
		if err != nil {
			return
		}

		commissionProvider, err := strconv.ParseFloat(read[8], 64)
		if err != nil {
			return
		}

		dateInput, err := time.Parse("2006-01-02 15:04:05", read[9])
		if err != nil {
			return
		}

		datePost, err := time.Parse("2006-01-02 15:04:05", read[10])
		if err != nil {
			return
		}

		serviceId, err := strconv.Atoi(read[14])
		if err != nil {
			return
		}

		payeeId, err := strconv.Atoi(read[16])
		if err != nil {
			return
		}

		payeeBankMfo, err := strconv.Atoi(read[18])
		if err != nil {
			return
		}

		parsedData = append(parsedData, models.Transaction{
			TransactionId:      transactionId,
			RequestId:          requestId,
			TerminalId:         terminalId,
			PartnerObjectId:    partnerObjectId,
			AmountTotal:        amountTotal,
			AmountOriginal:     amountOriginal,
			CommissionPS:       commissionPS,
			CommissionClient:   commissionClient,
			CommissionProvider: commissionProvider,
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
			PayeeBankAccount:   read[18],
			PaymentNarrative:   read[19],
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
// @Failure      400  {string}  string
// @Failure      500  {object}  object
// @Router       /api/transaction [get]
func (h *handler) GetTransaction(ctx *gin.Context) {
	transactionId, err := strconv.Atoi(ctx.Query("transactionId"))
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, "invalid transaction id")
		return
	}
	var terminalIdsInt []int
	terminalIds := strings.Split(ctx.Query("terminalIds"), ",")
	for _, id := range terminalIds {
		idInt, err := strconv.Atoi(id)
		if err != nil {
			return
		}

		terminalIdsInt = append(terminalIdsInt, idInt)
	}
	status := ctx.Query("status")
	paymentType := ctx.Query("paymentType")

	var datePostFrom time.Time
	var datePostTo time.Time

	datePost := strings.Split(ctx.Query("datePost"), ":")

	if len(datePost) >= 1 {
		datePostFrom, err = time.Parse("2006-01-02T15:04:05Z", datePost[0])
		if err != nil {
			newErrorResponse(ctx, http.StatusBadRequest, "invalid date post from time")
			return
		}
		if len(datePost) == 2 {
			datePostTo, err = time.Parse("2006-01-02T15:04:05Z", datePost[1])
			if err != nil {
				newErrorResponse(ctx, http.StatusBadRequest, "invalid date post from time")
				return
			}
		}
	}

	paymentNarrative := ctx.Query("paymentNarrative")

	transactions, err := h.service.GetTransaction(models.SearchParams{
		TransactionId:    transactionId,
		TerminalIds:      terminalIdsInt,
		Status:           status,
		PaymentType:      paymentType,
		DatePostFrom:     datePostFrom,
		DatePosTo:        datePostTo,
		PaymentNarrative: paymentNarrative,
	})

	ctx.JSON(http.StatusOK, transactions)
}
