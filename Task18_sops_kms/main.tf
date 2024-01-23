terraform {
  backend "gcs" {
    bucket  = "vit-secret"
    prefix  = "terraform/state"
  }
}

module "github_repository" {
  source                   = "./modules/github_repository"
  github_owner             = var.GITHUB_OWNER
  github_token             = var.GITHUB_TOKEN
  repository_name          = var.FLUX_GITHUB_REPO
  public_key_openssh       = module.tls_private_key.public_key_openssh
  public_key_openssh_title = "flux-ssh-pub"
}

module "gke_cluster" {
  source         = "./modules/gke_cluster"
  GOOGLE_REGION  = var.GOOGLE_REGION
  GOOGLE_PROJECT = var.GOOGLE_PROJECT
  GKE_NUM_NODES  = 2
}

module "tls_private_key" {
  source    = "./modules/tls_private_key"
  algorithm = "RSA"
}

module "flux_bootstrap" {
  source            = "./modules/flux_bootstrap"
  github_repository = "${var.GITHUB_OWNER}/${var.FLUX_GITHUB_REPO}"
  private_key       = module.tls_private_key.private_key_pem
  config_host       = module.gke_cluster.config_host
  config_token      = module.gke_cluster.config_token
  config_ca         = module.gke_cluster.config_ca
  github_token      = var.GITHUB_TOKEN
}

module "gke-workload-identity" {
  source              = "terraform-google-modules/kubernetes-engine/google//modules/workload-identity"
  use_existing_k8s_sa = true
  annotate_k8s_sa     = true
  name                = "kustomize-controller"
  namespace           = "flux-system"
  project_id          = var.GOOGLE_PROJECT
  location            = var.GOOGLE_REGION
  cluster_name        = "main"  
  roles               = ["roles/cloudkms.cryptoKeyEncrypterDecrypter"]

  module_depends_on = [
    module.flux_bootstrap
  ]
}


# data "google_kms_key_ring" "key_ring" {
#   name     = "sops-flux"
#   location = "global"
#   project  = var.GOOGLE_PROJECT
# }

# import {
#   to = google_kms_key_ring.key_ring
#   id = "projects/${var.GOOGLE_PROJECT}/locations/${data.google_kms_key_ring.key_ring.location}/keyRings/${data.google_kms_key_ring.key_ring.name}"
# }

# resource "google_kms_key_ring" "key_ring" {
#   count    = data.google_kms_key_ring.key_ring.name != null ? 0 : 1
#   name     = "sops-flux"
#   location = "global"
#   project  = var.GOOGLE_PROJECT

#   lifecycle {
#     prevent_destroy = false
#   }
# }

module "kms" {
  source             = "terraform-google-modules/kms/google"
  version            = "2.2.3"
  project_id         = var.GOOGLE_PROJECT
  keyring            = "sops-flux5"
  location           = "global"
  keys               = ["sops-key-flux"]
  prevent_destroy    = false
}


# module "google_kms" {
#   source          = "terraform-google-modules/kms/google"
#   version         = "2.2.3"
#   keyring         = coalesce(data.google_kms_key_ring.key_ring.name, try(google_kms_key_ring.key_ring[0].name, null))
#   keys            = ["sops-key-flux"]
#   location        = "global"
#   project_id      = var.GOOGLE_PROJECT
#   prevent_destroy = false
# }


resource "null_resource" "cluster_credentials" {
  depends_on = [
    module.gke_cluster,
    module.flux_bootstrap
  ]

  provisioner "local-exec" {
    command = <<EOF
      ${module.gke_cluster.cluster.gke_get_credentials_command}
    EOF
  }
}

resource "null_resource" "git_commit" {
  depends_on = [
    module.flux_bootstrap,
    module.kms
  ]

  provisioner "local-exec" {
    command = <<EOF
      if [ -d ${var.FLUX_GITHUB_REPO} ]; then
        rm -rf ${var.FLUX_GITHUB_REPO}
      fi
      git clone ${module.github_repository.values.http_clone_url}    
      ./sops -e -gcp-kms ${module.kms.keys.sops-key-flux} --encrypted-regex '^(token)$' secret.yaml > ./demo_app/demo/secret-enc.yaml
      cp -r demo_app/* ${var.FLUX_GITHUB_REPO}/${var.FLUX_GITHUB_TARGET_PATH}/     
      cd ${var.FLUX_GITHUB_REPO}
      git add .
      git commit -m"Added all manifests"
      git push
      cd ..
      rm -rf ${var.FLUX_GITHUB_REPO}
    EOF
  }
}

# resource "null_resource" "gitops_destroy" {
#   triggers = {
#     repo_name = module.github_repository.values.name
#   }

#   provisioner "local-exec" {
#     when    = destroy
#     command = <<EOF
#       if [ -d ${self.triggers.repo_name} ]; then
#         rm -rf ${self.triggers.repo_name}
#       fi
#     EOF
#   }
# }