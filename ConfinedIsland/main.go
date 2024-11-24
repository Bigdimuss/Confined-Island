/*package main

import (
	"confinedisland/generator/island"
	"fmt"
)

func main() {
	// Définir les dimensions de l'image
	width, height := 1500, 1500
	config := island.IslandConfig{Width: width, Height: height}
	island := island.NewIsland(config)
	island.Create_image()
	fmt.Println(island.Seed)
}*/

package main

import (
	"confinedisland/generator/island"
	"confinedisland/player"
	"confinedisland/screen"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

type Game struct {
	world  *island.Island
	scene  *screen.Screen
	player *player.Player
}

func NewGame(world *island.Island, scene *screen.Screen, player *player.Player) *Game {
	return &Game{world: world, scene: scene, player: player}
}

func (g *Game) Update() error {
	// Vérification des touches fléchées
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		if g.player.WorldY-1 >= 0 {
			g.player.WorldY -= 1
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		if g.player.WorldY+1 < g.world.Height {
			g.player.WorldY += 1
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		if g.player.WorldX-1 >= 0 {
			g.player.WorldX -= 1
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		if g.player.WorldX+1 < g.world.Width {
			g.player.WorldX += 1
		}
	}
	g.scene.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255}) // Fond noir
	// Dessine un carré rouge
	i := ebiten.NewImage(g.scene.Width, g.scene.Height)
	i.Fill(color.RGBA{R: 150, G: 50, B: 150, A: 250})
	o := &ebiten.DrawImageOptions{}
	o.GeoM.Translate(0, 0)
	screen.DrawImage(i, o)
	for _, item := range g.scene.Background {
		for _, v := range item {
			r := ebiten.NewImage(32, 32)
			r.Fill(v.Color) // Carré rouge
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(v.X, v.Y)
			screen.DrawImage(r, op)
			fmt.Println("Draw !")
		}
	}
	r := ebiten.NewImage(32, 32)
	r.Fill(color.RGBA{R: 250, G: 0, B: 0, A: 50}) // Carré rouge
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.player.X, g.player.Y)
	screen.DrawImage(r, op)
	active := g.world.Background[g.player.WorldY][g.player.WorldX]
	xperso, yperso := int(g.player.X/32), int(g.player.Y/32)
	activeS := g.scene.Background[yperso][xperso]
	fmt.Printf("You are on: %v", active)
	fmt.Printf("cord: %v - %v", g.player.WorldX, g.player.WorldY)
	fmt.Printf("You are on: %v", activeS)
	fmt.Printf("cord: %v - %v", xperso, yperso)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 960, 544 // Taille de la fenêtre
}

func main() {
	islandConf := island.IslandConfig{Width: 50, Height: 50}
	width, height := 960, 544
	world := island.NewIsland(islandConf)
	player := &player.Player{X: float64(width) / 2, Y: float64(height)/2 - 16, WorldX: int(islandConf.Width)/2 - 1, WorldY: int(islandConf.Height)/2 - 1}

	scene := screen.NewScreen(int(height), int(width), player, world)
	g := NewGame(world, scene, player) // Position initiale du carré

	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("Mon premier jeu Ebitengine")

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}

}
