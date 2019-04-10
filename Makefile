#!/usr/bin/make
build:
	go build -o ./bin/create_user ./cmd/cli/*
	go build -o ./bin/shorty ./cmd/web/*

clean:
	rm bin/create_user
	rm bin/shorty
