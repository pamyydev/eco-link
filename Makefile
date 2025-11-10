FRONT_SRC ?= ../Downloads/carbon-scan-now-main

.PHONY: web-build
web-build:
	@echo "Building frontend from $(FRONT_SRC) and copying to backend/frontend/dist"
	./scripts/copy_frontend_dist.sh $(FRONT_SRC)
