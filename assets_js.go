// +build netgo

package engo

import (
	"github.com/gopherjs/gopherjs/js"
)

type Image interface {
	Data() *js.Object
	Width() int
	Height() int
}
