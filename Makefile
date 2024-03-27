.DEFAULT_GOAL := build
API := ${API_NAME}
CLEARING := ${CLR_NAME}
DOCKER_COMPOSE := ${DOCKER_COMPOSE}
DOCKER_NETWORK := ${DOCKER_NETWORK}
ENV := ${PATH_SRC}/.env
REQUIRED_VARIABLES := PATH_BIN PATH_SRC PATH_WORKBENCH

all: init up run
build: build_init build_services
build_init: dep
	@echo "${GUM_PREFIX}building init"
	@cd init && go mod tidy
	@. $(ENV) && cd init && GOOS="" go build -o ${PATH_BIN}/init
build_services: dep
	@echo "${GUM_PREFIX}building $(API)"
	@. $(ENV) && cd $(API) && sed -i '' -e 's|//.*@BasePath.*$$|//	@BasePath	${API_SWAGGER_BASEPATH}|g' "$(API).go"
	@cd $(API) && swag fmt --generalInfo $(API).go
	@cd $(API) && swag init --generalInfo $(API).go --output ./swagger --parseDependency true --parseDependencyLevel 3 --parseInternal true
	@cd $(API) && go mod tidy
	@. $(ENV) && cd $(API) && GOOS="" go build -o ${PATH_BIN}/$(API)-native
	@. $(ENV) && cd $(API) && GOOS="linux" go build -o ${PATH_BIN}/$(API)-linux
	@. $(ENV) && cp -r $(API)/${API_SWAGGER_DIRPATH} ${PATH_WORKBENCH}

	@echo "${GUM_PREFIX}building $(CLEARING)"
	@cd $(CLEARING) && go mod tidy
	@. $(ENV) && cd $(CLEARING) && GOOS="" go build -o ${PATH_BIN}/$(CLEARING)-native
	@. $(ENV) && cd $(CLEARING) && GOOS="linux" go build -o ${PATH_BIN}/$(CLEARING)-linux
clean: dep down
	@. $(ENV) && gum style 'THIS TARGET CAN BE DESTRUCTIVE' 'IT SHOULD BE RUN WITH SPECIAL CARE'
	@. $(ENV) && gum confirm "${GUM_PREFIX}do you want to proceed?" || exit 1
	@. $(ENV) && gum confirm "${GUM_PREFIX}is PATH_WORKBENCH == ${PATH_WORKBENCH} correct?" && gum spin --title "removing ${PATH_WORKBENCH}/*" -- find "${PATH_WORKBENCH}" -mindepth 1 -delete || exit 1
	@. $(ENV) && echo "${GUM_PREFIX}${PATH_WORKBENCH} is clean"
dep:
	@. $(ENV)
	$(foreach var,$(REQUIRED_VARIABLES),$(if $(value $(var)),,$(error $(GUM_PREFIX) $(var) is not set, load $(ENV) before running make)))
	@. $(ENV) && echo "${GUM_PREFIX}checkig for dependencies"
	@go version
	@gum --version
	@swag --version
	@. $(ENV) && echo "${GUM_PREFIX}all good"
down:
	@. $(ENV) && echo "${GUM_PREFIX}kill all docker"
	-$(DOCKER_COMPOSE) -f ${PATH_WORKBENCH}/docker-compose.yaml -p $(DOCKER_NETWORK) down
	@if [ -n "$$(docker ps -q)" ]; then docker kill $(docker ps -q); fi
	@if [ -n "$$(docker ps -aq)" ]; then docker rm $(docker ps -aq); fi
	-docker network rm $(DOCKER_NETWORK)
init: clean build_init
	@. $(ENV) && echo "${GUM_PREFIX}executing ${PATH_BIN}/init templates"
	@. $(ENV) && ${PATH_BIN}/init
	-$(MAKE) up_infra
	@. $(ENV) && echo "${GUM_PREFIX}executing ${PATH_BIN}/init gorm"
	@. $(ENV) && gum spin --title "waiting for services for ${DB_INIT_TIMEOUT}s..." -- sleep ${DB_INIT_TIMEOUT}
	@. $(ENV) && ${PATH_BIN}/init gorm
	@. $(ENV) && echo "${GUM_PREFIX}deactivating ${PATH_BIN}/init"
	chmod -x ${PATH_BIN}/init 
	mv ${PATH_BIN}/init ${PATH_BIN}/init_deactivated
up: build_services up_infra
	@. $(ENV) && echo "${GUM_PREFIX}$(DOCKER_COMPOSE) up services"
	-docker network create $(DOCKER_NETWORK)
	$(DOCKER_COMPOSE) -f ${PATH_WORKBENCH}/docker-compose.yaml -p $(DOCKER_NETWORK) up -d ${API_NAME} ${CLR_NAME}
up_all: dep build_services
	@. $(ENV) && echo "${GUM_PREFIX}$(DOCKER_COMPOSE) up all"
	-docker network create $(DOCKER_NETWORK)
	$(DOCKER_COMPOSE) -f ${PATH_WORKBENCH}/docker-compose.yaml -p $(DOCKER_NETWORK) up -d
up_infra: dep
	@. $(ENV) && echo "${GUM_PREFIX}$(DOCKER_COMPOSE) up infra"
	-docker network create $(DOCKER_NETWORK)
	$(DOCKER_COMPOSE) -f ${PATH_WORKBENCH}/docker-compose.yaml -p $(DOCKER_NETWORK) up -d ${ADMINER_NAME} ${CDC_NAME} ${DB_NAME} ${MQ_NAME}
