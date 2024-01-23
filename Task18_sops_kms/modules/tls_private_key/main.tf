resource "tls_private_key" "this" {
  algorithm   = var.algorithm
  ecdsa_curve = var.ecdsa_curve
}