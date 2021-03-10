docker-build:
	./mvnw spring-boot:build-image -Dspring-boot.build-image.imageName=localregistry:5000/micropet/pets-spring -Dmaven.test.skip=true

docker-push: docker-build
	docker push  localregistry:5000/micropet/pets-spring:latest