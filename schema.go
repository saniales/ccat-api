package ccatapi

// schemaProperty contains the data about a single generic setting schema property.
type schemaProperty struct {
	Title   string `json:"title"`
	Type    string `json:"type"`
	Default *any   `json:"default,omitempty"`
}

// settingSchema contains the data about a single generic setting schema.
type settingSchema struct {
	Description       string                    `json:"description"`
	HumanReadableName string                    `json:"humanReadableName"`
	Link              string                    `json:"link"`
	Properties        map[string]schemaProperty `json:"properties"`
	Required          []string                  `json:"required"`
	Title             string                    `json:"title"`
	Type              string                    `json:"type"`
}
