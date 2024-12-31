.PHONY: create_users generate_jwt run_api

create_users:
	@bash scripts/create_users.sh

generate_jwt:
	go run cmd/cli/user/main.go generate-token --username llm@example.com
	go run cmd/cli/user/main.go generate-token --username nonllm@example.com



run_api:
	go run cmd/mass/main.go