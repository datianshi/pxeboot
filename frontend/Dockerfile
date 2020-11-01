FROM node:12.19.0-alpine3.10 AS builder
WORKDIR /app
COPY . .
RUN yarn run build

FROM node:12.19.0-alpine3.10
RUN yarn global add serve
WORKDIR /app
COPY --from=builder /app/build .
CMD ["serve", "-p", "8080", "-s", "."]
