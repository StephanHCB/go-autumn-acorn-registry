package mocka

import (
	auacornapi "github.com/StephanHCB/go-autumn-acorn-registry/api"
	"github.com/StephanHCB/go-autumn-acorn-registry/test/mockb"
	"github.com/StephanHCB/go-autumn-acorn-registry/test/mockc"
	"github.com/StephanHCB/go-autumn-acorn-registry/test/rec"
)

const MockAName = "mocka"

type MockA interface {
	IsA() bool
}

// Acorn

type MockAImpl struct {
	MockB mockb.MockB
	MockC mockc.MockC
}

func New() auacornapi.Acorn {
	rec.Add("a.New")
	return &MockAImpl{}
}

func (m *MockAImpl) AcornName() string {
	return MockAName
}

func (m *MockAImpl) AssembleAcorn(registry auacornapi.AcornRegistry) error {
	rec.Add("a.AssembleAcorn")
	m.MockB = registry.GetAcornByName(mockb.MockBName).(mockb.MockB)
	m.MockC = registry.GetAcornByName(mockc.MockCName).(mockc.MockC)

	return nil
}

func (m *MockAImpl) SetupAcorn(registry auacornapi.AcornRegistry) error {
	err := registry.SetupAfter(m.MockB.(auacornapi.Acorn))
	if err != nil {
		return err
	}

	err = registry.SetupAfter(m.MockC.(auacornapi.Acorn))
	if err != nil {
		return err
	}

	rec.Add("a.SetupAcorn")
	return nil
}

func (m *MockAImpl) TeardownAcorn(registry auacornapi.AcornRegistry) error {
	err := registry.TeardownAfter(m.MockC.(auacornapi.Acorn))
	if err != nil {
		return err
	}

	rec.Add("a.TeardownAcorn")
	return nil
}

// MockA

func (m *MockAImpl) IsA() bool {
	return m.MockB != nil && m.MockC != nil
}
