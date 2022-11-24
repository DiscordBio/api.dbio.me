package routes

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"go.dbio.me/cmd/middlewares"
	"go.dbio.me/cmd/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/oauth2"
)

var state = "random"

func Login(c *fiber.Ctx) error {
	AuthConfig := c.Locals("authConfig").(*oauth2.Config)
	sess := c.Locals("session").(*session.Session)

	var next string
	if c.Query("next") == "" {
		next = "/"
	} else {
		next = c.Query("next")
	}

	if err := sess.Save(); err != nil {
		return c.Redirect(AuthConfig.AuthCodeURL(state + "&next=" + next))
	}

	return c.Redirect(AuthConfig.AuthCodeURL(state + "&next=" + next))
}

func Callback(c *fiber.Ctx) error {
	AuthConfig := c.Locals("authConfig").(*oauth2.Config)
	config := c.Locals("config").(types.Config)
	sess := c.Locals("session").(*session.Session)

	var queryState string
	var next string

	total := strings.Split(c.Query("state"), "&")

	if len(total) == 2 {
		queryState = total[0]
		next = strings.Split(total[1], "=")[1]
	} else {
		queryState = total[0]
		next = "/"
	}

	if next == "" {
		next = "/"
	}

	if queryState != state {
		return c.Redirect(AuthConfig.AuthCodeURL(state + "&next=" + next))
	}

	token, err := AuthConfig.Exchange(c.Context(), c.Query("code"))

	if err != nil {
		return c.Redirect(AuthConfig.AuthCodeURL(state + "&next=" + next))
	}

	res, err := AuthConfig.Client(context.Background(), token).Get("https://discord.com/api/users/@me")

	if err != nil {
		return c.Redirect(AuthConfig.AuthCodeURL(state + "&next=" + next))
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return c.Redirect(AuthConfig.AuthCodeURL(state + "&next=" + next))
	}

	var user types.User

	err = json.Unmarshal(body, &user)

	if err != nil {
		return c.Redirect(AuthConfig.AuthCodeURL(state + "&next=" + next))
	}

	var db *mongo.Client = c.Locals("db").(*mongo.Client)
	users := db.Database("dbio").Collection("users")

	var result types.User

	users.FindOne(context.Background(), bson.M{"discord.id": user.ID}).Decode(&result)

	tokenSplit := strings.Split(middlewares.GenerateJWT(c, user.ID), " ")

	user.AccessToken = tokenSplit[1]
	user.Token = tokenSplit[0]

	opts := options.Update().SetUpsert(true)

	users.UpdateOne(context.Background(), bson.M{"id": user.ID}, bson.M{"$set": bson.M{
		"id":                user.ID,
		"username":          user.Username,
		"discriminator":     user.Discriminator,
		"avatar":            user.Avatar,
		"banner":            user.Banner,
		"accent_color":      user.AccentColor,
		"banner_color":      user.BannerColor,
		"locale":            user.Locale,
		"mfa_enabled":       user.MfaEnabled,
		"flags":             user.Flags,
		"premium_type":      user.PremiumType,
		"public_flags":      user.PublicFlags,
		"avatar_decoration": user.AvatarDecoration,
		"token":             user.Token,
		"access_token":      user.AccessToken,
		"loggedAt":          time.Now(),
	}}, opts)

	entities := db.Database("dbio").Collection("entities")
	var entity types.Entity

	entities.FindOne(context.Background(), bson.M{"discord.id": user.ID}).Decode(&entity)

	if entity.ID != "" {

		entities.UpdateOne(context.Background(), bson.M{"discord.id": user.ID}, bson.M{"$set": bson.M{
			"discord.id":            user.ID,
			"discord.username":      user.Username,
			"discord.discriminator": user.Discriminator,
			"discord.avatar":        user.Avatar,
		}})

	}

	gob.Register(types.User{})
	sess.Set("user", result)
	if err := sess.Save(); err != nil {
		fmt.Println(err.Error())
		return c.Redirect(AuthConfig.AuthCodeURL(state))
	}

	jwtAccessToken := user.AccessToken

	return c.Redirect(config.Web.ReturnUrl + "?next=" + next + "&access_token=" + jwtAccessToken + "&action=login")
}

func Logout(c *fiber.Ctx) error {
	sess := c.Locals("session").(*session.Session)
	config := c.Locals("config").(types.Config)
	sess.Destroy()

	back := c.Query("back")

	if back == "" {
		back = "/"
	}

	return c.Redirect(config.Web.ReturnUrl + "?next=" + back + "&action=logout")
}

func GetCurrentUser(c *fiber.Ctx) error {
	user := middlewares.GetUser(c)
	if user == nil {
		return c.Status(401).JSON(middlewares.NewError("You are not authorized to access this resource"))
	}

	return c.JSON(middlewares.NewSuccess("Successfully fetched user", user))
}
