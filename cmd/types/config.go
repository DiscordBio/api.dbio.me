package types

type Config struct {
	ApiVersion int `json:"apiVersion"`
	Database   `json:"database"`
	Web        `json:"web"`
	Client     `json:"client"`
	Collection string `json:"collection"`
	APIUrl     string `json:"apiUrl"`
}

type Database struct {
	Url string `json:"url"`
}

type Web struct {
	Port           string `json:"port"`
	ImageUploadKey string `json:"imageUploadKey"`
	ReturnUrl      string `json:"returnUrl"`
}

type Client struct {
	Id       string `json:"id"`
	Secret   string `json:"secret"`
	Token    string `json:"token"`
	Callback string `json:"callback"`
}

type Roles struct {
	Slug  string `json:"slug"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type Skills struct {
	Slug  string `json:"slug"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}
