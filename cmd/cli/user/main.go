package main

import (
	"fmt"
	"os"

	"github.com/brandonbraner/maas/external/usersapi"
	"github.com/brandonbraner/maas/pkg/permissions"
	"github.com/spf13/cobra"
)

func main() {
	// Initialize user service
	userService, err := usersapi.NewUserService()
	if err != nil {
		fmt.Printf("Failed to initialize user service: %v\n", err)
		os.Exit(1)
	}

	var rootCmd = &cobra.Command{
		Use:   "user",
		Short: "User management CLI",
	}

	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a new user",
		Run: func(cmd *cobra.Command, args []string) {
			username, _ := cmd.Flags().GetString("username")
			password, _ := cmd.Flags().GetString("password")
			firstname, _ := cmd.Flags().GetString("firstname")
			lastname, _ := cmd.Flags().GetString("lastname")

			// Create user with default permissions
			user, err := userService.NewUser(username, password, firstname, lastname, 0, permissions.Permissions{})
			if err != nil {
				fmt.Printf("Failed to create user: %v\n\n", err)
				return
			}

			// Save user to database
			createdUser, err := userService.CreateUser(user)
			if err != nil {
				fmt.Printf("Failed to save user: %v\n\n", err)
				return
			}

			fmt.Printf("Successfully created user: %+v\n\n", createdUser.Username)
		},
	}

	createCmd.Flags().StringP("username", "", "", "User's email address (required)")
	createCmd.Flags().StringP("password", "", "", "User's password (required)")
	createCmd.Flags().StringP("firstname", "", "", "User's first name")
	createCmd.Flags().StringP("lastname", "", "", "User's last name")
	createCmd.MarkFlagRequired("username")
	createCmd.MarkFlagRequired("password")

	var permissionsCmd = &cobra.Command{
		Use:   "permissions",
		Short: "Manage user permissions",
	}

	var setPermissionsCmd = &cobra.Command{
		Use:   "set",
		Short: "Set user permissions",
		Run: func(cmd *cobra.Command, args []string) {
			username, _ := cmd.Flags().GetString("username")
			permission, _ := cmd.Flags().GetString("permission")
			value, _ := cmd.Flags().GetBool("value")

			err := userService.UpdatePermission(username, permission, value)
			if err != nil {
				fmt.Printf("Failed to update permissions: %v\n", err)
				return
			}

			fmt.Printf("Successfully updated permissions for user %s\n", username)
		},
	}

	setPermissionsCmd.Flags().StringP("username", "", "", "User's username (required)")
	setPermissionsCmd.Flags().StringP("permission", "", "", "Name of Permission")
	setPermissionsCmd.Flags().BoolP("value", "", false, "Set to true to set the permission to true. Omit to set the permission to false.")
	_ = setPermissionsCmd.MarkFlagRequired("username")
	_ = setPermissionsCmd.MarkFlagRequired("permission")

	permissionsCmd.AddCommand(setPermissionsCmd)

	var generateTokenCmd = &cobra.Command{
		Use:   "generate-token",
		Short: "Generate a JWT token for a user",
		Run: func(cmd *cobra.Command, args []string) {
			username, _ := cmd.Flags().GetString("username")

			token, err := userService.GenerateJwt(username)
			if err != nil {
				fmt.Printf("Failed to generate token: %v\n", err)
				return
			}

			fmt.Printf("Successfully generated token for \x1b[31m user: %s\n\n \x1b[0m Token: %s\n\n", username, token)
		},
	}

	generateTokenCmd.Flags().StringP("username", "", "", "User's username (required)")
	_ = generateTokenCmd.MarkFlagRequired("username")

	var deleteAllCmd = &cobra.Command{
		Use:   "delete-all",
		Short: "Delete all users from the database",
		Run: func(cmd *cobra.Command, args []string) {
			count, err := userService.DeleteAllUsers()
			if err != nil {
				fmt.Printf("Failed to delete users: %v\n", err)
				return
			}

			fmt.Printf("Successfully deleted %d users\n", count)
		},
	}

	var addTokensCmd = &cobra.Command{
		Use:   "add-tokens",
		Short: "Add tokens to a user's account",
		Run: func(cmd *cobra.Command, args []string) {
			username, _ := cmd.Flags().GetString("username")
			amount, _ := cmd.Flags().GetInt("amount")

			err := userService.AddTokens(username, amount)
			if err != nil {
				fmt.Printf("Failed to add tokens: %v\n", err)
				return
			}

			fmt.Printf("Successfully added %d tokens to user %s\n", amount, username)
		},
	}

	addTokensCmd.Flags().StringP("username", "", "", "User's username (required)")
	addTokensCmd.Flags().IntP("amount", "", 0, "Number of tokens to add (required)")
	_ = addTokensCmd.MarkFlagRequired("username")
	_ = addTokensCmd.MarkFlagRequired("amount")

	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(permissionsCmd)
	rootCmd.AddCommand(generateTokenCmd)
	rootCmd.AddCommand(deleteAllCmd)
	rootCmd.AddCommand(addTokensCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
