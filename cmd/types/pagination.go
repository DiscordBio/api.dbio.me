package types

type Pagination struct {
	Page     string `json:"page"`
	Limit    string `json:"limit"`
	Sort     string `json:"sort"`
	Roles    string `json:"roles"`
	Skills   string `json:"skills"`
	Query    string `json:"query"`
	Language string `json:"language"`
}
