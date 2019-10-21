resource "runscope_bucket" "bucket" {
  name      = "${var.instance_name}"
  team_uuid = "${var.runscope_team_id}"
}

resource "runscope_test" "api" {
  name        = "${var.instance_name}"
  description = "Ensure echo service is up, running, and available"
  bucket_id   = "${runscope_bucket.bucket.id}"
}

resource "runscope_step" "main_page" {
  bucket_id = "${runscope_bucket.bucket.id}"
  test_id   = "${runscope_test.api.id}"
  step_type = "request"
  url       = "http://${google_compute_instance.instance.network_interface.0.access_config.0.nat_ip}:3333"
  note      = "Is our endpoint up?"
  method    = "POST"
  variables {
    name   = "httpStatus"
    source = "response_status"
  }

  assertions {
    source     = "response_status"
    comparison = "equal_number"
    value      = "200"
  }
}

resource "runscope_environment" "environment" {
  bucket_id = "${runscope_bucket.bucket.id}"
  name      = "production"
}

resource "runscope_schedule" "api" {
  bucket_id      = "${runscope_bucket.bucket.id}"
  test_id        = "${runscope_test.api.id}"
  interval       = "5m"
  note           = "Test API every 5 mins"
  environment_id = "${runscope_environment.environment.id}"
}
