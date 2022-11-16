package revcirclea

import (
	auacornapi "github.com/StephanHCB/go-autumn-acorn-registry/api"
	"github.com/StephanHCB/go-autumn-acorn-registry/test/rec"
	"github.com/StephanHCB/go-autumn-acorn-registry/test/revcircleint"
)

// Acorn

type RevCircleAImpl struct {
	RevCircleB revcircleint.RevCircleB
}

func New() auacornapi.Acorn {
	rec.Add("a.New")
	return &RevCircleAImpl{}
}

func (m *RevCircleAImpl) AcornName() string {
	return revcircleint.RevCircleAName
}

func (m *RevCircleAImpl) AssembleAcorn(registry auacornapi.AcornRegistry) error {
	rec.Add("a.AssembleAcorn")
	m.RevCircleB = registry.GetAcornByName(revcircleint.RevCircleBName).(revcircleint.RevCircleB)

	return nil
}

func (m *RevCircleAImpl) SetupAcorn(registry auacornapi.AcornRegistry) error {
	rec.Add("a.SetupAcorn")
	return nil
}

func (m *RevCircleAImpl) TeardownAcorn(registry auacornapi.AcornRegistry) error {
	rec.Add("a.PreTeardownAcorn")
	err := registry.TeardownAfter(m.RevCircleB.(auacornapi.Acorn))
	if err != nil {
		rec.Add("a.TeardownErr")
		return err
	}

	rec.Add("a.TeardownAcorn")
	return nil
}

// RevCircleA

func (m *RevCircleAImpl) IsA() bool {
	return m.RevCircleB != nil
}
