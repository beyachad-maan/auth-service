build:
	go build .

image:
	docker build -f ./Dockerfile . -t auth-service:latest

deploy:
	kubectl apply -f ./deployment.yaml

clean-deploy:
	kubectl delete -f ./deployment.yaml