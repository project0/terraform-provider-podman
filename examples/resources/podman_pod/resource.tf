# Create default local volume
resource "podman_volume" "vol" {}

# A pod with mounts
resource "podman_pod" "pod" {
  name     = "mypod"
  hostname = "superhost"
  mounts = [
    {
      destination = "/mount"
      bind = {
        path  = "/mnt/local"
        chown = true
      }
    },
    {
      destination = "/data"
      volume = {
        name      = podman_volume.vol.name
        read_only = true
      }
    },
  ]
}
