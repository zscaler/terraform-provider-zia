
resource "zia_url_categories" "blacklist_category"{
    super_category = "USER_DEFINED"
    configured_name = "Blacklist Category"
    description = "Blacklist Category"
    custom_category = true
    editable = true
    type = "URL_CATEGORY"
    urls = [
            ".malware-traffic-analysis.net",
            ".secure.eicar.org",
        ]
}


resource "zia_url_categories" "ssl_bypass"{
    super_category = "USER_DEFINED"
    configured_name = "SSL Bypass List"
    description = "SSL Bypass List"
    custom_category = true
    editable = true
    type = "URL_CATEGORY"
    urls = [
            ".zpath.net",
            ".zscaler.com",
            ".splunk.com",
            "splunkbase.splunk.com",
            "api.private.zscaler.com",
            "dist.private.zscaler.com",
            "zpa-updates.prod.zpath.net",
        ]
}

resource "zia_url_categories" "isolate_allow_download"{
    super_category = "USER_DEFINED"
    configured_name = "Isolate - Allow Download - No CopyPaste - Render Office Docs"
    description = "Isolate - Allow Download - No CopyPaste - Render Office Docs"
    custom_category = true
    type = "URL_CATEGORY"
    urls = [
            "file-examples.com/index.php/sample-documents-download/sample-ppt-file"
        ]
}

resource "zia_url_categories" "isolate_allow_paste"{
    super_category = "USER_DEFINED"
    configured_name = "Isolate - Allow Paste - Up Down - Render Docs"
    description = "Isolate - Allow Paste - Up Down - Render Docs"
    custom_category = true
    editable = true
    type = "URL_CATEGORY"
    urls = [
            ".malicious.safemarch.com"
        ]
}

resource "zia_url_categories" "isolate_no_transfer"{
    super_category = "USER_DEFINED"
    configured_name = "Isolate - No Transfer - No CopyPaste - Render Office Docs"
    description = "Isolate - No Transfer - No CopyPaste - Render Office Docs"
    custom_category = true
    editable = true
    type = "URL_CATEGORY"
    urls = [
            ".controlc.com",
            ".hastebin.com",
            "file-examples.com/index.php/sample-documents-download/sample-pdf-download",
        ]
}
/*
resource "zia_url_categories" "mcas_unsanctioned_apps" {
  super_category      = "USER_DEFINED"
  configured_name     = "MCAS Unsanctioned Apps"
  description         = "MCAS Unsanctioned Apps"
  keywords            = ["microsoft"]
  custom_category     = true
  db_categorized_urls = [".creditkarma.com", ".youku.com"]
  type                = "URL_CATEGORY"
  scopes {
    type = "LOCATION"
    scope_entities {
      id = [zia_location_management.ca_vancouver_ipsec.id]
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
*/