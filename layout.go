package main

import (
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
)

const sideWidth = 228

type mainLayout struct {
	top      fyne.CanvasObject
	bottom   fyne.CanvasObject
	left     fyne.CanvasObject
	right    fyne.CanvasObject
	textbox  fyne.CanvasObject
	content  fyne.CanvasObject
	dividers [3]fyne.CanvasObject
}

func newLayout(top, bottom, left, right, textbox, content fyne.CanvasObject, dividers [3]fyne.CanvasObject) fyne.Layout {
	return &mainLayout{top: top, bottom: bottom, left: left, right: layout.NewSpacer(), textbox: textbox, content: content, dividers: dividers}
}

func maxInt(arr []int) int {
	if len(arr) == 0 {
		// Handle empty slice case (you might want to return an error or a sentinel value)
		return math.MinInt
	}

	max := arr[0]
	for _, value := range arr[1:] {
		if value > max {
			max = value
		}
	}
	return max
}

// Function to remove an element from a slice by value
func removeElement(slice []int, value int) []int {
	index := -1
	// Find the index of the element to remove
	for i, v := range slice {
		if v == value {
			index = i
			break
		}
	}

	// If the value is not found, return the original slice
	if index == -1 {
		return slice
	}

	// Remove the element at the found index
	return append(slice[:index], slice[index+1:]...)
}

// type RelativePosition int

// const (
// 	Above RelativePosition = iota
// 	Below
// 	Left
// 	Right
// 	None
// )

// Function to check the relative position of widgetB with respect to widgetA
// func relativePosition(closeWidget, input fyne.CanvasObject) int {
// 	posA := closeWidget.Position()
// 	sizeA := closeWidget.Size()
// 	posB := input.Position()
// 	sizeB := input.Size()

// 	// Check if widgetB is directly above widgetA
// 	if posB.X == posA.X && posB.Y+sizeB.Height == posA.Y {
// 		return 1
// 	}

// 	// Check if widgetB is directly below widgetA
// 	if posB.X == posA.X && posB.Y == posA.Y+sizeA.Height {
// 		return 2
// 	}

// 	// Check if widgetB is directly to the left of widgetA
// 	if posB.Y == posA.Y && posB.X+sizeB.Width == posA.X {
// 		return 3
// 	}

// 	// Check if widgetB is directly to the right of widgetA
// 	if posB.Y == posA.Y && posB.X == posA.X+sizeA.Width {
// 		return 4
// 	}

//		return -1
//	}
func minInt(arr []int) int {
	min := arr[0]
	for _, val := range arr {
		if val < min {
			min = val
		}
	}
	return min
}

func adjustLayout(container *fyne.Container, input fyne.CanvasObject) {
	var distances []int
	var closestDistances []int

	inputPos := input.Position()
	inputSize := input.Size()

	// Calculate distances to input widget
	for _, obj := range container.Objects {
		if obj == input {
			continue
		}
		objPos := obj.Position()
		distance := int(math.Sqrt(math.Pow(float64(objPos.X)-float64(inputPos.X), 2) + math.Pow(float64(objPos.Y)-float64(inputPos.Y), 2)))
		distances = append(distances, distance)
	}

	// Find the four closest widgets
	for i := 0; i < 4; i++ {
		if len(distances) == 0 {
			break
		}
		widgetDist := minInt(distances)
		closestDistances = append(closestDistances, widgetDist)
		distances = removeElement(distances, widgetDist)
	}

	// Adjust the layout of the closest widgets
	for _, obj := range container.Objects {
		if obj == input {
			continue
		}
		objPos := obj.Position()
		objSize := obj.Size()

		distance := int(math.Sqrt(math.Pow(float64(objPos.X)-float64(inputPos.X), 2) + math.Pow(float64(objPos.Y)-float64(inputPos.Y), 2)))
		for _, closestDistance := range closestDistances {
			if distance == closestDistance {
				// Check if obj is horizontally aligned with input
				if objPos.Y < inputPos.Y+inputSize.Height && objPos.Y+objSize.Height > inputPos.Y {
					if objPos.X < inputPos.X {
						// Expand to the right to cover input space without shrinking original space
						newWidth := objSize.Width + (inputPos.X - (objPos.X + objSize.Width))
						obj.Resize(fyne.NewSize(newWidth, objSize.Height))
					} else {
						// Move and resize to cover input space without shrinking original space
						newWidth := inputPos.X + inputSize.Width - objPos.X
						obj.Move(fyne.NewPos(inputPos.X, objPos.Y))
						obj.Resize(fyne.NewSize(newWidth, objSize.Height))
					}
				}
				// Check if obj is vertically aligned with input
				if objPos.X < inputPos.X+inputSize.Width && objPos.X+objSize.Width > inputPos.X {
					if objPos.Y < inputPos.Y {
						// Expand downwards to cover input space without shrinking original space
						newHeight := objSize.Height + (inputPos.Y - (objPos.Y + objSize.Height))
						obj.Resize(fyne.NewSize(objSize.Width, newHeight))
					} else {
						// Move and resize to cover input space without shrinking original space
						newHeight := inputPos.Y + inputSize.Height - objPos.Y
						obj.Move(fyne.NewPos(objPos.X, inputPos.Y))
						obj.Resize(fyne.NewSize(objSize.Width, newHeight))
					}
				}
			}
		}
	}
}

