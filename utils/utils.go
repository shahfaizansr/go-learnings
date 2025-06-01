package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/remiges-tech/alya/config"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/logharbour/logharbour"
	sqldb "github.com/shahfaizansr/migrate/db"
	"github.com/shahfaizansr/models"
	"github.com/shahfaizansr/models/mycalc"
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

func StoreCalcLogToDB(
	ctx *gin.Context,
	srv *service.Service,
	request mycalc.CalcRequestModel,
	response mycalc.CalcResponse,
	requestTime, responseTime time.Time,
	duration float64,
) {
	inputJSON, err := json.Marshal(request.Input)
	if err != nil {
		srv.LogHarbour.WithModule("MYCALC").WithOp("StoreCalcLogToDB").Err().Log("Failed to marshal input: " + err.Error())
		return
	}

	resultJSON, err := json.Marshal(response.Result)
	if err != nil {
		srv.LogHarbour.WithModule("MYCALC").WithOp("StoreCalcLogToDB").Err().Log("Failed to marshal result: " + err.Error())
		return
	}

	log := mycalc.CalcResponseModel{
		Input:        string(inputJSON),
		Operation:    request.Operation,
		Result:       string(resultJSON),
		Error:        response.Error,
		RequestTime:  requestTime,
		ResponseTime: responseTime,
		DurationMs:   duration,
	}

	dbHandler, ok := srv.Database.(*sqldb.DBHandler)
	if !ok {
		srv.LogHarbour.WithModule("MYCALC").WithOp("StoreCalcLogToDB").Err().Log("Invalid database handler type in srv.Database")
		return
	}

	if err := dbHandler.DB.Create(&log).Error; err != nil {
		srv.LogHarbour.WithModule("MYCALC").WithOp("StoreCalcLogToDB").Err().Log("DB insert failed: " + err.Error())
	}
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
