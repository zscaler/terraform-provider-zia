resource "zia_dlp_engines" "this" {
    name = "Example1000"
    description = "Example1000"
    engine_expression = "((D63.S > 1))"
    custom_dlp_engine = true
}