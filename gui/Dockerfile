FROM node:12.10.0-alpine AS builder
COPY . ./builder
WORKDIR /builder

RUN npm i
RUN $(npm bin)/ng build --prod

FROM nginx:1.17.3-alpine
EXPOSE 80
COPY --from=builder /builder/dist/gui/ /usr/share/nginx/html