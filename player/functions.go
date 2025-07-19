package player

import "math"

func get_angle(radian bool, centerPersoX int, centerPersoY int, positionCibleX int, positionCibleY int) float64 {
	dx := float64(positionCibleX - centerPersoX)
	dy := float64(positionCibleY - centerPersoY)
	rads := math.Atan2(-dy, dx)
	if radian {
		return rads
	}
	return rads * (180.0 / math.Pi) // conversion en degrés
}
