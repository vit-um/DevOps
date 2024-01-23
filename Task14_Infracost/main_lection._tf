terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "4.52.0"
    }
  }
}

variable "GOOGLE_PROJECT" {
  type        = string
  description = "GCP project name"
}

variable "GOOGLE_REGION" {
  type        = string
  default     = "us-central1-c"
  description = "GCP region to use"
}


variable "GKE_MACHINE_TYPE" {
  type        = string
  default     = "g1-small"
  description = "Machine type"
}

variable "GKE_NUM_NODES" {
  type        = number
  default     = 2
  description = "GKE nodes number"
}

resource "google_container_cluster" "demo" {
  name     = "demo"
  location = var.GOOGLE_REGION

  remove_default_node_pool = true
  initial_node_count       = 1
}


resource "google_container_node_pool" "main" {
  name       = "main"
  project    = google_container_cluster.demo.project
  cluster    = google_container_cluster.demo.name
  location   = google_container_cluster.demo.location
  node_count = var.GKE_NUM_NODES

  node_config {
    machine_type = var.GKE_MACHINE_TYPE
  }
}

data "google_compute_instance_group" "node_instance_groups" {
  self_link = google_container_cluster.demo.node_pool[0].managed_instance_group_urls[0]
}

data "google_compute_instance" "nodes" {
  for_each = toset(data.google_compute_instance_group.node_instance_groups.instances[*])
  self_link = each.key
}

output "node_ip" {
  value = [for node in data.google_compute_instance.nodes : node.network_interface[0].access_config[0].nat_ip]
}

provider "google" {
  # Configuration options
  project = var.GOOGLE_PROJECT
  region  = var.GOOGLE_REGION
}


