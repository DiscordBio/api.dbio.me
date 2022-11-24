package middlewares

import (
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"go.dbio.me/cmd/configuration"
	"go.dbio.me/cmd/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var jwtKey = []byte("my_secret_key")

func GenerateJWT(c *fiber.Ctx, userID string) string {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, _ := token.SignedString(jwtKey)

	accessToken := jwt.New(jwt.SigningMethodHS256)

	claims = accessToken.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	accessTokenString, _ := accessToken.SignedString(jwtKey)

	return tokenString + " " + accessTokenString
}

func ValidateJWT(c *fiber.Ctx) interface{} {
	token := c.Get("Authorization")
	if token == "" {
		return false
	}
	token = parseJWT(c, token).(string)

	_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, nil
		}
		return jwtKey, nil
	})

	if err != nil {
		return false
	}

	db := c.Locals("db").(*mongo.Client)
	users := db.Database("dbio").Collection("users")

	var result types.User

	err = users.FindOne(context.Background(), bson.M{"access_token": token}).Decode(&result)

	if err != nil {
		return false
	}

	return token
}

func parseJWT(c *fiber.Ctx, string string) interface{} {
	parsedJWT := strings.Split(string, " ")

	if len(parsedJWT) == 2 {
		if parsedJWT[0] == "Bearer" {
			return parsedJWT[1]
		} else {
			return c.JSON(NewError("You are not authorized to access this resource"))
		}
	} else {
		return c.JSON(NewError("You are not authorized to access this resource"))
	}
}

func AuthRequired(c *fiber.Ctx) error {
	token := ValidateJWT(c)

	if token == false {
		return c.Status(401).JSON(NewError("You are not authorized to access this resource"))
	}

	return c.Next()
}

func GetUser(c *fiber.Ctx) interface{} {
	DConfig := configuration.GetConfig()

	token := ValidateJWT(c)

	if token == false {
		return c.Status(401).JSON(NewError("You are not authorized to access this resource"))
	}

	db := c.Locals("db").(*mongo.Client)
	users := db.Database("dbio").Collection("users")

	var result types.User

	err := users.FindOne(context.Background(), bson.M{"access_token": token}).Decode(&result)

	if err != nil {
		return c.Status(401).JSON(NewError("You are not authorized to access this resource"))
	}

	result.AccessToken = "FAILED TO GET TOKEN"
	result.Token = "FAILED TO GET TOKEN"

	result.AppID = nil
	entities := db.Database("dbio").Collection("entities")

	var entitiesResult types.Entity
	entities.FindOne(context.Background(), bson.M{"discord.id": result.ID}).Decode(&entitiesResult)

	if entitiesResult.URL != "" {
		result.AppID = entitiesResult.URL
		data := fiber.Map{
			"url":        entitiesResult.URL,
			"id":         entitiesResult.ID,
			"about":      entitiesResult.About,
			"socials":    entitiesResult.Socials,
			"roles":      entitiesResult.Roles,
			"occupation": entitiesResult.Occupation,
			"skills":     entitiesResult.Skills,
			"likes":      Count(entitiesResult.Likes),
			"createdAt":  entitiesResult.CreatedAt,
			"updatedAt":  entitiesResult.UpdatedAt,
			"isLiked":    entitiesResult.IsLiked,
			"isSelf":     entitiesResult.IsSelf,
			"views":      Count(entitiesResult.Views),
			"isPremium":  entitiesResult.Premium,
			"isVerified": entitiesResult.Verified,
			"privacy":    entitiesResult.Privacy,
			"email":      entitiesResult.Email,
			"gender":     entitiesResult.Gender,
			"pronouns":   entitiesResult.Pronouns,
			"birthday":   entitiesResult.Birthday,
			"location":   entitiesResult.Location,
			"language":   entitiesResult.Language,
		}

		if entitiesResult.Banner != "" {
			data["banner"] = DConfig.APIUrl + "/banners/" + entitiesResult.URL
		}

		if entitiesResult.Avatar != "" {
			data["avatar"] = DConfig.APIUrl + "/avatars/" + entitiesResult.URL
		}

		result.Entity = data
	}

	return result
}

func GetUserID(c *fiber.Ctx) string {
	token := ValidateJWT(c)

	if token == false {
		return ""
	}

	db := c.Locals("db").(*mongo.Client)
	users := db.Database("dbio").Collection("users")

	var result types.User

	err := users.FindOne(context.Background(), bson.M{"access_token": token}).Decode(&result)

	if err != nil {
		return ""
	}

	return result.ID
}

func GetDetailedUser(c *fiber.Ctx) types.User {
	token := ValidateJWT(c)

	if token == false {
		return types.User{}
	}

	db := c.Locals("db").(*mongo.Client)
	users := db.Database("dbio").Collection("users")

	var result types.User

	err := users.FindOne(context.Background(), bson.M{"access_token": token}).Decode(&result)

	if err != nil {
		return types.User{}
	}

	return result
}
