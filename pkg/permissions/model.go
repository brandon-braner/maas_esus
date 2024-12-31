package permissions

import (
	"reflect"
	"strings"
)

type Permissions struct {
	/*
		Permissions is a struct that contains all the permissions for a user.
		Right now every permission is a boolean.
	*/
	GenerateLlmMeme bool `bson:"generate_llm_meme"`
}

func ValidatePermission(permission string) bool {
	/*
		ValidatePermission takes a permission name and a permission struct
		Returns true if the permission exists in the struct as a bson value.
	*/
	val := reflect.ValueOf(Permissions{})
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		bsonTag := field.Tag.Get("bson")

		// Split the bson tag to get the actual name
		if parts := strings.Split(bsonTag, ","); len(parts) > 0 {
			if parts[0] == permission {
				return true
			}
		}
	}
	return false
}
