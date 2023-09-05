package mold

type Writer interface {
	Write(map[string]MoldTemplateVariable) error
}
