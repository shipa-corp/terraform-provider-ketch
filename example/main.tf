terraform {
  required_providers {
    ketch = {
      version = "0.0.1"
      source = "ketch.io/terraform/ketch"
    }
  }
}

provider "ketch" {}

# Create framework
resource "ketch_framework" "tf" {
  name = "tf-dev"
  ingress_controller {
    class_name = "istio"
    service_endpoint = "1.2.3.4"
    type = "istio"
  }
}

# Create app
resource "ketch_app" "tf" {
  name = "tf-app1"
  image = "docker.io/shipasoftware/bulletinboard:1.0"
  framework = ketch_framework.tf.name
  cnames = [
    "app1.ketch.io",
    "app2.ketch.io",
    "app3.ketch.io"]
  ports = [8081, 8082]
  units = 5
  processes {
    name = "web"
    cmd = [
      "docker-entrypoint.sh",
      "npm",
      "start"]
  }
  routing_settings {
    weight = 100
  }
}

# Create job
resource "ketch_job" "tf" {
  name = "tf-job1"
  framework = ketch_framework.tf.name
  containers {
    name = "pi"
    image = "perl"
    command = [
      "perl",
      "-Mbignum=bpi",
      "-wle",
      "print bpi(2000)"
    ]
  }
}
