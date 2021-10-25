package errors

const (
	invalidRequestFormat = 400400 + iota
	invalidID
	invalidStart
	invalidLimit
	invalidHubsID
	invalidPasscode
)

var (
	RoomInvalidRequestFormat = BadRequestError(invalidRequestFormat, "Invalid param: request param")
	RoomInvalidID            = BadRequestError(invalidID, "Invalid path: room_id")
	RoomInvalidStart         = BadRequestError(invalidStart, "Invalid param: start")
	RoomInvalidLimit         = BadRequestError(invalidLimit, "Invalid param: limit")
	RoomInvalidHubsID        = BadRequestError(invalidHubsID, "Invalid path: hubs_id")
	RoomInvalidPasscode      = BadRequestError(invalidPasscode, "Invalid content: passcode")
)
