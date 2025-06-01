package myhelloservice

import (
	"github.com/gin-gonic/gin"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	"github.com/shahfaizansr/models/myhello"
	"github.com/shahfaizansr/utils/cvlconstant"
)

func HelloHandler(ctx *gin.Context, srv *service.Service) {

	logger := srv.LogHarbour.WithModule(cvlconstant.MODULE_MYHELLO).WithClass(cvlconstant.MyHello.String()).WithOp(cvlconstant.MyHello.String())

	// Log the incoming request using LogHarbour Utilities
	logger.Info().Log("Hello request received")

	data := myhello.MyHelloModels{
		Message: "Hello, Alya",
	}

	// Create a standardized response using Alya's response utilities
	response := wscutils.Response{
		Status:   wscutils.SuccessStatus,
		Data:     data,
		Messages: []wscutils.ErrorMessage{},
	}

	logger.Debug2().LogActivity("Saying hello", map[string]any{
		"message": data.Message,
	})
	logger.Info().Log("Hello response sent")

	// Send standardized JSON response using Alya's response utilities
	wscutils.SendSuccessResponse(ctx, &response)
}
