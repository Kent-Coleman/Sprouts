package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 1000
	screenHeight = 480
)

var (
	running  = true
	bkgColor = rl.NewColor(147, 211, 196, 255)

	grassSprite    rl.Texture2D
	mountainSprite rl.Texture2D
	fenceSprite    rl.Texture2D
	houseSprite    rl.Texture2D
	waterSprite    rl.Texture2D
	tilledSprite   rl.Texture2D

	tex rl.Texture2D

	playerSprite rl.Texture2D

	//Player Movement Variables
	playerSrc                                     rl.Rectangle
	playerDest                                    rl.Rectangle
	playerMoving                                  bool
	playerDir                                     int
	playerUp, playerDown, playerRight, playerLeft bool
	// playerStandingUp, playerStandingDown, playerStandingRight, playerStandingLeft bool
	playerFrame int

	frameCount int

	//Map Variables
	tileDest   rl.Rectangle
	tileSrc    rl.Rectangle
	tileMap    []int
	srcMap     []string
	mapW, mapH int

	playerSpeed float32 = 1 //Players moving speed

	musicPause bool
	music      rl.Music

	cam rl.Camera2D //Camera Angle Variable
)

func drawScene() {
	//Set Up Map
	for i := 0; i < len(tileMap); i++ {
		if tileMap[i] != 0 {
			tileDest.X = tileDest.Width * float32(i%mapW)  //To Screen
			tileDest.Y = tileDest.Height * float32(i/mapW) //To Screen

			//Function to decide which tile map source file to use
			//Use MapEditor file to configure map
			if srcMap[i] == "g" {
				tex = grassSprite
			}
			if srcMap[i] == "m" {
				tex = mountainSprite
			}
			if srcMap[i] == "f" {
				tex = fenceSprite
			}
			if srcMap[i] == "h" {
				tex = houseSprite
			}
			if srcMap[i] == "w" {
				tex = waterSprite
			}
			if srcMap[i] == "t" {
				tex = tilledSprite
			}

			//Draws grass below fence and house tiles
			if srcMap[i] == "h" || srcMap[i] == "f" {
				tileSrc.X = 16
				tileSrc.Y = 16
				rl.DrawTexturePro(grassSprite, tileSrc, tileDest, rl.NewVector2(tileDest.Width, tileDest.Height), 0, rl.White)
			}

			tileSrc.X = tileSrc.Width * float32((tileMap[i]-1)%int(tex.Width/int32(tileSrc.Width)))                //From Image
			tileSrc.Y = tileSrc.Height * float32((tileMap[i]-1)/int(tex.Width/int32(tileSrc.Width)))               //From Image
			rl.DrawTexturePro(tex, tileSrc, tileDest, rl.NewVector2(tileDest.Width, tileDest.Height), 0, rl.White) //Draw Map
		}
	}

	rl.DrawTexturePro(playerSprite, playerSrc, playerDest, rl.NewVector2(playerDest.Width, playerDest.Height), 0, rl.White) //Draw Player
}

