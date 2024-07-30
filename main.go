//go:generate ./bin/fyne bundle -o data.go Icon.png

package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

// "fyne.io/fyne/v2/container"

// "fyne.io/fyne/v2/canvas"
// "fyne.io/fyne/v2/container"
// "fyne.io/fyne/v2/theme"
// "fyne.io/fyne/v2/widget"
// "github.com/fyne-io/examples/bugs"
// "github.com/fyne-io/examples/clock"
// "github.com/fyne-io/examples/fractal"
// "github.com/fyne-io/examples/img/icon"
// "github.com/fyne-io/examples/tictactoe"
// "github.com/fyne-io/examples/xkcd"

// type appInfo struct {
//      name string
//      icon fyne.Resource
//      canv bool
//      run  func(fyne.Window) fyne.CanvasObject
// }

//	var apps = []appInfo{
//	     {"Bugs", icon.BugBitmap, false, bugs.Show},
//	     {"XKCD", icon.XKCDBitmap, false, xkcd.Show},
//	     {"Clock", icon.ClockBitmap, true, clock.Show},
//	     {"Fractal", icon.FractalBitmap, true, fractal.Show},
//	     {"Tic Tac Toe", nil, true, tictactoe.Show},
//	}
// func createGrid() *fyne.Container {
// 	grid := container.NewGridWithColums(8)

// 	for
// }

func main() {

	a := app.New()
	a.Settings().SetTheme(newTheme())
	w := a.NewWindow("app")
	w.Resize(fyne.NewSize(1034, 768))

	// split := container.NewHSplit(appList, content)
	// split.Offset = 0.1
	// w.SetContent(split)

	// // Create a multi-line text box
	// multiLineEntry := w.NewMultiLineEntry()
	// multiLineEntry.SetPlaceHolder("Enter multi-line text...")

	w.SetContent(makeGUI(w))
	// w.SetContent(widget.NewLabel("app"))
	w.ShowAndRun()

	// grid := createGrid()
	// w.SetContent((grid))
	// w.ShowAndRun()
	// fmt.Println("hello world")
	// a := app.New()
	// w := a.NewWindow("Hello World")

	// w.SetContent(widget.NewLabel("Hello World!"))
	// w.Resize(fyne.NewSize(200, 200))

	// //app.New()

	// a.NewWindow("Here is the title for my app")
	// b := widget.NewButton("this is a button", func() {
	// 	fmt.Println("test")
	// 	checkbox := widget.NewCheck("title of the check box", func(b bool) {
	// 		fmt.Println(fmt.Sprintf("my check box value %t", b))
	// 	})

	// 	w.SetContent(checkbox)

	// 	w.ShowAndRun()
	// })
	// w.SetContent(b)
	// //w.ShowAndRun()

	// // check box widget

	// w.ShowAndRun()
	// a := app.New()
	//a.SetIcon(resourceIconPng)

	// content := container.NewMax()
	// w := a.NewWindow("Examples")

	// apps[4].icon = theme.RadioButtonIcon() // lazy load Fyne resource to avoid error
	// appList := widget.NewList(
	//      func() int {
	//              return len(apps)
	//      },
	//      func() fyne.CanvasObject {
	//              icon := &canvas.Image{}
	//              label := widget.NewLabel("Text Editor")
	//              labelHeight := label.MinSize().Height
	//              icon.SetMinSize(fyne.NewSize(labelHeight, labelHeight))
	//              return container.NewBorder(nil, nil, icon, nil,
	//                      label)
	//      },
	//      func(id widget.ListItemID, obj fyne.CanvasObject) {
	//              img := obj.(*fyne.Container).Objects[1].(*canvas.Image)
	//              text := obj.(*fyne.Container).Objects[0].(*widget.Label)
	//              img.Resource = apps[id].icon
	//              img.Refresh()
	//              text.SetText(apps[id].name)
	//      })
	// appList.OnSelected = func(id widget.ListItemID) {
	//      content.Objects = []fyne.CanvasObject{apps[id].run(w)}
	// }

	// split := container.NewHSplit(appList, content)
	// split.Offset = 0.1
	// w.SetContent(split)
	// w.Resize(fyne.NewSize(480, 360))
	// w.ShowAndRun()
}
