package engi

import (
	"log"
	"math"
)

type AABB struct {
	Min, Max Point
}

func IsIntersecting(rect1 AABB, rect2 AABB) bool {
	if rect1.Max.X > rect2.Min.X && rect1.Min.X < rect2.Max.X && rect1.Max.Y > rect2.Min.Y && rect1.Min.Y < rect2.Max.Y {
		return true
	}

	return false
}

func MinimumTranslation(rect1 AABB, rect2 AABB) Point {
	mtd := Point{}

	left := float64(rect2.Min.X - rect1.Max.X)
	right := float64(rect2.Max.X - rect1.Min.X)
	top := float64(rect2.Min.Y - rect1.Max.Y)
	bottom := float64(rect2.Max.Y - rect1.Min.Y)

	if left > 0 || right < 0 {
		log.Println("Box aint intercepting")
		return mtd
		//box doesnt intercept
	}

	if top > 0 || bottom < 0 {
		log.Println("Box aint intercepting")
		return mtd
		//box doesnt intercept
	}
	if math.Abs(left) < right {
		mtd.X = float32(left)
	} else {
		mtd.X = float32(right)
	}

	if math.Abs(top) < bottom {
		mtd.Y = float32(top)
	} else {
		mtd.Y = float32(bottom)
	}

	if math.Abs(float64(mtd.X)) < math.Abs(float64(mtd.Y)) {
		mtd.Y = 0
	} else {
		mtd.X = 0
	}

	return mtd
}
