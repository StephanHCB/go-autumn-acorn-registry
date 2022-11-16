package auacorn

import (
	"fmt"
	"github.com/StephanHCB/go-autumn-acorn-registry/test/circlea"
	"github.com/StephanHCB/go-autumn-acorn-registry/test/circleb"
	"github.com/StephanHCB/go-autumn-acorn-registry/test/circleint"
	"github.com/StephanHCB/go-autumn-acorn-registry/test/mocka"
	"github.com/StephanHCB/go-autumn-acorn-registry/test/mockb"
	"github.com/StephanHCB/go-autumn-acorn-registry/test/mockc"
	"github.com/StephanHCB/go-autumn-acorn-registry/test/rec"
	"github.com/StephanHCB/go-autumn-acorn-registry/test/revcirclea"
	"github.com/StephanHCB/go-autumn-acorn-registry/test/revcircleb"
	"github.com/StephanHCB/go-autumn-acorn-registry/test/revcircleint"
	"github.com/StephanHCB/go-autumn-acorn-registry/test/reversea"
	"github.com/StephanHCB/go-autumn-acorn-registry/test/reverseb"
	"github.com/StephanHCB/go-autumn-acorn-registry/test/reverseint"
	"os"
	"strings"
	"testing"
)

func assertRecording(t *testing.T, expected ...[]string) {
	actualArr := fmt.Sprintf("%v", rec.Get())
	allowedArrs := ""
	for _, allowed := range expected {
		allowedArrs += fmt.Sprintf("%v", allowed)
	}

	if !strings.Contains(allowedArrs, actualArr) {
		_, _ = os.Stderr.WriteString("actual : " + actualArr + "\n")
		_, _ = os.Stderr.WriteString("allowed: " + allowedArrs + "\n")
		t.FailNow()
	}
}

func TestRegistry_NormalLifecycle(t *testing.T) {
	Registry = New()

	Registry.Register(mocka.New)
	Registry.Register(mockb.New)
	Registry.Register(mockc.New)

	rec.Reset()
	Registry.Create()
	assertRecording(t, []string{"a.New", "b.New", "c.New"})

	rec.Reset()
	err := Registry.Assemble()
	if err != nil {
		t.FailNow()
	}
	// assemble has no order
	assertRecording(t,
		[]string{"a.AssembleAcorn", "b.AssembleAcorn", "c.AssembleAcorn"},
		[]string{"a.AssembleAcorn", "c.AssembleAcorn", "b.AssembleAcorn"},
		[]string{"b.AssembleAcorn", "a.AssembleAcorn", "c.AssembleAcorn"},
		[]string{"b.AssembleAcorn", "c.AssembleAcorn", "a.AssembleAcorn"},
		[]string{"c.AssembleAcorn", "a.AssembleAcorn", "b.AssembleAcorn"},
		[]string{"c.AssembleAcorn", "b.AssembleAcorn", "a.AssembleAcorn"},
	)

	rec.Reset()
	err = Registry.Setup()
	if err != nil {
		t.FailNow()
	}
	// dependencies make c set up first, then b, then a
	assertRecording(t, []string{"c.SetupAcorn", "b.SetupAcorn", "a.SetupAcorn"})

	a := Registry.GetAcornByName(mocka.MockAName).(mocka.MockA)
	b := Registry.GetAcornByName(mockb.MockBName).(mockb.MockB)
	c := Registry.GetAcornByName(mockc.MockCName).(mockc.MockC)
	// check internal pointers are correctly set
	if !a.IsA() || !b.IsB() || !c.IsC() {
		t.FailNow()
	}

	rec.Reset()
	err = Registry.Teardown()
	if err != nil {
		t.FailNow()
	}
	// dependencies make c tear down first, then a and b in any order (no dependencies)
	assertRecording(t,
		[]string{"c.TeardownAcorn", "a.TeardownAcorn", "b.TeardownAcorn"},
		[]string{"c.TeardownAcorn", "b.TeardownAcorn", "a.TeardownAcorn"})
}

