// +build 1.6,codegen

package api

import (
	"testing"
)

func TestUniqueInputAndOutputs(t *testing.T) {
	shamelist["FooService"] = map[string]struct {
		input  bool
		output bool
	}{}
	v := shamelist["FooService"]["OpOutputNoRename"]
	v.output = true
	shamelist["FooService"]["OpOutputNoRename"] = v
	v = shamelist["FooService"]["InputNoRename"]
	v.input = true
	shamelist["FooService"]["OpInputNoRename"] = v
	v = shamelist["FooService"]["BothNoRename"]
	v.input = true
	v.output = true
	shamelist["FooService"]["OpBothNoRename"] = v

	cases := [][]struct {
		expectedInput  string
		expectedOutput string
		operation      string
		input          string
		inputRef       string
		output         string
		outputRef      string
	}{
		{
			{
				expectedInput:  "FooOperationInput",
				expectedOutput: "FooOperationOutput",
				operation:      "FooOperation",
				input:          "FooInputShape",
				inputRef:       "FooInputShapeRef",
				output:         "FooOutputShape",
				outputRef:      "FooOutputShapeRef",
			},
			{
				expectedInput:  "BarOperationInput",
				expectedOutput: "BarOperationOutput",
				operation:      "BarOperation",
				input:          "FooInputShape",
				inputRef:       "FooInputShapeRef",
				output:         "FooOutputShape",
				outputRef:      "FooOutputShapeRef",
			},
		},
		{
			{
				expectedInput:  "FooOperationInput",
				expectedOutput: "FooOperationOutput",
				operation:      "FooOperation",
				input:          "FooInputShape",
				inputRef:       "FooInputShapeRef",
				output:         "FooOutputShape",
				outputRef:      "FooOutputShapeRef",
			},
			{
				expectedInput:  "OpOutputNoRenameInput",
				expectedOutput: "OpOutputNoRenameOutputShape",
				operation:      "OpOutputNoRename",
				input:          "OpOutputNoRenameInputShape",
				inputRef:       "OpOutputNoRenameInputRef",
				output:         "OpOutputNoRenameOutputShape",
				outputRef:      "OpOutputNoRenameOutputRef",
			},
		},
		{
			{
				expectedInput:  "FooOperationInput",
				expectedOutput: "FooOperationOutput",
				operation:      "FooOperation",
				input:          "FooInputShape",
				inputRef:       "FooInputShapeRef",
				output:         "FooOutputShape",
				outputRef:      "FooOutputShapeRef",
			},
			{
				expectedInput:  "OpInputNoRenameInputShape",
				expectedOutput: "OpInputNoRenameOutput",
				operation:      "OpInputNoRename",
				input:          "OpInputNoRenameInputShape",
				inputRef:       "OpInputNoRenameInputRef",
				output:         "OpInputNoRenameOutputShape",
				outputRef:      "OpInputNoRenameOutputRef",
			},
		},
		{
			{
				expectedInput:  "FooOperationInput",
				expectedOutput: "FooOperationOutput",
				operation:      "FooOperation",
				input:          "FooInputShape",
				inputRef:       "FooInputShapeRef",
				output:         "FooOutputShape",
				outputRef:      "FooOutputShapeRef",
			},
			{
				expectedInput:  "OpInputNoRenameInputShape",
				expectedOutput: "OpInputNoRenameOutputShape",
				operation:      "OpBothNoRename",
				input:          "OpInputNoRenameInputShape",
				inputRef:       "OpInputNoRenameInputRef",
				output:         "OpInputNoRenameOutputShape",
				outputRef:      "OpInputNoRenameOutputRef",
			},
		},
	}

	for _, c := range cases {
		a := &API{
			name:       "FooService",
			Operations: map[string]*Operation{},
		}

		expected := map[string][]string{}
		a.Shapes = map[string]*Shape{}
		for _, op := range c {
			a.Operations[op.operation] = &Operation{
				ExportedName: op.operation,
			}
			a.Operations[op.operation].Name = op.operation
			a.Operations[op.operation].InputRef = ShapeRef{
				API:       a,
				ShapeName: op.inputRef,
				Shape: &Shape{
					API:       a,
					ShapeName: op.input,
				},
			}
			a.Operations[op.operation].OutputRef = ShapeRef{
				API:       a,
				ShapeName: op.outputRef,
				Shape: &Shape{
					API:       a,
					ShapeName: op.output,
				},
			}

			a.Shapes[op.input] = &Shape{
				ShapeName: op.input,
			}
			a.Shapes[op.output] = &Shape{
				ShapeName: op.output,
			}

			expected[op.operation] = append(expected[op.operation], op.expectedInput)
			expected[op.operation] = append(expected[op.operation], op.expectedOutput)
		}

		a.fixStutterNames()
		a.renameToplevelShapes()
		for k, v := range expected {
			if a.Operations[k].InputRef.Shape.ShapeName != v[0] {
				t.Errorf("Error %d case: Expected %q, but received %q", k, v[0], a.Operations[k].InputRef.Shape.ShapeName)
			}
			if a.Operations[k].OutputRef.Shape.ShapeName != v[1] {
				t.Errorf("Error %d case: Expected %q, but received %q", k, v[1], a.Operations[k].OutputRef.Shape.ShapeName)
			}
		}

	}
}

func TestCollidingFields(t *testing.T) {
	cases := map[string]struct {
		Members map[string]*ShapeRef
		Expect  []string
	}{
		"SimpleMembers": {
			MemberRefs: map[string]*ShapeRef{
				"String":   &ShapeRef{},
				"GoString": &ShapeRef{},
				"Validate": &ShapeRef{},
				"Foo":      &ShapeRef{},
				"SetFoo":   &ShapeRef{},
				"Code":     &ShapeRef{},
				"Message":  &ShapeRef{},
				"OrigErr":  &ShapeRef{},
			},
			Expect: []string{
				"String_",
				"GoString_",
				"Validate_",
				"Foo",
				"SetFoo_",
				"Code",
				"Message",
				"OrigErr",
			},
		},
		"ExceptionShape": {
			MemberRefs: map[string]*ShapeRef{
				"Code":    &ShapeRef{Shape: &Shape{Exception: true}},
				"Message": &ShapeRef{Shape: &Shape{Exception: true}},
				"OrigErr": &ShapeRef{Shape: &Shape{Exception: true}},
				"String":  &ShapeRef{Shape: &Shape{Exception: true}},
				"Other":   &ShapeRef{Shape: &Shape{Exception: true}},
			},
			Expect: []string{
				"Code_",
				"Message_",
				"OrigErr_",
				"String_",
				"Other",
			},
		},
	}

	for _, c := range testCases {
		t.Run(k, func(t *testing.T) {
			a := &API{
				Shapes: []*Shape{
					ShapeName:  k,
					MemberRefs: c.Members,
				},
			}

			a.renameCollidingFields()

			for i, name := range a.Shapes[0].MemberNames() {
				if e, a := c.Expect[i], name; e != a {
					t.Errorf("expect %v, got %v", e, a)
				}
			}
		})
	}
}
