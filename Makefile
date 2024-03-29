APPNAME := any-method-api
APPIMAGE := asia-northeast1-docker.pkg.dev/$(GOOGLE_CLOUD_PROJECT)/my-app/$(APPNAME)
LOADINGIMAGE := asia-northeast1-docker.pkg.dev/$(GOOGLE_CLOUD_PROJECT)/my-app/loading-client

.PHONY: build-all
build-all: build-app build-loading-client

.PHONY: build-app
build-app:
	sh build.sh $(APPIMAGE) > /dev/null
	@echo "Image had been built as [$(APPIMAGE)]"

.PHONY: build-loading-client
build-loading-client:
	( cd loading-client; sh ./build.sh $(LOADINGIMAGE) > /dev/null )
	@echo "Image had been built as [$(LOADINGIMAGE)]"

.PHONY: prepare-repo
prepare-repo:
	gcloud artifacts repositories create my-app --location=asia-northeast1 --repository-format=docker

.PHONY: deploy-app
deploy-app:
	APPNAME=$(APPNAME) APPIMAGE=$(APPIMAGE) envsubst < manifests.yaml | kubectl apply -f -

.PHONY: deploy-loading-client
deploy-loading-client:
	IP=$(IP) LOADINGIMAGE=$(LOADINGIMAGE) envsubst < ./loading-client/manifests.yaml | kubectl apply -f -

.PHONY: deploy-all
deploy-all: deploy-app deploy-loading
