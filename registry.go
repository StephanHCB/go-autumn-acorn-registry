package auacorn

import (
	"errors"
	"fmt"
	auacornapi "github.com/StephanHCB/go-autumn-acorn-registry/api"
)

const (
	phaseCreateDone   = 1
	phaseAssembleDone = 2
	phaseSetupDone    = 3
	phaseTeardownDone = 4

	phaseInRecursiveSetup    = 93 // special phase value so we can detect circular setup dependencies
	phaseInRecursiveTeardown = 94 // special phase value so we can detect circular teardown dependencies
)

type AcornRegistryImpl struct {
	constructors    []auacornapi.Constructor
	instancesByName map[string]auacornapi.Acorn
	phase           uint8
	phaseByInstance map[auacornapi.Acorn]uint8
	setupBefore     map[auacornapi.Acorn][]auacornapi.Acorn // dependency -> prerequisites
}

// Registry is the singleton instance of AcornRegistry provided by this library.
//
// Note: you can create your own instances, but normally you should not need to.
var Registry auacornapi.AcornRegistry

func init() {
	Registry = New()
}

func New() auacornapi.AcornRegistry {
	return &AcornRegistryImpl{
		constructors:    make([]auacornapi.Constructor, 0),
		instancesByName: make(map[string]auacornapi.Acorn),
		phaseByInstance: make(map[auacornapi.Acorn]uint8),
		setupBefore:     make(map[auacornapi.Acorn][]auacornapi.Acorn),
	}
}

func (a *AcornRegistryImpl) Register(constructor auacornapi.Constructor) {
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

// CreateOverride lets you override an instance after create.
//
// MUST use before Assemble()
//
// useful for testing
func (a *AcornRegistryImpl) CreateOverride(name string, instance auacornapi.Acorn) {
	a.instancesByName[name] = instance
	a.phaseByInstance[instance] = phaseCreateDone
}

func (a *AcornRegistryImpl) lifecycleStep(step string, fromPhase uint8, toPhase uint8, receiver func(auacornapi.Acorn) error) error {
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

// SkipAssemble lets you mark an instance as already assembled, so it will be skipped during Assemble().
//
// useful for testing
func (a *AcornRegistryImpl) SkipAssemble(instance auacornapi.Acorn) {
	a.phaseByInstance[instance] = phaseAssembleDone
}

func (a *AcornRegistryImpl) Assemble() error {
	if a.phase != phaseCreateDone {
		return errors.New("wrong acorn registry phase order: Assemble() comes after Create()")
	}
	return a.lifecycleStep("assembly", phaseCreateDone, phaseAssembleDone, func(instance auacornapi.Acorn) error {
		return instance.AssembleAcorn(a)
	})
}

// SkipSetup lets you mark an instance as already set up, so it will be skipped during Setup().
//
// useful for testing
func (a *AcornRegistryImpl) SkipSetup(instance auacornapi.Acorn) {
	a.phaseByInstance[instance] = phaseSetupDone
}

func (a *AcornRegistryImpl) injectExtraSetupAfterCallsThenSetup(instance auacornapi.Acorn) error {
	extraPrerequisites, ok := a.setupBefore[instance]
	if ok {
		for _, prerequisite := range extraPrerequisites {
			err := a.SetupAfter(prerequisite)
			if err != nil {
				return err
			}
		}
	}
	return instance.SetupAcorn(a)
}

func (a *AcornRegistryImpl) Setup() error {
	if a.phase != phaseAssembleDone {
		return errors.New("wrong acorn registry phase order: Setup() comes after Assemble()")
	}
	return a.lifecycleStep("setup", phaseAssembleDone, phaseSetupDone, func(instance auacornapi.Acorn) error {
		return a.injectExtraSetupAfterCallsThenSetup(instance)
	})
}

// SkipTeardown lets you mark an instance as already torn down, so it will be skipped during Teardown().
//
// useful for testing
func (a *AcornRegistryImpl) SkipTeardown(instance auacornapi.Acorn) {
	a.phaseByInstance[instance] = phaseTeardownDone
}

func (a *AcornRegistryImpl) Teardown() error {
	// we allow teardown even for lower phase numbers, so partial setup can be cleaned up
	return a.lifecycleStep("teardown", phaseSetupDone, phaseTeardownDone, func(instance auacornapi.Acorn) error {
		return instance.TeardownAcorn(a)
	})
}

func (a *AcornRegistryImpl) GetAcornByName(acornName string) auacornapi.Acorn {
	return a.instancesByName[acornName]
}

func (a *AcornRegistryImpl) SetupAfter(otherAcorn auacornapi.Acorn) error {
	if a.phase != phaseAssembleDone {
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
	err := a.injectExtraSetupAfterCallsThenSetup(otherAcorn)
	a.phaseByInstance[otherAcorn] = phaseSetupDone
	return err
}

func (a *AcornRegistryImpl) TeardownAfter(otherAcorn auacornapi.Acorn) error {
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

func (a *AcornRegistryImpl) AddSetupOrderRule(prerequisite auacornapi.Acorn, dependency auacornapi.Acorn) error {
	if a.phase != phaseCreateDone {
		return errors.New("wrong acorn registry phase for call to AddSetupOrderRule() - only allowed during assembly phase")
	}
	if prerequisite == nil || dependency == nil {
		return errors.New("cannot add setup order rule for nil acorns")
	}

	currentSetupBefore, ok := a.setupBefore[dependency]
	if !ok {
		currentSetupBefore = make([]auacornapi.Acorn, 0)
	}

	a.setupBefore[dependency] = append(currentSetupBefore, prerequisite)
	return nil
}
