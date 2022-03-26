package goauacorn

import (
	"errors"
	"fmt"
	goauacornapi "github.com/StephanHCB/go-autumn-acorn-registry/api"
)

const (
	phaseCreateDone = 1
	phaseAssembleDone = 2
	phaseSetupDone = 3
	phaseTeardownDone = 4

	phaseInRecursiveSetup = 93 // special phase value so we can detect circular setup dependencies
	phaseInRecursiveTeardown = 94 // special phase value so we can detect circular teardown dependencies
)

type AcornRegistryImpl struct {
	constructors    []goauacornapi.Constructor
	instancesByName map[string]goauacornapi.Acorn
	phase           uint8
	phaseByInstance map[goauacornapi.Acorn]uint8
}

func init() {
	goauacornapi.Registry = New()
}

func New() goauacornapi.AcornRegistry {
	return &AcornRegistryImpl{
		constructors:    make([]goauacornapi.Constructor, 0),
		instancesByName: make(map[string]goauacornapi.Acorn),
		phaseByInstance: make(map[goauacornapi.Acorn]uint8),
	}
}

func (a *AcornRegistryImpl) Register(constructor goauacornapi.Constructor) {
	a.constructors = append(a.constructors, constructor)
}

func (a *AcornRegistryImpl) Create() {
	for _, constructor := range a.constructors {
		instance := constructor()
		name := instance.AcornName()
		a.instancesByName[name] = instance
		a.phaseByInstance[instance] = phaseCreateDone
	}
	a.phase = phaseCreateDone
}

func (a *AcornRegistryImpl) lifecycleStep(step string, fromPhase uint8, toPhase uint8, receiver func(goauacornapi.Acorn) error) error {
	for name, instance := range a.instancesByName {
		if a.phaseByInstance[instance] == fromPhase {
			// only do the phase if it hasn't already been done
			err := receiver(instance)
			if err != nil {
				return fmt.Errorf("error during %s of Acorn '%s': %s", step, name, err.Error())
			}
			a.phaseByInstance[instance] = toPhase
		}
	}
	a.phase = toPhase
	return nil
}

func (a *AcornRegistryImpl) Assemble() error {
	if a.phase != 1 {
		return errors.New("wrong acorn registry phase order: Assemble() comes after Create()")
	}
	return a.lifecycleStep("assembly", phaseCreateDone, phaseAssembleDone, func(instance goauacornapi.Acorn) error {
		return instance.AssembleAcorn(a)
	})
}

func (a *AcornRegistryImpl) Setup() error {
	if a.phase != 2 {
		return errors.New("wrong acorn registry phase order: Setup() comes after Assemble()")
	}
	return a.lifecycleStep("setup", phaseAssembleDone, phaseSetupDone, func(instance goauacornapi.Acorn) error {
		return instance.SetupAcorn(a)
	})
}

func (a *AcornRegistryImpl) Teardown() error {
	// we allow teardown even for lower phase numbers, so partial setup can be cleaned up
	return a.lifecycleStep("teardown", phaseSetupDone, phaseTeardownDone, func(instance goauacornapi.Acorn) error {
		return instance.TeardownAcorn(a)
	})
}

func (a *AcornRegistryImpl) GetAcornByName(acornName string) goauacornapi.Acorn {
	return a.instancesByName[acornName]
}

func (a *AcornRegistryImpl) SetupAfter(otherAcorn goauacornapi.Acorn) error {
	if a.phase != 2 {
		return errors.New("wrong acorn registry phase for call to SetupAfter() - only allowed during setup phase")
	}
	if a.phaseByInstance[otherAcorn] == phaseInRecursiveSetup {
		// circular dependency
		return fmt.Errorf("circular setup dependency involving Acorn %s - not allowed", otherAcorn.AcornName())
	}
	if a.phaseByInstance[otherAcorn] != phaseAssembleDone {
		// was already set up, that is ok
		return nil
	}

	a.phaseByInstance[otherAcorn] = phaseInRecursiveSetup
	err := otherAcorn.SetupAcorn(a)
	a.phaseByInstance[otherAcorn] = phaseSetupDone
	return err
}

func (a *AcornRegistryImpl) TeardownAfter(otherAcorn goauacornapi.Acorn) error {
	if a.phaseByInstance[otherAcorn] == phaseInRecursiveTeardown {
		// circular dependency
		return fmt.Errorf("circular teardown dependency involving Acorn %s - not allowed", otherAcorn.AcornName())
	}
	if a.phaseByInstance[otherAcorn] != phaseSetupDone {
		// was already torn down, or never set up, that is ok
		return nil
	}

	a.phaseByInstance[otherAcorn] = phaseInRecursiveTeardown
	err := otherAcorn.TeardownAcorn(a)
	a.phaseByInstance[otherAcorn] = phaseTeardownDone
	return err
}
