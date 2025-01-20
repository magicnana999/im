package protocol

type ApiRequest struct {
	Uri   string `json:"uri"`
	Param any    `json:"param"`
}

type ApiResponse struct {
	Uri    string `json:"uri"`
	Result any    `json:"result"`
}
