package mold

type Writer interface {
	Write([]MoldTemplateVariable) error
}
