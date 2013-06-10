// Copyright 2013 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package scene

import (
	"github.com/ajhager/eng"
)

type Vec2 struct {
	X, Y float32
}

type Actor struct {
	Region   *eng.Region
	Position *Vec2
	Origin   *Vec2
	Scale    *Vec2
	Rotation float32
	Color    *eng.Color
}
