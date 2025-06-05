package mybatchservice

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/remiges-tech/alya/jobs"
	"github.com/remiges-tech/alya/jobs/pg/batchsqlc"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/remiges-tech/logharbour/logharbour"
	"github.com/shahfaizansr/models/mybulkcalc"
	"github.com/shahfaizansr/models/mycalc"
	"github.com/shahfaizansr/utils"
	"github.com/shahfaizansr/utils/cvlconstant"
)

// MyCalcInitializer implements the jobs.Initializer interface to support batch initialization.
type MyCalcInitializer struct{}

// MyBulkCalcBathProcessor handles the batch job logic for bulk calculations.
type MyBulkCalcBathProcessor struct{}

// Init initializes any contextual data required before batch processing.
// Currently returns an empty InitBlock. Can be extended in future for setup needs (e.g., auth, config).
func (m *MyCalcInitializer) Init(arg string) (jobs.InitBlock, error) {
	fmt.Println("Initializing MyCalcInitializer with arg:", arg)
	var initBlock jobs.InitBlock
	return initBlock, nil
}

// DoBatchJob is invoked for each line of the batch job.
// It handles decoding input JSON, validating, performing the calculation, and formatting results.
func (p *MyBulkCalcBathProcessor) DoBatchJob(initBlock jobs.InitBlock, context jobs.JSONstr, line int, input jobs.JSONstr) (
	status batchsqlc.StatusEnum, result jobs.JSONstr, messages []wscutils.ErrorMessage, blobRows map[string]string, err error,
) {
	// Setup logger
	lctx := logharbour.NewLoggerContext(logharbour.DefaultPriority)
	logger := logharbour.NewLogger(lctx, "DoBatchJob", os.Stdout)

	// Parse input JSON into domain model
	inputStr := input.String()
	var inputStrJson mybulkcalc.BatchBulkCalcRequestModel

	fmt.Println("Inside MyBulkCalcBatchProcessor, Line:", line)
	fmt.Println("Raw Input:", inputStr)

	// Decode JSON input string
	err = json.Unmarshal([]byte(inputStr), &inputStrJson)
	if err != nil {
		logger.Error(err).LogActivity(cvlconstant.DOBATCHJOB, err.Error())
		result, _ = jobs.NewJSONstr(`{"message": "Invalid Input Provided"}`)
		return batchsqlc.StatusEnumFailed, result, nil, blobRows, err
	}

	fmt.Println("Parsed input JSON:", inputStrJson)

	// Reject empty inputs early
	if len(inputStrJson.Input) == 0 {
		logger.Error(err).LogActivity(cvlconstant.EMPTY_NUMBERS_LIST, nil)
		result, _ = jobs.NewJSONstr(`{"message": "Empty input provided"}`)
		return batchsqlc.StatusEnumFailed, result, nil, blobRows, nil
	}

	// Build calculation request and perform operation
	req := mycalc.CalcRequestModel{
		Operation: inputStrJson.Operation,
		Input:     inputStrJson.Input,
	}

	// Perform the calculation using the business logic utility
	resp, err := utils.CalculationProcess(req)
	if err != nil {
		logger.Error(err).LogActivity(cvlconstant.CALCULATION_ERROR, err.Error())
		result, _ = jobs.NewJSONstr(`{"message": "Invalid Operation"}`)
		return batchsqlc.StatusEnumFailed, result, nil, blobRows, nil
	}

	// Log success to stdout
	fmt.Printf("calculation successful for line %d. Operation: %s, Result: %v\n",
		line, req.Operation, resp.Result)

	// Store result in blobRows to be persisted in the database
	blobRows = map[string]string{
		"Result": fmt.Sprintf("%v", resp.Result),
	}

	// Set structured success response
	result, _ = jobs.NewJSONstr(`{"message": "Operation performed successfully"}`)
	return batchsqlc.StatusEnumSuccess, result, nil, blobRows, nil
}

// MarkDone is called once the batch completes processing.
// Used to log the final status and summary stats.
func (p *MyBulkCalcBathProcessor) MarkDone(initBlock jobs.InitBlock, context jobs.JSONstr, details jobs.BatchDetails_t) error {
	log.Printf("batch %s completed with status: %s", details.ID, details.Status)
	log.Printf("Results: Success=%d, Failed=%d, Aborted=%d",
		details.NSuccess, details.NFailed, details.NAborted)
	return nil
}
