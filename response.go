package middleware

type Response struct {
	Message     string `json:"message"`
	Data        any    `json:"data"`
	RequestUuid string `json:"request_uuid"`
}
