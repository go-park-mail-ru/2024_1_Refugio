# path to docker compose file
DCOMPOSE:=docker-compose.yml

# improve build time
DOCKER_BUILD_KIT:=COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1

down:
	docker compose -f ${DCOMPOSE} down --remove-orphans

build:
	${DOCKER_BUILD_KIT} docker compose build

up:
	docker compose up --build -d --remove-orphans

# Vendoring is useful for local debugging since you don't have to
# reinstall all packages again and again in docker
mod:
	go mod tidy -compat=1.22 && go mod vendor && go install ./...

tests:
	go test -json ./... -coverprofile coverprofile_.tmp -coverpkg=./... ; \
    grep -v -e '_easyjson.go' -e 'gen_notes.go' -e 'cmd' -e 'docs' -e 'db' -e 'monitoring' -e 'proto' -e 'mock' coverprofile_.tmp > coverprofile.tmp ; \
    rm coverprofile_.tmp ; \
    go tool cover -html coverprofile.tmp -o ../heatmap.html; \
    go tool cover -func coverprofile.tmp

mock:
	mockgen -source=internal/microservice/email/proto/email_grpc.pb.go -destination=internal/microservice/email/mock/email_grpc_mock.go -package=mock proto EmailServiceClient \
	&& mockgen -source=internal/microservice/email/interface/iemail_repo.go -destination=internal/microservice/email/mock/email_repository_mock.go -package=mock \
	&& mockgen -source=internal/microservice/email/interface/iemail_service.go -destination=internal/microservice/email/mock/email_service_mock.go -package=mock

swag:
	swag init -g cmd/mail/main.go

lint:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run