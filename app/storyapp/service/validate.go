package service

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/pkg/translation"
)

const (
	minTitleLength    = 3
	maxTitleLength    = 255
	minCaptionLength  = 3
	maxCaptionLength  = 2000
	minLinkURLLength  = 10
	maxLinkURLLength  = 255
	minLinkTextLength = 3
	maxLinkTextLength = 100
)

type Validate struct {
	t *translation.Translate
}

func NewValidate(t *translation.Translate) Validate {
	return Validate{
		t: t,
	}
}

func (v Validate) ValidateAddStoryRequest(req AddStoryRequest) error {
	const op = "service.validate.ValidateAddStoryRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.MediaFileID,
			validation.Required.Error(servermsg.MsgFieldRequired)),
		validation.Field(&req.Title,
			validation.Length(minTitleLength, maxTitleLength).Error(servermsg.MsgInvalidLengthOfStoryTitle)),
		validation.Field(&req.Caption,
			validation.Length(minCaptionLength, maxCaptionLength).Error(servermsg.MsgInvalidLengthOfStoryCaption)),
		validation.Field(&req.LinkURL,
			validation.Length(minLinkURLLength, maxLinkURLLength).Error(servermsg.MsgInvalidLengthOfStoryLinkURL)),
		validation.Field(&req.LinkText,
			validation.Length(minLinkTextLength, maxLinkTextLength).Error(servermsg.MsgInvalidLengthOfStoryLinkText)),
		validation.Field(&req.DurationSeconds,
			validation.Required.Error(servermsg.MsgFieldRequired)),
		validation.Field(&req.PublishAt,
			validation.Required.Error(servermsg.MsgFieldRequired)),
		validation.Field(&req.ExpiresAt,
			validation.Required.Error(servermsg.MsgFieldRequired)),
	); err != nil {
		fieldErrors := make(map[string]string)

		vErr := validation.Errors{}
		if errors.As(err, &vErr) {
			for key, value := range vErr {
				if value != nil {
					fieldErrors[key] = v.t.TranslateMessage(value.Error())
				}
			}

			return richerror.New(op).WithMessage(servermsg.MsgInvalidInput).WithKind(richerror.KindInvalid).
				WithErrorFields(fieldErrors).WithMeta(map[string]interface{}{"req": req})
		}
	}

	return nil
}
