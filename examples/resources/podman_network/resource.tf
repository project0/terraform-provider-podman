# Create default bridge network
resource "podman_network" "network" {}

# Full dual-stack example (netavark backend)
resource "podman_network" "dualstack" {
  name        = "dualstack"
  driver      = "bridge"
  ipam_driver = "dhcp"
  options = {
    mtu = 1500
  }
  internal = false
  dns      = true
  # enable dual stack
  ipv6 = true
  subnets = [
    {
      subnet  = "2001:db8::/64"
      gateway = "2001:db8::1"
    },
    {
      subnet  = "192.0.2.0/24/24"
      gateway = "192.0.2.1"
    }
  ]
}