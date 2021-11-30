# Audio Demo

## What does it do?
It demonstrates how one can play audio.

For doing so, it created bird noises as a background sound, and some weird man's voice as a sound effect, when pressing SPACE.
You can also pause/unpause the sounds at any point by pressing P.   

One can zoom in/out, and hear the sound effect increase/decrease in volume. This is because every non-background sound, is rendered in 3d.

## What are important aspects of the code?
These things are key in this demo:

* `w.AddSystem(&common.AudioSystem{})`, to add/enable the `AudioSystem`;
* The `moveSystem`, to enable the sound effects / graphical effects;
* `whoop.AudioComponent = common.AudioComponent{File: "326064.wav", Repeat: false})`, to add the audio component as a sound effect;
* `whoop.AudioComponent = common.AudioComponent{File: "326488.wav", Repeat: true, Background: true})`, to dd the audio component as background music.
