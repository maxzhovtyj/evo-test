package handler

import (
	_ "evo-test/docs"
	"evo-test/internal/apperror"
	"evo-test/internal/models"
	"evo-test/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"       // swagger embed files
	"github.com/swaggo/gin-swagger" // gin-swagger middleware
	"net/http"
)

const (
	paramsCtx      = "params"
	parsedDataCtx  = "parsedData"
	swaggerUrl     = "/swagger/*any"
	apiUrl         = "/api"
	loadDataUrl    = "/load-data"
	transactionUrl = "/transaction"
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

	router.GET(swaggerUrl, ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group(apiUrl)
	{
		api.POST(loadDataUrl, h.ParseFileData(h.LoadData))
		api.GET(transactionUrl, h.TransactionQueryParams(h.GetTransaction))
	}

	return router
}

// LoadData godoc
// @Summary      Load data to database from a csv file
// @Description  load data
// @Tags         api
// @Accept       multipart/form-data
// @Produce      json
// @Param		 file formData file true "CSV file with data"
// @Success      201  {string}  string
// @Failure      400  {object}  apperror.Error
// @Failure      500  {object}  apperror.Error
// @Router       /api/load-data [post]
func (h *handler) LoadData(ctx *gin.Context) {
	parsedData, exists := ctx.Get(parsedDataCtx)
	if !exists {
		newErrorResponse(ctx, http.StatusBadRequest, "data wasn't found")
		return
	}

	parsedDataTr := parsedData.([]models.Transaction)

	err := h.service.LoadData(parsedDataTr)
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
// @Param 		 terminalIds query string false "Array with terminal ids, example: 3506,3507"
// @Param		 status query string false "Transaction status"
// @Param		 paymentType query string false "Transaction payment type"
// @Param		 datePostFrom query string false "Transaction min date post"
// @Param		 datePostTo query string false "Transaction max date post"
// @Param		 paymentNarrative query string false "Transaction payment narrative"
// @Success      200  {array}  	models.Transaction
// @Failure      400  {object}  apperror.Error
// @Failure      500  {object}  apperror.Error
// @Router       /api/transaction [get]
func (h *handler) GetTransaction(ctx *gin.Context) {
	paramsStr, exists := ctx.Get(paramsCtx)
	if !exists {
		newErrorResponse(ctx, http.StatusBadRequest, "query params wasn't found")
		return
	}

	params := paramsStr.(models.SearchParams)

	transactions, err := h.service.GetTransactions(params)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, transactions)
}
