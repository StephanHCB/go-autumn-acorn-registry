package reverseb

import (
	auacornapi "github.com/StephanHCB/go-autumn-acorn-registry/api"
	"github.com/StephanHCB/go-autumn-acorn-registry/test/rec"
	"github.com/StephanHCB/go-autumn-acorn-registry/test/reverseint"
)

type ReverseBImpl struct {
}

func New() auacornapi.Acorn {
	rec.Add("b.New")
	return &ReverseBImpl{}
}

func (m *ReverseBImpl) AcornName() string {
	return reverseint.ReverseBName
}

func (m *ReverseBImpl) AssembleAcorn(registry auacornapi.AcornRegistry) error {
	rec.Add("b.AssembleAcorn")

	// this is the interesting part of this test case
	//   specify: "b" (this Acorn) must set up before "a"
	a := registry.GetAcornByName(reverseint.ReverseAName)
	err := registry.AddSetupOrderRule(m, a)
	if err != nil {
		rec.Add("b.AssembleErr")
		return err
	}

	return nil
}

func (m *ReverseBImpl) SetupAcorn(registry auacornapi.AcornRegistry) error {
	rec.Add("b.SetupAcorn")
	return nil
}

func (m *ReverseBImpl) TeardownAcorn(registry auacornapi.AcornRegistry) error {
	rec.Add("b.TeardownAcorn")
	return nil
}

// ReverseB

func (m *ReverseBImpl) IsB() bool {
	return true
}
