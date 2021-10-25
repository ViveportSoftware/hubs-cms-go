package dto

import "encoding/json"

type Page struct {
	Prev string `json:"prev"`
	Next string `json:"next"`
}

type DirectusMeta struct {
	FilterCount int64 `json:"filter_count"`
	TotalCount  int64 `json:"total_count"`
}

type PageRequestParam struct {
	Limit json.Number `form:"limit" binding:"omitempty,PageLimitValidator"`
	Start json.Number `form:"start" binding:"omitempty,PageStartValidator"`
}
