https://hackmd.io/@ryanjbaxter/spring-on-k8s-workshop

Start the service
./mvnw spring-boot:run

Build the image
./mvnw spring-boot:build-image
./mvnw spring-boot:build-image -Dspring-boot.build-image.imageName=registry.local:5000/micropet/pets-spring

Run image with docker
docker run --name k8s-demo-app -p 8080:8080  pets:0.0.1-SNAPSHOT

