# Variáveis para facilitar a configuração
IMAGE_NAME = go-postgres-app
PORT = 8085

build:
	docker build -t $(IMAGE_NAME) .

run:
	docker run -p $(PORT):$(PORT) $(IMAGE_NAME)

stop:
	docker stop $(shell docker ps -q --filter ancestor=$(IMAGE_NAME))

rm:
	docker rm $(shell docker ps -aq --filter ancestor=$(IMAGE_NAME))

push:
	docker tag go-postgres-app pedrosilva1/go-project
	docker push pedrosilva1/go-project