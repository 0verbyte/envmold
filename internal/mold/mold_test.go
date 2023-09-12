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

	if _, err := New(strings.NewReader(moldTemplate), nil); err != nil {
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

	mold, err := New(strings.NewReader(moldTemplate), nil)
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

	if _, err := New(strings.NewReader(moldTemplate), nil); err != ErrInvalidDataType {
		t.Errorf("Failed to create new mold: %v", err)
		return
	}
}

func TestMoldPromptReader(t *testing.T) {
	var moldTemplate = `
- name: foo
  value: bar
  type: string
  required: true

- name: debug
  value: true
  type: boolean
  required: false
`

	m, err := New(strings.NewReader(moldTemplate), nil)
	if err != nil {
		t.Errorf("Failed to create new mold: %v", err)
		return
	}

	var promptText = `yes
data`

	m.SetPromptReader(strings.NewReader(promptText))

	if err := m.Generate(); err != nil {
		t.Errorf("Failed to generate mold: %v", err)
		return
	}

	specs := []struct {
		input  string
		output string
	}{
		{input: "foo", output: "data"},
	}

	for _, spec := range specs {
		v, err := m.GetVariable(spec.input)
		if err != nil {
			t.Errorf("Failed to get mold get %s: %v", spec.input, err)
			return
		}

		if v.Value != spec.output {
			t.Errorf("Expected: %s, got %s", spec.output, v.Value)
			return
		}
	}
}

func TestHasTag(t *testing.T) {
	var moldTemplate = `
- name: foo
  value: "bar"
  type: string
  required: true
  tags: ["local", "dev"]

- name: debug
  value: true
  type: boolean
  required: false
`

	m, err := New(strings.NewReader(moldTemplate), nil)
	if err != nil {
		t.Errorf("Failed to create new mold: %v", err)
		return
	}

	specs := []struct {
		envVar string
		input  string
		output bool
	}{
		{envVar: "foo", input: "local", output: true},
		{envVar: "foo", input: "production", output: false},
	}

	for _, spec := range specs {
		envVar, err := m.GetVariable(spec.envVar)
		if err != nil {
			t.Errorf("Got error: %v", err)
			return
		}
		got := envVar.HasTag(spec.input)
		if got != spec.output {
			t.Errorf("Failed tag lookup for %s. Expected %t, got %t", spec.input, spec.output, got)
		}
	}
}

func TestTagsAll(t *testing.T) {
	var moldTemplate = `
- name: foo
  value: "bar"
  type: string
  required: true
  tags: ["local", "dev"]

- name: debug
  value: true
  type: boolean
  required: false
`

	m, err := New(strings.NewReader(moldTemplate), nil)
	if err != nil {
		t.Errorf("Failed to create new mold: %v", err)
		return
	}

	specs := []struct {
		envVar string
		output []string
	}{
		{envVar: "foo", output: []string{"local", "dev"}},
		{envVar: "debug", output: []string{}},
	}

	contains := func(items []string, key string) bool {
		for _, item := range items {
			if key == item {
				return true
			}
		}
		return false
	}

	for _, spec := range specs {
		envVar, err := m.GetVariable(spec.envVar)
		if err != nil {
			t.Errorf("Got error: %v", err)
			return
		}

		gotTags := envVar.AllTags()
		for _, expected := range spec.output {
			if !contains(gotTags, expected) {
				t.Errorf("Expected %s in tags, got tags %s", expected, strings.Join(gotTags, ", "))
			}
		}
	}
}
