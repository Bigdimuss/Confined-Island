package island

import (
	"confinedisland/generator"
	"confinedisland/sprite"
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
	"strings"
	"time"

	"github.com/aquilax/go-perlin"
	"golang.org/x/exp/rand"
)

type IslandConfig struct {
	Width  int
	Height int
	Seed   int64
}

type Island struct {
	Width      int
	Height     int
	Seed       int64
	Base       [][]float64
	Background [][]*sprite.Block
}

func NewIsland(config IslandConfig) *Island {
	// Initialiser le générateur de nombres aléatoires...
	background := make([][]*sprite.Block, config.Height)
	if config.Seed == 0 {
		// Initialiser le générateur de nombres aléatoires
		rand.Seed(uint64(time.Now().UnixNano()))
		// Définir la plage
		min := -100000 * 100000 // valeur minimale
		max := 100000 * 100000  // valeur maximale

		// Générer un nombre aléatoire entre min et max
		config.Seed = int64(rand.Intn(max-min+1) + min)
	}
	for i := range background {
		background[i] = make([]*sprite.Block, config.Width)
	}
	i := Island{Width: config.Width, Height: config.Height, Seed: config.Seed, Background: background}
	i.CreateIslandBase() // Crée l'île
	i.CreateIslandBackGround()
	i.GenerateBlock() // Remplit l'arrière-plan avec les couleurs
	return &i
}

func (i *Island) CreateIslandBase() {

	// Créer un générateur de bruit de Perlin
	p := perlin.NewPerlin(0.5, 0.5, 2, i.Seed)
	i.Base = make([][]float64, i.Height) // Initialiser Base ici
	for y := 0; y < i.Height; y++ {
		ligne := make([]float64, i.Width) // Créez la ligne avec la largeur

		for x := 0; x < i.Width; x++ {
			// Calculer la distance du pixel par rapport au centre
			distFromCenter := math.Sqrt(float64((x-i.Width/2)*(x-i.Width/2)+(y-i.Height/2)*(y-i.Height/2))) / (float64(i.Width) / 2)

			// Calculer la valeur de bruit de Perlin pour chaque pixel
			value := p.Noise2D(float64(x)/(float64(i.Width)/10), float64(y)/(float64(i.Height)/10))

			// Ajuster la valeur de bruit pour créer une île centrale avec différents reliefs
			value = (value + 1.0) / 2.0
			value = value * value * (3 - 2*value)
			value = value * (1.0 - distFromCenter)
			ligne[x] = value
		}
		i.Base[y] = ligne
	}
}

/*
def presence_voisin(mapi, cible, coo):
    enum_co = enumerate(coo)
    for k, val in enum_co:
        if mapi[coo[k][0]][coo[k][1]] == cible:
            return True
    return False
*/

func (i *Island) presenceVoisin(cible string, coord [][]int) bool {
	for _, val := range coord {
		if val[0] >= 0 && val[0] < i.Height && val[1] >= 0 && val[1] < i.Width {
			if i.Background[val[0]][val[1]].Name == cible {
				return true
			}
		}

	}
	return false
}
func (i *Island) OrientationBlock(cible string, coord [][]int) string {
	var blockCode []string = []string{"0", "0", "0", "0", "0", "0", "0", "0"}
	for k, val := range coord {
		//if val[0] >= 0 && val[0] < i.Height && val[1] >= 0 && val[1] < i.Width {
		if i.Background[val[0]][val[1]].Name == cible {
			blockCode[k] = "1"
		} else {
			blockCode[k] = "0"
		}
		//}
	}
	return strings.Join(blockCode, "")
}

func (i *Island) GenerateBlock() {

	for y := 0; y < i.Height; y++ {
		for x := 0; x < i.Width; x++ {
			coordonnees := [][]int{
				[]int{y - 1, x - 1},
				[]int{y - 1, x},
				[]int{y - 1, x + 1},
				[]int{y, x - 1},
				[]int{y, x + 1},
				[]int{y + 1, x - 1},
				[]int{y + 1, x},
				[]int{y + 1, x + 1},
			}
			if generator.TEMPLATE_GROUND_RESSOURCES.Biomes[i.Background[y][x].Name].Oriented {
				target := generator.TEMPLATE_GROUND_RESSOURCES.Biomes[i.Background[y][x].Name].Voisin
				voisin := i.presenceVoisin(target, coordonnees)
				if voisin {
					code := i.OrientationBlock(target, coordonnees)
					if code != "11111111" {
						for orientation, val := range generator.TEMPLATE_ORIENTATION {
							for _, v := range val {
								fmt.Printf("Comparaison: %s avec %s\n", code, v)
								if code == v {
									i.Background[y][x] = generator.TEMPLATE_GROUND_RESSOURCES.Biomes[i.Background[y][x].Name].OrientationSprites[orientation]
								}
							}
						}
					} else {
						i.Background[y][x] = generator.TEMPLATE_GROUND_RESSOURCES.Biomes[i.Background[y][x].Name].BaseImage
					}
				} else {
					i.Background[y][x] = generator.TEMPLATE_GROUND_RESSOURCES.Biomes[i.Background[y][x].Name].BaseImage
				}

			} else {
				i.Background[y][x] = generator.TEMPLATE_GROUND_RESSOURCES.Biomes[i.Background[y][x].Name].BaseImage

			}
		}
	}
}
func (i *Island) CreateIslandBackGround() {
	img := image.NewRGBA(image.Rect(0, 0, i.Width, i.Height))
	var name string
	for y := 0; y < i.Height; y++ {
		for x := 0; x < i.Width; x++ {
			value := i.Base[y][x]
			// Définir la couleur du pixel en fonction de la valeur de bruit
			if value < 0.04 {
				name = "ocean" // Bleu foncé pour l'eau profonde
			} else if value < 0.07 {
				name = "cote" // Bleu clair pour l'eau peu profonde
			} else if value < 0.10 {
				name = "plage" // Beige pour la plage
			} else if value < 0.15 {
				name = "desert" // Jaune pour le désert
			} else if value < 0.28 {
				name = "plaine" // Vert clair pour les plaines
			} else if value < 0.40 {
				name = "jungle" // Vert foncé pour la forêt
			} else if value < 0.5 {
				name = "forest"
			} else if value < 0.6 {
				name = "montain" // Gris pour les montagnes
			} else if value < 0.7 {
				name = "snow" // Blanc pour les sommets enneigés
			} else if value < 0.8 {
				name = "obsidian" // Obsidienne
			} else {
				name = "lave" // Lave
			}
			color := generator.TEMPLATE_GROUND_RESSOURCES.Biomes[name].Color
			img.Set(x, y, color)

			i.Background[y][x] = &sprite.Block{
				Name:      name,
				BaseColor: color,
			} // Remplissez Background correctement

		}
	}
	// Enregistrer l'image dans un fichier...
	// Enregistrer l'image dans un fichier
	file, err := os.Create("island.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}
}
