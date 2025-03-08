package executor

type ExecuteRequest struct {
	FuncID    string `json:"funcID"`
	Parameter string `json:"parameter"`
}

type ExecuteTask struct {
	UserID  string `json:"userID"`
	ReplyTo string `json:"replyTo"`
	Image   string `json:"image"`
	Params  string `json:"params"`
}

type ExecuteResponse struct {
	ExecutionTime string `json:"executionTime"`
	Result        string `json:"result"`
	Error         string `json:"error"`
}
