terraform {
  backend "gcs" {
    bucket = "vit-secret"
    prefix = "terraform/state"
  }
}

module "gke_cluster" {
  source         = "github.com/vit-um/tf-google-gke-cluster?ref=Task2W7"
  GOOGLE_REGION  = var.GOOGLE_REGION
  GOOGLE_PROJECT = var.GOOGLE_PROJECT
  GKE_NUM_NODES  = var.GKE_NUM_NODES
}

