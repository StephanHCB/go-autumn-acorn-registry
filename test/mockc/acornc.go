package mockc

import (
	auacornapi "github.com/StephanHCB/go-autumn-acorn-registry/api"
	"github.com/StephanHCB/go-autumn-acorn-registry/test/rec"
)

const MockCName = "mockc"

type MockC interface {
	IsC() bool
}

type MockCImpl struct {
}

func New() auacornapi.Acorn {
	rec.Add("c.New")
	return &MockCImpl{}
}

func (m *MockCImpl) AcornName() string {
	return MockCName
}

func (m *MockCImpl) AssembleAcorn(_ auacornapi.AcornRegistry) error {
	rec.Add("c.AssembleAcorn")
	return nil
}

func (m *MockCImpl) SetupAcorn(_ auacornapi.AcornRegistry) error {
	rec.Add("c.SetupAcorn")
	return nil
}

func (m *MockCImpl) TeardownAcorn(_ auacornapi.AcornRegistry) error {
	rec.Add("c.TeardownAcorn")
	return nil
}

// MockC

func (m *MockCImpl) IsC() bool {
	return true
}
