resource "zia_url_categories" "example" {
  super_category      = "USER_DEFINED"
  configured_name     = "MCAS Unsanctioned Apps2"
  description         = "MCAS Unsanctioned Apps2"
  keywords            = ["microsoft"]
  custom_category     = true
  db_categorized_urls = [".creditkarma.com", ".youku.com"]
  type                = "URL_CATEGORY"
  scopes {
    type = "LOCATION"
    scope_entities {
      id = [33079472]
    }
    scope_group_member_entities {
      id = []
    }
  }
  urls = [
    ".coupons.com",
    ".resource.alaskaair.net",
    ".techrepublic.com",
    ".dailymotion.com",
    ".osiriscomm.com",
    ".uefa.com",
    ".Logz.io",
    ".alexa.com",
    ".baidu.com",
    ".cnn.com",
    ".level3.com",
  ]
}

