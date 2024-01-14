# Terraform Flux Bootstrap Git Module

This Terraform module creates a Git repository to be used as a source for Flux Bootstrap.

## Usage

```hcl
module "flux_bootstrap" {
  source            = "github.com/den-vasyliev/tf-fluxcd-flux-bootstrap?ref=gke_auth"
  github_repository = "${var.GITHUB_OWNER}/${var.FLUX_GITHUB_REPO}"
  private_key       = module.tls_private_key.private_key_pem
  config_host       = module.gke_cluster.config_host
  config_token       = module.gke_cluster.config_token
  config_ca         = module.gke_cluster.config_ca
  github_token      = var.GITHUB_TOKEN
}
```
## Inputs
- github_repository - (Required) The name of the Git repository to be created.
- private_key - (Optional) The SSH private key to use for Git operations.
- config_host - (Required) Kubernetes APIServer endpoint.
- config_token - (Required) Client token for authenticating to cluster.
- config_ca - (Required) Client verifies the server certificate with this CA cert.
- github_token -(Required) The GitHub token ised by fluxcd/flux provider

## Outputs
None.

License
This module is licensed under the MIT License. See the LICENSE file for details.
