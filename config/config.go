package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	GOOGLE_GEOCODE_API_KEY  string
	OPENAI_API_KEY          string
	MONGODB_URI             string
	MONGO_DB_NAME           string
	TEXT_MEME_TOKEN_COST    int
	AI_TEXT_MEME_TOKEN_COST int
	REDIS_URL               string
	REDIS_CACHE_TTL         int

	//Mongo Collection Names
	USER_COLLECTION_NAME string

	GET_REAL_GEOLOCATION bool
}

var AppConfig Config

// TODO don't hardcode this, but the jwt libaray yells at me if I try to bring it in from an env var
const JWT_SECRET_KEY = "7df355d32f6ea0ab9d55d877eb2d9317e0de08236c1c1680281e007622221c9c374d358aa41c1b5ee430128fac55d8509ba5a8168a795eb3c853f000bbd721867e3bdb6470e1bf5055efcc2e899ecf4faa5769613c3fd90168d253ead614fd572ae42a4058cb50070a2ffd6ff71254ab6ae3367ffba17e2c1035e285df25ed5b1a2263c333c7b55e7742dba441c2bae60ebc01f24c4192a63170bf8aca913570c4a5d16a89af4056410f2889537eb64de62d597528fc74b367c8bd597bdabbc48971dcb264851a69f04e19429f2c797ff737de84b92ddaf4eb31f883004b0bef67cc6f4210b3ce746507a8110b0cdcfa79be88083cef28318aa9d02fd10a4150"

func init() {
	godotenv.Load()

	AppConfig.GOOGLE_GEOCODE_API_KEY = os.Getenv("GOOGLE_GEOCODE_API_KEY")
	AppConfig.OPENAI_API_KEY = os.Getenv("OPENAI_API_KEY")
	AppConfig.MONGODB_URI = generateMongoUri()
	AppConfig.MONGO_DB_NAME = os.Getenv("MONGO_DB_NAME")
	AppConfig.REDIS_URL = os.Getenv("REDIS_URL")

	AppConfig.GET_REAL_GEOLOCATION = true

	// Default cache TTL to 1 hour if not set
	redisTTL, _ := strconv.Atoi(os.Getenv("REDIS_CACHE_TTL"))
	if redisTTL == 0 {
		redisTTL = 3600
	}
	AppConfig.REDIS_CACHE_TTL = redisTTL

	//TODO ignoring errors for now as I am not sure what type of error to throw here
	textMemeTokenCost, _ := strconv.Atoi(os.Getenv("TEXT_MEME_TOKEN_COST"))
	aiTextMemeTokenCost, _ := strconv.Atoi(os.Getenv("AI_TEXT_MEME_TOKEN_COST"))

	AppConfig.TEXT_MEME_TOKEN_COST = textMemeTokenCost
	AppConfig.AI_TEXT_MEME_TOKEN_COST = aiTextMemeTokenCost

	//Mongo Collection Names
	AppConfig.USER_COLLECTION_NAME = "users"

}

func generateMongoUri() string {
	var uri string

	host := os.Getenv("MONGO_HOST")
	user := os.Getenv("MONGO_USER")
	password := os.Getenv("MONGO_PASSWORD")
	port := os.Getenv("MONGO_PORT")
	db := os.Getenv("MONGO_DB_NAME")

	uri = "mongodb://" + user + ":" + password + "@" + host + ":" + port + "/" + db + "?authSource=admin"

	return uri

}

func SetupMongoTestConfig() {
	// set the appconfig db name to the test db
	AppConfig.MONGO_DB_NAME = "testmaas"
	AppConfig.MONGODB_URI = "mongodb://root:password@localhost:27017"
}
