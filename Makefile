start-local-storage:
	docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=postgres -d postgres

start-local-file-storage:
	docker run --name s3 -p 4566:4566 -p 4571:4571 -e SERVICES=s3 -d localstack/localstack:1.1.0

run:
	export LOCAL_DEV=true && export CONFIG_DIR=config && go run cmd/main.go