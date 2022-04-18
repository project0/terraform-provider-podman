terraform {
  required_providers {
    podman = {
      source = "project0/podman"
    }
  }
}

# Per default connects to local unix socket
provider "podman" {
  // default
  uri = "unix:///run/podman/podman.sock"
}

# connect via ssh
provider "podman" {
  alias    = "ssh"
  uri      = "ssh://<user>@<host>[:port]/run/podman/podman.sock?secure=True"
  identity = "/tmp/ssh_identity_key"
}
