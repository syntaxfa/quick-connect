package service

import (
	"context"
	"fmt"

	"github.com/oklog/ulid/v2"
)

func (s Service) SendNotification(_ context.Context, req SendNotificationRequest) (SendNotificationResponse, error) {
	if vErr := s.vld.ValidateSendNotificationRequest(req); vErr != nil {
		return SendNotificationResponse{}, vErr
	}

	// validate user_id
	_, pErr := ulid.Parse(string(req.UserID))
	if pErr != nil {
		fmt.Println("not ok")
	}

	// save notification

	// check notification type if is critical, send it and return response

	return SendNotificationResponse{}, nil
}
