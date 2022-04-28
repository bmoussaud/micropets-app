run-database:
	docker run -p 5432:5432 --name bird-postgres -e POSTGRES_BD=birds -e POSTGRES_PASSWORD=mysecretpassword -d postgres


run-app:
	./mvnw spring-boot:run