// closestWidgets = append(closestWidgets, obj)
// if objPos.X < inputPos.X {
// 	// obj is to the left of input
// 	if objPos.Y < inputPos.Y {
// 		// obj is above input
// 		obj.Resize(fyne.NewSize(objSize.Width, objSize.Height+inputSize.Height))
// 	} else {
// 		// obj is below input
// 		obj.Move(fyne.NewPos(objPos.X, objPos.Y-inputSize.Height))
// 		obj.Resize(fyne.NewSize(objSize.Width, objSize.Height+inputSize.Height))
// 	}
// } else {
// 	// obj is to the right of input
// 	if objPos.Y < inputPos.Y {
// 		// obj is above input
// 		obj.Resize(fyne.NewSize(objSize.Width+inputSize.Width, objSize.Height))
// 	} else {
// 		// obj is below input
// 		obj.Move(fyne.NewPos(objPos.X-inputSize.Width, objPos.Y))
// 		obj.Resize(fyne.NewSize(objSize.Width+inputSize.Width, objSize.Height))
// 	}
// }

// pos := relativePosition(obj, input)
// fmt.Println(pos)
// switch pos {
// case 1:
//
//	fmt.Println("Above")
//	newHeight := obj.Size().Height + input.Size().Height
//	obj.Resize(fyne.NewSize(obj.Size().Width, newHeight))
//
// case 2:
//
//	fmt.Println("Below")
//	obj.Move(fyne.NewPos(obj.Position().X, obj.Position().Y-input.Size().Height))
//	newHeight := obj.Size().Height + input.Size().Height
//	obj.Resize(fyne.NewSize(obj.Size().Width, newHeight))
//
// case 3:
//
//	fmt.Println("Left")
//	newWidth := obj.Size().Width + input.Size().Width
//	obj.Resize(fyne.NewSize(newWidth, obj.Size().Height))
//
// case 4:
//
//		fmt.Println("Right")
//		obj.Move(fyne.NewPos(obj.Position().X-input.Size().Width, obj.Position().Y))
//		newWidth := obj.Size().Width + input.Size().Width
//		obj.Resize(fyne.NewSize(newWidth, obj.Size().Height))
//	}
func (l *mainLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	topHeight := l.top.MinSize().Height
	bottomHeight := l.bottom.MinSize().Height

	leftEdgeRight := float32(sideWidth)
	// rightEdgeLeft := size.Width - float32(sideWidth)
	// middleX := (leftEdgeRight + rightEdgeLeft) / 2

	l.top.Resize(fyne.NewSize(size.Width, topHeight))

	l.left.Move(fyne.NewPos(0, topHeight))
	l.left.Resize(fyne.NewSize(sideWidth, size.Height-topHeight))

	// l.right.Move(fyne.NewPos(0, 0))
	// l.right.Resize(fyne.NewSize(0, 0))
	l.right.Move(fyne.NewPos(size.Width-sideWidth, topHeight))
	l.right.Resize(fyne.NewSize(sideWidth, size.Height-topHeight))
	// l.right.Visible()
	l.right.Hide()

	l.content.Move(fyne.NewPos(sideWidth, topHeight))
	l.content.Resize(fyne.NewSize(size.Width-2*sideWidth, size.Height-topHeight))

	dividerThickness := theme.SeparatorThicknessSize()
	l.dividers[0].Move(fyne.NewPos(0, topHeight))
	l.dividers[0].Resize(fyne.NewSize(size.Width, dividerThickness))

	l.dividers[1].Move(fyne.NewPos(sideWidth, topHeight))
	l.dividers[1].Resize(fyne.NewSize(dividerThickness, size.Height-topHeight))

	l.dividers[2].Move(fyne.NewPos(size.Width-sideWidth, topHeight))
	l.dividers[2].Resize(fyne.NewSize(dividerThickness, size.Height-topHeight))

	// Resize and position the bottom component
	l.bottom.Move(fyne.NewPos(0, size.Height-bottomHeight))
	l.bottom.Resize(fyne.NewSize(size.Width, bottomHeight))

	// Resize and position the textbox component within the bottom area
	// textboxHeight := l.textbox.MinSize().Height
	// l.textbox.Move(fyne.NewPos(middleX/3, size.Height-bottomHeight-40))
	// l.textbox.Resize(fyne.NewSize(size.Width, textboxHeight))
	textboxHeight := l.textbox.MinSize().Height
	padding := float32(35) // adjust this value as needed for padding
	l.textbox.Move(fyne.NewPos(leftEdgeRight+padding, size.Height-bottomHeight-40))
	l.textbox.Resize(fyne.NewSize(size.Width, textboxHeight))

	// Create a container with all the widgets
	// container := container.NewWithoutLayout(objects...)

	// Call adjustLayout function
	//adjustLayout(container, l.right)
}

func (l *mainLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	// return fyne.NewSize(10, 10)

	borders := fyne.NewSize(
		sideWidth*2,
		l.top.MinSize().Height,
	)
	return borders.AddWidthHeight(100, 100)
}
