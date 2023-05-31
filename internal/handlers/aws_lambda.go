package handlers

import (
	"context"
	"github.com/BaronBonet/content-generator/internal/core/ports"
)

type AWSLambdaEventHandler struct {
	logger ports.Logger
	srv    ports.Service
}

func NewAWSLambdaEventHandler(logger ports.Logger, srv ports.Service) *AWSLambdaEventHandler {
	return &AWSLambdaEventHandler{logger: logger, srv: srv}
}

func (handler *AWSLambdaEventHandler) HandleEvent(ctx context.Context, request interface{}) {
	err := handler.srv.GenerateNewsContent(ctx)
	if err != nil {
		handler.logger.Fatal("Error while generating news content.", "error", err)
	}
}
