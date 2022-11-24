package middlewares

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"go.dbio.me/cmd/configuration"
	"go.dbio.me/cmd/types"
)

func UploadImage(name string, base64Image string) (*types.ImgbbResponse, error) {
	config := configuration.GetConfig()
	resp, err := http.PostForm("https://api.imgbb.com/1/upload?key="+config.ImageUploadKey+"&name="+name,
		url.Values(map[string][]string{
			"image": {ParseBase64(base64Image)},
		}),
	)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	imgbb_response := types.ImgbbResponse{}

	err = json.Unmarshal(body, &imgbb_response)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return &imgbb_response, nil
}

func ParseBase64(base64 string) string {
	var base64Image = strings.Split(base64, ",")

	if len(base64Image) == 2 {
		return base64Image[1]
	} else {
		return ""
	}
}
