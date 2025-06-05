package utils

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/remiges-tech/alya/config"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/logharbour/logharbour"
	sqldb "github.com/shahfaizansr/migrate/db"
	"github.com/shahfaizansr/models"
	"github.com/shahfaizansr/models/mybulkcalc"
	"github.com/shahfaizansr/models/mycalc"
	"github.com/shahfaizansr/utils/cvlconstant"
)

func SetConfigEnvironment(environment models.Environment) models.AppConfig {
	var appConfig models.AppConfig
	if !environment.IsValid() {
		log.Fatal("environment params is not valid")
	}
	switch environment {
	case models.DevEnv:
		config.LoadConfigFromFile("config/env/config_dev.json", &appConfig)
	case models.ProdEnv:
		config.LoadConfigFromFile("config/env/config_prod.json", &appConfig)
	case models.UATEnv:
		config.LoadConfigFromFile("config/env/config_uat.json", &appConfig)
	}
	return appConfig
}

// GetCurrentLoggerPriority returns the appropriate logharbour.LogPriority
// based on the provided logger priority string.
func GetCurrentLoggerPriority(loggerPriority string) logharbour.LogPriority {
	var currentPriority logharbour.LogPriority
	switch loggerPriority {
	case logharbour.LogPriorityDebug2:
		currentPriority = logharbour.Debug2
	case logharbour.LogPriorityDebug1:
		currentPriority = logharbour.Debug1
	case logharbour.LogPriorityDebug0:
		currentPriority = logharbour.Debug0
	case logharbour.LogPriorityInfo:
		currentPriority = logharbour.Info
	case logharbour.LogPriorityWarn:
		currentPriority = logharbour.Warn
	case logharbour.LogPriorityErr:
		currentPriority = logharbour.Err
	case logharbour.LogPriorityCrit:
		currentPriority = logharbour.Crit
	case logharbour.LogPrioritySec:
		currentPriority = logharbour.Sec
	case logharbour.LogPriorityUnknown:
		currentPriority = logharbour.LogPriority(logharbour.Unknown)
	}
	return currentPriority
}

func ArithmeticCalculation(req mycalc.CalcRequestModel, srv *service.Service) (mycalc.CalcResponse, error) {
	var result float64

	switch req.Operation {
	case "sum":
		for _, val := range req.Input {
			result += val
		}
		return mycalc.CalcResponse{
			Result: []float64{result},
			Error:  "",
		}, nil

	case "mean", "average":
		if len(req.Input) == 0 {
			return mycalc.CalcResponse{}, errors.New("cannot calculate mean/average on empty input")
		}
		for _, val := range req.Input {
			result += val
		}
		result /= float64(len(req.Input))
		return mycalc.CalcResponse{
			Result: []float64{result},
			Error:  "",
		}, nil

	case "sort":
		sorted := make([]float64, len(req.Input))
		copy(sorted, req.Input)
		sort.Float64s(sorted)
		return mycalc.CalcResponse{
			Result: sorted,
			Error:  "",
		}, nil

	default:
		return mycalc.CalcResponse{}, fmt.Errorf("unsupported operation: %s", req.Operation)
	}
}

func DecodeBase64CSV(input string) ([]string, error) {
	data, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(data), "\n"), nil
}

// ProcessCSVLine processes a single CSV line and returns numbers or error
func ProcessCSVLine(line string, lineNumber int) ([]float64, error) {
	record := strings.Split(line, ",")
	numbers := []float64{}

	// Convert and validate fields into float64 values
	for _, field := range record {
		field = strings.TrimSpace(field)
		if field == "" {
			continue // Skip empty fields
		}
		num, parseErr := ParseFloatSafe(field)
		if parseErr != nil {
			return nil, parseErr
		}
		numbers = append(numbers, num)
	}

	if len(numbers) == 0 {
		return nil, errors.New("empty record")
	}

	return numbers, nil
}

// ParseFloatSafe parses a float64 safely from a trimmed string.
func ParseFloatSafe(s string) (float64, error) {
	return strconv.ParseFloat(strings.TrimSpace(s), 64)
}

func ValidateBulkCalcRequest(request mybulkcalc.BulkCalcRequestModel) []wscutils.ErrorMessage {
	var validationErrors []wscutils.ErrorMessage

	// Validate input
	if len(request.Input) == 0 {
		validationErrors = append(validationErrors, wscutils.BuildErrorMessage(
			cvlconstant.EMPTY_MESSAGE_ID,
			cvlconstant.ERRORCODE_EMPTY,
			cvlconstant.INPUT,
			cvlconstant.INPUT,
		))
	}

	// Validate operation is provided
	if request.Operation == "" {
		validationErrors = append(validationErrors, wscutils.BuildErrorMessage(
			cvlconstant.EMPTY_MESSAGE_ID,
			cvlconstant.ERRORCODE_EMPTY,
			cvlconstant.OPERATION,
			cvlconstant.OPERATION,
		))
	} else {
		// Validate operation value
		validOps := map[string]bool{
			"sum":     true,
			"mean":    true,
			"average": true,
			"sort":    true,
		}
		if !validOps[strings.ToLower(request.Operation)] {
			validationErrors = append(validationErrors, wscutils.BuildErrorMessage(
				cvlconstant.INVALID_MESSAGE_ID,
				cvlconstant.ERRORCODE_INVALID,
				cvlconstant.OPERATION,
				cvlconstant.SUM,
				cvlconstant.MEAN,
				cvlconstant.SORT,
				cvlconstant.AVERAGE,
			))
		}
	}

	return validationErrors
}

