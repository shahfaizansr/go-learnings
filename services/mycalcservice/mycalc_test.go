package mycalcservice_test

// import (
// 	"bytes"
// 	"encoding/json"
// 	"flag"
// 	"log"
// 	"net/http"
// 	"net/http/httptest"
// 	"os"
// 	"testing"

// 	"github.com/gin-gonic/gin"
// 	"github.com/remiges-tech/alya/service"
// 	"github.com/remiges-tech/alya/wscutils"
// 	"github.com/remiges-tech/logharbour/logharbour"
// 	sqldb "github.com/shahfaizansr/migrate/db"
// 	"github.com/shahfaizansr/models"
// 	"github.com/shahfaizansr/models/mycalc"
// 	"github.com/shahfaizansr/services/mycalcservice"
// 	"github.com/shahfaizansr/utils"
// 	"github.com/shahfaizansr/utils/cvlconstant"
// 	"github.com/shahfaizansr/utils/rigel"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"gorm.io/gorm"
// )

// // MockDBHandler implements sqldb.DBHandler interface for testing
// type MockDBHandler struct {
// 	mock.Mock
// }

// func (m *MockDBHandler) GetDB() *gorm.DB {
// 	args := m.Called()
// 	return args.Get(0).(*gorm.DB)
// }

// func (m *MockDBHandler) Create(value interface{}) *gorm.DB {
// 	args := m.Called(value)
// 	return args.Get(0).(*gorm.DB)
// }

// // Add other required methods from sqldb.DBHandler interface here
// // For example:
// func (m *MockDBHandler) First(dest interface{}, conds ...interface{}) *gorm.DB {
// 	args := m.Called(dest, conds)
// 	return args.Get(0).(*gorm.DB)
// }

// // RequestWrapper wraps the request data
// type RequestWrapper struct {
// 	Data mycalc.CalcRequestModel `json:"data"`
// }

// var envFlag = flag.String("env", "dev_env", "setup the environment using CLI")

// func SetupTestContext(method, path string, body []byte, dbHandler sqldb.DBHandler) (*gin.Context, *httptest.ResponseRecorder, *service.Service) {
// 	gin.SetMode(gin.TestMode)
// 	w := httptest.NewRecorder()
// 	c, _ := gin.CreateTestContext(w)
// 	req := httptest.NewRequest(method, path, bytes.NewReader(body))
// 	req.Header.Set("Content-Type", "application/json")
// 	c.Request = req

// 	currentLogPriority := utils.GetCurrentLoggerPriority("info")
// 	lctx := logharbour.NewLoggerContext(currentLogPriority)
// 	l := logharbour.NewLogger(lctx, cvlconstant.APP_NAME, os.Stdout)

// 	appConfig := utils.SetConfigEnvironment(models.Environment(*envFlag))
// 	rigelClient := rigel.SetRigel(appConfig)

// 	if rigelClient == nil {
// 		log.Fatal("Failed to create Rigel instance")
// 	} else {
// 		log.Println("Rigel instance created successfully")
// 	}

// 	svc := service.NewService(nil)
// 	svc.LogHarbour = l
// 	svc.RigelConfig = rigelClient
// 	svc.Database = dbHandler

// 	return c, w, svc
// }

// func TestSuccessForAverage(t *testing.T) {
// 	// Setup test data
// 	req := RequestWrapper{
// 		Data: mycalc.CalcRequestModel{
// 			Input:     []float64{1, 2, 3, 4, 5},
// 			Operation: "average",
// 		},
// 	}

// 	body, err := json.Marshal(req)
// 	assert.NoError(t, err, "failed to marshal request body")

// 	// Create mock DB handler
// 	mockDB := &MockDBHandler{}

// 	// Setup expectations
// 	mockDB.On("GetDB").Return(&gorm.DB{})
// 	mockDB.On("Create", mock.AnythingOfType("*mycalc.CalcLog")).Return(&gorm.DB{})

// 	c, w, s := SetupTestContext("POST", "/mycalc", body, mockDB)

// 	// Call the handler
// 	mycalcservice.CalculatorHandler(c, s)

// 	// Verify response
// 	assert.Equal(t, http.StatusOK, w.Code, "Expected status 200 OK")

// 	var response wscutils.Response
// 	err = json.Unmarshal(w.Body.Bytes(), &response)
// 	assert.NoError(t, err, "failed to unmarshal response")

// 	assert.Equal(t, "success", response.Status, "Expected status 'success'")

// 	dataMap, ok := response.Data.(map[string]interface{})
// 	assert.True(t, ok, "Expected data to be a map")

// 	result := dataMap["result"].([]interface{})
// 	assert.Len(t, result, 1, "Expected result length 1")
// 	assert.Equal(t, 3.0, result[0].(float64), "Expected result [3]")

// 	// Verify mock expectations
// 	mockDB.AssertExpectations(t)
// }
