# Shapes Demo

## What does it do?
It demonstrates how one can efficiently draw various shapes in solid colors (with/without border).

## What are important aspects of the code?

Each shape has to have both a `SpaceComponent` and a `RenderComponent`:

```go
triangle1 := MyShape{BasicEntity: ecs.NewBasic()}
triangle1.SpaceComponent = common.SpaceComponent{Width: 100, Height: 100}
triangle1.RenderComponent = common.RenderComponent{Drawable: common.Triangle{}, Color: color.RGBA{255, 0, 0, 255}}
```

In this example, the `Width` and `Height` are being used to compute the size of the shape, we set the `Drawable`
to an instance of a `common.Triangle`, and set the `Color` to the fill color we want to use.

Finally, we add it to the `RenderSystem`:

```go
for _, system := range w.Systems() {
    switch sys := system.(type) {
    case *common.RenderSystem:
        sys.Add(&triangle1.BasicEntity, &triangle1.RenderComponent, &triangle1.SpaceComponent)
    }
}
```
