// Copyright 2013 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eng

// An Effect is a type that returns a shader and is able to set that
// shader up on its own.
type Effect interface {
	Setup()
	Shader() *Shader
}

// A Composer composes a series of render passes.
type Composer struct {
}
