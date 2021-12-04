package data

// Datastructure descripting Planet
type Planet struct {
	Mass     float64 `json:"mass"`
	Location Vector3 `json:"location"`
	Speed    Vector3 `json:"speed"`
	force    Vector3
}

// Clear force for next speed update
func (a *Planet) clearForce() {
	a.force.X = 0
	a.force.Y = 0
	a.force.Z = 0
}

// Update force by interact with other planet
func (a *Planet) UpdateForce(b *Planet, g float64, softning float64) {
	// Calculate F = G * m_1 * m_2 /r^2
	delta := Sub(b.Location, a.Location)
	dist := Distance(a.Location, b.Location)
	invDist2 := 1.0 / (dist + softning)
	invDist2 *= invDist2
	force := Mul(delta, a.Mass*b.Mass*g*invDist2)

	// Update speed
	a.force = Add(a.force, force)
}

// Update Speed
func (a *Planet) UpdateSpeed(dt float64) {
	a.force = Mul(a.force, dt)
	a.Speed = Add(a.Speed, a.force)
	a.clearForce()
}

// Update Location
func (a *Planet) UpdateLocation(dt float64) {
	a.Location = Add(a.Location, Mul(a.Speed, dt))
}
