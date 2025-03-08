package functions

type RegisterRequest struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type RegisterResponse struct {
	FunctionIdentifier string `json:"functionIdentifier"`
}
