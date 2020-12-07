FROM node:15.3.0-alpine3.10 AS builder
COPY . ./builder
WORKDIR /builder

RUN npm i
RUN $(npm bin)/ng build --prod

FROM nginx:1.19.5-alpine
EXPOSE 80
COPY --from=builder /builder/dist/gui/ /usr/share/nginx/html