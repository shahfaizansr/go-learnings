package mycalcservice

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/shahfaizansr/models/mycalc"
	"github.com/shahfaizansr/utils"
	"github.com/shahfaizansr/utils/cvlconstant"
)

func CalculatorHandler(ctx *gin.Context, srv *service.Service) {
	logger := srv.LogHarbour.WithModule(cvlconstant.MODULE_MYCALC).WithClass(cvlconstant.MyCalc.String()).WithOp(cvlconstant.MyCalc.String())
	logger.Info().Log("Calculator request received")

	var (
		request      mycalc.CalcRequestModel
		requestTime  = time.Now()
		responseTime time.Time
		duration     float64
		calcResult   mycalc.CalcResponse
	)

	// Parse request
	if err := wscutils.BindJSON(ctx, &request); err != nil {
		responseTime = time.Now()
		duration = responseTime.Sub(requestTime).Seconds() * 1000
		errorMsg := fmt.Sprintf("Invalid request format: %v", err)

		utils.StoreCalcLogToDB(ctx, srv, request, mycalc.CalcResponse{
			Result: nil,
			Error:  errorMsg,
		}, requestTime, responseTime, duration)

		return
	}

	// Validate input
	var validationErrors []wscutils.ErrorMessage
	if len(request.Input) == 0 {
		validationErrors = append(validationErrors, wscutils.BuildErrorMessage(
			cvlconstant.EMPTY_MESSAGE_ID,
			cvlconstant.ERRORCODE_EMPTY,
			cvlconstant.INPUT,
			cvlconstant.INPUT,
		))
	}
	if request.Operation == "" {
		validationErrors = append(validationErrors, wscutils.BuildErrorMessage(
			cvlconstant.EMPTY_MESSAGE_ID,
			cvlconstant.ERRORCODE_EMPTY,
			cvlconstant.OPERATION,
			cvlconstant.OPERATION,
		))
	}

	if len(validationErrors) > 0 {
		responseTime = time.Now()
		duration = responseTime.Sub(requestTime).Seconds() * 1000

		// Convert validation errors to JSON string for storage
		errorBytes, err := json.Marshal(validationErrors)
		errorMsg := string(errorBytes)
		if err != nil {
			errorMsg = fmt.Sprintf("Failed to marshal validation errors: %v", err)
			logger.Err().Log(errorMsg)
		}

		utils.StoreCalcLogToDB(ctx, srv, request, mycalc.CalcResponse{
			Result: nil,
			Error:  errorMsg,
		}, requestTime, responseTime, duration)

		wscutils.SendErrorResponse(ctx, wscutils.NewResponse(
			wscutils.ErrorStatus,
			nil,
			validationErrors,
		))
		return
	}

	// Perform calculation
	calcResult, calcErr := utils.ArithmeticCalculation(request, srv)
	if calcErr != nil {
		responseTime = time.Now()
		duration = responseTime.Sub(requestTime).Seconds() * 1000

		errorMsg := calcErr.Error()
		logger.Err().LogActivity(cvlconstant.CALCULATION_ERROR, errorMsg)

		utils.StoreCalcLogToDB(ctx, srv, request, mycalc.CalcResponse{
			Result: nil,
			Error:  errorMsg,
		}, requestTime, responseTime, duration)

		wscutils.SendErrorResponse(ctx, &wscutils.Response{
			Status: wscutils.ErrorStatus,
			Messages: []wscutils.ErrorMessage{
				wscutils.BuildErrorMessage(
					cvlconstant.INVALID_MESSAGE_ID,
					cvlconstant.ERRORCODE_INVALID,
					cvlconstant.OPERATION,
					request.Operation,
					cvlconstant.SUM,
					cvlconstant.SORT,
					cvlconstant.AVERAGE,
					cvlconstant.MEAN,
				),
			},
		})
		return
	}

	// Success case
	responseTime = time.Now()
	duration = responseTime.Sub(requestTime).Seconds() * 1000

	utils.StoreCalcLogToDB(ctx, srv, request, mycalc.CalcResponse{
		Result: calcResult.Result,
		Error:  "",
	}, requestTime, responseTime, duration)

	logger.Info().LogActivity("Calculation completed", map[string]interface{}{
		cvlconstant.OPERATION: request.Operation,
		cvlconstant.INPUT:     request.Input,
		cvlconstant.RESULT:    calcResult.Result,
	})

	wscutils.SendSuccessResponse(ctx, &wscutils.Response{
		Status:   wscutils.SuccessStatus,
		Data:     calcResult,
		Messages: []wscutils.ErrorMessage{},
	})
}
