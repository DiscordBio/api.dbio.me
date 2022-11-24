package types

type Social []struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	URL     string `json:"url"`
	Color   string `json:"color"`
	Icon    Icon   `json:"icon"`
	Enabled bool   `json:"enabled"`
}

type Icon struct {
	Label string `json:"label"`
	Value string `json:"value"`
	URL   string `json:"url"`
}
