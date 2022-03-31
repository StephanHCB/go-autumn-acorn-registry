package circlea

import (
	auacornapi "github.com/StephanHCB/go-autumn-acorn-registry/api"
	"github.com/StephanHCB/go-autumn-acorn-registry/test/circleint"
	"github.com/StephanHCB/go-autumn-acorn-registry/test/rec"
)

// Acorn

type MockAImpl struct {
	MockB circleint.MockB
}

func New() auacornapi.Acorn {
	rec.Add("a.New")
	return &MockAImpl{}
}

func (m *MockAImpl) AcornName() string {
	return circleint.MockAName
}

func (m *MockAImpl) AssembleAcorn(registry auacornapi.AcornRegistry) error {
	rec.Add("a.AssembleAcorn")
	m.MockB = registry.GetAcornByName(circleint.MockBName).(circleint.MockB)

	return nil
}

func (m *MockAImpl) SetupAcorn(registry auacornapi.AcornRegistry) error {
	rec.Add("a.PreSetupAcorn")
	err := registry.SetupAfter(m.MockB.(auacornapi.Acorn))
	if err != nil {
		rec.Add("a.SetupErr")
		return err
	}

	rec.Add("a.SetupAcorn")
	return nil
}

func (m *MockAImpl) TeardownAcorn(registry auacornapi.AcornRegistry) error {
	rec.Add("a.PreTeardownAcorn")
	err := registry.TeardownAfter(m.MockB.(auacornapi.Acorn))
	if err != nil {
		rec.Add("a.TeardownErr")
		return err
	}

	rec.Add("a.TeardownAcorn")
	return nil
}

// MockA

func (m *MockAImpl) IsA() bool {
	return m.MockB != nil
}
