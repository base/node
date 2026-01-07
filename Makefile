.PHONY: doctor validate-config help

help:
	@echo "Available targets:"
	@echo "  doctor          - Validate configuration (alias for validate-config)"
	@echo "  validate-config - Validate environment configuration file"
	@echo ""
	@echo "Usage:"
	@echo "  make doctor [ENV_FILE=.env.mainnet]"
	@echo "  make validate-config [ENV_FILE=.env.mainnet]"

doctor: validate-config

validate-config:
	@if [ -z "$(ENV_FILE)" ]; then \
		if [ -f ".env.mainnet" ]; then \
			./scripts/validate-config.sh .env.mainnet; \
		elif [ -f ".env.sepolia" ]; then \
			./scripts/validate-config.sh .env.sepolia; \
		else \
			echo "Error: No .env file found. Please specify ENV_FILE=path/to/.env"; \
			exit 1; \
		fi \
	else \
		./scripts/validate-config.sh $(ENV_FILE); \
	fi