func ValidateMyCalcRequest(request mycalc.CalcRequestModel) []wscutils.ErrorMessage {
	var validationErrors []wscutils.ErrorMessage

	// Validate input
	if len(request.Input) == 0 {
		validationErrors = append(validationErrors, wscutils.BuildErrorMessage(
			cvlconstant.EMPTY_MESSAGE_ID,
			cvlconstant.ERRORCODE_EMPTY,
			cvlconstant.INPUT,
			cvlconstant.INPUT,
		))
	}

	// Validate operation is provided
	if request.Operation == "" {
		validationErrors = append(validationErrors, wscutils.BuildErrorMessage(
			cvlconstant.EMPTY_MESSAGE_ID,
			cvlconstant.ERRORCODE_EMPTY,
			cvlconstant.OPERATION,
			cvlconstant.OPERATION,
		))
	} else {
		// Validate operation value
		validOps := map[string]bool{
			"sum":     true,
			"mean":    true,
			"average": true,
			"sort":    true,
		}
		if !validOps[strings.ToLower(request.Operation)] {
			validationErrors = append(validationErrors, wscutils.BuildErrorMessage(
				cvlconstant.INVALID_MESSAGE_ID,
				cvlconstant.ERRORCODE_INVALID,
				cvlconstant.OPERATION,
				cvlconstant.SUM,
				cvlconstant.MEAN,
				cvlconstant.SORT,
				cvlconstant.AVERAGE,
			))
		}
	}

	return validationErrors
}

// StoreCalcLogToDB logs the calculation request and response details to the database.
// It also measures duration and handles any marshalling or DB errors during the logging process.
func StoreCalcLogToDB(
	ctx context.Context,
	srv *service.Service,
	request any,
	response any,
	requestTime, responseTime time.Time,
	durationMs float64,
) {
	dbHandler, ok := srv.Database.(*sqldb.DBHandler)
	if !ok {
		srv.LogHarbour.Err().Log("Failed to assert *sqldb.DBHandler")
		return
	}

	reqJSON, err := json.Marshal(request)
	if err != nil {
		srv.LogHarbour.Err().LogActivity(cvlconstant.MARSHALING_ERROR, err.Error())
		return
	}

	resJSON, err := json.Marshal(response)
	if err != nil {
		srv.LogHarbour.Err().LogActivity(cvlconstant.MARSHALING_ERROR, err.Error())
		return
	}

	// Extract fields from response if available
	var errMsg string
	switch resp := response.(type) {
	case mycalc.CalcResponse:
		errMsg = resp.Error
	case *mycalc.CalcResponse:
		errMsg = resp.Error
	default:
		srv.LogHarbour.Err().
			LogActivity("TypeAssertionFailed", fmt.Sprintf("Unknown response type: %T", response))
	}

	// Extract fields from request if available
	var operation, input string
	switch req := request.(type) {
	case mycalc.CalcRequestModel:
		operation = req.Operation
		input = fmt.Sprintf("%v", req.Input)
	case mybulkcalc.BulkCalcRequestModel:
		operation = req.Operation
		input = fmt.Sprintf("%v", req.Input)
	default:
		srv.LogHarbour.Warn().Log(fmt.Sprintf("Request fields not extracted. Unknown type: %T", request))
	}

	entry := mycalc.CalcLogEntry{
		RequestTime:  requestTime,
		ResponseTime: responseTime,
		DurationMs:   durationMs,
		RequestData:  string(reqJSON),
		ResponseData: string(resJSON),
		Error:        errMsg,
		Operation:    operation,
		Input:        input,
	}

	query := os.Getenv(cvlconstant.CALC_LOG_INSERT_QUERY)

	_, err = dbHandler.DB.NamedExecContext(ctx, query, entry)
	if err != nil {
		srv.LogHarbour.Err().LogActivity(cvlconstant.ERROR_INSERTING_LOG_ENTRY, err.Error())
		return
	}

	srv.LogHarbour.Info().Log(cvlconstant.DB_STORAGE_SUCCESS)
}

// CalculationProcess performs basic arithmetic operations (sum, mean, average, sort)
// on the input slice provided in the request and returns the result.
func CalculationProcess(request mycalc.CalcRequestModel) (mycalc.CalcResponse, error) {
	input := request.Input
	var result []float64
	sum := 0.0

	switch request.Operation {
	case "sum":
		for _, num := range input {
			sum += num
		}
		result = append(result, sum)

	case "mean", "average":
		for _, num := range input {
			sum += num
		}
		mean := sum / float64(len(input))
		result = append(result, mean)

	case "sort":
		sort.Float64s(input)
		result = input

	default:
		return mycalc.CalcResponse{}, errors.New(cvlconstant.UNSUPPORTED_OPERATION + request.Operation)
	}

	return mycalc.CalcResponse{Result: result}, nil
}
