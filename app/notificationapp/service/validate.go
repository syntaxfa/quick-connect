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

func (v Validate) ValidateSendNotificationRequest(req SendNotificationRequest) error {
	const op = "service.validate.ValidateSendNotificationRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.UserID,
			validation.Required.Error(servermsg.MsgFieldRequired),
			validation.Length(1, 255).Error(servermsg.MsgInvalidLengthOfUserID),
		),
		validation.Field(&req.Type,
			validation.Required.Error(servermsg.MsgFieldRequired),
			validation.By(v.ValidateNotificationType),
		),
		validation.Field(&req.Title,
			validation.Required.Error(servermsg.MsgFieldRequired),
			validation.Length(1, 255).Error(servermsg.MsgInvalidLengthOfNotificationTitle),
		),
		validation.Field(&req.Body,
			validation.Required.Error(servermsg.MsgFieldRequired),
		),
		validation.Field(&req.ChannelDeliveries,
			validation.Required.Error(servermsg.MsgFieldRequired),
			validation.By(v.ValidateNotificationChannelDeliveries),
		),
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

func (v Validate) ValidateNotificationType(value interface{}) error {
	notificationType, ok := value.(string)
	if !ok {
		return errors.New(servermsg.MsgInvalidNotificationType)
	}

	if !IsValidNotificationType(notificationType) {
		return errors.New(servermsg.MsgInvalidNotificationType)
	}

	return nil
}

func (v Validate) ValidateNotificationChannelDeliveries(value interface{}) error {
	channelDeliveries, ok := value.([]ChannelDelivery)
	if !ok {
		return errors.New(servermsg.MsgInvalidNotificationChannelDelivery)
	}

	if len(channelDeliveries) < 1 {
		return errors.New(servermsg.MsgInvalidNotificationChannelDelivery)
	}

	for _, channel := range channelDeliveries {
		if !IsValidChannelType(string(channel.Channel)) {
			return errors.New(servermsg.MsgInvalidNotificationChannelDelivery)
		}
	}

	return nil
}
