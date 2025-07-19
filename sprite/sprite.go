package sprite

import (
	"confinedisland/config"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite interface {
	Draw(*ebiten.Image, *ebiten.DrawImageOptions)
}

type Block struct {
	Name         string
	Orientation  string
	Sprite       Sprite
	BaseColor    color.RGBA
	X            int
	Y            int
	doDamage     bool
	damageValue  int64
	damageSpeed  float64
	elementState []string // Attendu : swiming, burn

}

func (b *Block) Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) {

	if b.Sprite != nil {
		b.Sprite.Draw(screen, op)
	} else {
		i := ebiten.NewImage(config.UNITE, config.UNITE)
		i.Fill(b.BaseColor)
		screen.DrawImage(i, op)
	}

}
func (b *Block) SetDamageValue(damage int64) {
	b.damageValue = damage
}

func (b *Block) SetDamageSpeed(speed float64) {
	b.damageSpeed = speed
}

func (b *Block) SetDoDamage(do bool) {
	b.doDamage = do
}

func (b *Block) AddElementState(state string) {
	b.elementState = append(b.elementState, state)
}

func (b *Block) RemoveElementState(state string) {
	var newlist []string
	for _, st := range b.elementState {
		if st != state {
			newlist = append(newlist, st)
		}
	}
	b.elementState = newlist
}

/*
func newBlocksList(path string, startX int64, nbinteration int) map[string]Block {
	var blocks map[string]Block
	img, _, err := ebitenutil.NewImageFromFile(path)
	images := AnimatedSprites{}
	for i := 0; i < nbinteration; i++ {

	}
	return blocks
}
*/
/*
	    img, err := ebitenutil.NewImageFromFile(path)
	    if err != nil {
	        return nil, err
	    }

	    var sprites []*ebiten.Image
	    for y := 0; y < img.Bounds().Dy(); y += spriteHeight {
	        for x := 0; x < img.Bounds().Dx(); x += spriteWidth {
	            rect := image.Rect(x, y, x+spriteWidth, y+spriteHeight)
	            sprite := img.SubImage(rect).(*ebiten.Image)
	            sprites = append(sprites, sprite)
	        }
	    }
	    return sprites, nil
	}
*/
func Lerp(start, end float64, t float64) float64 {
	return start + (end-start)*t
}

type StaticSprite struct {
	Width, Height float64
	Image         *ebiten.Image
}

func (sp *StaticSprite) Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	if sp.Image != nil {
		screen.DrawImage(sp.Image, op)
	}
}

type AnimatedSprites struct {
	Sprites  []*StaticSprite
	Duration float64
	Type     string
	isRun    bool
	cursor   int
}

func (a *AnimatedSprites) GetRunState() bool {
	return a.isRun
}
func (a *AnimatedSprites) Start() {
	a.isRun = true
}
func (a *AnimatedSprites) Stop() {
	a.isRun = false
	a.cursor = 0
}
func (a *AnimatedSprites) Pause() {
	a.isRun = false
}
func (a *AnimatedSprites) Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	if a.Sprites != nil {
		screen.DrawImage(a.Sprites[0].Image, op)
	}

}

func NewAnimation(sprites []*StaticSprite, duration float64, isrun bool) *AnimatedSprites {
	a := AnimatedSprites{Sprites: sprites, Duration: duration, isRun: isrun}
	return &a
}
