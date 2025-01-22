resource "zia_url_filtering_and_cloud_app_settings" "this" {
    block_skype                             = true
    consider_embedded_sites                 = false
    enable_block_override_for_non_auth_user = false
    enable_chatgpt_prompt                   = false
    enable_cipa_compliance                  = false
    enable_dynamic_content_cat              = true
    enable_gemini_prompt                    = false
    enable_meta_prompt                      = false
    enable_microsoft_copilot_prompt         = false
    enable_msft_o365                        = false
    enable_newly_registered_domains         = false
    enable_office365                        = true
    enable_per_plexity_prompt               = false
    enable_poep_prompt                      = false
    enable_ucaas_logmein                    = false
    enable_ucaas_ring_central               = false
    enable_ucaas_talkdesk                   = false
    enable_ucaas_webex                      = false
    enable_ucaas_zoom                       = false
    enforce_safe_search                     = false
}