# go-autumn-acorn-registry

A registry for singleton components that comprise an application. Also helps you manage dependencies
between the components.

## About go-autumn

A collection of libraries for [enterprise microservices](https://github.com/StephanHCB/go-mailer-service/blob/master/README.md) in golang that
- is heavily inspired by Spring Boot / Spring Cloud
- is very opinionated
- names modules by what they do
- unlike Spring Boot avoids certain types of auto-magical behaviour
- is not a library monolith, that is every part only depends on the api parts of the other components
  at most, and the api parts do not add any dependencies.  

Fall is my favourite season, so I'm calling it go-autumn.

## About go-autumn-acorn-registry

One of the most flexible ways to assemble an application is building it out of singleton instances that
implement an interface. You create them at program startup. Because this is go-autumn, let's call them **Acorn**s. 

The various instances that build your application keep references to each other, only referring to their respective 
interfaces. This allows easy mocking and switching between multiple implementations.

But then wiring up the application manually can become cumbersome. This is where go-autumn-acorn-registry steps in,
giving you a very easy way to code the application assembly without relying on any kind of reflection framework
or auto-magical behavior.

As far as I am concerned, this has most benefits of dependency injection without some of its major drawbacks.

### Features

- zero dependencies
- no reflection
- you have full control

## Acorn lifecycle

### 1. creation

The registry calls the constructor of your Acorn. It can do that because you have given it a reference to it.

The constructor must not take any arguments, and must not assume any of the other Acorns are already non-nil.
You must not rely on creation order.

You can set up non-Acorn fields if you wish, but under no circumstances should you try accessing the registry
for any other Acorn.

You must also avoid doing anything that has side effects, or you'll lose a very convenient approach to mocking
(the constructor gets called before the mocks can be inserted in tests).

### 2. assembly

Now that all the Acorns have been instantiated, it's time to wire up the references between them.

The registry then calls the `AssembleAcorn()` method of all Acorns in unspecified order. It is given a reference to the registry
which you can use to look up other Acorns by their name. You should cast each reference to its primary 
interface type in order to store it in your instance.

You MUST NOT access the other Acorns at this time, but it is perfectly alright to do Acorn-internal early setup tasks.

One typical example is loading the application configuration. This lets you get out of a very typical circular
dependency between configuration and logging.

### 3. setup

Now that all Acorns are wired up, the registry picks one to set up first. It calls its `SetupAcorn()` method, 
again passing in a reference to the registry.

If you need another component set up first, there's a method `registry.SetupAfter(otherAcorn Acorn)` which you can 
use to specify the dependency tree. Just call it before using the reference to the other Acorn. If it has already been
set up, this call does nothing, but if it hasn't, it will recurse into its `SetupAcorn()` method.

It is a run time error to set up circular dependencies in Setup(), and this is detected by the registry and
an error is raised.

### 4. teardown

When it comes to tearing the application down and doing cleanup, once again the registry will call your
`TeardownAcorn()` method, giving you a reference to the registry.

If you need another component torn down first, there's a method `registry.TeardownAfter(otherAcorn Acorn)`
which you can use to specify the teardown dependency tree. Just call it before proceeding with your
own teardown. If the other Acorn has already been torn down, this call does nothing, but if it hasn't,
it will recurse into its `TeardownAcorn()` method.

## Usage

### Making your interface implementation an Acorn

1. Provide a parameterless constructor that returns your instance cast to Acorn:

`func New() Acorn`

It does not have to be called `New`, but that is a convention I find useful.

2. Implement the Acorn interface. Its methods are 
   
  - `AcornName() string` which should normally return the package and primary interface name of your Acorn separated by a `.`, 
    for use by other Acorn's `AssembleAcorn()` to look it up.
    If you have multiple implementations, and wish to use them at the same time, disambiguate the name by adding another `.`
    followed by the implementation. 
  - `AssembleAcorn(registry AcornRegistry) error`
  - `SetupAcorn(registry AcornRegistry) error` 
  - `TeardownAcorn(registry AcornRegistry) error`

You should do so using a pointer receiver.

### Informing the registry about an Acorn

Call `registry.Register(mypackage.New)` with a reference to your constructor.

I find it useful to collect all these calls in a top level application class.

_Note: since you are registering the Acorn constructors yourself, you can change the constructor that you
are registering. This allows very easy switching to mocks for tests. You can even overwrite already
registered constructors by providing an implementation with the same return value of `AcornName()`,
because the registry remembers registration order, and the last one wins._

### Testing

During test scenarios, you have several methods that you can call between the major lifecycle phases

#### Before Create()

You can register extra testing-only Acorns, which may later do things like populate caches, etc.

#### Between Create() and Assemble()

`CreateOverride(name string, instance Acorn)` lets you replace an instance before it is wired into
all the other instances during Assemble(). You need to take care to implement the same interfaces as the
original Acorn, and you'll need to provide the correct name.

`SkipAssemble(instance Acorn)` lets you mark an Acorn as already assembled, so it will be skipped.
Of course this means you're not getting a call to your AssembleAcorn() method, so you'll need to do
any reference wiring yourself. Since your implementation is likely a mock anyway, this may be exactly what you want.

#### Between Assemble() and Setup()

`SkipSetup(instance Acorn)` lets you mark an Acorn as already set up, so it will be skipped.
This works even when another acorn requests your Acorn to be set up first, the logic just thinks
it's already set up and does nothing.

#### Before Teardown()

`SkipTeardown(instance Acorn)` lets you mark an Acorn as already torn down.

### Library Authors

Implement the `Acorn` interface for any class that an application author might wish to directly wire up as
a singleton (an "Acorn") in their application.

Otherwise, do not interact with the registry, that is for the application author to do. You can't know if the
author wants to use your component as a singleton.

### Application Authors

Implement the `Acorn` interface in the singleton components of your application.

Put some place in the code where you `Register()` all your acorns with the registry.

Once that's done, call 

  - `registry.Create()`
  - `registry.Assemble()`
  - `registry.Setup()`

in that order.

During teardown, call `registry.Teardown()`.

That's it.
