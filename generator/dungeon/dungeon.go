package dungeon

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
)

type Room struct {
	row, col, height, width int
}

type DungeonGenerator struct {
	width   int
	height  int
	dungeon [][]float64 // Changez pour utiliser float64
	rooms   []Room
}

func NewDungeonGenerator(w, h int) *DungeonGenerator {
	dg := &DungeonGenerator{
		width:   w,
		height:  h,
		dungeon: make([][]float64, h),
		rooms:   []Room{},
	}
	for i := range dg.dungeon {
		dg.dungeon[i] = make([]float64, w)
		for j := range dg.dungeon[i] {
			dg.dungeon[i][j] = 1 // Fond noir
		}
	}
	return dg
}

// Vérifier si la nouvelle salle chevauche une salle existante
func (dg *DungeonGenerator) isRoomOverlap(newRoom Room) bool {
	for _, room := range dg.rooms {
		if newRoom.row < room.row+room.height+1 && newRoom.row+newRoom.height > room.row-1 &&
			newRoom.col < room.col+room.width+1 && newRoom.col+newRoom.width > room.col-1 {
			return true // Superposition détectée
		}
	}
	return false
}

// Créer le donjon circulaire
func (dg *DungeonGenerator) createCircularDungeon(centerX, centerY, radius int) {
	numRooms := rand.Intn(20) + 2 // Ajuster le nombre de pièces
	for i := 0; i < numRooms; i++ {
		angle := rand.Float64() * 2 * math.Pi
		r := rand.Intn(radius - 5)
		roomX := int(float64(centerX) + float64(r)*math.Cos(angle))
		roomY := int(float64(centerY) + float64(r)*math.Sin(angle))

		roomWidth := rand.Intn(18) + 6
		roomHeight := rand.Intn(18) + 6

		newRoom := Room{roomY, roomX, roomHeight, roomWidth}

		if roomX < 1 || roomX+roomWidth >= dg.width-1 || roomY < 1 || roomY+roomHeight >= dg.height-1 || dg.isRoomOverlap(newRoom) {
			continue
		}

		// Carver la salle et ajouter des murs autour
		for r := roomY; r < roomY+roomHeight; r++ {
			for c := roomX; c < roomX+roomWidth; c++ {
				dg.dungeon[r][c] = 0 // Sol (blanc)
			}
		}

		// Ajouter des murs autour de la salle
		for r := roomY - 1; r <= roomY+roomHeight; r++ {
			for c := roomX - 1; c <= roomX+roomWidth; c++ {
				if dg.dungeon[r][c] != 0 {
					dg.dungeon[r][c] = 0.5 // Murs (gris)
				}
			}
		}

		dg.rooms = append(dg.rooms, newRoom)
	}
}

// Carver le couloir entre deux salles
func (dg *DungeonGenerator) carveCorridorBetweenRooms(room1, room2 Room) {
	startRow := room1.row + rand.Intn(room1.height-2)
	startCol := room1.col + rand.Intn(room1.width-2)

	if rand.Intn(2) == 0 {
		startCol = room1.col // Début à gauche
	} else {
		startCol = room1.col + room1.width - 1 // Début à droite
	}

	endRow := room2.row + rand.Intn(room2.height)
	endCol := room2.col + rand.Intn(room2.width)

	if rand.Intn(2) == 0 {
		endCol = room2.col // Fin à gauche
	} else {
		endCol = room2.col + room2.width - 1 // Fin à droite
	}
	if rand.Float32() < 0.5 {
		if startCol < endCol {
			for c := startCol; c <= endCol; c++ {
				dg.dungeon[startRow][c] = 0 // Sol (blanc)
			}
		} else {
			for c := endCol; c <= startCol; c++ {
				dg.dungeon[startRow][c] = 0 // Sol (blanc)
			}
		}
		if startRow < endRow {
			for r := startRow; r <= endRow; r++ {
				dg.dungeon[r][endCol] = 0 // Sol (blanc)
			}
		} else {
			for r := endRow; r <= startRow; r++ {
				dg.dungeon[r][endCol] = 0 // Sol (blanc)
			}
		}
	} else {
		if startRow < endRow {
			for r := startRow; r <= endRow; r++ {
				dg.dungeon[r][startCol] = 0 // Sol (blanc)
			}
		} else {
			for r := endRow; r <= startRow; r++ {
				dg.dungeon[r][startCol] = 0 // Sol (blanc)
			}
		}
		if startCol < endCol {
			for c := startCol; c <= endCol; c++ {
				dg.dungeon[endRow][c] = 0 // Sol (blanc)
			}
		} else {
			for c := endCol; c <= startCol; c++ {
				dg.dungeon[endRow][c] = 0 // Sol (blanc)
			}
		}
	}
}

// Connecter les salles
func (dg *DungeonGenerator) connectRooms() {
	for i := 0; i < len(dg.rooms)-1; i++ {
		dg.carveCorridorBetweenRooms(dg.rooms[i], dg.rooms[i+1])
	}
}

// Générer la carte du donjon
func (dg *DungeonGenerator) generateMap() {
	centerX := dg.width / 2
	centerY := dg.height / 2
	radius := int(math.Min(float64(dg.width), float64(dg.height)) / 2)
	dg.createCircularDungeon(centerX, centerY, radius)
	dg.connectRooms()
}

// Dessiner l'image du donjon
func (dg *DungeonGenerator) drawToImage() {
	img := image.NewRGBA(image.Rect(0, 0, dg.width, dg.height))
	for y := 0; y < dg.height; y++ {
		for x := 0; x < dg.width; x++ {
			var col color.Color
			switch dg.dungeon[y][x] {
			case 0:
				col = color.White // Sol
			case 0.5:
				col = color.RGBA{128, 128, 128, 255} // Murs en gris
			case 1:
				col = color.Black // Fond noir
			}
			img.Set(x, y, col)
		}
	}

	file, err := os.Create("dungeon.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		panic(err)
	}
}

/*
func main() {
	rand.Seed(time.Now().UnixNano())
	dg := NewDungeonGenerator(150, 150) // Taille de l'image
	dg.generateMap()
	dg.drawToImage()

	fmt.Println("Dungeon map generated as dungeon.png")
}*/
