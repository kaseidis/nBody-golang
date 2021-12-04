package data

// Data structrue holding results
type Result struct {
	TimeStamp float64  `json:"timestamp"`
	Planets   []Planet `json:"planets"`
}
