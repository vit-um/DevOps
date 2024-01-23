variable "github_repository" {
  type        = string
  description = "GitHub repository to store Flux manifests"
}

variable "flux_namespace" {
  type        = string
  default     = "flux-system"
  description = "GitHub repository to store Flux manifests"
}

variable "target_path" {
  type        = string
  default     = "clusters"
  description = "Flux manifests subdirectory"
}

variable "github_token" {
  type        = string
  default     = ""
  description = "The token used to authenticate with the Git repository"
}

variable "private_key" {
  type        = string
  description = "The private key used to authenticate with the Git repository"
}

variable "config_host" {
  type        = string
  default     = "gke"
  description = "The url for gke"
}

variable "config_token" {
  type        = string
  default     = "token"
  description = "The token for gke"
}

variable "config_ca" {
  type        = string
  default     = "ca"
  description = "The ca for gke"
}