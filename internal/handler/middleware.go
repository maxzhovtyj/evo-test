package handler

import (
	"encoding/csv"
	"evo-test/internal/models"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (h *handler) ParseFileData(handlerFunc func(ctx *gin.Context)) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
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
			var tr models.Transaction
			read, err := r.Read()

			if skipFirst == 0 {
				skipFirst++
				continue
			}

			if err == io.EOF {
				break
			} else if err != nil {
				newErrorResponse(ctx, http.StatusBadRequest, err.Error())
				return
			}

			tr.TransactionId, err = strconv.Atoi(read[0])
			if err != nil {
				newErrorResponse(ctx, http.StatusBadRequest, err.Error())
				return
			}

			tr.RequestId, err = strconv.Atoi(read[1])
			if err != nil {
				newErrorResponse(ctx, http.StatusBadRequest, err.Error())
				return
			}

			tr.TerminalId, err = strconv.Atoi(read[2])
			if err != nil {
				newErrorResponse(ctx, http.StatusBadRequest, err.Error())
				return
			}

			tr.PartnerObjectId, err = strconv.Atoi(read[3])
			if err != nil {
				newErrorResponse(ctx, http.StatusBadRequest, err.Error())
				return
			}

			amountTotal, err := strconv.ParseFloat(read[4], 32)
			if err != nil {
				newErrorResponse(ctx, http.StatusBadRequest, err.Error())
				return
			}
			tr.AmountTotal = float32(amountTotal)

			amountOriginal, err := strconv.ParseFloat(read[5], 32)
			if err != nil {
				newErrorResponse(ctx, http.StatusBadRequest, err.Error())
				return
			}
			tr.AmountOriginal = float32(amountOriginal)

			commissionPS, err := strconv.ParseFloat(read[6], 32)
			tr.CommissionPS = float32(commissionPS)
			if err != nil {
				newErrorResponse(ctx, http.StatusBadRequest, err.Error())
				return
			}

			commissionClient, err := strconv.ParseFloat(read[7], 32)
			if err != nil {
				newErrorResponse(ctx, http.StatusBadRequest, err.Error())
				return
			}
			tr.CommissionClient = float32(commissionClient)

			commissionProvider, err := strconv.ParseFloat(read[8], 32)
			if err != nil {
				newErrorResponse(ctx, http.StatusBadRequest, err.Error())
				return
			}
			tr.CommissionProvider = float32(commissionProvider)

			tr.DateInput, err = time.Parse("2006-01-02 15:04:05", read[9])
			if err != nil {
				newErrorResponse(ctx, http.StatusBadRequest, err.Error())
				return
			}

			tr.DatePost, err = time.Parse("2006-01-02 15:04:05", read[10])
			if err != nil {
				newErrorResponse(ctx, http.StatusBadRequest, err.Error())
				return
			}

			tr.Status = read[11]
			tr.PaymentType = read[12]
			tr.PaymentNumber = read[13]

			tr.ServiceId, err = strconv.Atoi(read[14])
			if err != nil {
				newErrorResponse(ctx, http.StatusBadRequest, err.Error())
				return
			}

			tr.Service = read[15]

			tr.PayeeId, err = strconv.Atoi(read[16])
			if err != nil {
				newErrorResponse(ctx, http.StatusBadRequest, err.Error())
				return
			}

			tr.PayeeName = read[17]
			tr.PayeeBankMFO = read[18]
			tr.PayeeBankAccount = read[19]
			tr.PaymentNarrative = read[20]

			parsedData = append(parsedData, tr)
		}

		ctx.Set(parsedDataCtx, parsedData)

		handlerFunc(ctx)
	}
}

func (h *handler) TransactionQueryParams(handlerFunc func(ctx *gin.Context)) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
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

		terminalIdsTrim := strings.ReplaceAll(ctx.Query("terminalIds"), " ", "")
		terminalIds := strings.Split(terminalIdsTrim, ",")
		if terminalIds[0] != "" {
			for _, id := range terminalIds {
				_, err = strconv.Atoi(id)
				if err != nil {
					newErrorResponse(ctx, http.StatusBadRequest, "invalid terminal ids")
					return
				}
			}
		}

		params := models.SearchParams{
			TransactionId:    transactionIdInt,
			TerminalIds:      terminalIds,
			Status:           ctx.Query("status"),
			PaymentType:      ctx.Query("paymentType"),
			DatePostFrom:     ctx.Query("datePostFrom"),
			DatePostTo:       ctx.Query("datePostTo"),
			PaymentNarrative: ctx.Query("paymentNarrative"),
		}

		ctx.Set(paramsCtx, params)

		handlerFunc(ctx)
	}
}
