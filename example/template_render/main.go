package main

import (
	"fmt"
	"github.com/syntaxfa/quick-connect/app/notificationapp/service"
	"time"
)

var TextOtpCodeEn = `login code is: {{.otp_code}}

quick connect`

var TextOtpCodeFa = `کد ورود شما است: {{.otp_code}}

کوئیک کانکت`

var EmailOtpCodeEn string = `
<!DOCTYPE html>
<html>
<body>

<h1>otp code:</h1>
<p>{{.otp_code}}</p>

</body>
</html>`

var EmailOtpCodeFa string = `
<!DOCTYPE html>
<html>
<body>

<h1>کد ورود:</h1>
<p>{{.otp_code}}</p>

</body>
</html>`

var otpCodeTemplate = service.Template{
	ID:   "ssmsms",
	Name: "otp_code",
	Contents: []service.TemplateContent{
		{Channel: service.ChannelTypeSMS, Bodies: []service.ContentBody{
			{
				Lang:  "en",
				Body:  TextOtpCodeEn,
				Title: "your opt code",
			},
			{
				Lang:  "fa",
				Body:  TextOtpCodeFa,
				Title: "کد ورود",
			},
		}},
		{Channel: service.ChannelTypeEmail, Bodies: []service.ContentBody{
			{
				Lang:  "en",
				Body:  EmailOtpCodeEn,
				Title: "your opt code",
			},
			{
				Lang:  "fa",
				Body:  EmailOtpCodeFa,
				Title: "کد ورود",
			},
		}},
	},
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
}

var textTemplate2 string = `hello dear {{.username}}`

var htmlTemplate2 string = `
<!DOCTYPE html>
<html>
<body>

<h1>Hello mr {{.username}}</h1>

</body>
</html>`

func main() {
	render := service.NewRenderService("fa")

	notification := service.Notification{
		ID:                "ssse",
		UserID:            "ewvwr2",
		Type:              service.NotificationTypeInfo,
		Data:              nil,
		TemplateName:      "otp_code",
		DynamicBodyData:   map[string]string{"otp_code": "123458"},
		DynamicTitleData:  nil,
		IsRead:            false,
		IsInApp:           false,
		CreatedAt:         time.Time{},
		ChannelDeliveries: nil,
		OverallStatus:     "",
	}

	for i := 1; i < 100; i++ {
		render, rErr := render.RenderTemplate(otpCodeTemplate, service.ChannelTypeEmail, "fa", notification.DynamicTitleData, notification.DynamicBodyData)
		if rErr != nil {
			panic(rErr)
		}

		fmt.Println("title:\n", render.Title)
		fmt.Println("body:\n", render.Body)
	}
}
