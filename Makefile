IMAGE_NAME := johma_site_reborn
CONTAINER_NAME := johma_site_reborn_container
PROJECT_ID ?= demo-johma-site

.PHONY: build run lrun shutdown restart clean logs help

.DEFAULT_GOAL := help

build:
	@docker build --no-cache -t $(IMAGE_NAME) .

run: build
	@docker rm -f $(CONTAINER_NAME) 2>/dev/null || true
	@docker run -d -p 8080:8080 --name $(CONTAINER_NAME) \
		-e GOOGLE_CLOUD_PROJECT=$(PROJECT_ID) \
		-e GOOGLE_APPLICATION_CREDENTIALS=/app/gcloud/application_default_credentials.json \
		-v "$(GOOGLE_APPLICATION_CREDENTIALS):/app/gcloud/application_default_credentials.json:ro" \
		$(IMAGE_NAME)

lrun: build
	@docker rm -f $(CONTAINER_NAME) 2>/dev/null || true
	@docker run -d -p 8080:8080 --name $(CONTAINER_NAME) \
		-e FIRESTORE_EMULATOR_HOST=host.docker.internal:8081 \
		-e FIREBASE_AUTH_EMULATOR_HOST=host.docker.internal:9099 \
		-e GOOGLE_CLOUD_PROJECT=$(PROJECT_ID) \
		$(IMAGE_NAME)
	@echo "http://localhost:8080"

shutdown:
	@docker rm -f $(CONTAINER_NAME) 2>/dev/null || true

restart: shutdown run
clean: shutdown

logs:
	@docker logs -f $(CONTAINER_NAME)

help:
	@echo "make build / run / lrun / shutdown / restart / logs"
