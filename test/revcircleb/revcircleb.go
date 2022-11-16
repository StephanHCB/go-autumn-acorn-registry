package revcircleb

import (
	auacornapi "github.com/StephanHCB/go-autumn-acorn-registry/api"
	"github.com/StephanHCB/go-autumn-acorn-registry/test/rec"
	"github.com/StephanHCB/go-autumn-acorn-registry/test/revcircleint"
)

type RevCircleBImpl struct {
	RevCircleA revcircleint.RevCircleA
}

func New() auacornapi.Acorn {
	rec.Add("b.New")
	return &RevCircleBImpl{}
}

func (m *RevCircleBImpl) AcornName() string {
	return revcircleint.RevCircleBName
}

func (m *RevCircleBImpl) AssembleAcorn(registry auacornapi.AcornRegistry) error {
	rec.Add("b.AssembleAcorn")
	a := registry.GetAcornByName(revcircleint.RevCircleAName)
	m.RevCircleA = a.(revcircleint.RevCircleA)

	// specifies: "b" (this Acorn) must set up before "a"
	err := registry.AddSetupOrderRule(m, a)
	if err != nil {
		rec.Add("b.AssembleErr")
		return err
	}

	return nil
}

func (m *RevCircleBImpl) SetupAcorn(registry auacornapi.AcornRegistry) error {
	rec.Add("b.PreSetupAcorn")

	// specifies: "b" (this Acorn) must set up after "a"
	err := registry.SetupAfter(m.RevCircleA.(auacornapi.Acorn))
	if err != nil {
		rec.Add("b.SetupErr")
		return err
	}

	rec.Add("b.SetupAcorn")
	return nil
}

func (m *RevCircleBImpl) TeardownAcorn(registry auacornapi.AcornRegistry) error {
	rec.Add("b.PreTeardownAcorn")
	err := registry.TeardownAfter(m.RevCircleA.(auacornapi.Acorn))
	if err != nil {
		rec.Add("b.TeardownErr")
		return err
	}
	rec.Add("b.TeardownAcorn")
	return nil
}

// RevCircleB

func (m *RevCircleBImpl) IsB() bool {
	return m.RevCircleA != nil
}
