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

func (v Validate) ValidateListConversationsRequest(req ListConversationsRequest) error {
	const op = "service.validate.ValidateListConversationsRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Statuses,
			validation.By(v.validateConversationStatus)),
	); err != nil {
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

	return nil
}

func (v Validate) validateConversationStatus(value interface{}) error {
	statuses, ok := value.([]ConversationStatus)
	if !ok {
		return errors.New(servermsg.MsgInvalidUserRole)
	}

	for _, status := range statuses {
		if !IsValidConversationStatus(status) {
			return errors.New(v.t.TranslateMessage(servermsg.MsgInvalidConversationStatus))
		}
	}

	return nil
}
