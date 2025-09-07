package world

import (
	"fmt"
	"math"

	"github.com/OpenBiohazard2/OpenBiohazard2/fileio"
	"github.com/go-gl/mathgl/mgl32"
)

func RemoveCollisionEntity(collisionEntities []fileio.CollisionEntity, entityId int) {
	for i, entity := range collisionEntities {
		if entity.ScaIndex == entityId {
			collisionEntities = append(collisionEntities[:i], collisionEntities[i+1:]...)
			fmt.Println("Removing collision entity id ", entityId)
			return
		}
	}
}

func CheckCollision(newPosition mgl32.Vec3, collisionEntities []fileio.CollisionEntity) *fileio.CollisionEntity {
	playerFloorNum := int(math.Round(float64(newPosition.Y()) / fileio.FLOOR_HEIGHT_UNIT))
	for _, entity := range collisionEntities {
		// The boundary is on a different floor than the player
		if !entity.FloorCheck[playerFloorNum] {
			continue
		}

		switch entity.Shape {
		case 0:
			// Rectangle
			corner1 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z)}
			corner2 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z) + float32(entity.Density)}
			corner3 := mgl32.Vec3{float32(entity.X) + float32(entity.Width), 0, float32(entity.Z) + float32(entity.Density)}
			corner4 := mgl32.Vec3{float32(entity.X) + float32(entity.Width), 0, float32(entity.Z)}
			if isPointInRectangle(newPosition, corner1, corner2, corner3, corner4) {
				return &entity
			}
		case 1:
			// Triangle \\|
			vertex1 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z + entity.Density)}
			vertex2 := mgl32.Vec3{float32(entity.X + entity.Width), 0, float32(entity.Z + entity.Density)}
			vertex3 := mgl32.Vec3{float32(entity.X + entity.Width), 0, float32(entity.Z)}
			if isPointInTriangle(newPosition, vertex1, vertex2, vertex3) {
				return &entity
			}
		case 2:
			// Triangle |/
			vertex1 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z)}
			vertex2 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z + entity.Density)}
			vertex3 := mgl32.Vec3{float32(entity.X + entity.Width), 0, float32(entity.Z + entity.Density)}
			if isPointInTriangle(newPosition, vertex1, vertex2, vertex3) {
				return &entity
			}
		case 3:
			// Triangle /|
			vertex1 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z)}
			vertex2 := mgl32.Vec3{float32(entity.X + entity.Width), 0, float32(entity.Z + entity.Density)}
			vertex3 := mgl32.Vec3{float32(entity.X + entity.Width), 0, float32(entity.Z)}
			if isPointInTriangle(newPosition, vertex1, vertex2, vertex3) {
				return &entity
			}
		case 6:
			// Circle
			radius := float32(entity.Width) / 2.0
			center := mgl32.Vec3{float32(entity.X) + radius, 0, float32(entity.Z) + radius}
			if isPointInCircle(newPosition, center, radius) {
				return &entity
			}
		case 7:
			// Ellipse, rectangle with rounded corners on the x-axis
			majorAxis := float32(entity.Width) / 2.0
			minorAxis := float32(entity.Density) / 2.0
			center := mgl32.Vec3{float32(entity.X) + majorAxis, 0, float32(entity.Z) + minorAxis}
			if isPointInEllipseXAxisMajor(newPosition, center, majorAxis, minorAxis) {
				return &entity
			}
		case 8:
			// Ellipse, rectangle with rounded corners on the z-axis
			majorAxis := float32(entity.Density) / 2.0
			minorAxis := float32(entity.Width) / 2.0
			center := mgl32.Vec3{float32(entity.X) + minorAxis, 0, float32(entity.Z) + majorAxis}
			if isPointInEllipseZAxisMajor(newPosition, center, majorAxis, minorAxis) {
				return &entity
			}
		case 9:
			// Rectangle climb up
			corner1 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z)}
			corner2 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z) + float32(entity.Density)}
			corner3 := mgl32.Vec3{float32(entity.X) + float32(entity.Width), 0, float32(entity.Z) + float32(entity.Density)}
			corner4 := mgl32.Vec3{float32(entity.X) + float32(entity.Width), 0, float32(entity.Z)}
			if isPointInRectangle(newPosition, corner1, corner2, corner3, corner4) {
				return &entity
			}
		case 10:
			// Rectangle jump down
			corner1 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z)}
			corner2 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z) + float32(entity.Density)}
			corner3 := mgl32.Vec3{float32(entity.X) + float32(entity.Width), 0, float32(entity.Z) + float32(entity.Density)}
			corner4 := mgl32.Vec3{float32(entity.X) + float32(entity.Width), 0, float32(entity.Z)}
			if isPointInRectangle(newPosition, corner1, corner2, corner3, corner4) {
				return &entity
			}
		case fileio.SCA_TYPE_SLOPE: // 11
			corner1 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z)}
			corner2 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z) + float32(entity.Density)}
			corner3 := mgl32.Vec3{float32(entity.X) + float32(entity.Width), 0, float32(entity.Z) + float32(entity.Density)}
			corner4 := mgl32.Vec3{float32(entity.X) + float32(entity.Width), 0, float32(entity.Z)}
			if isPointInRectangle(newPosition, corner1, corner2, corner3, corner4) {
				return &entity
			}
		case fileio.SCA_TYPE_STAIRS: // 12
			corner1 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z)}
			corner2 := mgl32.Vec3{float32(entity.X), 0, float32(entity.Z) + float32(entity.Density)}
			corner3 := mgl32.Vec3{float32(entity.X) + float32(entity.Width), 0, float32(entity.Z) + float32(entity.Density)}
			corner4 := mgl32.Vec3{float32(entity.X) + float32(entity.Width), 0, float32(entity.Z)}
			if isPointInRectangle(newPosition, corner1, corner2, corner3, corner4) {
				return &entity
			}
		}
	}
	return nil
}

