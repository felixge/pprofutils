.PHONY: README.md
README.md:
	go run ./scripts/generate_readme.go < README.template.md > README.md

.PHONY: deploy
deploy:
	go run ./scripts/deploy_pprofutils.bash

.PHONY: deploy-agent
deploy-agent:
	go run ./scripts/deploy_agent.bash
