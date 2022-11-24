package configuration

import (
	"go.dbio.me/cmd/types"
)

func getConfig() types.Config {
	return types.Config{
		ApiVersion: 1,
		Database: types.Database{
			Url: "localhost:5432",
		},
		Web: types.Web{
			Port:           "8080",
			ImageUploadKey: "1234567890",
			ReturnUrl:      "https://dbio.me",
		},
		Client: types.Client{
			Id:       "1234567890",
			Secret:   "1234567890",
			Token:    "",
			Callback: "https://dbio.me/api/v1/auth/callback",
		},
		Collection: "entities",
		APIUrl:     "https://dbio.me/api/v1",
	}
}

func GetConfig() types.Config {
	return getConfig()
}
