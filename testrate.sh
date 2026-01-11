#!/bin/bash


echo "health test"
for i in {1..100}; do
	curl http://localhost:8080/health
	done

echo "404 test"
for i in {1..100}; do
	curl http://localhost:8080/
	done