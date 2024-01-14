output "values" {
  value = {
    url            = github_repository.this.html_url
    name           = github_repository.this.name
    full_name      = github_repository.this.full_name
    git_clone_url  = github_repository.this.git_clone_url
    ssh_clone_url  = github_repository.this.ssh_clone_url
    http_clone_url = github_repository.this.http_clone_url
  }
  description = "The GitHub repository details"
}