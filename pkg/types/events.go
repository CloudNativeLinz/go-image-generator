package types

// Talk represents a talk with title and speaker
type Talk struct {
	Title   string `yaml:"title"`
	Speaker string `yaml:"speaker"`
}

// Event represents an event with talks, host, and date
type Event struct {
	ID    int    `yaml:"id"`
	Date  string `yaml:"date"`
	Title string `yaml:"title"`
	Talks []Talk `yaml:"talks"`
	Host  string `yaml:"host"`
}

// EventsYAML represents a list of events
type EventsYAML []Event

// EventData represents extracted event information for text rendering
type EventData struct {
	Speaker1Title string
	Speaker1Name  string
	Speaker2Title string
	Speaker2Name  string
	Sponsor       string
	Date          string
}