func input() {
	if rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp) {
		playerMoving = true
		playerDir = 5
		playerUp = true
	} else if rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown) {
		playerMoving = true
		playerDir = 4
		playerDown = true
	} else if rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft) {
		playerMoving = true
		playerDir = 7
		playerLeft = true
	} else if rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight) {
		playerMoving = true
		playerDir = 6
		playerRight = true
	} else {
		playerMoving = false
		playerDir = 0
	}

	if rl.IsKeyPressed(rl.KeySpace) {
		musicPause = !musicPause
	}
}
func update() {
	running = !rl.WindowShouldClose()

	playerSrc.X = playerSrc.Width * float32(playerFrame)
	if playerMoving {
		if playerUp {
			playerDest.Y -= playerSpeed
		}
		if playerDown {
			playerDest.Y += playerSpeed
		}
		if playerRight {
			playerDest.X += playerSpeed
		}
		if playerLeft {
			playerDest.X -= playerSpeed
		}
		if frameCount%8 == 7 { //Frame count when player is walking
			playerFrame++
		}
	} else if frameCount%18 == 7 { //Frame count when player is idle
		playerFrame++
	}

	frameCount++
	if playerFrame > 7 {
		playerFrame = 0
	}
	if !playerMoving && playerFrame > 7 {
		playerFrame = 0
	}

	playerSrc.X = playerSrc.Width * float32(playerFrame)
	playerSrc.Y = playerSrc.Height * float32(playerDir)

	rl.UpdateMusicStream(music)
	if musicPause {
		rl.PauseMusicStream(music)
	} else {
		rl.ResumeMusicStream(music)
	}

	cam.Target = rl.NewVector2(float32(playerDest.X-(playerDest.Width/2)),
		float32(playerDest.Y-(playerDest.Height/2)))

	playerMoving = false
	playerUp, playerDown, playerRight, playerLeft = false, false, false, false
}
func render() {
	rl.BeginDrawing()
	rl.ClearBackground(bkgColor)
	rl.BeginMode2D(cam)

	drawScene()

	rl.EndMode2D()
	rl.EndDrawing()
}

func loadMap(mapFile string) {
	file, err := os.ReadFile(mapFile) //Reads text file for map
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	remNewLines := strings.Replace(string(file), "\n", " ", -1)
	sliced := strings.Split(remNewLines, " ")
	mapW = -1
	mapH = -1
	for i := 0; i < len(sliced); i++ {
		s, _ := strconv.ParseInt(sliced[i], 10, 64)
		m := int(s)
		if mapW == -1 {
			mapW = m
		} else if mapH == -1 {
			mapH = m
		} else if i < mapW*mapH+2 {
			tileMap = append(tileMap, m)
		} else {
			srcMap = append(srcMap, sliced[i])
		}
	}
	if len(tileMap) > mapW*mapH {
		tileMap = tileMap[:len(tileMap)-1]
	}
}

// Method to load images into GUI
func init() {
	rl.InitWindow(screenWidth, screenHeight, "Sprouts")
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	//Map source files
	grassSprite = rl.LoadTexture("Resources/TileSets/GroundTiles/NewTiles/GrassTiles.png")
	mountainSprite = rl.LoadTexture("Resources/TileSets/GroundTiles/NewTiles/MountainTiles.png")
	fenceSprite = rl.LoadTexture("Resources/TileSets/BuildingParts/FenceTiles.png")
	houseSprite = rl.LoadTexture("Resources/TileSets/BuildingParts/HouseTiles.png")
	waterSprite = rl.LoadTexture("Resources/TileSets/WaterTiles.png")
	tilledSprite = rl.LoadTexture("Resources/TileSets/GroundTiles/NewTiles/TilledTiles.png")

	tileDest = rl.NewRectangle(0, 0, 10, 10) //Size of displayed image
	tileSrc = rl.NewRectangle(0, 0, 8, 8)    //Based on pixels of source image

	playerSprite = rl.LoadTexture("Resources/Characters/PremiumCharacterSpriteSheet.png")

	playerSrc = rl.NewRectangle(0, 0, 48, 48)
	playerDest = rl.NewRectangle(200, 200, 31, 31) //Player size

	rl.InitAudioDevice()
	music = rl.LoadMusicStream("Resources/AverysFarm.mp3")
	musicPause = false
	rl.PlayMusicStream(music)

	cam = rl.NewCamera2D(rl.NewVector2(float32(screenWidth/2), float32(screenHeight/2)), rl.NewVector2(float32(playerDest.X-(playerDest.Width/2)),
		float32(playerDest.Y-(playerDest.Height/2))), 0.0, 1.0)

	cam.Zoom = 5 //Camera zoom

	loadMap("MapEditor.map") //Source of array that generates map
}

func quit() {
	rl.UnloadTexture(grassSprite)
	rl.UnloadTexture(playerSprite)
	rl.UnloadMusicStream(music)
	rl.CloseAudioDevice()
	rl.CloseWindow()
}

func main() {

	for running {
		input()
		update()
		render()
	}

	quit()
}
