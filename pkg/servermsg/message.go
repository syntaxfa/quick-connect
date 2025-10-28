package servermsg

const (
	MsgSomethingWentWrong                  = "something went wrong"
	MsgInvalidToken                        = "invalid token"
	MsgInvalidTokenAlgorithm               = "invalid token algorithm"
	MsgFieldRequired                       = "this field is required"
	MsgInvalidLengthOfUsername             = "the username must be between 6 and 191 characters"
	MsgConflictUsername                    = "This username already exists"
	MsgInvalidInput                        = "invalid input"
	MsgInvalidLengthOfPassword             = "the password must be between 8 and 64 characters"
	MsgInvalidLengthOfFullname             = "the fullname must be between 3 and 191 characters"
	MsgInvalidUserRole                     = "invalid user role"
	MsgRecordNotFound                      = "record not found"
	MsgInvalidLengthOfUserID               = "the user id must be less than 255 characters"
	MsgInvalidNotificationType             = "invalid notification type"
	MsgInvalidNotificationChannelDelivery  = "invalid notification channel delivery"
	MsgConflictNotificationChannelDelivery = "channel delivery has conflict"
	MsgPageSizeMin                         = "page size must be greater than 0"
	MsgPageMin                             = "page must be greater than 0"
	MsgInvalidLengthOfTemplateName         = "the name must be between 1 and 255"
	MsgInvalidTemplateContent              = "invalid template contents"
	MsgConflictTemplateChannel             = "template channel has conflict"
	MsgConflictTemplateChannelLang         = "template content channel language has conflict"
	MsgConflictTemplate                    = "template is already exists"
	MsgTemplateNotFound                    = "this template does not exist"
	MsgInvalidIgnoreChannel                = "invalid ignore channel"

	// Admin app

	MsgUsernameAndPasswordAreRequired = "username and password are required"
)
