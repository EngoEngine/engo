# Text Demo

## What does it do?
It demonstrates how one can efficiently draw text onto the screen.

## What are important aspects of the code?

First off, we have to load the `ttf` (TrueType Font) file within the `Preload` method:

```go
err := engo.Files.Load("Roboto-Regular.ttf")
```

Then, whenever we want to use it, we can load it like this:

```go
fnt := &common.Font{
    URL:  "Roboto-Regular.ttf", // This is the sme URL as defined with the Preload method
    FG:   color.Black,          // This is the color of the text
    Size: 64,                   // This is the size
}
err := fnt.CreatePreloaded()    // This is required to load the preloaded file into this struct
```

Next, whenever we want to draw something, we set the `Drawable` of the `RenderComponent`:

```go
label1.RenderComponent.Drawable = common.Text{
    Font: fnt,              // This is the `fnt` in the snippet above
    Text: "Hello world !",  // This can be anything, and may include newlines
}
```

And finally, we add the whole thing to the `RenderSystem`:

```go
for _, system := range w.Systems() {
    switch sys := system.(type) {
    case *common.RenderSystem:
        sys.Add(&label1.BasicEntity, &label1.RenderComponent, &label1.SpaceComponent)
    }
}
```
