package types

type ImgbbImage struct {
	Filename  string `json:"filename"`
	Name      string `json:"name"`
	Mime      string `json:"mime"`
	Extension string `json:"extension"`
	Url       string `json:"url"`
}

type ImgbbResponse struct {
	Data    ImgbbResponseData `json:"data"`
	Success bool              `json:"success"`
	Status  int               `json:"status"`
}

type ImgbbResponseData struct {
	ID string `json:"id"`

	DisplayURL string `json:"display_url"`
	DeleteURL  string `json:"delete_url"`

	Expiration string `json:"expiration"`

	Height string `json:"height"`
	Width  string `json:"width"`

	Image  ImgbbImage `json:"image"`
	Thumb  ImgbbImage `json:"thumb"`
	Medium ImgbbImage `json:"medium"`
}
