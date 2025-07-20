package generator

import (
	"confinedisland/config"
	"confinedisland/sprite"
	"fmt"
	"image"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var UNITE int = 32

var TEMPLATE_ORIENTATION map[string][]string = map[string][]string{
	"N":  []string{"11100000", "01000000", "11000000", "01100000", "10111000", "10100000"},
	"S":  []string{"00000111", "00000010", "00000110", "00000011", "00011101", "00000101"},
	"O":  []string{"10010100", "00010000", "00010100", "10010000", "11000110", "11000110", "10000100"},
	"E":  []string{"00101001", "00001000", "00101000", "00001001", "01100011", "00100001"},
	"NE": []string{"00100000"},
	"NO+": []string{
		"11010000", "11110000", "11010100", "11110100", "01110100", "11000100",
		"11010100", "11010000", "10001001", "10110000", "10110100", "01010100",
		"11000100"},
	"NO": []string{"10000000"},
	"NE+": []string{
		"11100100", "01101000", "11101000", "01101001", "11101001", "10101001",
		"10101000", "00110100", "01100001", "11101000", "11101100", "01101000",
		"01100100", "00101100", "11001001", "11001000", "11100001"},
	"SE": []string{"00000001"},
	"SO+": []string{
		"00010110", "10010110", "00010111", "10010111", "10010101", "10000011",
		"00001101", "00010111", "00010110", "10010001", "10000111", "10010011",
		"00010011", "00010101", "10000110"},
	"SO": []string{"00000100"},
	"SE+": []string{
		"00001011", "00101011", "00001111", "00101111", "00101101", "00101100",
		"00101011", "00001011", "00100110", "00101101", "00101110", "00101010",
		"00100011"},
	"ALONE": []string{
		"10011011", "01111111", "10111111", "11011111", "11101111", "11110111",
		"11111011", "11111101", "11111110", "11111111", "11111000", "00011111",
		"11010110", "01101011", "11111001", "11110110", "10011111", "01101111",
		"00110110", "11010001", "01101100", "10001011", "00110111", "11010101",
		"11101100", "10101011", "11101011", "11111100", "11010111", "00111111",
		"10111101", "11100111", "10111101", "1010111", "11101101", "11110101",
		"10110111", "11000110", "10111000", "01100011", "00011101", "11100111",
		"10111101",
	},
}

type Damage struct {
	Value      int64
	State      string // burn Gel etc...
	Speed      float64
	Animations map[string]sprite.AnimatedSprites
}

type Biome struct {
	Name               string
	Description        string
	Color              color.RGBA
	StartX             int
	Damage             Damage
	Animated           bool
	IterationNumber    int
	AnimationType      string
	AnimationDuration  float64
	SpeedReduce        float64
	State              []string // Attendu : swiming, burn
	Oriented           bool
	OrientationSprites map[string]*sprite.Block
	Voisin             string
	BaseImage          *sprite.Block
}

type Coord struct {
	X int
	Y int
}

type Template_Ground struct {
	Chemin string
	Image  ebiten.Image
	Biomes map[string]Biome
}

var TemplateGroundPosition map[string][]Coord = map[string][]Coord{
	"SE": []Coord{Coord{X: 0, Y: 0}},
	"SO": []Coord{Coord{X: 1, Y: 0}},
	"NE": []Coord{Coord{X: 0, Y: 1}},
	"NO": []Coord{Coord{X: 1, Y: 1}},
	"ALONE": []Coord{
		Coord{X: 2, Y: 0},
		Coord{X: 2, Y: 1},
	},
	"NO+": []Coord{Coord{X: 0, Y: 2}},
	"N":   []Coord{Coord{X: 1, Y: 2}},
	"NE+": []Coord{Coord{X: 2, Y: 2}},
	"O":   []Coord{Coord{X: 0, Y: 3}},
	"E":   []Coord{Coord{X: 2, Y: 3}},
	"SO+": []Coord{Coord{X: 0, Y: 4}},
	"S":   []Coord{Coord{X: 1, Y: 4}},
	"SE+": []Coord{Coord{X: 2, Y: 4}},
	"BUG": []Coord{Coord{X: 0, Y: 5}},
	"RANDOM": []Coord{
		Coord{X: 1, Y: 5},
		Coord{X: 2, Y: 5},
		Coord{X: 1, Y: 6},
		Coord{X: 2, Y: 6},
	},
}
var StaticSpritePool = sync.Pool{
	New: func() interface{} {
		return &sprite.StaticSprite{}
	},
}

var AnimatedSpritePool = sync.Pool{
	New: func() interface{} {
		return &sprite.AnimatedSprites{}
	},
}

var BlockPool = sync.Pool{
	New: func() interface{} {
		return &sprite.Block{}
	},
}

func (templateGround *Template_Ground) GenerateOrientation(templatePosition map[string][]Coord) {
	img, _, err := ebitenutil.NewImageFromFile(templateGround.Chemin)
	if err != nil {
		fmt.Println("Erreur chargement image")
	} else {
		for v, b := range templateGround.Biomes {
			if b.Oriented {
				for k, tp := range templatePosition {
					for _, coord := range tp {
						//block := BlockPool.Get().(*sprite.Block)
						block := &sprite.Block{}
						block.Name = b.Name
						block.Orientation = k
						block.Animated = false
						if b.IterationNumber == 0 {
							x := (b.StartX + coord.X) * UNITE
							y := coord.Y * int(UNITE)
							sprite := StaticSpritePool.Get().(*sprite.StaticSprite)
							sprite.Image = img.SubImage(image.Rect(x, y, x+int(UNITE), y+int(UNITE))).(*ebiten.Image)
							block.Sprite = sprite
						} else {
							var sprites []*sprite.StaticSprite
							for i := 0; i < b.IterationNumber; i++ {
								x := (b.StartX + coord.X) * UNITE
								if i > 0 {
									x += (UNITE * 3) * i
								}
								y := coord.Y * int(UNITE)
								sprite := StaticSpritePool.Get().(*sprite.StaticSprite)
								sprite.Image = img.SubImage(image.Rect(x, y, x+int(UNITE), y+int(UNITE))).(*ebiten.Image)
								sprites = append(sprites, sprite)

							}
							animated := AnimatedSpritePool.Get().(*sprite.AnimatedSprites)
							animated.Sprites = sprites
							animated.Type = b.AnimationType
							animated.Duration = b.AnimationDuration
							animated.Start()
							block.Sprite = animated
							block.Animated = true
						}
						templateGround.Biomes[v].OrientationSprites[k] = block

					}
				}
			}
			biome := templateGround.Biomes[v]
			img := ebiten.NewImage(config.UNITE, config.UNITE)
			img.Fill(b.Color)
			block := BlockPool.Get().(*sprite.Block)
			block.Name = b.Name
			sprite := StaticSpritePool.Get().(*sprite.StaticSprite)
			sprite.Image = img
			block.Sprite = sprite
			biome.BaseImage = block
			templateGround.Biomes[v] = biome
		}

	}
}

var TEMPLATE_GROUND_RESSOURCES Template_Ground = Template_Ground{
	Chemin: "ressources/TileMapV2-Sheet.png",
	Biomes: map[string]Biome{
		"limit": Biome{
			Name:        "limite",
			Description: "Vide galactique",
			Color:       color.RGBA{R: 0, G: 0, B: 0, A: 255},
			Oriented:    false,
		},
		"ocean": Biome{
			Name:        "ocean",
			Description: "Eau tres profonde",
			SpeedReduce: 0.2,
			Animated:    false,
			Color:       color.RGBA{R: 11, G: 94, B: 101, A: 255},
			State:       []string{"Swiming"},
			// Initialisation ici
			Oriented: false,
		},
		"cote": Biome{
			Name:               "cote",
			Description:        "Mer peu profonde",
			Color:              color.RGBA{R: 43, G: 209, B: 194, A: 255},
			Animated:           true,
			AnimationType:      "loop",
			AnimationDuration:  0.5,
			IterationNumber:    5,
			StartX:             0,
			OrientationSprites: make(map[string]*sprite.Block), // Initialisation ici
			Oriented:           true,
			Voisin:             "ocean",
		},
		"plage": Biome{
			Name:               "plage",
			Description:        "Plage",
			Color:              color.RGBA{R: 247, G: 243, B: 183, A: 255},
			Animated:           true,
			AnimationType:      "ping-pong",
			AnimationDuration:  0.5,
			IterationNumber:    5,
			StartX:             15,
			OrientationSprites: make(map[string]*sprite.Block), // Initialisation ici
			Oriented:           true,
			Voisin:             "cote",
		},
		"desert": Biome{
			Name:               "desert",
			Description:        "Territoire aride ! attention au la chaleur et au manque de ressource",
			Color:              color.RGBA{R: 251, G: 185, B: 84, A: 255},
			Animated:           false,
			StartX:             30,
			OrientationSprites: make(map[string]*sprite.Block), // Initialisation ici
			Oriented:           true,
			Voisin:             "plage",
		},
		"plaine": Biome{
			Name:               "plaine",
			Description:        "Plaine Arboré",
			Color:              color.RGBA{R: 130, G: 178, B: 28, A: 255},
			Animated:           false,
			StartX:             33,
			OrientationSprites: make(map[string]*sprite.Block), // Initialisation ici
			Oriented:           true,
			Voisin:             "desert",
		},
		"jungle": Biome{
			Name:               "jungle",
			Description:        "Jungle tropicale",
			Color:              color.RGBA{R: 43, G: 127, B: 7, A: 255},
			Animated:           false,
			StartX:             36,
			OrientationSprites: make(map[string]*sprite.Block), // Initialisation ici
			Oriented:           true,
			Voisin:             "plaine",
		},
		"forest": Biome{
			Name:               "forest",
			Description:        "Foret dense",
			Color:              color.RGBA{R: 20, G: 95, B: 43, A: 255},
			Animated:           false,
			StartX:             39,
			OrientationSprites: make(map[string]*sprite.Block), // Initialisation ici
			Oriented:           true,
			Voisin:             "jungle",
		},
		"montain": Biome{
			Name:               "montain",
			Description:        "Montagne",
			Color:              color.RGBA{R: 70, G: 59, B: 51, A: 255},
			Animated:           false,
			StartX:             42,
			OrientationSprites: make(map[string]*sprite.Block), // Initialisation ici
			Oriented:           true,
			Voisin:             "forest",
		},
		"snow": Biome{
			Name:               "snow",
			Description:        "Snow",
			Color:              color.RGBA{R: 239, G: 255, B: 244, A: 255},
			Animated:           false,
			StartX:             45,
			OrientationSprites: make(map[string]*sprite.Block), // Initialisation ici
			Oriented:           true,
			Voisin:             "montain",
		},
		"obsidian": Biome{
			Name:               "obsidian",
			Description:        "obsidian",
			Color:              color.RGBA{R: 42, G: 35, B: 34, A: 255},
			Animated:           false,
			StartX:             48,
			State:              []string{"Burning"},
			OrientationSprites: make(map[string]*sprite.Block), // Initialisation ici
			Oriented:           true,
			Voisin:             "snow",
		},
		"lave": Biome{
			Name:               "lave",
			Description:        "Lave attention ça brule !",
			Color:              color.RGBA{R: 255, G: 80, B: 32, A: 255},
			Animated:           true,
			AnimationType:      "loop",
			AnimationDuration:  0.5,
			IterationNumber:    5,
			StartX:             51,
			State:              []string{"Burning"},
			OrientationSprites: make(map[string]*sprite.Block), // Initialisation ici
			Oriented:           true,
			Voisin:             "obsidian",
		},
	},
}

func (tg *Template_Ground) Update(fps float64) {
	for v := range tg.Biomes {
		for _, value := range tg.Biomes[v].OrientationSprites {
			if value.Animated {
				value.Sprite.Update(fps)
			}
		}

	}
}
