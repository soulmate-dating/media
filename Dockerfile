FROM golang:1.20-alpine as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server cmd/main/main.go


FROM ubuntu
LABEL authors="ageev"

COPY --from=build ./app/server ./server
COPY --from=build ./app/s3-policy.json ./s3-policy.json

EXPOSE 8082
CMD ./server