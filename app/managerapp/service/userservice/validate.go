package userservice

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

func (v Validate) ValidateLoginRequest(req UserLoginRequest) error {
	const op = "service.validate.ValidateLoginRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Username,
			validation.Required.Error(servermsg.MsgFieldRequired),
			validation.Length(6, 191).Error(servermsg.MsgInvalidLengthOfUsername)),
		validation.Field(&req.Password,
			validation.Required,
			validation.Length(8, 191).Error(servermsg.MsgInvalidLengthOfPassword)),
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

		return richerror.New(op).WithMessage(servermsg.MsgInvalidInput).WithKind(richerror.KindUnexpected).
			WithMeta(map[string]interface{}{"req": req}).WithErrorFields(fieldErrors)
	}

	return nil
}
