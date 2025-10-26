package http

import "github.com/syntaxfa/quick-connect/pkg/translation"

type Handler struct {
	t *translation.Translate
}

func NewHandler(t *translation.Translate) Handler {
	return Handler{
		t: t,
	}
}
