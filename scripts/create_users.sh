#!/bin/bash

NONLLM_EMAIL="nonllm@example.com"
LLM_EMAIL="llm@example.com"

# Drop all users
go run cmd/cli/user/main.go delete-all

# Create user nonllm
go run cmd/cli/user/main.go create --username $NONLLM_EMAIL --password password

# Create user llm
go run cmd/cli/user/main.go create --username $LLM_EMAIL --password password

# Set llm generate_llm_meme to true
go run cmd/cli/user/main.go permissions set --username $LLM_EMAIL --permission generate_llm_meme --value true

# Add Tokens to both users. Adding less to nonllm because I assume their requests will be cheaper as they aren't using llm
go run cmd/cli/user/main.go add-tokens --username $LLM_EMAIL --amount 100
go run cmd/cli/user/main.go add-tokens --username $NONLLM_EMAIL --amount 50
