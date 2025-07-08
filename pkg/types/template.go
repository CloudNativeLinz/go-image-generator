package types

// Position represents a position with X and Y coordinates
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// TextElement represents a text element with font, size, color, position and box width
type TextElement struct {
	Text     string   `json:"text"`
	Font     string   `json:"font"`
	FontSize float64  `json:"fontSize"`
	Color    string   `json:"color"`
	Position Position `json:"position"`
	BoxWidth float64  `json:"boxWidth"`
}

// BackgroundConfig represents background image configuration
type BackgroundConfig struct {
	Image    string `json:"image"`
	Position struct {
		X int `json:"x"`
		Y int `json:"y"`
	} `json:"position"`
	Size struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"size"`
}

// Template represents the complete template configuration
type Template struct {
	Background    BackgroundConfig `json:"background"`
	Speaker1title TextElement      `json:"speaker1title"`
	Speaker1name  TextElement      `json:"speaker1name"`
	Speaker2title TextElement      `json:"speaker2title"`
	Speaker2name  TextElement      `json:"speaker2name"`
	Sponsor       TextElement      `json:"sponsor"`
	Date          TextElement      `json:"date"`
}
