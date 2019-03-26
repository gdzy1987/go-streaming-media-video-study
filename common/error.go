package common

// 自定义错误返回
type ApiBody struct {
	Url     string `json:"url"`
	Method  string `json:"method"`
	ReqBody string `json:"req_body"`
}

type ErrorResponse struct {
	HttpSC int `json:"http_status_code"`
	Error  Err `json:"error"`
}

type Err struct {
	Error     string `json:"error"`
	ErrorCode string `json:"error_code"`
}

var (
	ErrorRequestBodyParseFailed = ErrorResponse{
		HttpSC: 400,
		Error: Err{
			Error:     "Request body is not correct",
			ErrorCode: "001",
		},
	}

	ErrorNotAuthUser = ErrorResponse{
		HttpSC: 401,
		Error: Err{
			Error:     "User authentication failed.",
			ErrorCode: "002",
		},
	}

	ErrorDBError = ErrorResponse{
		HttpSC: 500,
		Error: Err{
			Error:     "DB ops failed",
			ErrorCode: "003",
		},
	}

	ErrorInternalFaults = ErrorResponse{
		HttpSC: 500,
		Error: Err{
			Error:     "Internal service error",
			ErrorCode: "004",
		},
	}

	ErrorRequestNotRecognized = Err{
		Error:     "api not recognized, bad request",
		ErrorCode: "001",
	}
)
