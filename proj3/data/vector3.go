package data

import "math"

// Simple vector3 data structure
type Vector3 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// Add vectors
func Add(a Vector3, b Vector3) Vector3 {
	return Vector3{
		X: a.X + b.X,
		Y: a.Y + b.Y,
		Z: a.Z + b.Z,
	}
}

// Subtract vectors
func Sub(a Vector3, b Vector3) Vector3 {
	return Vector3{
		X: a.X - b.X,
		Y: a.Y - b.Y,
		Z: a.Z - b.Z,
	}
}

// Multiply vectors
func Mul(a Vector3, b float64) Vector3 {
	return Vector3{
		X: a.X * b,
		Y: a.Y * b,
		Z: a.Z * b,
	}
}

// Distance function
func Distance(a Vector3, b Vector3) float64 {
	delta := Sub(a, b)
	return math.Sqrt(delta.X*delta.X + delta.Y*delta.Y + delta.Z*delta.Z)
}
