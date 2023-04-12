package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type FullLayout struct {

}


func (d *FullLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		minSize = minSize.Max(child.MinSize())
	}

	return minSize
}

func (d *FullLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	padding := theme.DefaultTheme().Size(theme.SizeNamePadding)
	topLeft := fyne.NewPos(-padding, -padding)
	for _, child := range objects {
		child.Resize(fyne.NewSize(containerSize.Width + 2 * padding, containerSize.Height + 2 * padding))
		child.Move(topLeft)
	}
}