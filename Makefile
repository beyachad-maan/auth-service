build:
	go build .

image:
	docker build -f ./Dockerfile . -t auth-service:latest

image-push:
	docker login -u kubeadmin -p $$(oc whoami -t) default-route-openshift-image-registry.apps-crc.testing
	docker tag auth-service:latest   default-route-openshift-image-registry.apps-crc.testing/beyachad-maan/auth-service:latest
	docker push default-route-openshift-image-registry.apps-crc.testing/beyachad-maan/auth-service:latest

deploy:
	kubectl apply -f ./deployment.yaml

clean-deploy:
	kubectl delete -f ./deployment.yaml