package template

type Template struct {
	ID          int64
	Type        string
	Name        string
	TitleFormat string
	BodyFormat  string
	IsDefault   bool
}
