package dto

type BearerTokenMiddlewareRequest struct {
	Token string `header:"Authorization" binding:"required,BearerTokenValidator"`
}
