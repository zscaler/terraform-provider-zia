resource "zia_dlp_notification_templates" "terraform_dlp_template" {
  name               = "Terraform DLP Template"
  attach_content     = true
  tls_enabled 		 = true
  subject            = local.subject
  plain_text_message = local.msg_plain_text
  html_message       = local.msg_html
}

locals {

subject = <<SUBJECT
  DLP Violation: $${TRANSACTION_ID} $${RULENAME}
SUBJECT

msg_plain_text = <<MSGPLAINTEXT
    The attached content triggered a Web DLP rule for your organization.

    Transaction ID: $${TRANSACTION_ID}
    User Accessing the URL: $${USER}
    URL Accessed: $${URL}
    Posting Type: $${TYPE}
    DLP MD5: $${DLPMD5}
    Triggered DLP Violation Engines (assigned to the hit rule): $${ENGINES_IN_RULE}
    Triggered DLP Violation Dictionaries (assigned to the hit rule): $${DICTIONARIES}

    No action is required on your part.
MSGPLAINTEXT

msg_html = <<MSGHTML
  <!DOCTYPE html>
<html>
	<head>
		<style>
			.user {color: rgb(1, 81, 152);}
			.url {color: rgb(1, 81, 152);}
			.postingtype {color: rgb(1, 81, 152);}
			.engines {color: rgb(1, 81, 152);}
			.dictionaries {color: rgb(1, 81, 152);}
		</style>
	</head>
	<body>
		The attached content triggered a Web DLP rule for your organization.
		<br/><br/>
		Transaction ID: <span class="transaction_id">$${TRANSACTION_ID}</span>
		<br/>
		User Accessing the URL: <span class="user">$${USER}</span>
		<br/>
		URL Accessed: <span class="url">$${URL}</span>
		<br/>
		Posting Type: <span class="postingtype">$${TYPE}</span>
		<br/>
		DLP MD5: <span class="dlpmd5">$${DLPMD5}</span>
		<br/>
		Triggered DLP Violation Engines (assigned to the hit rule): <span class="engines">$${ENGINES_IN_RULE}</span>
		<br/>
		Triggered DLP Violation Dictionaries (assigned to the hit rule): <span class="dictionaries">$${DICTIONARIES}</span>
		<br/><br/>
		No action is required on your part.
		<br/><br/>
	</body>
</html>
MSGHTML
}