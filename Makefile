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
    grep -v -e '_easyjson.go' -e 'cmd' -e 'docs' -e 'db' -e 'monitoring' -e 'proto' -e 'mock' -e 'gmail' coverprofile_.tmp > coverprofile.tmp ; \
    rm coverprofile_.tmp ; \
    go tool cover -html coverprofile.tmp -o ../heatmap.html; \
    go tool cover -func coverprofile.tmp

mock:
	mockgen -source=internal/microservice/auth/proto/auth_grpc.pb.go -destination=internal/microservice/auth/mock/auth_grpc_mock.go -package=mock proto AuthServiceClient \
	&& mockgen -source=internal/microservice/email/proto/email_grpc.pb.go -destination=internal/microservice/email/mock/email_grpc_mock.go -package=mock proto EmailServiceClient \
	&& mockgen -source=internal/microservice/email/interface/iemail_repo.go -destination=internal/microservice/email/mock/email_repository_mock.go -package=mock \
	&& mockgen -source=internal/microservice/email/interface/iemail_service.go -destination=internal/microservice/email/mock/email_service_mock.go -package=mock \
	&& mockgen -source=internal/microservice/folder/interface/ifolder_repo.go -destination=internal/microservice/folder/mock/folder_repo_mock.go -package=mock \
	&& mockgen -source=internal/microservice/folder/interface/ifolder_service.go -destination=internal/microservice/folder/mock/folder_service_mock.go -package=mock \
	&& mockgen -source=internal/microservice/folder/proto/folder_grpc.pb.go -destination=internal/microservice/folder/mock/folder_grpc_mock.go -package=mock proto FolderServiceClient \
	&& mockgen -source=internal/microservice/questionnaire/interface/iquestion_repo.go -destination=internal/microservice/questionnaire/mock/question_repository_mock.go -package=mock \
	&& mockgen -source=internal/microservice/questionnaire/interface/iquestion_service.go -destination=internal/microservice/questionnaire/mock/question_service_mock.go -package=mock \
	&& mockgen -source=internal/microservice/questionnaire/proto/question-answer_grpc.pb.go -destination=internal/microservice/questionnaire/mock/question-answer_grpc_mock.go -package=mock proto QuestionServiceClient \
	&& mockgen -source=internal/microservice/session/interface/isession_repo.go -destination=internal/microservice/session/mock/session_repository_mock.go -package=mock \
	&& mockgen -source=internal/microservice/session/interface/isession_service.go -destination=internal/microservice/session/mock/session_service_mock.go -package=mock \
	&& mockgen -source=internal/microservice/session/proto/session_grpc.pb.go -destination=internal/microservice/session/mock/session_grpc_mock.go -package=mock proto SessionServiceClient \
	&& mockgen -source=internal/microservice/user/interface/iuser_repo.go -destination=internal/microservice/user/mock/user_repository_mock.go -package=mock \
	&& mockgen -source=internal/microservice/user/interface/iuser_service.go -destination=internal/microservice/user/mock/user_service_mock.go -package=mock \
	&& mockgen -source=internal/microservice/user/proto/user_grpc.pb.go -destination=internal/microservice/user/mock/user_grpc_mock.go -package=mock proto UserServiceClient

swag:
	swag init -g cmd/mail/main.go

lint:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run

tree:
	tree  > file_tree.txt