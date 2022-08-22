variable "region" {
  type        = string
  default     = "eastus" # us east
  description = "Region"
}

variable "env" {
  type        = string
  default     = "development"
  description = "Environment"
}

variable "resource_group" {
  type        = string
  description = "Resource Group"
}

variable "registry_url" {
  type        = string
  description = "ACR Registry URL"
}

variable "registry_username" {
  type        = string
  description = "ACR Registry Username"
  sensitive   = true
}

variable "registry_password" {
  type        = string
  description = "ACR Registry Password"
  sensitive   = true
}

variable "container_name" {
  type        = string
  default     = "http-simple-api"
  description = "Container Name"
}

variable "container_cpu" {
  type        = string
  default     = "1"
  description = "Container CPU Amount"
}

variable "container_memory" {
  type        = string
  default     = "1"
  description = "Container Memory Amount"
}

variable "container_port" {
  type        = number
  default     = 9000
  description = "Container TCP Port"
}

variable "image_name" {
  type        = string
  default     = "http-simple-api"
  description = "Docker Image Name"
}

variable "image_tag" {
  type        = string
  default     = "v1.0.0"
  description = "Docker Image Tag"
}
