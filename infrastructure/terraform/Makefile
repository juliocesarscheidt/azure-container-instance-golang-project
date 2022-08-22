# env config
ENV?=development
# azure config
REGION?=eastus
RESOURCE_GROUP?=
# backend config
STORAGE_ACCOUNT_BACKEND?=
STORAGE_ACCOUNT_BACKEND_LOCATION?=
# registry oonfig
REGISTRY_USERNAME?=
REGISTRY_URL?=$(REGISTRY_USERNAME).azurecr.io
# image config
IMAGE_NAME?=http-simple-api-$(ENV)
IMAGE_TAG?=v1.0.0
# docker and terraform config
DOCKER_BUILDKIT=1
BUILDKIT_PROGRESS=plain
PLAN_FILE?=tfplan

all: push-image init apply

local-run:
	docker container run --rm --name $(IMAGE_NAME) -d -p 80:80 "$(REGISTRY_URL)/$(IMAGE_NAME):$(IMAGE_TAG)"

create-repo:
	-@az acr create --resource-group $(RESOURCE_GROUP) --name $(REGISTRY_USERNAME) --sku Basic
	-@az acr update --admin-enabled true --resource-group $(RESOURCE_GROUP) --name $(REGISTRY_USERNAME)

docker-login:
	export REGISTRY_PASSWORD="$$(az acr credential show --resource-group $(RESOURCE_GROUP) --name $(REGISTRY_USERNAME) --query passwords[0].value | sed 's/^"//; s/"$$//')"; \
	echo "$${REGISTRY_PASSWORD}" | \
		docker login "$(REGISTRY_URL)" \
			--username "$(REGISTRY_USERNAME)" \
			--password-stdin

build-image:
	docker image build --tag "$(REGISTRY_URL)/$(IMAGE_NAME):$(IMAGE_TAG)" ../../src

push-image: create-repo build-image docker-login
	docker image push "$(REGISTRY_URL)/$(IMAGE_NAME):$(IMAGE_TAG)"

create-backend-storage-account:
	az storage account show --resource-group $(RESOURCE_GROUP) --name $(STORAGE_ACCOUNT_BACKEND) || \
		az storage account create --resource-group $(RESOURCE_GROUP) --name $(STORAGE_ACCOUNT_BACKEND) \
		--location $(STORAGE_ACCOUNT_BACKEND_LOCATION) --sku Standard_LRS

create-workspace:
	-@terraform workspace new development 2> /dev/null
	-@terraform workspace new test 2> /dev/null
	-@terraform workspace new production 2> /dev/null
	terraform workspace select $(ENV)

init: create-backend-storage-account
	terraform init -upgrade=true \
		-backend-config="resource_group_name=$(RESOURCE_GROUP)" \
		-backend-config="storage_account_name=$(STORAGE_ACCOUNT_BACKEND)" \
		-backend-config="container_name=terraform" \
		-backend-config="key=$(env).terraform.tfstate"
	make create-workspace
	make plan

tf-fmt:
	terraform fmt -write=true -recursive

tf-validate:
	terraform validate

tf-lint:
	docker container run --rm -t \
		--name tflint \
		--env TFLINT_LOG=debug \
		-v $$(pwd):/data \
		ghcr.io/terraform-linters/tflint

output:
	terraform output

plan: tf-validate tf-lint tf-fmt
	export REGISTRY_PASSWORD="$$(az acr credential show --resource-group $(RESOURCE_GROUP) --name $(REGISTRY_USERNAME) --query passwords[0].value | sed 's/^"//; s/"$$//')"; \
	terraform plan \
		-out=$(PLAN_FILE) \
		-var-file=$(ENV).tfvars \
		-var region="$(REGION)" \
		-var env="$(ENV)" \
		-var resource_group="$(RESOURCE_GROUP)" \
		-var registry_url="$(REGISTRY_URL)" \
		-var registry_username="$(REGISTRY_USERNAME)" \
		-var registry_password="$${REGISTRY_PASSWORD}" \
		-var image_name="$(IMAGE_NAME)" \
		-var image_tag="$(IMAGE_TAG)" \
		-input=false

apply:
	terraform apply $(PLAN_FILE)
	make output

container-logs:
	az container logs --resource-group $(RESOURCE_GROUP) --name $(IMAGE_NAME) --only-show-errors --follow

container-attach:
	az container attach --resource-group $(RESOURCE_GROUP) --name $(IMAGE_NAME)

destroy:
	export REGISTRY_PASSWORD="$$(az acr credential show --resource-group $(RESOURCE_GROUP) --name $(REGISTRY_USERNAME) --query passwords[0].value | sed 's/^"//; s/"$$//')"; \
	terraform destroy \
		-var-file=$(ENV).tfvars \
		-var region="$(REGION)" \
		-var env="$(ENV)" \
		-var resource_group="$(RESOURCE_GROUP)" \
		-var registry_url="$(REGISTRY_URL)" \
		-var registry_username="$(REGISTRY_USERNAME)" \
		-var registry_password="$${REGISTRY_PASSWORD}" \
		-var image_name="$(IMAGE_NAME)" \
		-var image_tag="$(IMAGE_TAG)" \
		-auto-approve
