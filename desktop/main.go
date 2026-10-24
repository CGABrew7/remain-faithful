package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()
	_ = wails.Run(&options.App{
		Title:            "Remain Faithful",
		Width:            1100,
		Height:           760,
		MinWidth:         960,
		MinHeight:        680,
		BackgroundColour: &options.RGBA{R: 18, G: 33, B: 76, A: 255},
		AssetServer:      &assetserver.Options{Assets: assets},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Bind:             []interface{}{app},
		Mac: &mac.Options{
			TitleBar: mac.TitleBarDefault(),
			About: &mac.AboutInfo{
				Title:   "Remain Faithful",
				Message: "Accountability between trusted friends",
			},
		},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
		},
	})
}
