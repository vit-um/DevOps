# ${path.module} is an interpolated string in Terraform, which references the path to the current
# module. In this case, it returns the path to the directory containing the current module.

# output "kubeconfig" {
#  value       = "${path.module}/kubeconfig"
#  description = "The path to the kubeconfig file"
#}

output "config_host" {
  value = "https://${data.google_container_cluster.main.endpoint}"
}

output "config_token" {
  value = data.google_client_config.current.access_token
}

output "config_ca" {
  value = base64decode(
    data.google_container_cluster.main.master_auth[0].cluster_ca_certificate,
  )
}

output "name" {
  value = google_container_cluster.this.name
}

output "cluster" {
  value = {
    name    = google_container_cluster.this.name
    zone    = google_container_cluster.this.location
    project = google_container_cluster.this.project

    gke_get_credentials_command = "gcloud container clusters get-credentials ${google_container_cluster.this.name} --zone ${google_container_cluster.this.location} --project ${google_container_cluster.this.project}"
  }
  description = "The GKE cluster details"
}