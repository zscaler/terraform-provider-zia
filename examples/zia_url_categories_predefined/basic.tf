resource "zia_url_categories_predefined" "education" {
  name = "EDUCATION"
  urls = [
    ".internal-learning.example.com",
    ".corporate-training.example.com",
  ]
}

resource "zia_url_categories_predefined" "finance" {
  name = "FINANCE"
  keywords = [
    "internal-trading",
    "corporate-finance",
  ]
  ip_ranges = [
    "10.0.0.0/8",
    "172.16.0.0/12",
  ]
}