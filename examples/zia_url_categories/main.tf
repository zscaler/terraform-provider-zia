terraform {
  required_providers {
    zia = {
      version = "1.0.0"
      source  = "zscaler.com/zia/zia"
    }
  }
}
provider "zia" {}

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
    ".alaskaair.com",
    ".filedropper.com",
    ".nationalgeographic.com",
    ".mini900.cn",
    ".lemonde.fr",
    ".telegram.org",
    ".extratorrent.cc",
    ".strava.com",
    ".slashdot.org",
    ".demonoid.pw",
    ".sozcu.com.tr",
    ".bittorrent.com",
    ".ft.com",
    ".thepiratebay.org",
    ".theguardian.com",
    ".accuweather.com",
    ".aa.com",
    ".agoda.com",
    ".centurylink.com",
    ".singaporeair.com",
    ".savvis.com",
    ".thepiratebay.se"
  ]
}

data "zia_url_categories" "example"{
    //id = "SOCIAL_NETWORKING"
    id = zia_url_categories.example.id
    //configured_name = "Custom_Category"
}

output "zia_url_categories"{
    value = data.zia_url_categories.example
}
