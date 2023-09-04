package mold

import (
	"errors"
	"io"

	"gopkg.in/yaml.v3"
)

var (
	ErrMissingVariableName             = errors.New("missing environment variable name")
	ErrEnvironmentVariableDoesNotExist = errors.New("environment variable does not exist")
	ErrInvalidDataType                 = errors.New("value does not implement the required type")
)

const (
	MoldDataTypeString  = "string"
	MoldDataTypeNumber  = "number"
	MoldDataTypeBoolean = "boolean"
)

// MoldTemplateVariable data representation for a mold template variable
type MoldTemplateVariable struct {
	Name     string      `yaml:"name"`
	Value    interface{} `yaml:"value"`
	Type     string      `yaml:"type"`
	Required bool        `yaml:"required"`
}

// MoldTemplate data representation for the MoldTemplate
type MoldTemplate struct {
	variables map[string]MoldTemplateVariable
}

// New creates a new MoldTemplate from an io.Reader. Use the helper functions to read from the respective input.
func New(r io.Reader) (*MoldTemplate, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	moldTemplate := []MoldTemplateVariable{}
	if err := yaml.Unmarshal(b, &moldTemplate); err != nil {
		return nil, err
	}

	moldTemplateVariables := make(map[string]MoldTemplateVariable)
	for _, moldTemplateVariable := range moldTemplate {
		if moldTemplateVariable.Name == "" {
			return nil, ErrMissingVariableName
		}

		// Check for type constraint on the value field
		switch moldTemplateVariable.Value.(type) {
		case int, float32, float64:
			if moldTemplateVariable.Type != MoldDataTypeNumber {
				return nil, ErrInvalidDataType
			}
		case string:
			if moldTemplateVariable.Type != MoldDataTypeString {
				return nil, ErrInvalidDataType
			}
		case bool:
			if moldTemplateVariable.Type != MoldDataTypeBoolean {
				return nil, ErrInvalidDataType
			}
		}

		moldTemplateVariables[moldTemplateVariable.Name] = moldTemplateVariable
	}

	return &MoldTemplate{
		variables: moldTemplateVariables,
	}, nil
}

// GetVariable gets a MoldTemplateVariable by key. If the key does not exist an error will be returned.
func (m *MoldTemplate) GetVariable(key string) (*MoldTemplateVariable, error) {
	if v, ok := m.variables[key]; ok {
		return &v, nil
	}
	return nil, ErrEnvironmentVariableDoesNotExist
}
