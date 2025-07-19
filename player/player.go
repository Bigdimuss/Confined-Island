package player

import "confinedisland/config"

type Player struct {
	Name             string
	X, Y             float64
	WorldX           int
	WorldY           int
	MoveSpeed        float64
	ProcessMoveSpeed float64
	Moving           bool // Indique si le joueur est en mouvement
	MoveTargetX      int  // Cible de position X dans le monde
	MoveTargetY      int  // Cible de position Y dans le monde
	Orientation      string
}

func NewPlayer(name string) *Player {
	p := &Player{Name: name}
	return p
}

func (p *Player) GetCenterCoord() (int, int) {
	return int(p.X) + config.UNITE/2, int(p.Y) + config.UNITE/2
}

func (p *Player) UpdateOrientation(mouseX int, mouseY int) {
	centerX, centerY := p.GetCenterCoord()
	angle := get_angle(false, centerX, centerY, mouseX, mouseY)
	switch {
	case angle <= 112.5 && angle >= 67.5:
		p.Orientation = "nord"
	case angle <= 67.5 && angle > 22.5:
		p.Orientation = "nord-est"
	case angle < 22.5 && angle > -22.5:
		p.Orientation = "est"
	case angle < -22.5 && angle > -67.5:
		p.Orientation = "sud-est"
	case angle < -67.5 && angle > -112.5:
		p.Orientation = "sud"
	case angle < -112.5 && angle > -157.5:
		p.Orientation = "sud-ouest"
	case angle < 157.5 && angle > 112.5:
		p.Orientation = "nord-ouest"
	case (angle < -157.5 && angle >= -180) || (angle > 157.5 && angle <= 180):
		p.Orientation = "ouest"
	default:
		p.Orientation = "inconnu" // Au cas où l'angle ne correspond à aucun cas
	}
}
func (p *Player) Update(tps float64) {
	p.ProcessMoveSpeed += p.MoveSpeed * 60 / float64(tps) // Ajustez la vitesse ici
	if p.ProcessMoveSpeed >= 1 {
		p.ProcessMoveSpeed = 1
		p.Moving = false
	}
	if !p.Moving {
		p.WorldX = p.MoveTargetX
		p.WorldY = p.MoveTargetY
		// Mettre à jour l'écran avec les nouvelles coordonnées
	}
}
