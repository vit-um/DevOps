provider "flux" {
  kubernetes = {
    host                   = var.config_host
    client_key             = var.config_client_key
    cluster_ca_certificate = var.config_ca
    client_certificate     = var.config_crt
  }
  git = {
    url = "https://github.com/${var.github_repository}.git"
    http = {
      username = "git"
      password = var.github_token
    }
  }
}

resource "flux_bootstrap_git" "this" {
  path = var.target_path
}
