package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"image/color"
)

type OpaqueTheme struct {

}
var _ fyne.Theme = (*OpaqueTheme)(nil)
var ifOpenOpaque = false

var themeBGColor color.Color

func init()  {
	themeBGColor = color.NRGBA{
		R: 255,
		G: 255,
		B: 255,
		A: 255,
	}
}

func (te OpaqueTheme)Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if ifOpenOpaque == true {
		if name == theme.ColorNameBackground  {
			return color.NRGBA{
				R: 0,
				G: 0,
				B: 0,
				A: 0,
			}
		}
	}else {
		if name == theme.ColorNameBackground  {
			return themeBGColor
		}
	}

	if name == theme.ColorNameInputBorder || name == theme.ColorNameSeparator {
		if variant == theme.VariantLight{
			return color.NRGBA{
				R: 0,
				G: 0,
				B: 0,
				A: 0,
			}
		}
	}


	return theme.DefaultTheme().Color(name, variant)
}

func (te OpaqueTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	//if name == theme.IconNameHome {
	//	fyne.NewStaticResource("myHome", homeBytes)
	//}

	return theme.DefaultTheme().Icon(name)
}

func (m OpaqueTheme) Font(style fyne.TextStyle) fyne.Resource {
	return resourceHeiTtf
	//return theme.DefaultTheme().Font(style)
}

func (m OpaqueTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameText {
		return 14
	}
	//if name == theme.SizeNamePadding {
	//	return 0
	//}
	return theme.DefaultTheme().Size(name)
}
