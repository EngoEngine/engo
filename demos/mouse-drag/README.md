# Mouse Dragging Demo

## What does it do?
It demonstrates how one can use mouse dragging in `engo`

## What are important aspects of the code?
These lines are key in this demo:

* Get the state of your mouseable entity with `e.MouseComponent.Clicked` and `e.MouseComponent.Released`

* Create a clickable shape using lines to make a polygon

```go
bananaOutline := []engo.Line{}
points := []float32{0, 95, 11, 93, 47, 89, 60.5, 82.3, 76.7, 66.2, 81.9, 57.0, 85.3, 31.5,
  78.9, 14.6, 78.6, 6.4, 72.2, 0, 88.1, 0, 93.0, 4.1, 98.8, 20.1, 110.2,
  42.6, 110.2, 60.4, 97.9, 91.7, 84.8, 105.4, 71.4, 113.6, 50.1, 119.8,
  28.7, 119.8, 28.7, 119.8, 28.7, 119.8, 12, 112.5, 0, 103.7, 0, 95.0}
for i := 0; i < len(points)-2; i += 2 {
  line := engo.Line{
    P1: engo.Point{
      X: points[i],
      Y: points[i+1],
    },
    P2: engo.Point{
      X: points[i+2],
      Y: points[i+3],
    },
  }
  bananaOutline = append(bananaOutline, line)
}
banana.SpaceComponent.AddShape(common.Shape{Lines: bananaOutline})
```

* Create a clickable shape using an ellipse

```go
watermelon.SpaceComponent.AddShape(common.Shape{Ellipse: common.Ellipse{Cx: 61, Cy: 50, Rx: 61, Ry: 50}})
```

* Create a clicable shape using multiple shapes

```go
cherry.SpaceComponent.AddShape(common.Shape{Ellipse: common.Ellipse{Cx: 25, Cy: 57.5, Rx: 25, Ry: 25.5}})
cherry.SpaceComponent.AddShape(common.Shape{Ellipse: common.Ellipse{Cx: 59, Cy: 75, Rx: 26, Ry: 25}})
cherryOutline := []engo.Line{}
points = []float32{36.2, 37.5, 46.4, 22.7, 29.9, 18.2, 28.5, 16.4, 44.1, 3.5,
  50, 0, 61.2, 0, 69.7, 5.0, 83.1, 5.1, 90.2, 9.0, 98.8, 22.6, 106.0, 45.9,
  106, 49.4, 100.7, 49.3, 79.8, 39.1, 70, 28.9, 69, 24.7, 65.7, 37.8, 65.7,
  54.2, 59.5, 54.1, 59.2, 35.2, 62.9, 22.7, 54.2, 24.0, 40.5, 42.0, 36.2, 37.5}
for i := 0; i < len(points)-2; i += 2 {
  line := engo.Line{
    P1: engo.Point{
      X: points[i],
      Y: points[i+1],
    },
    P2: engo.Point{
      X: points[i+2],
      Y: points[i+3],
    },
  }
  cherryOutline = append(cherryOutline, line)
}
cherry.SpaceComponent.AddShape(common.Shape{Lines: cherryOutline})
```
