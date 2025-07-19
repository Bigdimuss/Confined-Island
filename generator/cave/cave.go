package cave

import (
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
)

type Cave struct {
	grid   [][]int
	width  int
	height int
}

func NewCave(width int, height int) *Cave {
	cave := &Cave{
		width:  width,
		height: height,
	}
	cave.InitCave()

	return cave
}
func (c *Cave) InitCave() {
	c.grid = make([][]int, c.height)
	for i := range c.grid {
		c.grid[i] = make([]int, c.width)
	}
}
func (c *Cave) Generate() {
	// Générer des murs organiques sur les bords
	for y := 0; y < c.height; y++ {
		for x := 0; x < c.width; x++ {
			if x < 5 || x > c.width-6 || y < 5 || y > c.height-6 {
				if rand.Float64() < 0.6 {
					c.grid[y][x] = 1 // Mur
				} else {
					c.grid[y][x] = 0 // Vide
				}
			} else {
				c.grid[y][x] = 0 // Toujours vide au centre
			}
		}
	}

	// Ajouter des murs internes de manière organique
	c.addInternalWalls()

	// Ajouter des taches organiques
	c.addOrganicSpots()

	// Appliquer un bruit pour créer des parois organiques
	c.addOrganicEdges()

	// Application de l'automate cellulaire pour lisser les murs
	for i := 0; i < 7; i++ {
		c.applyCellularAutomata()
	}
}

func (c *Cave) addInternalWalls() {
	for y := 5; y < c.height-5; y++ {
		for x := 5; x < c.width-5; x++ {
			if rand.Float64() < 0.05 { // Ajustez la probabilité pour plus ou moins de murs internes
				c.grid[y][x] = 1 // Créer des murs internes
			}
		}
	}
}

func (c *Cave) addOrganicSpots() {
	for y := 10; y < c.height-10; y++ {
		for x := 10; x < c.width-10; x++ {
			if rand.Float64() < 0.01 { // Réduisez la probabilité pour générer moins de taches
				// Créer une tache de taille aléatoire
				size := rand.Intn(5) + 3 // Taille entre 3 et 7
				for dy := -size; dy <= size; dy++ {
					for dx := -size; dx <= size; dx++ {
						if dx*dx+dy*dy < size*size { // Garder une forme circulaire
							nx, ny := x+dx, y+dy
							if nx >= 0 && nx < c.width && ny >= 0 && ny < c.height {
								c.grid[ny][nx] = 1 // Mur pour la tache
							}
						}
					}
				}
			}
		}
	}
}

func (c *Cave) addOrganicEdges() {
	for y := 0; y < c.height; y++ {
		for x := 0; x < c.width; x++ {
			if (x < 5 || x > c.width-6 || y < 5 || y > c.height-6) && c.grid[y][x] == 1 {
				if rand.Float64() < 0.15 {
					c.grid[y][x] = 0 // Créer des espaces aléatoires
				}
			}
		}
	}
}

func (c *Cave) applyCellularAutomata() {
	newGrid := make([][]int, c.height)
	for i := range newGrid {
		newGrid[i] = make([]int, c.width)
	}

	for y := 0; y < c.height; y++ {
		for x := 0; x < c.width; x++ {
			aliveNeighbors := c.countAliveNeighbors(x, y)

			if c.grid[y][x] == 1 {
				if aliveNeighbors < 4 {
					newGrid[y][x] = 0 // Meurt
				} else {
					newGrid[y][x] = 1 // Reste mur
				}
			} else {
				if aliveNeighbors > 4 {
					newGrid[y][x] = 1 // Devient mur
				} else {
					newGrid[y][x] = 0 // Reste vide
				}
			}
		}
	}

	// Assurer que les bords restent des murs
	for y := 0; y < c.height; y++ {
		newGrid[y][0] = 1
		newGrid[y][c.width-1] = 1
	}
	for x := 0; x < c.width; x++ {
		newGrid[0][x] = 1
		newGrid[c.height-1][x] = 1
	}

	c.grid = newGrid
}

func (c *Cave) countAliveNeighbors(x, y int) int {
	count := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx, ny := x+dx, y+dy
			if nx >= 0 && nx < c.width && ny >= 0 && ny < c.height && c.grid[ny][nx] == 1 {
				count++
			}
		}
	}
	return count
}

func (c *Cave) SaveImage(filename string) error {
	img := image.NewRGBA(image.Rect(0, 0, c.width, c.height))
	black := color.RGBA{0, 0, 0, 255}       // Couleur pour les murs
	white := color.RGBA{255, 255, 255, 255} // Couleur pour les espaces vides

	for y := 0; y < c.height; y++ {
		for x := 0; x < c.width; x++ {
			if c.grid[y][x] == 1 {
				img.Set(x, y, black) // Mur
			} else {
				img.Set(x, y, white) // Vide
			}
		}
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}

/*
func main() {
	rand.Seed(time.Now().UnixNano())
	cave := NewCave(1500, 1500)
	cave.Generate()

	err := cave.SaveImage("cave.png")
	if err != nil {
		panic(err)
	}
}*/
