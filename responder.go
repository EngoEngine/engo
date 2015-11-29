// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package engi

import (
	"github.com/paked/engi/ecs"
)

type CustomGame interface {
	Preload()
	Setup(*ecs.World)
}
