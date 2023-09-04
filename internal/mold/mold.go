package mold

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	ErrMissingVariableName             = errors.New("missing environment variable name")
	ErrEnvironmentVariableDoesNotExist = errors.New("environment variable does not exist")
	ErrInvalidDataType                 = errors.New("value does not implement the required type")
	ErrEmptyMold                       = errors.New("mold variables are empty")
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

func (m *MoldTemplateVariable) String() string {
	return fmt.Sprintf("%s = %v (type=%s, required=%t)", m.Name, m.Value, m.Type, m.Required)
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

// Generate runs the main logic to check for any required fields in the mold template and fills the values
func (m *MoldTemplate) Generate() error {
	if m.variables == nil {
		return ErrEmptyMold
	}

	for k, v := range m.variables {
		if !v.Required {
			continue
		}

		if v.Value != "" && v.Value != nil {
			fmt.Printf("'%s' is a required field, with the value of '%s'. Would you like to overwrite this value (yes/no)? ", v.Name, v.Value)
			reader := bufio.NewReader(os.Stdin)
			answer, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			answer = strings.TrimSpace(answer)
			if answer == "no" || answer == "n" {
				fmt.Println("Skipping...")
				continue
			}
		}

		fmt.Printf("Enter a value for %s (type=%s): ", v.Name, v.Type)
		reader := bufio.NewReader(os.Stdin)
		value, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		value = strings.TrimSpace(value)
		switch v.Type {
		case MoldDataTypeNumber:
			if strings.Contains(value, ".") {
				parsed, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return err
				}
				v.Value = parsed
			} else {
				parsed, err := strconv.Atoi(value)
				if err != nil {
					return err
				}
				v.Value = parsed
			}

		case MoldDataTypeBoolean:
			parsed, err := strconv.ParseBool(value)
			if err != nil {
				return err
			}
			v.Value = parsed

		case MoldDataTypeString:
			fallthrough
		default:
			v.Value = value
		}
		m.variables[k] = v
	}
	return nil
}

// GetVariable gets a MoldTemplateVariable by key. If the key does not exist an error will be returned.
func (m *MoldTemplate) GetVariable(key string) (*MoldTemplateVariable, error) {
	if v, ok := m.variables[key]; ok {
		return &v, nil
	}
	return nil, ErrEnvironmentVariableDoesNotExist
}

// GetAllVariables returns all the variables in the mold
func (m *MoldTemplate) GetAllVariables() []MoldTemplateVariable {
	variables := []MoldTemplateVariable{}
	for _, v := range m.variables {
		variables = append(variables, v)
	}
	return variables
}
