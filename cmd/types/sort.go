package types

type Sort []struct {
	Default bool   `json:"default"`
	Label   string `json:"label"`
	Value   string `json:"value"`
	Icon    Icon   `json:"icon"`
}