func CheckRamp(entity *fileio.CollisionEntity) bool {
	return entity.Shape == fileio.SCA_TYPE_SLOPE || entity.Shape == fileio.SCA_TYPE_STAIRS
}

func CheckNearbyBoxClimb(playerPosition mgl32.Vec3, collisionEntities []fileio.CollisionEntity) bool {
	for _, entity := range collisionEntities {
		switch entity.Shape {
		case 9:
			// Rectangle climb up
			rectMinX := math.Min(float64(entity.X), float64(entity.X)+float64(entity.Width))
			rectMaxX := math.Max(float64(entity.X), float64(entity.X)+float64(entity.Width))
			rectMinZ := math.Min(float64(entity.Z), float64(entity.Z)+float64(entity.Density))
			rectMaxZ := math.Max(float64(entity.Z), float64(entity.Z)+float64(entity.Density))

			dx := math.Max(rectMinX-float64(playerPosition.X()), float64(playerPosition.X())-rectMaxX)
			dz := math.Max(rectMinZ-float64(playerPosition.Z()), float64(playerPosition.Z())-rectMaxZ)
			dist := math.Sqrt(dx*dx + dz*dz)
			if dist <= 1000 {
				return true
			}
		case 10:
			// Rectangle climb down
			rectMinX := math.Min(float64(entity.X), float64(entity.X)+float64(entity.Width))
			rectMaxX := math.Max(float64(entity.X), float64(entity.X)+float64(entity.Width))
			rectMinZ := math.Min(float64(entity.Z), float64(entity.Z)+float64(entity.Density))
			rectMaxZ := math.Max(float64(entity.Z), float64(entity.Z)+float64(entity.Density))

			dx := math.Max(rectMinX-float64(playerPosition.X()), float64(playerPosition.X())-rectMaxX)
			dz := math.Max(rectMinZ-float64(playerPosition.Z()), float64(playerPosition.Z())-rectMaxZ)
			dist := math.Sqrt(dx*dx + dz*dz)
			if dist <= 1000 {
				return true
			}
		}
	}

	return false
}

func isPointInTriangle(point mgl32.Vec3, corner1 mgl32.Vec3, corner2 mgl32.Vec3, corner3 mgl32.Vec3) bool {
	// area of triangle ABC
	area := triangleArea(corner1, corner2, corner3)
	// area of PBC
	area1 := triangleArea(point, corner2, corner3)
	// area of APC
	area2 := triangleArea(corner1, point, corner3)
	// area of ABP
	area3 := triangleArea(corner1, corner2, point)

	// areas should be equal if point is in triangle
	areaDifference := area - (area1 + area2 + area3)
	return math.Abs(float64(areaDifference)) <= 0.01
}

// Find the area of triangle formed by p1, p2 and p3
func triangleArea(p1 mgl32.Vec3, p2 mgl32.Vec3, p3 mgl32.Vec3) float32 {
	return float32(math.Abs(float64((p1.X()*(p2.Z()-p3.Z()) + p2.X()*(p3.Z()-p1.Z()) + p3.X()*(p1.Z()-p2.Z())) / 2.0)))
}

func isPointInRectangle(point mgl32.Vec3, corner1 mgl32.Vec3, corner2 mgl32.Vec3, corner3 mgl32.Vec3, corner4 mgl32.Vec3) bool {
	x := point.X()
	z := point.Z()
	x1 := corner1.X()
	z1 := corner1.Z()

	x2 := corner2.X()
	z2 := corner2.Z()

	x3 := corner3.X()
	z3 := corner3.Z()

	x4 := corner4.X()
	z4 := corner4.Z()

	a := (x2-x1)*(z-z1) - (z2-z1)*(x-x1)
	b := (x3-x2)*(z-z2) - (z3-z2)*(x-x2)
	c := (x4-x3)*(z-z3) - (z4-z3)*(x-x3)
	d := (x1-x4)*(z-z4) - (z1-z4)*(x-x4)

	if (a > 0 && b > 0 && c > 0 && d > 0) ||
		(a < 0 && b < 0 && c < 0 && d < 0) {
		return true
	}
	return false
}

func isPointInCircle(point mgl32.Vec3, circleCenter mgl32.Vec3, radius float32) bool {
	distance := point.Sub(circleCenter).Len()
	return distance <= radius
}

func isPointInEllipseXAxisMajor(point mgl32.Vec3, ellipseCenter mgl32.Vec3, majorAxis float32, minorAxis float32) bool {
	xDistance := math.Pow(float64(point.X()-ellipseCenter.X()), 2) / float64(majorAxis*majorAxis)
	zDistance := math.Pow(float64(point.Z()-ellipseCenter.Z()), 2) / float64(minorAxis*minorAxis)
	return xDistance+zDistance <= 1.0
}

func isPointInEllipseZAxisMajor(point mgl32.Vec3, ellipseCenter mgl32.Vec3, majorAxis float32, minorAxis float32) bool {
	xDistance := math.Pow(float64(point.X()-ellipseCenter.X()), 2) / float64(minorAxis*minorAxis)
	zDistance := math.Pow(float64(point.Z()-ellipseCenter.Z()), 2) / float64(majorAxis*majorAxis)
	return xDistance+zDistance <= 1.0
}
