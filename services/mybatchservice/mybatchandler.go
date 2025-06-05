package mybatchservice

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/remiges-tech/alya/jobs"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/shahfaizansr/models/mybulkcalc"
	"github.com/shahfaizansr/models/mycalc"
	"github.com/shahfaizansr/utils"
	"github.com/shahfaizansr/utils/cvlconstant"
)

func HandleMyBulkCalcBatch(ctx *gin.Context, srv *service.Service) {
	logger := srv.LogHarbour.
		WithModule(cvlconstant.MODULE_MYBULKCALC_BATCH).
		WithClass(cvlconstant.MyBulkCalcBatch.String()).
		WithOp(cvlconstant.MyBulkCalcBatch.String())

	logger.Info().Log("HandleMyBulkCalcBatch request received")

	var (
		request      mybulkcalc.BulkCalcRequestModel
		jm           = srv.Dependencies["jobmanager"].(*jobs.JobManager)
		batchInputs  []jobs.BatchInput_t
		requestTime  = time.Now()
		responseTime time.Time
		duration     float64
	)

	// STEP 1: Bind JSON
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

	// STEP 2: Validate business rules
	if validationErrors := utils.ValidateBulkCalcRequest(request); len(validationErrors) > 0 {
		responseTime = time.Now()
		duration = responseTime.Sub(requestTime).Seconds() * 1000

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

	// STEP 3: Decode base64 CSV
	decodedCSV, err := utils.DecodeBase64CSV(request.Input)
	if err != nil {
		responseTime = time.Now()
		duration = responseTime.Sub(requestTime).Seconds() * 1000
		errorMsg := fmt.Sprintf("Base64 decode failed: %v", err)

		utils.StoreCalcLogToDB(ctx, srv, request, mycalc.CalcResponse{
			Result: nil,
			Error:  errorMsg,
		}, requestTime, responseTime, duration)

		logger.Error(err).LogActivity(cvlconstant.ERRORCODE_BASE64, nil)

		wscutils.SendErrorResponse(ctx, &wscutils.Response{
			Status: wscutils.ErrorStatus,
			Data:   nil,
			Messages: []wscutils.ErrorMessage{
				wscutils.BuildErrorMessage(cvlconstant.INVALID_BASE64, cvlconstant.ERRORCODE_BASE64, cvlconstant.INPUTFILE, err.Error()),
			},
		})
		return
	}

	// STEP 4: Convert each CSV line into job input
	for lineIdx, line := range decodedCSV {
		lineNumber := lineIdx + 1

		numbers, err := utils.ProcessCSVLine(line, lineNumber)
		if err != nil {
			logger.Warn().LogActivity(cvlconstant.CSV_PARSE_ERROR, map[string]interface{}{
				"line":   lineNumber,
				"record": line,
				"error":  err.Error(),
			})
			continue
		}

		calcReq := mycalc.CalcRequestModel{
			Input:     numbers,
			Operation: request.Operation,
		}

		inputBytes, err := json.Marshal(calcReq)
		if err != nil {
			logger.Warn().LogActivity(cvlconstant.JSON_MARSHAL_ERROR, map[string]interface{}{
				"line": lineNumber,
				"err":  err.Error(),
			})
			continue
		}

		inputJSON, err := jobs.NewJSONstr(string(inputBytes))
		if err != nil {
			logger.Warn().LogActivity(cvlconstant.INVALID_JSON_CONVEN, map[string]interface{}{
				"error": err.Error(),
				"json":  string(inputBytes),
			})
			continue
		}

		batchInputs = append(batchInputs, jobs.BatchInput_t{
			Line:  lineNumber,
			Input: inputJSON,
		})
	}

	if len(batchInputs) == 0 {
		responseTime = time.Now()
		duration = responseTime.Sub(requestTime).Seconds() * 1000
		errMsg := "No valid CSV lines to process"

		utils.StoreCalcLogToDB(ctx, srv, request, mycalc.CalcResponse{
			Result: nil,
			Error:  errMsg,
		}, requestTime, responseTime, duration)

		logger.Err().LogActivity("no_valid_batch_lines", errMsg)
		wscutils.SendErrorResponse(ctx, &wscutils.Response{
			Status: wscutils.ErrorStatus,
			Messages: []wscutils.ErrorMessage{
				wscutils.BuildErrorMessage(11111, "EMPTY_BATCH", "inputfile", errMsg),
			},
		})
		return
	}

	// STEP 5: Marshal context
	jsonBytes, err := json.Marshal(request)
	if err != nil {
		responseTime = time.Now()
		duration = responseTime.Sub(requestTime).Seconds() * 1000
		errorMsg := fmt.Sprintf("Request marshal failed: %v", err)

		utils.StoreCalcLogToDB(ctx, srv, request, mycalc.CalcResponse{
			Result: nil,
			Error:  errorMsg,
		}, requestTime, responseTime, duration)

		wscutils.SendErrorResponse(ctx, &wscutils.Response{
			Status: wscutils.ErrorStatus,
			Messages: []wscutils.ErrorMessage{
				wscutils.BuildErrorMessage(22222, "JSON_ERROR", "context", err.Error()),
			},
		})
		return
	}

	jsonStr, err := jobs.NewJSONstr(string(jsonBytes))
	if err != nil {
		responseTime = time.Now()
		duration = responseTime.Sub(requestTime).Seconds() * 1000
		errorMsg := fmt.Sprintf("jobs.NewJSONstr failed: %v", err)

		utils.StoreCalcLogToDB(ctx, srv, request, mycalc.CalcResponse{
			Result: nil,
			Error:  errorMsg,
		}, requestTime, responseTime, duration)

		wscutils.SendErrorResponse(ctx, &wscutils.Response{
			Status: wscutils.ErrorStatus,
			Messages: []wscutils.ErrorMessage{
				wscutils.BuildErrorMessage(3333, "JSON_FORMAT", "context", err.Error()),
			},
		})
		return
	}

	// STEP 6: Submit job
	jobID, err := jm.BatchSubmit(cvlconstant.MyBulkCalc.String(), cvlconstant.BULK_OPERATION, jsonStr, batchInputs, false)
	responseTime = time.Now()
	duration = responseTime.Sub(requestTime).Seconds() * 1000

	if err != nil {
		errorMsg := fmt.Sprintf("BatchSubmit failed: %v", err)
		logger.Error(err).LogActivity(cvlconstant.BATCH_SUBMISSION_FAILED, err.Error())

		utils.StoreCalcLogToDB(ctx, srv, request, mycalc.CalcResponse{
			Result: nil,
			Error:  errorMsg,
		}, requestTime, responseTime, duration)

		wscutils.SendErrorResponse(ctx, &wscutils.Response{
			Status: wscutils.ErrorStatus,
			Messages: []wscutils.ErrorMessage{
				wscutils.BuildErrorMessage(4444, "BATCH_SUBMIT", "job", err.Error()),
			},
		})
		return
	}

	// STEP 7: Success log
	utils.StoreCalcLogToDB(ctx, srv, request, mycalc.CalcResponseBUlk{
		Result: jobID,
		Error:  "",
	}, requestTime, responseTime, duration)

	logger.Info().LogActivity("Batch submitted successfully", map[string]interface{}{
		"jobID": jobID,
		"count": len(batchInputs),
	})

	wscutils.SendSuccessResponse(ctx, &wscutils.Response{
		Status:   wscutils.SuccessStatus,
		Data:     jobID,
		Messages: nil,
	})
}
