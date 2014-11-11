package engi

import (
	"log"
	"testing"
)

func TestAddEntity(t *testing.T) {
	world := World{}
	entityOne := Entity{}
	world.AddEntity(&entityOne)
	entityTwo := Entity{}
	world.AddEntity(&entityTwo)
	if len(world.Entities()) == 0 {
		t.Fail()
	}
}
