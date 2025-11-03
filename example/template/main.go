package main

import (
	"bytes"
	"fmt"
	htmlTmpl "html/template"
	textTmpl "text/template"
)

var textTemplate = `login code is: {{.otp_code}}

quick connect`

var htmlTemplate = `
<!DOCTYPE html>
<html>
<body>

<h1>otp code:</h1>
<p>{{.otp_code}}</p>

</body>
</html>`

func main() {
	fmt.Println(textTemplate)

	data := map[string]string{
		"otp_code": "123456",
	}

	text, err := textTmpl.New("template").Parse(textTemplate)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	if err = text.Execute(&buf, data); err != nil {
		panic(err)
	}

	fmt.Println(buf.String())

	html, err := htmlTmpl.New("template").Parse(htmlTemplate)
	if err != nil {
		panic(err)
	}

	var htmlBuf bytes.Buffer
	if err := html.Execute(&htmlBuf, data); err != nil {
		panic(err)
	}

	fmt.Println(htmlBuf.String())
}
