variable "github_owner" {
  type        = string
  description = "The GitHub owner"
}

variable "github_token" {
  type        = string
  description = "GitHub personal access token"
}

variable "repository_name" {
  type        = string
  default     = "flux-gitops"
  description = "GitHub repository"
}

variable "repository_visibility" {
  type        = string
  default     = "private"
  description = "The visibility of the GitOps repository"
}

variable "branch" {
  type        = string
  default     = "main"
  description = "GitHub branch"
}

variable "public_key_openssh" {
  type        = string
  description = "OpenSSH public key repository access"
}

variable "public_key_openssh_title" {
  type        = string
  description = "The title for OpenSSH public key"
}
