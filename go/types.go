package naturalvoid

// DB Models
type StoryData struct {
	Name        string
	ShortName   string
	Slug        string
	Description []string
}

type IndexData struct {
	Stories []StoryData
}
