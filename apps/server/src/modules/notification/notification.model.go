package notification

type Model struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	Active    bool    `json:"active"`
	IsDefault bool    `json:"is_default"`
	Config    *string `json:"config"`
}

type UpdateModel struct {
	ID        *string `json:"id"`
	Name      *string `json:"name"`
	Type      *string `json:"type"`
	Active    *bool   `json:"active"`
	IsDefault *bool   `json:"is_default"`
	Config    *string `json:"config"`
}
