package main

import (
	"fmt"
	"github.com/syntaxfa/quick-connect/app/notificationapp/service"
)

var textTemplate string = `login code is: {{.otp_code}}

quick connect`

var textTemplate2 string = `hello dear {{.username}}`

var htmlTemplate string = `
<!DOCTYPE html>
<html>
<body>

<h1>otp code:</h1>
<p>{{.otp_code}}</p>

</body>
</html>`

var htmlTemplate2 string = `
<!DOCTYPE html>
<html>
<body>

<h1>Hello mr {{.username}}</h1>

</body>
</html>`

func main() {
	render := service.NewRenderService()

	for i := 1; i < 100; i++ {
		body, rErr := render.RenderTemplate("opt_code:sms", textTemplate, service.TemplateTypeText, map[string]string{"otp_code": "123458"})
		if rErr != nil {
			panic(rErr)
		}

		fmt.Println(body)
	}

	for i := 1; i < 100; i++ {
		body, rErr := render.RenderTemplate("opt_code:email", htmlTemplate, service.TemplateTypeHTML, map[string]string{"otp_code": "1234569"})
		if rErr != nil {
			panic(rErr)
		}

		fmt.Println(body)
	}

	for i := 1; i < 100; i++ {
		body, rErr := render.RenderTemplate("welcome:sms", textTemplate2, service.TemplateTypeText, map[string]string{"username": "alireza"})
		if rErr != nil {
			panic(rErr)
		}

		fmt.Println(body)
	}

	for i := 1; i < 100; i++ {
		body, rErr := render.RenderTemplate("welcome:email", htmlTemplate2, service.TemplateTypeHTML, map[string]string{"username": "alireza"})
		if rErr != nil {
			panic(rErr)
		}

		fmt.Println(body)
	}
}
