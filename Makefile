APP_NAME = backend
DOCKER_IMAGE_NAME = learn-devops-$(APP_NAME)
DOCKERHUB_USERNAME = <your-dockerhub-username>

.PHONY: run build-app build-docker deploy-dockerhub

run:
	go run main.go

build-app:
	go build -o $(APP_NAME) main.go

build-docker:
	docker build -t $(DOCKER_IMAGE_NAME) .

deploy-dockerhub:
	# Replace <your-dockerhub-username> with your actual Docker Hub username
	docker tag $(DOCKER_IMAGE_NAME) $(DOCKERHUB_USERNAME)/$(DOCKER_IMAGE_NAME):latest
	docker push $(DOCKERHUB_USERNAME)/$(DOCKER_IMAGE_NAME):latest
