provider "google" {
  credentials = "${file("/Users/den/code/gcp/cred-gke.json")}"
  project     = "smartcity-gl"
  region      = "europe-west2"
  zone        = "europe-west2-a"
}

resource "google_container_cluster" "primary" {
  name                     = "demo"
  location                 = "europe-west2-a"
  remove_default_node_pool = true
  initial_node_count       = 1


}

resource "google_container_node_pool" "primary-pool" {
  name       = "primary-pool"
  cluster    = "${google_container_cluster.primary.name}"
  location   = "europe-west2-a"
  node_count = "2"

  node_config {
    preemptible  = true
    machine_type = "n1-standard-2"

    oauth_scopes = [
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]
  }

  autoscaling {
    min_node_count = 2
    max_node_count = 5
  }

  management {
    auto_repair  = true
    auto_upgrade = true
  }
}

