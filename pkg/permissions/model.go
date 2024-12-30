package permissions

type Permissions struct {
	/*
		Permissions is a struct that contains all the permissions for a user.
		Right now every permission is a boolean.
	*/
	GenerateLlmMeme bool `bson:"generate_llm_meme"`
}
