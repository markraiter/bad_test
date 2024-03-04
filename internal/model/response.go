package model

type Response struct {
	Message string `json:"message" example:"response message"`
}

type TaskResult struct {
	Max              int     `json:"max" example:"100"`
	Min              int     `json:"min" example:"1"`
	Median           float64 `json:"median" example:"50.5"`
	Avg              float64 `json:"avg" example:"50.5"`
	MaxIncreasingSeq []int   `json:"max_increasing_seq" example:"[-4390, -503, 3, 16, 5032]"`
	MaxDecreasingSeq []int   `json:"max_decreasing_seq" example:"[5032, 16, 3, -503, -4390]"`
	Time             string  `json:"time" example:"1s"`
}
