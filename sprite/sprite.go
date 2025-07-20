package sprite

import (
	"confinedisland/config"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite interface {
	Draw(*ebiten.Image, *ebiten.DrawImageOptions)
	GetImage() *ebiten.Image
	Update(fps float64)
}

type Block struct {
	Name         string
	Orientation  string
	Sprite       Sprite
	BaseColor    color.RGBA
	X            int
	Y            int
	Animated     bool
	doDamage     bool
	damageValue  int64
	damageSpeed  float64
	elementState []string // Attendu : swiming, burn

}

func UpdateCursorAnimation(fps float64, animation *AnimatedSprites) {
	animation.Update(fps)
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

func (b *Block) Image() *ebiten.Image {
	if b.Sprite != nil {
		return b.Sprite.GetImage()
	} else {
		i := ebiten.NewImage(config.UNITE, config.UNITE)
		i.Fill(b.BaseColor)
		return i
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

func (sp *StaticSprite) Update(fps float64) {
}

func (sp *StaticSprite) GetImage() *ebiten.Image {
	return sp.Image
}

type AnimatedSprites struct {
	Sprites     []*StaticSprite
	Duration    float64
	Type        string
	isRun       bool
	cursor      int
	id          int
	timeElapsed float64
}

func (a *AnimatedSprites) SetCursor(cursor int) {
	a.cursor = cursor
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

func (a *AnimatedSprites) Update(tps float64) {
	if !a.isRun || len(a.Sprites) == 0 {
		return
	}

	// Calculer l'intervalle pour changer d'image basé sur la durée totale
	interval := a.Duration / float64(len(a.Sprites)) // Durée par image

	// Incrémenter le temps écoulé
	a.timeElapsed += 1.0 / tps // Simule le temps écoulé

	// Vérifier si le temps écoulé dépasse l'intervalle
	if a.timeElapsed >= interval {
		a.timeElapsed = 0 // Réinitialiser le temps écoulé

		// Logique pour changer d'image selon le type d'animation
		if a.Type == "loop" {
			a.cursor++
			if a.cursor >= len(a.Sprites) {
				a.cursor = 0 // Recommencer l'animation
			}
		} else if a.Type == "ping-pong" {
			if a.cursor == len(a.Sprites)-1 {
				a.id = -1 // Inverser la direction
			} else if a.cursor == 0 {
				a.id = 1 // Inverser la direction
			}

			// Mise à jour du curseur
			a.cursor += a.id
		}
	}
}

func (a *AnimatedSprites) GetImage() *ebiten.Image {
	return a.Sprites[a.cursor].Image
}
func NewAnimation(sprites []*StaticSprite, duration float64, isrun bool) *AnimatedSprites {
	a := AnimatedSprites{Sprites: sprites, Duration: duration, isRun: isrun}
	return &a
}
