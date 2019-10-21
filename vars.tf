variable "instance_name" {
  default = "tech-test"
}

variable "runscope_team_id" {
  type        = string
  default     = "47991038-2c27-4ddb-bc20-42c8992d13ec"
  description = "Runscope Team UUID; this can only be generated and found via the UI"
}
