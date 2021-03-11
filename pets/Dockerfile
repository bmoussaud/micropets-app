FROM golang:latest AS build
WORKDIR /src
COPY . .

RUN go build -o /out/pets .
#FROM scratch AS bin
#FROM alpine:latest  
EXPOSE 7004
#COPY --from=build /out/pets /bin
CMD ["/out/pets"]