package reversea

import (
	auacornapi "github.com/StephanHCB/go-autumn-acorn-registry/api"
	"github.com/StephanHCB/go-autumn-acorn-registry/test/rec"
	"github.com/StephanHCB/go-autumn-acorn-registry/test/reverseint"
)

// Acorn

type ReverseAImpl struct {
	ReverseB reverseint.ReverseB
}

func New() auacornapi.Acorn {
	rec.Add("a.New")
	return &ReverseAImpl{}
}

func (m *ReverseAImpl) AcornName() string {
	return reverseint.ReverseAName
}

func (m *ReverseAImpl) AssembleAcorn(registry auacornapi.AcornRegistry) error {
	rec.Add("a.AssembleAcorn")
	m.ReverseB = registry.GetAcornByName(reverseint.ReverseBName).(reverseint.ReverseB)

	return nil
}

func (m *ReverseAImpl) SetupAcorn(registry auacornapi.AcornRegistry) error {
	rec.Add("a.SetupAcorn")
	return nil
}

func (m *ReverseAImpl) TeardownAcorn(registry auacornapi.AcornRegistry) error {
	rec.Add("a.TeardownAcorn")
	return nil
}

// ReverseA

func (m *ReverseAImpl) IsA() bool {
	return m.ReverseB != nil
}
