FRONT_SRC ?= ../Downloads/carbon-scan-now-main

.PHONY: web-build
web-build:
	@echo "Building frontend from $(FRONT_SRC) and copying to backend/frontend/dist"
	./scripts/copy_frontend_dist.sh $(FRONT_SRC)

.PHONY: docker-build docker-run
docker-build:
	@echo "Building docker image eco-link:latest (tries docker, falls back to sudo)"
	@docker build -t eco-link:latest . || sudo docker build -t eco-link:latest .

docker-run:
	@echo "Running docker container eco-link:latest (maps 8081)"
	@docker run -p 8081:8081 eco-link:latest || sudo docker run -p 8081:8081 eco-link:latest
