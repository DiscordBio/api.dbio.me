package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Entity struct {
	ID      string `json:"id"`
	Discord struct {
		Id            string `json:"id"`
		Username      string `json:"username"`
		Discriminator string `json:"discriminator"`
	} `json:"discord"`
	URL        string      `json:"url"`
	Banner     string      `json:"banner"`
	Avatar     string      `json:"avatar"`
	About      interface{} `json:"about"`
	Occupation interface{} `json:"occupation"`
	Birthday   interface{} `json:"birthday"`
	Location   interface{} `json:"location"`
	Gender     interface{} `json:"gender"`
	Pronouns   interface{} `json:"pronouns"`
	Language   string      `json:"language"`
	Website    interface{} `json:"website"`
	Like       int         `json:"like"`
	Email      interface{} `json:"email"`
	Views      []string    `json:"views"`
	Premium    bool        `json:"isPremium" bson:"isPremium"`
	Verified   bool        `json:"isVerified" bson:"isVerified"`
	Privacy    struct {
		IsShow            bool `json:"isShow"`
		IsEmailPrivate    bool `json:"isEmailPrivate"`
		IsBirthdayPrivate bool `json:"isBirthdayPrivate"`
		IsLocationPrivate bool `json:"isLocationPrivate"`
		IsGenderPrivate   bool `json:"isGenderPrivate"`
		IsPronounsPrivate bool `json:"isPronounsPrivate"`
	} `json:"privacy"`
	Roles   []string `json:"roles"`
	Likes   []string `json:"likes"`
	Skills  []string `json:"skills"`
	Socials []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		URL  string `json:"url"`
		Icon struct {
			Label string `json:"label"`
			Value string `json:"value"`
			URL   string `json:"url"`
		} `json:"icon"`
	} `json:"socials"`
	CreatedAt    primitive.DateTime `json:"createdAt"`
	UpdatedAt    primitive.DateTime `json:"updatedAt"`
	DeletedAt    interface{}        `json:"deletedAt"`
	IsLiked      bool               `json:"isLiked"`
	IsSelf       bool               `json:"isSelf"`
	IsTeamMember bool               `json:"isTeamMember"`
}
