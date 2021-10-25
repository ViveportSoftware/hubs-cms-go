package utils

import (
	"fmt"
)

func GetPageParam(start int64, limit int64) string {

	pageParam := ""

	if start > 0 {
		pageParam = fmt.Sprintf("offset=%d", start)
	}
	if limit > 0 {
		pageParam = fmt.Sprintf("%s&limit=%d", pageParam, limit)
	}

	return pageParam
}
