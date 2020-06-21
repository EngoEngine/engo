//+build !js

package common

func (s *basicShader) setupBatchSizeDefaults() {
	if s.BatchSize <= 0 {
		s.BatchSize = MaxSprites
	}
}
