package service

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/pkg/translation"
)

type Validate struct {
	t *translation.Translate
}

func NewValidate(t *translation.Translate) Validate {
	return Validate{
		t: t,
	}
}

// ValidateListConversationsRequest validates the request for listing conversations.
func (v Validate) ValidateListConversationsRequest(req ListConversationsRequest) error {
	const op = "service.validate.ValidateListConversationsRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Statuses,
			validation.By(v.validateConversationStatus)),
	); err != nil {
		return v.formatValidationErrors(op, err, req)
	}

	return nil
}

// ValidateClientMessage validates incoming websocket messages.
func (v Validate) ValidateClientMessage(req ClientMessage) error {
	const op = "service.validate.ValidateClientMessage"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Type,
			validation.Required.Error(servermsg.MsgFieldRequired),
			validation.By(v.validateMessageType)),
		validation.Field(&req.ConversationID, validation.Required.Error(servermsg.MsgFieldRequired)),
		validation.Field(&req.Content,
			validation.When(req.Type == MessageTypeText, validation.Required.Error(servermsg.MsgFieldRequired))),
		validation.Field(&req.SubType,
			validation.When(req.Type == MessageTypeSystem, validation.Required.Error(servermsg.MsgFieldRequired))),
	); err != nil {
		return v.formatValidationErrors(op, err, req)
	}

	return nil
}

// formatValidationErrors converts ozzo-validation errors into richerror.
func (v Validate) formatValidationErrors(op string, err error, req interface{}) error {
	fieldErrors := make(map[string]string)
	vErr := validation.Errors{}
	if errors.As(err, &vErr) {
		for key, value := range vErr {
			if value != nil {
				fieldErrors[key] = v.t.TranslateMessage(value.Error())
			}
		}
	}

	return richerror.New(op).WithMessage(servermsg.MsgInvalidInput).WithKind(richerror.KindInvalid).
		WithErrorFields(fieldErrors).WithMeta(map[string]interface{}{"req": req})
}

// validateConversationStatus is a custom validation rule for ConversationStatus slices.
func (v Validate) validateConversationStatus(value interface{}) error {
	statuses, ok := value.([]ConversationStatus)
	if !ok {
		return errors.New(servermsg.MsgInvalidInput)
	}

	for _, status := range statuses {
		if !IsValidConversationStatus(status) {
			return errors.New(v.t.TranslateMessage(servermsg.MsgInvalidConversationStatus))
		}
	}

	return nil
}

// validateMessageType is a custom validation rule for MessageType.
func (v Validate) validateMessageType(value interface{}) error {
	msgType, ok := value.(MessageType)
	if !ok {
		return errors.New(servermsg.MsgInvalidInput)
	}

	if !IsValidMessageType(msgType) {
		return errors.New(v.t.TranslateMessage(servermsg.MsgInvalidMessageType))
	}

	return nil
}
