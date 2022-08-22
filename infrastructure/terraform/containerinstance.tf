resource "azurerm_container_group" "api" {
  name                = "${var.container_name}-${var.env}"
  location            = data.azurerm_resource_group.default_resource_group.location
  resource_group_name = data.azurerm_resource_group.default_resource_group.name
  ip_address_type     = "Public"
  dns_name_label      = "${var.container_name}-${var.env}"
  os_type             = "Linux"
  restart_policy      = "Always"
  image_registry_credential {
    server   = var.registry_url
    username = var.registry_username
    password = var.registry_password
  }
  container {
    name   = var.container_name
    image  = "${var.registry_url}/${var.image_name}:${var.image_tag}"
    cpu    = var.container_cpu
    memory = var.container_memory
    liveness_probe {
      http_get {
        path   = "/healthcheck"
        port   = var.container_port
        scheme = "Http"
      }
      initial_delay_seconds = 10
      period_seconds        = 30
      failure_threshold     = 3
      success_threshold     = 1
      timeout_seconds       = 30
    }
    environment_variables = {
      API_PORT = var.container_port
      MESSAGE  = "Hello World From Azure - ${var.env}"
    }
    ports {
      port     = var.container_port
      protocol = "TCP"
    }
  }
  tags = {
    environment = var.env
  }
}

output "api_url" {
  value = "http://${azurerm_container_group.api.fqdn}:${var.container_port}"
}

output "api_ip" {
  value = azurerm_container_group.api.ip_address
}
