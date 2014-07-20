// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eng

type Responder interface {
	Load()
	Init()
	Setup()
	Draw()
	Resize(width, height int)
	Update(delta float32)
	Mouse(x, y float32, action Action)
	Scroll(amount float32)
	Key(key Key, modifier Modifier, action Action)
	Type(char rune)
}
