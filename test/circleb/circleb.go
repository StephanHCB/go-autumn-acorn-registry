package circleb

import (
	auacornapi "github.com/StephanHCB/go-autumn-acorn-registry/api"
	"github.com/StephanHCB/go-autumn-acorn-registry/test/circleint"
	"github.com/StephanHCB/go-autumn-acorn-registry/test/rec"
)

type MockBImpl struct {
	MockA circleint.MockA
}

func New() auacornapi.Acorn {
	rec.Add("b.New")
	return &MockBImpl{}
}

func (m *MockBImpl) AcornName() string {
	return circleint.MockBName
}

func (m *MockBImpl) AssembleAcorn(registry auacornapi.AcornRegistry) error {
	rec.Add("b.AssembleAcorn")
	m.MockA = registry.GetAcornByName(circleint.MockAName).(circleint.MockA)

	return nil
}

func (m *MockBImpl) SetupAcorn(registry auacornapi.AcornRegistry) error {
	rec.Add("b.PreSetupAcorn")
	err := registry.SetupAfter(m.MockA.(auacornapi.Acorn))
	if err != nil {
		rec.Add("b.SetupErr")
		return err
	}

	rec.Add("b.SetupAcorn")
	return nil
}

func (m *MockBImpl) TeardownAcorn(registry auacornapi.AcornRegistry) error {
	rec.Add("b.PreTeardownAcorn")
	err := registry.TeardownAfter(m.MockA.(auacornapi.Acorn))
	if err != nil {
		rec.Add("b.TeardownErr")
		return err
	}
	rec.Add("b.TeardownAcorn")
	return nil
}

// MockB

func (m *MockBImpl) IsB() bool {
	return m.MockA != nil
}
