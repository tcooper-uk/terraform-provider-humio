resource "humio_alert" "example_alert_with_labels" {
  repository  = humio_action.example_email.repository
  name        = "example_alert_with_labels"

  actions   = [humio_action.example_email.action_id]

  labels               = ["terraform", "ops"]
  throttle_time_millis = 300000
  enabled              = true
  query                = "count()"
  start                = "24h"
}

resource "humio_alert" "example_alert_without_labels" {
  repository  = humio_action.example_email_body.repository
  name        = "example_alert_without_labels"

  actions = [humio_action.example_email_body.action_id]

  throttle_time_millis = 300000
  enabled              = true
  query                = "count()"
  start                = "24h"
}

resource "humio_alert" "example_alert_with_description" {
  repository  = humio_action.example_email_body.repository
  name        = "example_alert_with_description"
  description = "lorem ipsum...."

  actions = [
    humio_action.example_email_body.action_id,
    humio_action.example_email_subject.action_id,
  ]

  link_url             = "http://localhost:8080/humio/search?query=count()&live=true&start=24h&fullscreen=false"
  labels               = ["terraform", "ops"]
  throttle_time_millis = 300000
  enabled              = true
  query                = "count()"
  start                = "24h"
}
