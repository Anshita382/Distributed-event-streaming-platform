SHELL := /bin/bash

.PHONY: up down logs topics build run-java run-producer run-consumer fmt clean smoke

up:
	docker compose up -d --build

down:
	docker compose down -v

logs:
	docker compose logs -f --tail=100

topics:
	bash scripts/create-topics.sh

build:
	cd java-core-service && mvn -q -DskipTests package
	cd go-producer && go build ./cmd/producer
	cd go-consumer && go build ./cmd/consumer

run-java:
	cd java-core-service && mvn spring-boot:run

run-producer:
	cd go-producer && go run ./cmd/producer

run-consumer:
	cd go-consumer && go run ./cmd/consumer

fmt:
	cd go-producer && gofmt -w $$(find . -name '*.go')
	cd go-consumer && gofmt -w $$(find . -name '*.go')

smoke:
	curl -s http://localhost:8080/api/events/health

clean:
	find . -type f -name 'producer' -delete
	find . -type f -name 'consumer' -delete
