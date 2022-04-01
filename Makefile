TAG := $(shell git rev-parse --verify --short HEAD)
# Artifact registry
REGISTRY := us-west2-docker.pkg.dev
PROJECT_ID := wham-ad61b
REPOSITORY := platform


build-push-image:
	docker build -t $(REGISTRY)/$(PROJECT_ID)/$(REPOSITORY)/image:$(TAG) .
	docker push $(REGISTRY)/$(PROJECT_ID)/$(REPOSITORY)/image:$(TAG)