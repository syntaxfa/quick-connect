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
		validation.Field(&req.ExternalUserID,
			validation.Required.Error(servermsg.MsgFieldRequired),
			validation.Length(1, 255).Error(servermsg.MsgInvalidLengthOfUserID),
		),
		validation.Field(&req.Type,
			validation.Required.Error(servermsg.MsgFieldRequired),
			validation.By(v.ValidateNotificationType),
		),
		validation.Field(&req.TemplateName,
			validation.Required.Error(servermsg.MsgFieldRequired),
			validation.Length(1, 255).Error(servermsg.MsgInvalidLengthOfTemplateName),
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
	notificationType, ok := value.(NotificationType)
	if !ok {
		return errors.New(servermsg.MsgInvalidNotificationType)
	}

	if !IsValidNotificationType(notificationType) {
		return errors.New(servermsg.MsgInvalidNotificationType)
	}

	return nil
}

func (v Validate) ValidateNotificationChannelDeliveries(value interface{}) error {
	channelDeliveries, ok := value.([]ChannelDeliveryRequest)
	if !ok {
		return errors.New(servermsg.MsgInvalidNotificationChannelDelivery)
	}

	if len(channelDeliveries) < 1 {
		return errors.New(servermsg.MsgInvalidNotificationChannelDelivery)
	}

	for _, channel := range channelDeliveries {
		if !IsValidChannelType(channel.Channel) {
			return errors.New(servermsg.MsgInvalidNotificationChannelDelivery)
		}
	}

	return nil
}

func (v Validate) ValidateListNotificationRequest(req ListNotificationRequest) error {
	const op = "validate.ValidateListNotificationRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.ExternalUserID,
			validation.Required.Error(servermsg.MsgFieldRequired),
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

func (v Validate) ValidateAddTemplateRequest(req AddTemplateRequest) error {
	const op = "validate.ValidateAddTemplateRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Name,
			validation.Required.Error(servermsg.MsgFieldRequired),
			validation.Length(1, 255).Error(servermsg.MsgInvalidLengthOfTemplateName),
		),
		validation.Field(&req.Bodies,
			validation.Required.Error(servermsg.MsgFieldRequired),
			validation.By(v.ValidateTemplateBodies),
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

func (v Validate) ValidateTemplateBodies(value interface{}) error {
	bodies, ok := value.([]TemplateBody)
	if !ok {
		return errors.New(servermsg.MsgInvalidTemplateBody)
	}

	for index, body := range bodies {
		for i := index; i < (len(bodies) - 1); i++ {
			if body.Lang == bodies[i+1].Lang && body.Channel == bodies[i+1].Channel {
				return errors.New(servermsg.MsgConflictTemplateBody)
			}
		}
	}

	for _, body := range bodies {
		if !IsValidChannelType(body.Channel) {
			return errors.New(servermsg.MsgInvalidNotificationChannelDelivery)
		}
	}

	return nil
}

func (v Validate) ValidateUpdateUserSettingsRequest(req UpdateUserSettingRequest) error {
	const op = "validate.ValidateUpdateUserNotificationSettingsRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Lang,
			validation.Required.Error(servermsg.MsgFieldRequired)),
		validation.Field(&req.IgnoreChannels,
			validation.Required.Error(servermsg.MsgFieldRequired),
			validation.By(v.validateIgnoreChannel),
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

func (v Validate) validateIgnoreChannel(value interface{}) error {
	ignoreChannels, ok := value.([]IgnoreChannel)
	if !ok {
		return errors.New(servermsg.MsgInvalidIgnoreChannel)
	}

	for _, ignore := range ignoreChannels {
		if !IsValidChannelType(ignore.Channel) {
			return errors.New(servermsg.MsgInvalidNotificationChannelDelivery)
		}

		for _, nt := range ignore.NotificationTypes {
			if !IsValidNotificationType(nt) {
				return errors.New(servermsg.MsgInvalidNotificationType)
			}
		}
	}

	return nil
}
