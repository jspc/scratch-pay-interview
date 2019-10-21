// GCP auth json comes from $GOOGLE_CLOUD_KEYFILE_JSON
provider "google" {
  project = "scratch-pay-tech-test"
  region  = "us-central1"
  zone    = "us-central1-c"
}

// Runscope auth token comes from $RUNSCOPE_ACCESS_TOKEN
provider "runscope" {}
