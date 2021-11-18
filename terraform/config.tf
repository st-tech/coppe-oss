terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "3.5.0"
    }
  }

  backend "gcs" {
    bucket = "coppe-tfstate-bucket"
  }
}

provider "google" {
  project = var.project_id
  region  = "us-central1"
}
