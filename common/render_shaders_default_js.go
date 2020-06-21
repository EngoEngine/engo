//+build js

package common

func (s *basicShader) setupBatchSizeDefaults() {
	if s.BatchSize <= 0 {
		s.BatchSize = 2048 //js can't seem to handle the whole buffer size
	}
}
