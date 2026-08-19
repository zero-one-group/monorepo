resource "time_static" "placeholder" {
  count = var.enabled ? 1 : 0
}
