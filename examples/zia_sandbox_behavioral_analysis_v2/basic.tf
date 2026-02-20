resource "zia_sandbox_behavioral_analysis_v2" "this" {
  md5_hash_value_list {
    url         = "4EE43B71BB89CB9CBF7784495AE8D0DF"
    url_comment = "4EE43B71BB89CB9CBF7784495AE8D0DF"
    type        = "CUSTOM_FILEHASH_ALLOW"
  }

  md5_hash_value_list {
    url         = "8350dED6D39DF158E51D6CFBE36FB012"
    url_comment = "8350dED6D39DF158E51D6CFBE36FB012"
    type        = "CUSTOM_FILEHASH_DENY"
   }
}