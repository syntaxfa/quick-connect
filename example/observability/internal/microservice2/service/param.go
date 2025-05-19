package service

type GetCommentResponse struct {
	ID   uint64 `json:"id"`
	Body string `json:"body"`
}
