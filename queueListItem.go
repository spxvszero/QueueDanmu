package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"image/color"
)

var QueueListItemColor = &struct {
	PrefixColor		color.Color
	ContentColor	color.Color
}{
	PrefixColor: color.NRGBA{
		R: 255,
		G: 102,
		B: 153,
		A: 255,
	},
	ContentColor: color.NRGBA{
		R: 255,
		G: 102,
		B: 153,
		A: 255,
	},
}

type QueueListItem struct {
	fyne.Container
	//Prefix 		canvas.Text
	//Content		canvas.Text
}

func NewQueueListItem(prefix string, content string) fyne.CanvasObject {
	return container.NewHBox(
			canvas.NewText(prefix, QueueListItemColor.PrefixColor),
			canvas.NewText(content, QueueListItemColor.ContentColor))
}