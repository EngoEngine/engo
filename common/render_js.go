//+build js

package common

import "fmt"

func getShadersPtr(a, b Shader) (string, string) {
	return fmt.Sprintf("%p", a), fmt.Sprintf("%p", b)
}

func compareShaders(a, b Shader) bool {
	return a == b
}
