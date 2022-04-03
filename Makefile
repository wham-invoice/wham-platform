TAG := $(shell git rev-parse --verify --short HEAD)
# Artifact registry
REGISTRY := us-west2-docker.pkg.dev
REPOSITORY := platform


build-push-image:
	docker build -t $(REGISTRY)/$(PROJECT_ID)/$(REPOSITORY)/$(TAG) .
	docker push $(REGISTRY)/$(PROJECT_ID)/$(REPOSITORY)/$(TAG)

deploy-image:
	kustomize edit set image REGISTRY/PROJECT_ID/REPO/IMAGE=$(REGISTRY)/$(PROJECT_ID)/$(REPOSITORY)/$(TAG)
	kustomize build . | kubectl apply -f -
	kubectl rollout status deployment/$DEPLOYMENT_NAME
	kubectl get services -o wide