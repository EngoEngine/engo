// +build !netgo

package engo

import (
	"engo.io/gl"
)

type Image interface {
	Data() interface{}
	Width() int
	Height() int
}
