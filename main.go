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
	"confinedisland/config"
	"confinedisland/generator"
	"confinedisland/generator/island"
	"confinedisland/player"
	"confinedisland/screen"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

type Game struct {
	world                        *island.Island
	scene                        *screen.Screen
	player                       *player.Player
	redSquareImage               *ebiten.Image
	backgroundImages             [][]*ebiten.Image
	cibledX, cibledY             int
	mouseCoordX, mouseCoordY     int
	playerCenterX, playerCenterY int

	tps float64
	fps float64

	MenuState  bool
	GameState  bool
	PauseState bool
}

func NewGame(world *island.Island, scene *screen.Screen, player *player.Player) *Game {
	g := &Game{world: world, scene: scene, player: player}
	// Dans la fonction NewGame ou lors de l'initialisation
	g.redSquareImage = ebiten.NewImage(32, 32)
	g.redSquareImage.Fill(color.RGBA{R: 250, G: 150, B: 100, A: 250})
	return g
}

func (g *Game) Draw_Cursor(screen *ebiten.Image) {
	vector.StrokeLine(screen, float32(g.playerCenterX), float32(g.playerCenterY), float32(g.mouseCoordX), float32(g.mouseCoordY), 2, color.White, true)
}

func (g *Game) Update() error {
	g.fps = ebiten.ActualFPS()
	g.tps = ebiten.ActualTPS()
	generator.TEMPLATE_GROUND_RESSOURCES.Update(g.fps)

	g.mouseCoordX, g.mouseCoordY = ebiten.CursorPosition()
	g.playerCenterX, g.playerCenterY = g.player.GetCenterCoord()

	if g.player.Moving {
		g.player.Update(g.fps)
		return nil
	}

	if g.handleInput() && !g.player.Moving {
		g.player.Moving = true
		g.player.ProcessMoveSpeed = 0
	}
	g.scene.Update(g.player.MoveTargetX, g.player.MoveTargetY, g.player.MoveSpeed, g.fps)
	return nil
}

func (g *Game) handleInput() bool {
	moveX, moveY := 0, 0

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		moveY = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		moveY = 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		moveX = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		moveX = 1
	}

	if moveX != 0 || moveY != 0 {
		g.player.MoveTargetX = g.player.WorldX + int(moveX)
		g.player.MoveTargetY = g.player.WorldY + int(moveY)
		g.clampPlayerPosition()
		return true
	}
	return false
}

func (g *Game) clampPlayerPosition() {
	if g.player.MoveTargetY < 0 {
		g.player.MoveTargetY = g.player.WorldY
	}
	if g.player.MoveTargetY >= g.world.Height {
		g.player.MoveTargetY = g.player.WorldY
	}
	if g.player.MoveTargetX < 0 {
		g.player.MoveTargetX = g.player.WorldX
	}
	if g.player.MoveTargetX >= g.world.Width {
		g.player.MoveTargetX = g.player.WorldX
	}
}

func (g *Game) drawVisibleBlocks(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	maxY, maxX := g.scene.Size_y+4, g.scene.Size_x+4
	for j := 0; j < maxY; j++ {
		for i := 0; i < maxX; i++ {
			if g.scene.Background[j][i] != nil {
				op.GeoM.Reset()
				op.GeoM.Translate(float64((i-2)*config.UNITE), float64((j-2)*config.UNITE))
				screen.DrawImage(g.scene.Background[j][i].Image(), op)
			}

		}

	}

}
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})
	// Créer une image de 32x32 pour chaque sprite
	op := &ebiten.DrawImageOptions{}
	g.drawVisibleBlocks(screen, op)
	op.GeoM.Reset()
	op.GeoM.Translate(g.player.X, g.player.Y)
	//screen.DrawImage(screen, nil)
	g.Draw_Cursor(screen)
	screen.DrawImage(g.redSquareImage, op)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS : %v \nTPS : %v \nBIOME : %v", g.fps, g.tps, g.scene.GetGroundUnderPlayer().Name))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 1920, 1080 //960, 544 // Taille de la fenêtre
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	islandConf := island.IslandConfig{Width: 1500, Height: 1500}
	width, height := 1920, 1080 //960, 544

	generator.TEMPLATE_GROUND_RESSOURCES.GenerateOrientation(generator.TemplateGroundPosition)

	world := island.NewIsland(islandConf)
	player := &player.Player{
		Name:        "Bigdimuss",
		X:           float64(width)/2 - float64(config.UNITE),
		Y:           float64(height)/2 - float64(config.UNITE/2),
		WorldX:      int(islandConf.Width) / 2,
		WorldY:      int(islandConf.Height) / 2,
		MoveSpeed:   0.25,
		MoveTargetX: int(islandConf.Width) / 2,
		MoveTargetY: int(islandConf.Height) / 2,
		Moving:      true,
	}

	scene := screen.NewScreen(int(height), int(width), player, world)
	g := NewGame(world, scene, player) // Position initiale du carré

	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("ConfinedIsland")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	//ebiten.SetScreenClearedEveryFrame(true)
	//ebiten.SetVsyncEnabled(true)
	ebiten.SetTPS(config.TICKS)

	if err := ebiten.RunGameWithOptions(g, &ebiten.RunGameOptions{GraphicsLibrary: ebiten.GraphicsLibraryAuto}); err != nil {
		log.Fatal(err)

	}

}
