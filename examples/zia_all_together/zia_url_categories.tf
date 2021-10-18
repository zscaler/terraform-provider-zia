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

resource "zia_url_categories" "mcas_unsanctioned_apps"{
    super_category = "USER_DEFINED"
    configured_name = "MCAS Unsanctioned Apps"
    description = "MCAS Unsanctioned Apps"
    custom_category = true
    editable = true
    type = "URL_CATEGORY"
    keywords_retaining_parent_category = [".creditkarma.com", ".youku.com"]
    urls = [
        ".aa.com",
        ".accuweather.com",
        ".agoda.com",
        ".alaskaair.com",
        ".alexa.com",
        ".baidu.com",
        ".bittorrent.com",
        ".centurylink.com",
        ".cnn.com",
        ".controlc.com",
        ".coupons.com",
        ".dailymotion.com",
        ".demonoid.pw",
        ".extratorrent.cc",
        ".filedropper.com",
        ".ft.com",
        ".gigaom.com",
        ".hubpages.com",
        ".lemonde.fr",
        ".level3.com",
        ".Logz.io",
        ".mini900.cn",
        ".nationalgeographic.com",
        ".osiriscomm.com",
        ".pasted.co",
        ".resource.alaskaair.net",
        ".savvis.com",
        ".singaporeair.com",
        ".slashdot.org",
        ".sozcu.com.tr",
        ".strava.com",
        ".techrepublic.com",
        ".telegram.org",
        ".theguardian.com",
        ".thepiratebay.org",
        ".thepiratebay.se",
        ".uefa.com"
    ]
}