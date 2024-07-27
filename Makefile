PROJECT_NAME=rate_limiter
SERVICES=redis app

build:
	docker compose build
	
up:
	docker compose up -d
	
up-build: build up
	
down:
	docker compose down
	
up-%:
	docker compose up -d $*
	
up-build-%:
	docker compose build $*
	docker compose up -d $*
	
down-%:
	docker compose stop $*
	docker compose rm -f $*