resource "humio_notifier" "example_email" {
  repository = "sandbox"
  name       = "example_email"
  type     = "EmailAction"

  email {
    recipients = ["ops@example.com"]
  }
}

resource "humio_notifier" "example_email_body" {
  repository = "humio"
  name       = "example_email_body"
  type     = "EmailAction"

  email {
    recipients    = ["ops@example.com"]
    body_template = "{event_count}"
  }
}

resource "humio_notifier" "example_email_subject" {
  repository = "humio"
  name       = "example_email_subject"
  type     = "EmailAction"

  email {
    recipients       = ["ops@example.com"]
    subject_template = "{alert_name}"
  }
}

resource "humio_notifier" "example_email_body_subject" {
  repository = "humio"
  name       = "example_email_body_subject"
  type     = "EmailAction"

  email {
    recipients       = ["ops@example.com"]
    body_template    = "{event_count}"
    subject_template = "{alert_name}"
  }
}
