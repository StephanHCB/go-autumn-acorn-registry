package revcircleint

// circular dependency test with reverse dependency specification

const RevCircleAName = "revcirclea"

type RevCircleA interface {
	IsA() bool
}

const RevCircleBName = "revcircleb"

type RevCircleB interface {
	IsB() bool
}
