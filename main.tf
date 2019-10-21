resource "google_compute_network" "network" {
  name                    = "${var.instance_name}"
  auto_create_subnetworks = "true"
}

resource "google_compute_firewall" "default" {
  name    = "app-firewall"
  network = "${google_compute_network.network.name}"

  allow {
    protocol = "icmp"
  }

  allow {
    protocol = "tcp"
    ports    = ["3333"]
  }

  source_ranges = ["0.0.0.0/0"]
}

resource "google_compute_instance" "instance" {
  name         = "${var.instance_name}"
  machine_type = "f1-micro"

  service_account {
    scopes = ["compute-ro", "logging-write", "monitoring-write"]
  }

  metadata = {
    user-data           = "${file("script/cloud-init.yml")}"
    cos-metrics-enabled = true
    cos-update-strategy = "disabled" // Allowing auto-updates will put instances on different versions
  }

  boot_disk {
    initialize_params {
      image = "gce-uefi-images/cos-stable"
    }
  }

  network_interface {
    network = "${google_compute_network.network.self_link}"
    access_config {
    }
  }
}
