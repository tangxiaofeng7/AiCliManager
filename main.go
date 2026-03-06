package main

import (
	"embed"

	"AiCliManager/internal/app"
	"AiCliManager/internal/db"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	database := db.Init()
	appService := app.New(database)

	err := wails.Run(&options.App{
		Title:     "AiCliManager",
		Width:     1280,
		Height:    800,
		MinWidth:  960,
		MinHeight: 600,
		Frameless: true, // 无边框，使用自定义标题栏
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 18, G: 18, B: 18, A: 1},
		Bind: []interface{}{
			appService,
		},
		OnStartup:  appService.Startup,
		OnShutdown: appService.Shutdown,
	})
	if err != nil {
		panic(err)
	}
}
