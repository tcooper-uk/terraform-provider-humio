resource "humio_notifier" "example_pagerduty" {
  repository = "humio"
  name       = "example_pagerduty"
  type     = "PagerDutyAction"

  pagerduty {
    routing_key = "XXXXXXXXXXXXXXX"
    severity    = "critical"
  }
}
