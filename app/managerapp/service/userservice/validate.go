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

const (
	minUsernameLength = 4
	maxUsernameLength = 191
	minFullnameLength = 3
	maxFullnameLength = 191
	minEmailLength    = 10
	maxEmailLength    = 255
	minPasswordLength = 7
	maxPasswordLength = 64
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
			validation.Length(minUsernameLength, maxUsernameLength).Error(servermsg.MsgInvalidLengthOfUsername),
			validation.Match(usernameRegex).Error(servermsg.MsgInvalidUsernameFormat)),
		validation.Field(&req.Password,
			validation.Required,
			validation.Length(minPasswordLength, maxPasswordLength).Error(servermsg.MsgInvalidLengthOfPassword)),
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
			validation.Length(minUsernameLength, maxUsernameLength).Error(servermsg.MsgInvalidLengthOfUsername),
			validation.Match(usernameRegex).Error(servermsg.MsgInvalidUsernameFormat)),
		validation.Field(&req.Password,
			validation.Required.Error(servermsg.MsgFieldRequired),
			validation.Length(minPasswordLength, maxPasswordLength).Error(servermsg.MsgInvalidLengthOfPassword)),
		validation.Field(&req.Fullname,
			validation.Required.Error(servermsg.MsgFieldRequired),
			validation.Length(minFullnameLength, maxFullnameLength).Error(servermsg.MsgInvalidLengthOfFullname)),
		validation.Field(&req.Email,
			validation.Required.Error(servermsg.MsgFieldRequired),
			validation.Length(minEmailLength, maxEmailLength).Error(servermsg.MsgInvalidLengthOfEmail),
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
			validation.Length(minUsernameLength, maxUsernameLength).Error(servermsg.MsgInvalidLengthOfUsername)),
		validation.Field(&req.Roles,
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

func (v Validate) UserUpdateFromSuperuserRequest(req UserUpdateFromSuperuserRequest) error {
	const op = "service.validate.UserUpdateFromSuperuserRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Username,
			validation.Required.Error(servermsg.MsgFieldRequired),
			validation.Length(minUsernameLength, maxUsernameLength).Error(servermsg.MsgInvalidLengthOfUsername),
			validation.Match(usernameRegex).Error(servermsg.MsgInvalidUsernameFormat)),
		validation.Field(&req.Fullname,
			validation.Required.Error(servermsg.MsgFieldRequired),
			validation.Length(minFullnameLength, maxFullnameLength).Error(servermsg.MsgInvalidLengthOfFullname)),
		validation.Field(&req.Email,
			validation.Required.Error(servermsg.MsgFieldRequired),
			validation.Length(minEmailLength, maxEmailLength).Error(servermsg.MsgInvalidLengthOfEmail),
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

func (v Validate) ChangePasswordRequest(req ChangePasswordRequest) error {
	const op = "service.validate.ChangePasswordRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.OldPassword,
			validation.Required.Error(servermsg.MsgFieldRequired),
			validation.Length(minPasswordLength, maxPasswordLength).Error(servermsg.MsgInvalidLengthOfPassword)),
		validation.Field(&req.NewPassword,
			validation.Required.Error(servermsg.MsgFieldRequired),
			validation.Length(minPasswordLength, maxPasswordLength).Error(servermsg.MsgInvalidLengthOfPassword)),
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
