resource "zia_atp_malicious_urls" "this" {
    malicious_urls = [
        "test1.malicious.com",
        "test2.malicious.com",
        "test3.malicious.com",
    ]
}