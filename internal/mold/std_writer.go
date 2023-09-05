package mold

import (
	"fmt"
	"os"
	"strings"
)

const (
	WriterStdout = "stdout"
)

type StdoutWriter struct{}

func (w *StdoutWriter) Write(envVars map[string]MoldTemplateVariable) error {
	lines := []string{}
	for _, v := range envVars {
		valueFmt := "%s"
		switch v := v.Value.(type) {
		case string:
			if strings.Contains(v, " ") {
				valueFmt = "\"%v\""
			}
		case bool:
			valueFmt = "%t"
		case int:
			valueFmt = "%d"
		case float64:
			valueFmt = "%.2f"
		}

		line := fmt.Sprintf("export %s=%s", strings.ToUpper(v.Name), fmt.Sprintf(valueFmt, v.Value))
		lines = append(lines, line)
	}

	if _, err := fmt.Fprint(os.Stdout, strings.Join(lines, "\n")); err != nil {
		return err
	}
	return nil
}
