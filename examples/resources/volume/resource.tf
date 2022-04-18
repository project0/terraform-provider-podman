# Create default local volume
resource "podman_volume" "vol" {}

# A local volume with mount options
resource "podman_volume" "tmpfs" {
  name   = "myvol"
  driver = "local"
  # driver specific options
  options = {
    # mount device
    device = "tmpfs"
    # mount -t
    type = "tmpfs"
    # mount options
    o = "nodev,noexec"
  }
}