package auacornapi

// Registry is the singleton instance of AcornRegistry provided by this library.
//
// Note: you can create your own instances, but normally you should not need to.
var Registry AcornRegistry

type Constructor func() Acorn

type AcornRegistry interface {
	// --- methods to be called by the top level application ---

	// Register an Acorn's constructor with the registry.
	//
	// During the first phase, creation, the registry calls them all in unspecified order.
	//
	// Your constructor MUST NOT assume that any other acorn is present. This is just to get non-nil
	// instance pointers for all Acorns.
	//
	// You can do internal setup of variable values, but you should relegate any time-consuming
	// activities to phase three, setup.
	Register(constructor Constructor)

	// Create should be called after all Acorns have been registered with Register().
	//
	// It will use the registered constructors to create uninitialized instances of all registered Acorns.
	//
	// This does phase one, creation.
	Create()

	// Assemble should be called after Create.
	//
	// It will call AssembleAcorn on each Acorn.
	//
	// This does phase two, assembly.
	Assemble() error

	// Setup should be called after Assemble.
	//
	// It will call SetupAcorn on each Acorn.
	//
	// This does phase three, setup.
	Setup() error

	// Teardown should be called during application shutdown.
	//
	// It will call TeardownAcorn on each Acorn.
	//
	// This does phase four, teardown.
	Teardown() error

	// --- methods to be called by Acorns ---

	// GetAcornByName gives you a reference to another Acorn.
	//
	// Should ONLY be used during the second phase, assembly, meaning in your implementation of AssembleAcorn().
	//
	// You have no guarantee that the other acorn has been assembled yet, so all you are allowed to do at this
	// point is store the reference in your instance. You are guaranteed that the return value is not nil.
	//
	// DO NOT call any methods on the other Acorn yet. It may not be ready!
	GetAcornByName(acornName string) Acorn

	// SetupAfter allows you to specify that your SetupAcorn() method depends on another Acorn being set up first.
	//
	// Should ONLY be used during the third phase, setup, typically at the beginning of your SetupAcorn().
	//
	// When it returns, you can rely on the other Acorn being set up.
	//
	// It is an error to create a circular dependency. The registry will detect this.
	SetupAfter(otherAcorn Acorn) error

	// TeardownAfter allows you to specify that your TeardownAcorn() method depends on another Acorn being torn down first.
	//
	// Should ONLY be used during the teardown phase, typically at the beginning of your TeardownAcorn().
	//
	// When it returns, you can rely on the other Acorn being torn down.
	//
	// It is an error to create a circular dependency. The registry will detect this.
	TeardownAfter(otherAcorn Acorn) error

	// --- methods useful for testing ---

	// CreateOverride lets you override an instance after create.
	//
	// MUST use before Assemble()
	//
	// useful for testing
	CreateOverride(name string, instance Acorn)

	// SkipAssemble lets you mark an instance as already assembled, so it will be skipped during Assemble().
	//
	// useful for testing
	SkipAssemble(instance Acorn)

	// SkipSetup lets you mark an instance as already set up, so it will be skipped during Setup().
	//
	// useful for testing
	SkipSetup(instance Acorn)

	// SkipTeardown lets you mark an instance as already torn down, so it will be skipped during Teardown().
	//
	// useful for testing
	SkipTeardown(instance Acorn)
}

type Acorn interface {
	// AcornName should return the package name, followed by the Acorn's primary interface name, and possibly
	// a third part for disambiguation, separated by a dot (".").
	//
	// Example:
	// Assume you have an interface config.Configuration, and its implementations implement Acorn.
	// Then you should return "config.Configuration", possibly followed by further disambiguation, in case
	// you wish to use multiple implementations at the same time.
	AcornName() string

	// AssembleAcorn gets called during the second phase, assembly.
	//
	// Use registry.GetAcornByName to obtain references to other Acorns you depend on, and store these references
	// in your instance. You should type-cast them to their primary interface so they're convenient to use.
	// DO NOT call methods on them yet, they may not be assembled yet.
	//
	// You can do internal setup here, but you MUST NOT call methods on another Acorn.
	//
	// One typical internal setup done here is reading in the application configuration to avoid a circular
	// dependency with the logging subsystem, which needs the configuration to know what log level to set etc.
	AssembleAcorn(registry AcornRegistry) error

	// SetupAcorn gets called during the third phase, setup.
	//
	// If you need another Acorn set up first, start the implementation with a call to registry.SetupAfter(),
	// giving it the reference to the other Acorn you stored during the assembly phase. Then you can rely on
	// the other Acorn being already set up.
	//
	// Note: registry.SetupAfter() will give you an error if you create a circular dependency. If you have one of those,
	// you will need to relegate part of the setup to the assembly phase so that you can resolve one of the dependencies.
	SetupAcorn(registry AcornRegistry) error

	// TeardownAcorn gets called during application shutdown.
	//
	// If you need another Acorn torn down first, start the implementation with a call to registry.TeardownAfter(),
	// giving it the reference to the other Acorn you stored during the assembly phase. Then you can rely on
	// the other Acorn being already torn down.
	//
	// Note: registry.TeardownAfter() will give you an error if you create a circular dependency. You will have to
	// decide which dependency to sever, because we must tear down in SOME order.
	TeardownAcorn(registry AcornRegistry) error
}
