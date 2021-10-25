package errors

const (
	eventInvalidRequestFormat = 400100 + iota
	eventInvalidID
	eventInvalidStart
	eventInvalidLimit
)

var (
	EventInvalidRequestFormat = BadRequestError(eventInvalidRequestFormat, "Invalid param: request param")
	EventInvalidID            = BadRequestError(eventInvalidID, "Invalid path: event_id")
	EventInvalidStart         = BadRequestError(eventInvalidStart, "Invalid param: start")
	EventInvalidLimit         = BadRequestError(eventInvalidLimit, "Invalid param: limit")
)
