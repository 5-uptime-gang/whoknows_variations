package main

type StandardResponse struct {
	Data any `json:"data"`
}

type SearchResponse struct {
	Data []Page `json:"data"`
}

type AuthResponse struct {
	StatusCode *int    `json:"statusCode"`
	Message    *string `json:"message"`
}

type ValidationError struct {
	Loc  []any  `json:"loc"`
	Msg  string `json:"msg"`
	Type string `json:"type"`
}

type HTTPValidationError struct {
	Detail []ValidationError `json:"detail"`
}

type RequestValidationError struct {
	StatusCode int     `json:"statusCode" default:"422"`
	Message    *string `json:"message"`
}
