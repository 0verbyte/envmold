package mold

import (
	"strings"
	"testing"
)

func TestMoldNew(t *testing.T) {
	var moldTemplate = `
- name: foo
  value: "bar"
  type: string
  required: true

- name: debug
  value: true
  type: boolean
  required: false
`

	if _, err := New(strings.NewReader(moldTemplate)); err != nil {
		t.Errorf("Failed to create new mold: %v", err)
		return
	}
}

func TestCheckKeysInMold(t *testing.T) {
	var moldTemplate = `
- name: foo
  value: "bar"
  type: string
  required: true

- name: debug
  value: true
  type: boolean
  required: false
`

	mold, err := New(strings.NewReader(moldTemplate))
	if err != nil {
		t.Errorf("Failed to create new mold: %v", err)
		return
	}

	specs := []struct {
		input    string
		expected string
	}{
		{input: "foo", expected: "foo"},
		{input: "debug", expected: "debug"},
	}

	for _, spec := range specs {
		envVar, err := mold.GetVariable(spec.input)
		if err != nil {
			t.Errorf("Failed to get '%s' environment variable: %v", spec.input, err)
		}
		if envVar.Name != spec.expected {
			t.Errorf("Expected '%s', got '%s'", spec.expected, envVar.Name)
		}
	}
}

func TestVerifyTypeConstraint(t *testing.T) {
	var moldTemplate = `
- name: foo
  value: "bar"
  type: boolean
  required: true

- name: debug
  value: true
  type: boolean
  required: false
`

	if _, err := New(strings.NewReader(moldTemplate)); err != ErrInvalidDataType {
		t.Errorf("Failed to create new mold: %v", err)
		return
	}
}
