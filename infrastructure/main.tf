terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "5.7.0"
    }
  }
}

variable "region" {
  default = "europe-central2"
}

variable "project" {
  type = string
  nullable = false
}

variable "billing_account" {
  type = string
  nullable = false
}

provider "google" {
  project = var.project
  region  = var.region
}

resource "google_project" "project" {
  project_id = "${var.project}"
  name       = var.project
  billing_account = var.billing_account
}

resource "time_sleep" "wait_30_seconds" {
  depends_on = [google_project.project]

  create_duration = "30s"
}

resource "google_project_service" "firestore" {
  project = google_project.project.project_id
  service = "firestore.googleapis.com"

  depends_on = [ 
    google_project.project,
    time_sleep.wait_30_seconds 
  ]
}

resource "google_project_service" "cloudbuild" {
  project = google_project.project.project_id
  service = "cloudbuild.googleapis.com"

  depends_on = [ 
    google_project.project,
    time_sleep.wait_30_seconds 
  ]
}

resource "google_project_service" "cloudfunctions" {
  project = google_project.project.project_id
  service = "cloudfunctions.googleapis.com"

  depends_on = [ 
    google_project.project,
    time_sleep.wait_30_seconds 
  ]
}

resource "google_firestore_database" "database" {
  project     = google_project.project.project_id
  name        = "(default)"
  location_id = "eur3"
  type        = "FIRESTORE_NATIVE"

  depends_on = [google_project_service.firestore]
}

resource "random_id" "default" {
  byte_length = 8
}

resource "google_storage_bucket" "bucket" {
  name     = "${random_id.default.hex}-function-bucket"
  location = "EU"
  depends_on = [ 
    google_project.project,
  ]
}

data "archive_file" "function-archive" {
  type        = "zip"
  output_path = "/tmp/function-source.zip"
  source_dir  = "../"
}

resource "google_storage_bucket_object" "api-function-archive" {
  name   = "api-function.zip"
  bucket = google_storage_bucket.bucket.name
  source = data.archive_file.function-archive.output_path
  depends_on = [ 
    google_project.project,
    google_storage_bucket.bucket,
    data.archive_file.function-archive
  ]
}

resource "google_cloudfunctions_function" "function" {
  name                  = "api-function"
  description           = "The main entrypoint"
  runtime               = "go121"
  source_archive_bucket = google_storage_bucket.bucket.name
  source_archive_object = google_storage_bucket_object.api-function-archive.name
  trigger_http          = true
  entry_point           = "Entrypoint"
  available_memory_mb   = 128
  timeout               = 30
  ingress_settings      = "ALLOW_INTERNAL_AND_GCLB"
  max_instances         = 4
  min_instances         = 0
  environment_variables = {
    "GCP_PROJECT" = google_project.project.project_id
  }
  depends_on = [
    google_project.project,
    google_project_service.cloudbuild,
    google_project_service.cloudfunctions
  ]
}

# IAM entry for all users to invoke the function
resource "google_cloudfunctions_function_iam_member" "invoker" {
  project        = google_cloudfunctions_function.function.project
  region         = google_cloudfunctions_function.function.region
  cloud_function = google_cloudfunctions_function.function.name

  role   = "roles/cloudfunctions.invoker"
  member = "allUsers"
  depends_on = [ 
    google_project.project,
    google_cloudfunctions_function.function,
  ]
}