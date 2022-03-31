package mockb

import (
	auacornapi "github.com/StephanHCB/go-autumn-acorn-registry/api"
	"github.com/StephanHCB/go-autumn-acorn-registry/test/mockc"
	"github.com/StephanHCB/go-autumn-acorn-registry/test/rec"
)

const MockBName = "mockb"

type MockB interface {
	IsB() bool
}

type MockBImpl struct {
	MockC mockc.MockC
}

func New() auacornapi.Acorn {
	rec.Add("b.New")
	return &MockBImpl{}
}

func (m *MockBImpl) AcornName() string {
	return MockBName
}

func (m *MockBImpl) AssembleAcorn(registry auacornapi.AcornRegistry) error {
	rec.Add("b.AssembleAcorn")
	m.MockC = registry.GetAcornByName(mockc.MockCName).(mockc.MockC)

	return nil
}

func (m *MockBImpl) SetupAcorn(registry auacornapi.AcornRegistry) error {
	err := registry.SetupAfter(m.MockC.(auacornapi.Acorn))
	if err != nil {
		return err
	}

	rec.Add("b.SetupAcorn")
	return nil
}

func (m *MockBImpl) TeardownAcorn(registry auacornapi.AcornRegistry) error {
	rec.Add("b.TeardownAcorn")
	return nil
}

// MockB

func (m *MockBImpl) IsB() bool {
	return m.MockC != nil
}
