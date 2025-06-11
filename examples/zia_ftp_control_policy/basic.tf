resource "zia_ftp_control_policy" "this" {
    ftp_enabled = true
    ftp_over_http_enabled = true
    url_categories = ["HOBBIES_AND_LEISURE","HEALTH","HISTORY","INSURANCE","IMAGE_HOST","INTERNET_SERVICES","GOVERNMENT"]
    urls = ["test1.acme.com", "test10.acme.com"]
}
