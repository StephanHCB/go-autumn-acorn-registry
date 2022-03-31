package circleint

// circular dependency test

const MockAName = "mocka"

type MockA interface {
	IsA() bool
}

const MockBName = "mockb"

type MockB interface {
	IsB() bool
}

