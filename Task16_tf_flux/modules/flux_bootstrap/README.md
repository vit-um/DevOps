# Terraform Flux Bootstrap Git Module

This Terraform module creates a Git repository to be used as a source for Flux Bootstrap.

## Usage

```hcl
module "flux_bootstrap" {
  source            = "github.com/den-vasyliev/tf-fluxcd-flux-bootstrap?ref=kind_auth
  github_repository = "${var.GITHUB_OWNER}/${var.FLUX_GITHUB_REPO}"
  private_key       = module.tls_private_key.private_key_pem
  config_host       = module.kind_cluster.endpoint
  config_client_key = module.kind_cluster.client_key
  config_ca         = module.kind_cluster.ca
  config_crt        = module.kind_cluster.crt
  github_token      = var.GITHUB_TOKEN
}
}
```
## Inputs
- github_repository - (Required) The name of the Git repository to be created.
- target_path - (Optional) The path to clone the Git repository into. Default value is clusters.
- private_key - (Required) The SSH private key to use for Git operations.
- config_host - (Required) Kubernetes APIServer endpoint.
- config_client_key - (Required) Client key for authenticating to cluster.
- config_ca - (Required) Client verifies the server certificate with this CA cert.
- config_crt - (Required) Client certificate for authenticating to cluster.
- github_token - (Required) The GitHub token ised by fluxcd/flux provider

## Outputs
None.

License
This module is licensed under the MIT License. See the LICENSE file for details.
