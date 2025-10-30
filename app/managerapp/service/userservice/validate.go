package userservice

import (
	"errors"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/pkg/translation"
	"github.com/syntaxfa/quick-connect/types"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
var usernameRegex = regexp.MustCompile(`^[\x21-\x7E]+$`)

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
			validation.Length(4, 191).Error(servermsg.MsgInvalidLengthOfUsername)),
		validation.Field(&req.Password,
			validation.Required,
			validation.Length(7, 64).Error(servermsg.MsgInvalidLengthOfPassword)),
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
			WithMeta(map[string]interface{}{"req": req}).WithErrorFields(fieldErrors)
	}

	return nil
}

func (v Validate) ValidateUserCreateRequest(req UserCreateRequest) error {
	const op = "service.validate.ValidateLoginRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Username,
			validation.Required.Error(servermsg.MsgFieldRequired),
			validation.Length(4, 191).Error(servermsg.MsgInvalidLengthOfUsername),
			validation.Match(usernameRegex).Error(servermsg.MsgInvalidUsernameFormat)),
		validation.Field(&req.Password,
			validation.Required.Error(servermsg.MsgFieldRequired),
			validation.Length(7, 64).Error(servermsg.MsgInvalidLengthOfPassword)),
		validation.Field(&req.Fullname,
			validation.Required.Error(servermsg.MsgFieldRequired),
			validation.Length(3, 191).Error(servermsg.MsgInvalidLengthOfFullname)),
		validation.Field(&req.Email,
			validation.Required.Error(servermsg.MsgFieldRequired),
			validation.Length(10, 255).Error(servermsg.MsgInvalidLengthOfEmail),
			validation.Match(emailRegex).Error(servermsg.MsgInvalidEmailFormat)),
		validation.Field(&req.PhoneNumber,
			validation.Required.Error(servermsg.MsgFieldRequired)),
		validation.Field(&req.Roles,
			validation.Required.Error(servermsg.MsgFieldRequired),
			validation.By(v.validateUserRole)),
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

func (v Validate) validateUserRole(value interface{}) error {
	roles, ok := value.([]types.Role)
	if !ok {
		return errors.New(servermsg.MsgInvalidUserRole)
	}

	for _, role := range roles {
		if !types.IsValidRole(role) {
			return errors.New(servermsg.MsgInvalidUserRole)
		}
	}

	return nil
}

func (v Validate) ValidateListUserRequest(req ListUserRequest) error {
	const op = "validate.ValidateListUserRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Username,
			validation.Length(4, 191).Error(servermsg.MsgInvalidLengthOfUsername)),
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

func (v Validate) UserUpdateFromSuperuserRequest(req UserUpdateFromSuperuserRequest) error {
	const op = "service.validate.UserUpdateFromSuperuserRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Username,
			validation.Required.Error(servermsg.MsgFieldRequired),
			validation.Length(4, 191).Error(servermsg.MsgInvalidLengthOfUsername),
			validation.Match(usernameRegex).Error(servermsg.MsgInvalidUsernameFormat)),
		validation.Field(&req.Fullname,
			validation.Required.Error(servermsg.MsgFieldRequired),
			validation.Length(3, 191).Error(servermsg.MsgInvalidLengthOfFullname)),
		validation.Field(&req.Email,
			validation.Required.Error(servermsg.MsgFieldRequired),
			validation.Length(10, 255).Error(servermsg.MsgInvalidLengthOfEmail),
			validation.Match(emailRegex).Error(servermsg.MsgInvalidEmailFormat)),
		validation.Field(&req.PhoneNumber,
			validation.Required.Error(servermsg.MsgFieldRequired)),
		validation.Field(&req.Roles,
			validation.Required.Error(servermsg.MsgFieldRequired),
			validation.By(v.validateUserRole)),
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
