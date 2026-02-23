data "zia_datacenters" "filtered" {
    name = "CA Client Node DC"
}

## Example Usage - Filter by Multiple Criteria

data "zia_datacenters" "filtered" {
    city            = "San Jose"
    dc_provider     = "Zscaler Internal"
    gov_only        = false
    third_party_cloud = false
    virtual         = false
}