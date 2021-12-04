package data

// Datastructure holding input task
type Task struct {
	Step       float64  `json:"step"`
	G          float64  `json:"g"`
	Iterations int      `json:"iterations"`
	Softning   float64  `json:"softning"`
	Planets    []Planet `json:"planets"`
}
