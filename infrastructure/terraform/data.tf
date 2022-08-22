data "azurerm_subscription" "subscription" {}

data "azurerm_client_config" "client_config" {}

data "azurerm_resource_group" "default_resource_group" {
  name = var.resource_group
}
