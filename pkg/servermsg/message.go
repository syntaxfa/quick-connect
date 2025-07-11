package servermsg

const (
	MsgSomethingWentWrong                 = "something went wrong"
	MsgInvalidToken                       = "invalid token"
	MsgInvalidTokenAlgorithm              = "invalid token algorithm"
	MsgFieldRequired                      = "this field is required"
	MsgInvalidLengthOfUsername            = "the username must be between 6 and 191 characters"
	MsgInvalidInput                       = "invalid input"
	MsgInvalidLengthOfPassword            = "the password must be between 8 and 191 characters"
	MsgRecordNotFound                     = "record not found"
	MsgInvalidLengthOfUserID              = "the user id must be less than 255 characters"
	MsgInvalidLengthOfNotificationTitle   = "the title must be less than 255 characters"
	MsgInvalidNotificationType            = "invalid notification type"
	MsgInvalidNotificationChannelDelivery = "invalid notification channel delivery"
	MsgPageSizeMin                        = "page size must be greater than 0"
	MsgPageMin                            = "page must be greater than 0"
	MsgInvalidLengthOfTemplateName        = "the name must be between 1 and 255"
	MsgInvalidTemplateBody                = "invalid template bodies"
	MsgConflictTemplateBody               = "body language and channel cannot be the same as another body"
	MsgConflictTemplate                   = "template is already exists"
	MsgTemplateNotFound                   = "this template does not exist"
	MsgInvalidIgnoreChannel               = "invalid ignore channel"
)
