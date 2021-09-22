.PHONY: generate
generate:
	go generate ./cmd/pprofutils/

.PHONY: README.md
README.md:
	go run ./scripts/generate_readme.go < README.template.md > README.md

.PHONY: deploy
deploy: generate
	./scripts/deploy_pprofutils.bash

.PHONY: deploy-agent
deploy-agent:
	./scripts/deploy_agent.bash
