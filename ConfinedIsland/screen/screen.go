package screen

import (
	"confinedisland/generator/island"
	"confinedisland/player"
	"confinedisland/sprite"
	"image/color"
)

type Screen struct {
	X, Y           int
	Width, Height  int
	size_x, size_y int
	Background     [][]sprite.Sprite
	player         *player.Player
	world          *island.Island
}

func NewScreen(height int, width int, player *player.Player, world *island.Island) *Screen {
	size_x, size_y := int(width/32), int(height/32)
	background := make([][]sprite.Sprite, size_y) // Correction : utilisez size_y pour les lignes
	for i := range background {
		background[i] = make([]sprite.Sprite, size_x) // Correction : utilisez size_x pour les colonnes
	}
	s := Screen{X: 0, Y: 0, Height: height, Width: width, size_x: size_x, size_y: size_y, Background: background, player: player, world: world}
	return &s
}
func (s *Screen) Update() {
	height := s.size_y // Utilisez size_y pour la hauteur
	width := s.size_x  // Utilisez size_x pour la largeur
	cord_y := int(float64(s.player.WorldY) - float64(height)/2)
	cord_x := int(float64(s.player.WorldX) - float64(width)/2)

	y_map, x_map := cord_y, cord_x

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if y_map >= 0 && y_map < s.world.Height && x_map >= 0 && x_map < s.world.Width {
				s.Background[i][j] = sprite.Sprite{
					Color:  s.world.Background[y_map][x_map],
					X:      float64(j * 32),
					Y:      float64(i * 32),
					Width:  32,
					Height: 32,
				}
			} else {
				s.Background[i][j] = sprite.Sprite{
					Color:  color.RGBA{R: 0, B: 0, G: 0, A: 250},
					X:      float64(j * 32),
					Y:      float64(i * 32),
					Width:  32,
					Height: 32,
				}

			}
			x_map++ // Incrémente x_map
		}
		y_map++        // Incrémente y_map
		x_map = cord_x // Réinitialise x_map pour la prochaine ligne
	}
}
