package mybulkcalcservice

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/shahfaizansr/models/mybulkcalc"
	"github.com/shahfaizansr/models/mycalc"
	"github.com/shahfaizansr/utils"
	"github.com/shahfaizansr/utils/cvlconstant"
)

func HandleMyBulkCalc(ctx *gin.Context, srv *service.Service) {
	logger := srv.LogHarbour.
		WithModule(cvlconstant.MODULE_MYBULKCALC).
		WithClass(cvlconstant.MyBulkCalc.String()).
		WithOp(cvlconstant.MyBulkCalc.String())
	logger.Info().Log("Bulk calculator request received")

	var (
		request      mybulkcalc.BulkCalcRequestModel
		response     mybulkcalc.BulkCalcResponseModel
		requestTime  = time.Now()
		responseTime time.Time
		duration     float64
		allErrors    []wscutils.ErrorMessage
		results      [][]float64
		successLines []int
		failedLines  []int
		allLines     []string
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

	// STEP 2: Validate
	if valErrors := utils.ValidateBulkCalcRequest(request); len(valErrors) > 0 {
		allErrors = append(allErrors, valErrors...)
	}

	// STEP 3: Decode Base64 CSV
	decodedCSV, err := utils.DecodeBase64CSV(request.Input)
	if err != nil {
		allErrors = append(allErrors,
			wscutils.BuildErrorMessage(
				cvlconstant.INVALID_BASE64,
				cvlconstant.ERRORCODE_BASE64,
				cvlconstant.INPUTFILE,
				err.Error(),
			),
		)
		logger.Error(err).LogActivity(cvlconstant.ERRORCODE_BASE64, nil)
	}

	for lineIdx, line := range decodedCSV {
		if strings.TrimSpace(line) == "" {
			continue
		}
		response.NLines++
		allLines = append(allLines, line)

		numbers, err := utils.ProcessCSVLine(line, lineIdx+1)
		if err != nil {
			response.NFailed++
			failedLines = append(failedLines, lineIdx+1)

			logger.Error(err).LogActivity(cvlconstant.CSV_PARSE_ERROR, map[string]interface{}{
				"line":   lineIdx + 1,
				"record": line,
			})

			utils.StoreCalcLogToDB(ctx, srv, request, mycalc.CalcResponse{
				Result: nil,
				Error:  err.Error(),
			}, requestTime, time.Now(), 0)

			continue
		}

		calcReq := mycalc.CalcRequestModel{
			Input:     numbers,
			Operation: request.Operation,
		}

		calcResult, calcErr := utils.ArithmeticCalculation(calcReq, srv)
		if calcErr != nil {
			response.NFailed++
			failedLines = append(failedLines, lineIdx+1)

			logger.Error(calcErr).LogActivity(cvlconstant.UNSUPPORTED_OPERATION, map[string]interface{}{
				"operation": request.Operation,
			})

			utils.StoreCalcLogToDB(ctx, srv, calcReq, mycalc.CalcResponse{
				Result: nil,
				Error:  calcErr.Error(),
			}, requestTime, time.Now(), 0)

			continue
		}

		response.NSuccess++
		successLines = append(successLines, lineIdx+1)
		results = append(results, calcResult.Result)
	}

	// STEP 5: Final response assembly
	response.Results = results
	responseTime = time.Now()
	duration = responseTime.Sub(requestTime).Seconds() * 1000

	// Final Logging
	utils.StoreCalcLogToDB(ctx, srv, request, response, requestTime, responseTime, duration)

	// STEP 6: Send Response
	if response.NSuccess == 0 {
		wscutils.SendErrorResponse(ctx, &wscutils.Response{
			Status:   wscutils.ErrorStatus,
			Messages: allErrors,
		})
	} else if response.NFailed > 0 {
		// Mixed success
		wscutils.SendSuccessResponse(ctx, &wscutils.Response{
			Status:   wscutils.SuccessStatus,
			Data:     response,
			Messages: allErrors,
		})
	} else {
		// Full success
		wscutils.SendSuccessResponse(ctx, &wscutils.Response{
			Status: wscutils.SuccessStatus,
			Data:   response,
		})
	}

	logger.Info().LogActivity("Bulk calculation completed", map[string]interface{}{
		"nlines":   response.NLines,
		"nsuccess": response.NSuccess,
		"nfailed":  response.NFailed,
	})
}
