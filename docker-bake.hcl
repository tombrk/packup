variable "PLATFORMS" {
  default = ["linux/amd64", "linux/arm64"]
}

group "default" {
  targets = ["server", "agent"]
}

target "_common" {
  context    = "."
  dockerfile = "Dockerfile"
  platforms  = PLATFORMS
}

target "server" {
  inherits = ["_common"]
  target   = "server"
  tags     = ["ghcr.io/tombrk/packup/server:latest"]
}

target "agent" {
  inherits = ["_common"]
  target   = "agent"
  tags     = ["ghcr.io/tombrk/packup/agent:latest"]
}
