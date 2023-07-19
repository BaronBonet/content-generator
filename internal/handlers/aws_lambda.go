package handlers

import (
	"context"

	"github.com/BaronBonet/content-generator/internal/core/ports"
	"github.com/BaronBonet/go-logger/logger"
)

type AWSLambdaEventHandler struct {
	logger logger.Logger
	srv    ports.Service
}

func NewAWSLambdaEventHandler(logger logger.Logger, srv ports.Service) *AWSLambdaEventHandler {
	return &AWSLambdaEventHandler{logger: logger, srv: srv}
}

func (handler *AWSLambdaEventHandler) HandleEvent(ctx context.Context, request interface{}) {
	err := handler.srv.GenerateNewsContent(ctx)
	if err != nil {
		handler.logger.Fatal("Error while generating news content.", "error", err)
	}
}