func TestRegistry_CircleDetection(t *testing.T) {
	Registry = New()

	Registry.Register(circlea.New)
	Registry.Register(circleb.New)

	rec.Reset()
	Registry.Create()
	assertRecording(t, []string{"a.New", "b.New"})

	rec.Reset()
	err := Registry.Assemble()
	if err != nil {
		t.FailNow()
	}
	// assemble has no order
	assertRecording(t,
		[]string{"a.AssembleAcorn", "b.AssembleAcorn"},
		[]string{"b.AssembleAcorn", "a.AssembleAcorn"},
	)

	rec.Reset()
	err = Registry.Setup()
	if err == nil {
		t.FailNow()
	}
	assertRecording(t,
		[]string{"a.PreSetupAcorn", "b.PreSetupAcorn", "a.PreSetupAcorn", "a.SetupErr", "b.SetupErr", "a.SetupErr"},
		[]string{"b.PreSetupAcorn", "a.PreSetupAcorn", "b.PreSetupAcorn", "b.SetupErr", "a.SetupErr", "b.SetupErr"},
	)

	a := Registry.GetAcornByName(circleint.MockAName).(circleint.MockA)
	b := Registry.GetAcornByName(circleint.MockBName).(circleint.MockB)
	// check internal pointers are correctly set
	if !a.IsA() || !b.IsB() {
		t.FailNow()
	}

	rec.Reset()
	err = Registry.Teardown()
	if err == nil {
		t.FailNow()
	}
	assertRecording(t,
		[]string{"a.PreTeardownAcorn", "b.PreTeardownAcorn", "a.PreTeardownAcorn", "a.TeardownErr", "b.TeardownErr", "a.TeardownErr"},
		[]string{"b.PreTeardownAcorn", "a.PreTeardownAcorn", "b.PreTeardownAcorn", "b.TeardownErr", "a.TeardownErr", "b.TeardownErr"},
	)
}

func TestRegistry_NormalLifecycle_WithReverseDependency(t *testing.T) {
	Registry = New()

	Registry.Register(reversea.New)
	Registry.Register(reverseb.New)

	rec.Reset()
	Registry.Create()
	assertRecording(t, []string{"a.New", "b.New"})

	rec.Reset()
	err := Registry.Assemble()
	if err != nil {
		t.FailNow()
	}
	// assemble has no order
	assertRecording(t,
		[]string{"a.AssembleAcorn", "b.AssembleAcorn"},
		[]string{"b.AssembleAcorn", "a.AssembleAcorn"},
	)

	rec.Reset()
	err = Registry.Setup()
	if err != nil {
		t.FailNow()
	}
	// reverse specified dependencies make b set up first, then a
	assertRecording(t, []string{"b.SetupAcorn", "a.SetupAcorn"})

	a := Registry.GetAcornByName(reverseint.ReverseAName).(reverseint.ReverseA)
	b := Registry.GetAcornByName(reverseint.ReverseBName).(reverseint.ReverseB)
	// check internal pointers are correctly set
	if !a.IsA() || !b.IsB() {
		t.FailNow()
	}

	rec.Reset()
	err = Registry.Teardown()
	if err != nil {
		t.FailNow()
	}
	// tear down a and b in any order (no dependencies)
	assertRecording(t,
		[]string{"a.TeardownAcorn", "b.TeardownAcorn"},
		[]string{"b.TeardownAcorn", "a.TeardownAcorn"})
}

func TestRegistry_CircleDetection_WithReverseDependency(t *testing.T) {
	Registry = New()

	Registry.Register(revcirclea.New)
	Registry.Register(revcircleb.New)

	rec.Reset()
	Registry.Create()
	assertRecording(t, []string{"a.New", "b.New"})

	rec.Reset()
	err := Registry.Assemble()
	if err != nil {
		t.FailNow()
	}
	// assemble has no order
	assertRecording(t,
		[]string{"a.AssembleAcorn", "b.AssembleAcorn"},
		[]string{"b.AssembleAcorn", "a.AssembleAcorn"},
	)

	rec.Reset()
	err = Registry.Setup()
	if err == nil {
		t.FailNow()
	}
	assertRecording(t,
		// it picks a first, and the circle is detected the second time the additional order rule triggers
		[]string{"b.PreSetupAcorn", "b.SetupErr"},
		// it picks b first, and the circle is detected the second time it enters b.SetupAcorn()
		[]string{"b.PreSetupAcorn", "b.PreSetupAcorn", "b.SetupErr", "b.SetupErr"},
	)

	a := Registry.GetAcornByName(revcircleint.RevCircleAName).(revcircleint.RevCircleA)
	b := Registry.GetAcornByName(revcircleint.RevCircleBName).(revcircleint.RevCircleB)
	// check internal pointers are correctly set
	if !a.IsA() || !b.IsB() {
		t.FailNow()
	}

	rec.Reset()
	err = Registry.Teardown()
	if err == nil {
		t.FailNow()
	}
	assertRecording(t,
		[]string{"a.PreTeardownAcorn", "b.PreTeardownAcorn", "a.PreTeardownAcorn", "a.TeardownErr", "b.TeardownErr", "a.TeardownErr"},
		[]string{"b.PreTeardownAcorn", "a.PreTeardownAcorn", "b.PreTeardownAcorn", "b.TeardownErr", "a.TeardownErr", "b.TeardownErr"},
	)
}
