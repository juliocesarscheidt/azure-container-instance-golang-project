# HTTP Simple API with Azure Container Instances

A simple HTTP API made with Golang, to use the service of Azure Container Instances

## Up and Running with Terraform

```bash
# env config
export ENV='development'
# azure config
export REGION='REGION'
export RESOURCE_GROUP='RESOURCE_GROUP'
# backend config
export STORAGE_ACCOUNT_BACKEND='STORAGE_ACCOUNT_BACKEND'
export STORAGE_ACCOUNT_BACKEND_LOCATION='STORAGE_ACCOUNT_BACKEND_LOCATION'
# registry oonfig
export REGISTRY_USERNAME="REGISTRY_USERNAME"
export REGISTRY_URL="$REGISTRY_USERNAME.azurecr.io"

cd infrastructure/terraform/

# create the ACR repository, login into the ACR, build the image and pushes the image to the repository
make push-image

# create the backend storageaccount on Azure (if doesn't exist, this could take a few minutes), initializes the terraform, create the workspaces, validate and do the plan
make init

# apply the previous plan
make apply
```
