resource "humio_notifier" "example_opsgenie" {
  repository = "humio"
  name       = "example_opsgenie"
  type     = "OpsGenieAction"

  opsgenie {
    api_url   = "https://api.opsgenie.com"
    genie_key = "XXXXXXXXXXXXXXX"
  }
}
