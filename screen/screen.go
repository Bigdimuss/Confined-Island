package screen

import (
	"confinedisland/config"
	"confinedisland/generator"
	"confinedisland/generator/island"
	"confinedisland/player"
	"confinedisland/sprite"
)

type Screen struct {
	X, Y           int
	Width, Height  int
	Size_x, Size_y int
	Background     [][]*sprite.Block
	player         *player.Player
	world          *island.Island
}

func NewScreen(height int, width int, player *player.Player, world *island.Island) *Screen {
	size_x, size_y := int(width/config.UNITE), int(height/config.UNITE)
	background := make([][]*sprite.Block, size_y+4) // +4 pour les 2 blocs supplémentaires en haut et en bas
	for i := range background {
		background[i] = make([]*sprite.Block, size_x+4) // +4 pour l 2 blocs supplémentaires à gauche et à droite
	}
	s := Screen{X: 0, Y: 0, Height: height, Width: width, Size_x: size_x, Size_y: size_y, Background: background, player: player, world: world}
	return &s
}

func (s *Screen) GetGroundUnderPlayer() sprite.Block {
	y := int(s.player.Y/32) + 2
	x := int(s.player.X/32) + 2
	if s.Background[y][x] != nil {
		return *(s.Background)[y][x]
	} else {
		return sprite.Block{Name: "NO BLOCK !"}
	}

}

func (s *Screen) Update(targetX int, targetY int, moveSpeed float64) {
	// Mettre à jour les coordonnées de la caméra
	if s.player.Moving {
		s.X = int(sprite.Lerp(float64(s.X), float64(targetX*32), moveSpeed))
		s.Y = int(sprite.Lerp(float64(s.Y), float64(targetY*32), moveSpeed))
	}

	y_map, x_map := s.player.WorldY-s.Size_y/2, s.player.WorldX-s.Size_x/2

	// Charger les sprites, incluant 2 blocs hors champ
	for i := -2; i < s.Size_y+2; i++ {
		yScreen := i + 2
		for j := -2; j < s.Size_x+2; j++ {
			xScreen := j + 2
			mapY := y_map + i
			mapX := x_map + j + 1
			if mapY >= 0 && mapY < s.world.Height && mapX >= 0 && mapX < s.world.Width {
				s.Background[yScreen][xScreen] = s.world.Background[mapY][mapX]

			} else {
				s.Background[yScreen][xScreen] = &sprite.Block{
					Name: "limite",
					Sprite: &sprite.StaticSprite{
						Width:  float64(config.UNITE),
						Height: float64(config.UNITE),
					},
					BaseColor: generator.TEMPLATE_GROUND_RESSOURCES.Biomes["limite"].Color,
				}
			}
		}
	}
}
