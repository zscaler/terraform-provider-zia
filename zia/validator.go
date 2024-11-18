package zia

import (
	"context"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/biter777/countries"
	"github.com/fabiotavarespr/iso3166"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/dlp/dlp_web_rules"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/urlfilteringpolicies"
)

// Validate URL Filtering Category Options
var supportedURLCategories = []string{
	"ANY", "NONE",
	"OTHER_ADULT_MATERIAL", "ADULT_THEMES", "LINGERIE_BIKINI", "NUDITY", "PORNOGRAPHY", "SEXUALITY", "ADULT_SEX_EDUCATION", "K_12_SEX_EDUCATION", "SOCIAL_ADULT", "OTHER_BUSINESS_AND_ECONOMY", "CORPORATE_MARKETING", "FINANCE", "PROFESSIONAL_SERVICES", "CLASSIFIEDS", "TRADING_BROKARAGE_INSURANCE", "OTHER_DRUGS", "MARIJUANA", "OTHER_EDUCATION", "CONTINUING_EDUCATION_COLLEGES", "HISTORY", "K_12", "REFERENCE_SITES", "SCIENCE_AND_TECHNOLOGY", "OTHER_ENTERTAINMENT_AND_RECREATION", "ENTERTAINMENT", "TELEVISION_AND_MOVIES", "MUSIC", "STREAMING_MEDIA", "RADIO_STATIONS", "GAMBLING", "OTHER_GAMES", "SOCIAL_NETWORKING_GAMES", "OTHER_GOVERNMENT_AND_POLITICS", "GOVERNMENT", "POLITICS", "HEALTH", "OTHER_ILLEGAL_OR_QUESTIONABLE", "COPYRIGHT_INFRINGEMENT", "COMPUTER_HACKING", "QUESTIONABLE", "PROFANITY", "MATURE_HUMOR", "ANONYMIZER", "OTHER_INFORMATION_TECHNOLOGY", "TRANSLATORS", "IMAGE_HOST", "FILE_HOST", "SHAREWARE_DOWNLOAD", "WEB_BANNERS", "WEB_HOST", "WEB_SEARCH", "PORTALS", "SAFE_SEARCH_ENGINE", "CDN", "OSS_UPDATES", "DNS_OVER_HTTPS", "OTHER_INTERNET_COMMUNICATION", "INTERNET_SERVICES", "DISCUSSION_FORUMS", "ONLINE_CHAT", "EMAIL_HOST", "BLOG", "P2P_COMMUNICATION", "REMOTE_ACCESS", "WEB_CONFERENCING", "ZSPROXY_IPS", "JOB_SEARCH", "MILITANCY_HATE_AND_EXTREMISM", "OTHER_MISCELLANEOUS", "MISCELLANEOUS_OR_UNKNOWN", "NEWLY_REG_DOMAINS", "NON_CATEGORIZABLE", "NEWS_AND_MEDIA", "OTHER_RELIGION", "TRADITIONAL_RELIGION", "CULT", "ALT_NEW_AGE", "OTHER_SECURITY", "ADWARE_OR_SPYWARE", "ENCR_WEB_CONTENT", "MALICIOUS_TLD", "OTHER_SHOPPING_AND_AUCTIONS", "SPECIALIZED_SHOPPING", "REAL_ESTATE", "ONLINE_AUCTIONS", "OTHER_SOCIAL_AND_FAMILY_ISSUES", "SOCIAL_ISSUES", "FAMILY_ISSUES", "OTHER_SOCIETY_AND_LIFESTYLE", "ART_CULTURE", "ALTERNATE_LIFESTYLE", "HOBBIES_AND_LEISURE", "DINING_AND_RESTAURANT", "ALCOHOL_TOBACCO", "SOCIAL_NETWORKING", "OTHER_SHOPPING_AND_AUCTIONS", "SPECIALIZED_SHOPPING", "REAL_ESTATE", "ONLINE_AUCTIONS", "SPECIAL_INTERESTS", "SPORTS", "TASTELESS", "TRAVEL", "USER_DEFINED", "VEHICLES", "VIOLENCE", "WEAPONS_AND_BOMBS", "OTHER_SECURITY", "ADWARE_OR_SPYWARE", "P2P_COMMUNICATION", "MISCELLANEOUS_OR_UNKNOWN", "SOCIAL_ADULT", "SOCIAL_NETWORKING_GAMES", "REMOTE_ACCESS", "NEWLY_REG_DOMAINS", "CDN", "NON_CATEGORIZABLE", "WEB_CONFERENCING", "ZSPROXY_IPS", "ENCR_WEB_CONTENT",
	"OSS_UPDATES", "TRADING_BROKARAGE_INSURANCE", "DNS_OVER_HTTPS", "MARIJUANA", "DYNAMIC_DNS", "MILITARY", "AI_ML_APPS", "NEWLY_REVIVED_DOMAINS", "CUSTOM_00", "CUSTOM_01", "CUSTOM_02", "CUSTOM_03", "CUSTOM_04", "CUSTOM_05", "CUSTOM_06", "CUSTOM_07", "CUSTOM_08", "CUSTOM_09", "CUSTOM_10", "CUSTOM_11", "CUSTOM_12", "CUSTOM_13", "CUSTOM_14", "CUSTOM_15", "CUSTOM_16", "CUSTOM_17", "CUSTOM_18", "CUSTOM_19", "CUSTOM_20", "CUSTOM_21", "CUSTOM_22", "CUSTOM_23", "CUSTOM_24", "CUSTOM_25", "CUSTOM_26", "CUSTOM_27", "CUSTOM_28", "CUSTOM_29", "CUSTOM_30", "CUSTOM_31", "CUSTOM_32", "CUSTOM_33", "CUSTOM_34", "CUSTOM_35", "CUSTOM_36", "CUSTOM_37", "CUSTOM_38", "CUSTOM_39", "CUSTOM_40", "CUSTOM_41", "CUSTOM_42", "CUSTOM_43", "CUSTOM_44", "CUSTOM_45", "CUSTOM_46", "CUSTOM_47", "CUSTOM_48", "CUSTOM_49", "CUSTOM_50", "CUSTOM_51", "CUSTOM_52", "CUSTOM_53", "CUSTOM_54", "CUSTOM_55", "CUSTOM_56", "CUSTOM_57", "CUSTOM_58", "CUSTOM_59", "CUSTOM_60", "CUSTOM_61", "CUSTOM_62", "CUSTOM_63", "CUSTOM_64", "CUSTOM_65", "CUSTOM_66", "CUSTOM_67", "CUSTOM_68", "CUSTOM_69", "CUSTOM_70", "CUSTOM_71", "CUSTOM_72", "CUSTOM_73", "CUSTOM_74", "CUSTOM_75", "CUSTOM_76", "CUSTOM_77", "CUSTOM_78", "CUSTOM_79", "CUSTOM_80", "CUSTOM_81", "CUSTOM_82", "CUSTOM_83", "CUSTOM_84", "CUSTOM_85", "CUSTOM_86", "CUSTOM_87", "CUSTOM_88", "CUSTOM_89", "CUSTOM_90", "CUSTOM_91", "CUSTOM_92", "CUSTOM_93", "CUSTOM_94", "CUSTOM_95", "CUSTOM_96", "CUSTOM_97", "CUSTOM_98", "CUSTOM_99", "CUSTOM_100", "CUSTOM_101", "CUSTOM_102", "CUSTOM_103", "CUSTOM_104", "CUSTOM_105", "CUSTOM_106", "CUSTOM_107", "CUSTOM_108", "CUSTOM_109", "CUSTOM_110", "CUSTOM_111", "CUSTOM_112", "CUSTOM_113", "CUSTOM_114", "CUSTOM_115", "CUSTOM_116", "CUSTOM_117", "CUSTOM_118", "CUSTOM_119", "CUSTOM_120", "CUSTOM_121", "CUSTOM_122", "CUSTOM_123", "CUSTOM_124", "CUSTOM_125", "CUSTOM_126", "CUSTOM_127", "CUSTOM_128", "CUSTOM_129", "CUSTOM_130", "CUSTOM_131", "CUSTOM_132", "CUSTOM_133", "CUSTOM_134", "CUSTOM_135", "CUSTOM_136", "CUSTOM_137", "CUSTOM_138", "CUSTOM_139", "CUSTOM_140", "CUSTOM_141", "CUSTOM_142", "CUSTOM_143", "CUSTOM_144", "CUSTOM_145", "CUSTOM_146", "CUSTOM_147", "CUSTOM_148", "CUSTOM_149", "CUSTOM_150", "CUSTOM_151", "CUSTOM_152", "CUSTOM_153", "CUSTOM_154", "CUSTOM_155", "CUSTOM_156", "CUSTOM_157", "CUSTOM_158", "CUSTOM_159", "CUSTOM_160", "CUSTOM_161", "CUSTOM_162", "CUSTOM_163", "CUSTOM_164", "CUSTOM_165", "CUSTOM_166", "CUSTOM_167", "CUSTOM_168", "CUSTOM_169", "CUSTOM_170", "CUSTOM_171", "CUSTOM_172", "CUSTOM_173", "CUSTOM_174", "CUSTOM_175", "CUSTOM_176", "CUSTOM_177", "CUSTOM_178", "CUSTOM_179", "CUSTOM_180", "CUSTOM_181", "CUSTOM_182", "CUSTOM_183", "CUSTOM_184", "CUSTOM_185", "CUSTOM_186", "CUSTOM_187", "CUSTOM_188", "CUSTOM_189", "CUSTOM_190", "CUSTOM_191", "CUSTOM_192", "CUSTOM_193", "CUSTOM_194", "CUSTOM_195", "CUSTOM_196", "CUSTOM_197", "CUSTOM_198", "CUSTOM_199", "CUSTOM_200", "CUSTOM_201", "CUSTOM_202", "CUSTOM_203", "CUSTOM_204", "CUSTOM_205", "CUSTOM_206", "CUSTOM_207", "CUSTOM_208", "CUSTOM_209", "CUSTOM_210", "CUSTOM_211", "CUSTOM_212", "CUSTOM_213", "CUSTOM_214", "CUSTOM_215", "CUSTOM_216", "CUSTOM_217", "CUSTOM_218", "CUSTOM_219", "CUSTOM_220", "CUSTOM_221", "CUSTOM_222", "CUSTOM_223", "CUSTOM_224", "CUSTOM_225", "CUSTOM_226", "CUSTOM_227", "CUSTOM_228", "CUSTOM_229", "CUSTOM_230", "CUSTOM_231", "CUSTOM_232", "CUSTOM_233", "CUSTOM_234", "CUSTOM_235", "CUSTOM_236", "CUSTOM_237", "CUSTOM_238", "CUSTOM_239", "CUSTOM_240", "CUSTOM_241", "CUSTOM_242", "CUSTOM_243", "CUSTOM_244", "CUSTOM_245", "CUSTOM_246", "CUSTOM_247", "CUSTOM_248", "CUSTOM_249", "CUSTOM_250", "CUSTOM_251", "CUSTOM_252", "CUSTOM_253", "CUSTOM_254", "CUSTOM_255", "CUSTOM_256",
}

func validateURLFilteringCategories() schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		value, ok := i.(string)
		if !ok {
			return diag.Diagnostics{
				{
					Severity: diag.Error,
					Summary:  "Expected type to be string",
					Detail:   "Type assertion failed, expected string type for URL Filtering Categories validation",
				},
			}
		}

		// Convert the cty.Path to a string representation
		pathStr := fmt.Sprintf("%+v", path)

		// Use StringInSlice from helper/validation package
		var diags diag.Diagnostics
		if _, errs := validation.StringInSlice(supportedURLCategories, false)(value, pathStr); len(errs) > 0 {
			for _, err := range errs {
				diags = append(diags, diag.FromErr(err)...)
			}
		}

		return diags
	}
}

var supportedURLFilteringRequestMethods = []string{
	"OPTIONS", "GET", "HEAD", "POST", "PUT", "DELETE", "TRACE", "CONNECT", "OTHER",
}

func validateURLFilteringRequestMethods() schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		value, ok := i.(string)
		if !ok {
			return diag.Diagnostics{
				{
					Severity: diag.Error,
					Summary:  "Expected type to be string",
					Detail:   "Type assertion failed, expected string type for URL Filtering Request Methods validation",
				},
			}
		}

		// Convert the cty.Path to a string representation
		pathStr := fmt.Sprintf("%+v", path)

		// Use StringInSlice from helper/validation package
		var diags diag.Diagnostics
		if _, errs := validation.StringInSlice(supportedURLFilteringRequestMethods, false)(value, pathStr); len(errs) > 0 {
			for _, err := range errs {
				diags = append(diags, diag.FromErr(err)...)
			}
		}

		return diags
	}
}

var supportedURLFilteringProtocols = []string{
	"SMRULEF_ZPA_BROKERS_RULE", "ANY_RULE", "TCP_RULE", "UDP_RULE", "DOHTTPS_RULE", "TUNNELSSL_RULE",
	"HTTP_PROXY", "FOHTTP_RULE", "FTP_RULE", "HTTPS_RULE", "HTTP_RULE", "SSL_RULE", "TUNNEL_RULE", "WEBSOCKETSSL_RULE",
	"WEBSOCKET_RULE",
}

func validateURLFilteringProtocols() schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		value, ok := i.(string)
		if !ok {
			return diag.Diagnostics{
				{
					Severity: diag.Error,
					Summary:  "Expected type to be string",
					Detail:   "Type assertion failed, expected string type for URL Filtering Protocols validation",
				},
			}
		}

		// Convert the cty.Path to a string representation
		pathStr := fmt.Sprintf("%+v", path)

		// Use StringInSlice from helper/validation package
		var diags diag.Diagnostics
		if _, errs := validation.StringInSlice(supportedURLFilteringProtocols, false)(value, pathStr); len(errs) > 0 {
			for _, err := range errs {
				diags = append(diags, diag.FromErr(err)...)
			}
		}

		return diags
	}
}

var supportedURLSuperCategories = []string{
	"ANY", "NONE", "ADVANCED_SECURITY", "ENTERTAINMENT_AND_RECREATION", "NEWS_AND_MEDIA", "USER_DEFINED", "EDUCATION", "BUSINESS_AND_ECONOMY", "JOB_SEARCH", "INFORMATION_TECHNOLOGY", "INTERNET_COMMUNICATION", "OFFICE_365", "CUSTOM_SUPERCATEGORY", "CUSTOM_BP", "CUSTOM_BW", "MISCELLANEOUS", "TRAVEL", "VEHICLES", "GOVERNMENT_AND_POLITICS", "GLOBAL_INT", "GLOBAL_INT_BP", "GLOBAL_INT_BW", "GLOBAL_INT_OFC365", "ADULT_MATERIAL", "DRUGS", "GAMBLING", "VIOLENCE", "WEAPONS_AND_BOMBS", "TASTELESS", "MILITANCY_HATE_AND_EXTREMISM", "ILLEGAL_OR_QUESTIONABLE", "SOCIETY_AND_LIFESTYLE", "HEALTH", "SPORTS", "SPECIAL_INTERESTS_SOCIAL_ORGANIZATIONS", "GAMES", "SHOPPING_AND_AUCTIONS", "SOCIAL_AND_FAMILY_ISSUES", "RELIGION", "SECURITY",
}

var supportedUserRiskScoreLevels = []string{
	"LOW", "MEDIUM", "HIGH", "CRITICAL",
}

func validateUserRiskScoreLevels() schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		value, ok := i.(string)
		if !ok {
			return diag.Diagnostics{
				{
					Severity: diag.Error,
					Summary:  "Expected type to be string",
					Detail:   "Type assertion failed, expected string type for User Risk Score levels",
				},
			}
		}

		// Convert the cty.Path to a string representation
		pathStr := fmt.Sprintf("%+v", path)

		// Use StringInSlice from helper/validation package
		var diags diag.Diagnostics
		if _, errs := validation.StringInSlice(supportedUserRiskScoreLevels, false)(value, pathStr); len(errs) > 0 {
			for _, err := range errs {
				diags = append(diags, diag.FromErr(err)...)
			}
		}

		return diags
	}
}

var supportedUserAgentTypes = []string{
	"OPERA", "FIREFOX", "MSIE", "MSEDGE", "CHROME", "SAFARI", "MSCHREDGE", "OTHER",
}

func validateUserAgentTypes() schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		value, ok := i.(string)
		if !ok {
			return diag.Diagnostics{
				{
					Severity: diag.Error,
					Summary:  "Expected type to be string",
					Detail:   "Type assertion failed, expected string type for User Agent Types",
				},
			}
		}

		// Convert the cty.Path to a string representation
		pathStr := fmt.Sprintf("%+v", path)

		// Use StringInSlice from helper/validation package
		var diags diag.Diagnostics
		if _, errs := validation.StringInSlice(supportedUserAgentTypes, false)(value, pathStr); len(errs) > 0 {
			for _, err := range errs {
				diags = append(diags, diag.FromErr(err)...)
			}
		}

		return diags
	}
}

func validateURLSuperCategories() schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		value, ok := i.(string)
		if !ok {
			return diag.Diagnostics{
				{
					Severity: diag.Error,
					Summary:  "Expected type to be string",
					Detail:   "Type assertion failed, expected string type for URL Super Categories validation",
				},
			}
		}

		// Convert the cty.Path to a string representation
		pathStr := fmt.Sprintf("%+v", path)

		// Use StringInSlice from helper/validation package
		var diags diag.Diagnostics
		if _, errs := validation.StringInSlice(supportedURLSuperCategories, false)(value, pathStr); len(errs) > 0 {
			for _, err := range errs {
				diags = append(diags, diag.FromErr(err)...)
			}
		}

		return diags
	}
}

var supportedAppControlType = []string{
	"AI_ML", "BUSINESS_PRODUCTIVITY", "CONSUMER", "CUSTOM_CAPP", "DNS_OVER_HTTPS", "ENTERPRISE_COLLABORATION", "FILE_SHARE",
	"FINANCE", "HEALTH_CARE", "HOSTING_PROVIDER", "HUMAN_RESOURCES", "INSTANT_MESSAGING", "IT_SERVICES", "LEGAL", "SALES_AND_MARKETING",
	"SOCIAL_NETWORKING", "STREAMING_MEDIA", "SYSTEM_AND_DEVELOPMENT", "WEBMAIL",
}

func validateAppControlType() schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		value, ok := i.(string)
		if !ok {
			return diag.Diagnostics{
				{
					Severity: diag.Error,
					Summary:  "Expected type to be string",
					Detail:   "Type assertion failed, expected string type for User Agent Types",
				},
			}
		}

		// Convert the cty.Path to a string representation
		pathStr := fmt.Sprintf("%+v", path)

		// Use StringInSlice from helper/validation package
		var diags diag.Diagnostics
		if _, errs := validation.StringInSlice(supportedAppControlType, false)(value, pathStr); len(errs) > 0 {
			for _, err := range errs {
				diags = append(diags, diag.FromErr(err)...)
			}
		}

		return diags
	}
}

var validActionsForType = map[string][]string{
	"AI_ML":                    {"ALLOW_AI_ML_WEB_USE", "CAUTION_AI_ML_WEB_USE", "DENY_AI_ML_WEB_USE", "ISOLATE_AI_ML_WEB_USE"},
	"BUSINESS_PRODUCTIVITY":    {"ALLOW_BUSINESS_PRODUCTIVITY_APPS", "BLOCK_BUSINESS_PRODUCTIVITY_APPS", "CAUTION_BUSINESS_PRODUCTIVITY_APPS", "ISOLATE_BUSINESS_PRODUCTIVITY_APPS"},
	"CONSUMER":                 {"ALLOW_CONSUMER_APPS", "BLOCK_CONSUMER_APPS", "CAUTION_CONSUMER_APPS", "ISOLATE_CONSUMER_APPS"},
	"DNS_OVER_HTTPS":           {"ALLOW_DNS_OVER_HTTPS_USE", "DENY_DNS_OVER_HTTPS_USE"},
	"ENTERPRISE_COLLABORATION": {"ALLOW_ENTERPRISE_COLLABORATION_APPS", "BLOCK_ENTERPRISE_COLLABORATION_APPS", "CAUTION_ENTERPRISE_COLLABORATION_APPS", "ISOLATE_ENTERPRISE_COLLABORATION_APPS"},
	"FILE_SHARE":               {"ALLOW_FILE_SHARE_VIEW", "ALLOW_FILE_SHARE_UPLOAD", "CAUTION_FILE_SHARE_VIEW", "DENY_FILE_SHARE_VIEW", "DENY_FILE_SHARE_UPLOAD", "ISOLATE_FILE_SHARE_VIEW"},
	"FINANCE":                  {"ALLOW_FINANCE_USE", "CAUTION_FINANCE_USE", "DENY_FINANCE_USE", "ISOLATE_FINANCE_USE"},
	"HEALTH_CARE":              {"ALLOW_HEALTH_CARE_USE", "CAUTION_HEALTH_CARE_USE", "DENY_HEALTH_CARE_USE", "ISOLATE_HEALTH_CARE_USE"},
	"HOSTING_PROVIDER":         {"ALLOW_HOSTING_PROVIDER_USE", "CAUTION_HOSTING_PROVIDER_USE", "DENY_HOSTING_PROVIDER_USE", "ISOLATE_HOSTING_PROVIDER_USE"},
	"HUMAN_RESOURCES":          {"ALLOW_HUMAN_RESOURCES_USE", "CAUTION_HUMAN_RESOURCES_USE", "DENY_HUMAN_RESOURCES_USE", "ISOLATE_HUMAN_RESOURCES_USE"},
	"INSTANT_MESSAGING":        {"ALLOW_CHAT", "ALLOW_FILE_TRANSFER_IN_CHAT", "BLOCK_CHAT", "BLOCK_FILE_TRANSFER_IN_CHAT", "CAUTION_CHAT", "ISOLATE_CHAT"},
	"IT_SERVICES":              {"ALLOW_IT_SERVICES_USE", "CAUTION_LEGAL_USE", "DENY_IT_SERVICES_USE", "ISOLATE_IT_SERVICES_USE"},
	"LEGAL":                    {"ALLOW_LEGAL_USE", "DENY_DNS_OVER_HTTPS_USE", "DENY_LEGAL_USE", "ISOLATE_LEGAL_USE"},
	"SALES_AND_MARKETING":      {"ALLOW_SALES_MARKETING_APPS", "BLOCK_SALES_MARKETING_APPS", "CAUTION_SALES_MARKETING_APPS", "ISOLATE_SALES_MARKETING_APPS"},
	"STREAMING_MEDIA":          {"ALLOW_STREAMING_VIEW_LISTEN", "ALLOW_STREAMING_UPLOAD", "BLOCK_STREAMING_UPLOAD", "CAUTION_STREAMING_VIEW_LISTEN", "ISOLATE_STREAMING_VIEW_LISTEN"},
	"SOCIAL_NETWORKING":        {"ALLOW_SOCIAL_NETWORKING_VIEW", "ALLOW_SOCIAL_NETWORKING_POST", "BLOCK_SOCIAL_NETWORKING_VIEW", "BLOCK_SOCIAL_NETWORKING_POST", "CAUTION_SOCIAL_NETWORKING_VIEW"},
	"SYSTEM_AND_DEVELOPMENT":   {"ALLOW_SYSTEM_DEVELOPMENT_APPS", "ALLOW_SYSTEM_DEVELOPMENT_UPLOAD", "BLOCK_SYSTEM_DEVELOPMENT_APPS", "BLOCK_SYSTEM_DEVELOPMENT_UPLOAD", "CAUTION_SYSTEM_DEVELOPMENT_APPS", "ISOLATE_SALES_MARKETING_APPS"},
	"WEBMAIL":                  {"ALLOW_WEBMAIL_VIEW", "ALLOW_WEBMAIL_ATTACHMENT_SEND", "ALLOW_WEBMAIL_SEND", "CAUTION_WEBMAIL_VIEW", "BLOCK_WEBMAIL_VIEW", "BLOCK_WEBMAIL_ATTACHMENT_SEND", "BLOCK_WEBMAIL_SEND"},
}

/*
// Add valid applications for types that support tenancy_profile_ids
var validAppsForTypeWithTenancy = map[string][]string{
	"BUSINESS_PRODUCTIVITY":    {"GOOGLEANALYTICS"},
	"ENTERPRISE_COLLABORATION": {"GOOGLECALENDAR", "GOOGLEKEEP", "GOOGLEMEET", "GOOGLESITES", "WEBEX", "SLACK", "WEBEX_TEAMS", "ZOOM"},
	"FILE_SHARE":               {"DROPBOX", "GDRIVE", "GPHOTOS"},
	"HOSTING_PROVIDER":         {"GCLOUDCOMPUTE", "AWS", "IBMSMARTCLOUD", "GAPPENGINE", "GOOGLE_CLOUD_PLATFORM"},
	"IT_SERVICES":              {"MSLOGINSERVICES", "GOOGLOGINSERVICE", "WEBEX_LOGIN_SERVICES", "ZOHO_LOGIN_SERVICES"},
	"SOCIAL_NETWORKING":        {"GOOGLE_GROUPS", "GOOGLE_PLUS"},
	"STREAMING_MEDIA":          {"YOUTUBE", "GOOGLE_STREAMING"},
	"SYSTEM_AND_DEVELOPMENT":   {"GOOGLE_DEVELOPERS", "GOOGLEAPPMAKER"},
	"WEBMAIL":                  {"GOOGLE_WEBMAIL"},
}
*/

func validateActionsCustomizeDiff(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	ruleType := diff.Get("type").(string)
	actions := diff.Get("actions").(*schema.Set).List()

	validActions, typeExists := validActionsForType[ruleType]
	if !typeExists {
		return fmt.Errorf("invalid type: %s", ruleType)
	}

	validActionsMap := make(map[string]struct{}, len(validActions))
	for _, action := range validActions {
		validActionsMap[action] = struct{}{}
	}

	var invalidActions []string
	var containsIsolate bool
	var containsDenyOrBlock bool
	for _, action := range actions {
		actionStr, ok := action.(string)
		if !ok {
			return fmt.Errorf("expected action to be a string, got: %T", action)
		}
		if _, valid := validActionsMap[actionStr]; !valid {
			invalidActions = append(invalidActions, actionStr)
		}
		if strings.Contains(actionStr, "ISOLATE_") {
			containsIsolate = true
		}
		if strings.Contains(actionStr, "DENY_") || strings.Contains(actionStr, "BLOCK_") {
			containsDenyOrBlock = true
		}
	}

	if len(invalidActions) > 0 {
		return fmt.Errorf(
			"invalid actions %v for type %s. Valid actions are: %v. Please adjust the type or actions accordingly",
			invalidActions, ruleType, validActions)
	}

	// Additional validation logic for enforce_time_validity, validity_start_time, and validity_end_time
	if enforceTimeValidity, ok := diff.GetOk("enforce_time_validity"); ok && enforceTimeValidity.(bool) {
		if _, ok := diff.GetOk("validity_start_time"); !ok {
			return fmt.Errorf("validity_start_time must be set when enforce_time_validity is true")
		} else {
			validityStartTimeStr := diff.Get("validity_start_time").(string)
			if isSingleDigitDay(validityStartTimeStr) {
				return fmt.Errorf("validity_start_time must have a two-digit day (e.g., 02 instead of 2)")
			}
		}
		if _, ok := diff.GetOk("validity_end_time"); !ok {
			return fmt.Errorf("validity_end_time must be set when enforce_time_validity is true")
		} else {
			validityEndTimeStr := diff.Get("validity_end_time").(string)
			if isSingleDigitDay(validityEndTimeStr) {
				return fmt.Errorf("validity_end_time must have a two-digit day (e.g., 02 instead of 2)")
			}
		}
		if _, ok := diff.GetOk("validity_time_zone_id"); !ok {
			return fmt.Errorf("validity_time_zone_id must be set when enforce_time_validity is true")
		}
	} else {
		// If enforce_time_validity is false, ensure validity attributes are not set
		if _, ok := diff.GetOk("validity_start_time"); ok {
			return fmt.Errorf("validity_start_time can only be set when enforce_time_validity is true")
		}
		if _, ok := diff.GetOk("validity_end_time"); ok {
			return fmt.Errorf("validity_end_time can only be set when enforce_time_validity is true")
		}
		if _, ok := diff.GetOk("validity_time_zone_id"); ok {
			return fmt.Errorf("validity_time_zone_id can only be set when enforce_time_validity is true")
		}
	}

	// Validation 1: When the "actions" value contains the string "ISOLATE_", the attribute cbi_profile must be set.
	if containsIsolate {
		if _, ok := diff.GetOk("cbi_profile"); !ok {
			return fmt.Errorf("cbi_profile attribute must be set when actions contain ISOLATE_")
		}
		cbiProfileList := diff.Get("cbi_profile").([]interface{})
		if len(cbiProfileList) == 0 || cbiProfileList[0] == nil {
			return fmt.Errorf("cbi_profile attribute must be set when actions contain ISOLATE_")
		}
		cbiProfile := cbiProfileList[0].(map[string]interface{})
		if cbiProfile["id"] == "" && cbiProfile["name"] == "" && cbiProfile["url"] == "" {
			return fmt.Errorf("cbi_profile attribute must be properly set when actions contain ISOLATE_")
		}
	}

	// Validation 2: The attribute cascading_enabled can only be set when the actions value contain the string "DENY_" or "BLOCK_".
	if cascadingEnabled, ok := diff.GetOk("cascading_enabled"); ok && cascadingEnabled.(bool) && !containsDenyOrBlock {
		return fmt.Errorf("cascading_enabled can only be set when actions contain DENY_ or BLOCK_")
	}

	// Validation 3: When the "actions" value contains the string "ISOLATE_", the attribute user_agent_types must be set.
	if containsIsolate {
		if _, ok := diff.GetOk("user_agent_types"); !ok {
			return fmt.Errorf("user_agent_types attribute must be set when actions contain ISOLATE_")
		}
		userAgentTypes := diff.Get("user_agent_types").(*schema.Set).List()
		if len(userAgentTypes) == 0 {
			return fmt.Errorf("user_agent_types attribute must be set when actions contain ISOLATE_")
		}
		for _, userAgent := range userAgentTypes {
			if userAgent == "OTHER" {
				return fmt.Errorf("user_agent_types should not contain 'OTHER' when actions contain ISOLATE_. Valid options are: CHROME, FIREFOX, MSIE, MSEDGE, MSCHREDGE, OPERA, SAFARI")
			}
		}
	}

	return nil
}

func validateLocationManagementCountries() schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		name, ok := i.(string)
		if !ok {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Expected type to be string",
				Detail:   fmt.Sprintf("Expected string type for %s but got: %T", path, i),
			})
			return diags
		}

		// Special values
		if name == "NONE" {
			return diags
		}

		if country := countries.ByName(name); country == countries.Unknown {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid country name",
				Detail:   fmt.Sprintf("'%s' is not a valid country name. Please refer to ISO 3166-1 for a list of valid country names", name),
			})
		}

		return diags
	}
}

var supportedLocationManagemeTimeZones = []string{
	"NOT_SPECIFIED", "GMT_12_00_DATELINE", "GMT_11_00_SAMOA", "GMT_10_00_US_HAWAIIAN_TIME", "GMT_09_30_MARQUESAS", "GMT_09_00_US_ALASKA_TIME", "GMT_08_30_PITCARN", "GMT_08_00_PACIFIC_TIME", "GMT_07_00_US_MOUNTAIN_TIME", "GMT_07_00_US_MOUNTAIN_TIME_ARIZONA", "GMT_06_00_US_CENTRAL_TIME", "GMT_06_00_MEXICO", "GMT_05_00_US_EASTERN_TIME", "GMT_05_00_US_EASTERN_TIME_INDIANA", "GMT_05_00_COLUMBIA_PERU_SOUTH_AMERICA", "GMT_04_00_ATLANTIC_TIME", "GMT_03_30_NEWFOUNDLAND_CANADA", "GMT_03_00_ARGENTINA", "GMT_03_00_BRAZIL", "GMT_02_00_MID_ATLANTIC", "GMT_01_00_AZORES", "GMT", "GMT_01_00_WESTERN_EUROPE_GMT_01_00", "GMT_02_00_EASTERN_EUROPE_GMT_02_00", "GMT_02_00_EGYPT_GMT_02_00", "GMT_02_00_ISRAEL_GMT_02_00", "GMT_03_00_RUSSIA_GMT_03_00", "GMT_03_00_SAUDI_ARABIA_GMT_03_00", "GMT_03_30_IRAN_GMT_03_30", "GMT_04_00_ARABIAN_GMT_04_00", "GMT_04_30_AFGHANISTAN_GMT_04_30", "GMT_05_00_PAKISTAN_WEST_ASIA_GMT_05_00", "GMT_05_30_INDIA_GMT_05_30", "GMT_06_00_BANGLADESH_CENTRAL_ASIA_GMT_06_00", "GMT_06_30_BURMA_GMT_06_30", "GMT_07_00_BANGKOK_HANOI_JAKARTA_GMT_07_00", "GMT_08_00_CHINA_TAIWAN_GMT_08_00", "GMT_08_00_SINGAPORE_GMT_08_00", "GMT_08_00_AUSTRALIA_WT_GMT_08_00", "GMT_09_00_JAPAN_GMT_09_00", "GMT_09_00_KOREA_GMT_09_00", "GMT_09_30_AUSTRALIA_CT_GMT_09_30", "GMT_10_00_AUSTRALIA_ET_GMT_10_00", "GMT_10_30_AUSTRALIA_LORD_HOWE_GMT_10_30", "GMT_11_00_CENTRAL_PACIFIC_GMT_11_00", "GMT_11_30_NORFOLK_ISLANDS_GMT_11_30", "GMT_12_00_FIJI_NEW_ZEALAND_GMT_12_00", "AFGHANISTAN_ASIA_KABUL", "ALAND_ISLANDS_EUROPE_MARIEHAMN", "ALBANIA_EUROPE_TIRANE", "ALGERIA_AFRICA_ALGIERS", "AMERICAN_SAMOA_PACIFIC_PAGO_PAGO", "ANDORRA_EUROPE_ANDORRA", "ANGOLA_AFRICA_LUANDA", "ANGUILLA_AMERICA_ANGUILLA", "ANTARCTICA_CASEY", "ANTARCTICA_DAVIS", "ANTARCTICA_DUMONTDURVILLE", "ANTARCTICA_MAWSON", "ANTARCTICA_MCMURDO", "ANTARCTICA_PALMER", "ANTARCTICA_ROTHERA", "ANTARCTICA_SOUTH_POLE", "ANTARCTICA_SYOWA", "ANTARCTICA_VOSTOK", "ANTIGUA_AND_BARBUDA_AMERICA_ANTIGUA", "ARGENTINA_AMERICA_ARGENTINA_BUENOS_AIRES", "ARGENTINA_AMERICA_ARGENTINA_CATAMARCA", "ARGENTINA_AMERICA_ARGENTINA_CORDOBA", "ARGENTINA_AMERICA_ARGENTINA_JUJUY", "ARGENTINA_AMERICA_ARGENTINA_LA_RIOJA", "ARGENTINA_AMERICA_ARGENTINA_MENDOZA", "ARGENTINA_AMERICA_ARGENTINA_RIO_GALLEGOS", "ARGENTINA_AMERICA_ARGENTINA_SAN_JUAN", "ARGENTINA_AMERICA_ARGENTINA_TUCUMAN", "ARGENTINA_AMERICA_ARGENTINA_USHUAIA", "ARMENIA_ASIA_YEREVAN", "ARUBA_AMERICA_ARUBA", "AUSTRALIA_ADELAIDE", "AUSTRALIA_BRISBANE", "AUSTRALIA_BROKEN_HILL", "AUSTRALIA_CURRIE", "AUSTRALIA_DARWIN", "AUSTRALIA_EUCLA", "AUSTRALIA_HOBART", "AUSTRALIA_LINDEMAN", "AUSTRALIA_LORD_HOWE", "AUSTRALIA_MELBOURNE", "AUSTRALIA_PERTH", "AUSTRALIA_SYDNEY", "AUSTRIA_EUROPE_VIENNA", "AZERBAIJAN_ASIA_BAKU", "BAHAMAS_AMERICA_NASSAU", "BAHRAIN_ASIA_BAHRAIN", "BANGLADESH_ASIA_DHAKA", "BARBADOS_AMERICA_BARBADOS", "BELARUS_EUROPE_MINSK", "BELGIUM_EUROPE_BRUSSELS", "BELIZE_AMERICA_BELIZE", "BENIN_AFRICA_PORTO_NOVO", "BERMUDA_ATLANTIC_BERMUDA", "BHUTAN_ASIA_THIMPHU", "BOLIVIA_AMERICA_LA_PAZ", "BOSNIA_AND_HERZEGOVINA_EUROPE_SARAJEVO", "BOTSWANA_AFRICA_GABORONE", "BRAZIL_AMERICA_ARAGUAINA", "BRAZIL_AMERICA_BAHIA", "BRAZIL_AMERICA_BELEM", "BRAZIL_AMERICA_BOA_VISTA", "BRAZIL_AMERICA_CAMPO_GRANDE", "BRAZIL_AMERICA_CUIABA", "BRAZIL_AMERICA_EIRUNEPE", "BRAZIL_AMERICA_FORTALEZA", "BRAZIL_AMERICA_MACEIO", "BRAZIL_AMERICA_MANAUS", "BRAZIL_AMERICA_NORONHA", "BRAZIL_AMERICA_PORTO_VELHO", "BRAZIL_AMERICA_RECIFE", "BRAZIL_AMERICA_RIO_BRANCO", "BRAZIL_AMERICA_SAO_PAULO", "BRITISH_INDIAN_OCEAN_TERRITORY_INDIAN_CHAGOS", "BRUNEI_DARUSSALAM_ASIA_BRUNEI", "BULGARIA_EUROPE_SOFIA", "BURKINA_FASO_AFRICA_OUAGADOUGOU", "BURUNDI_AFRICA_BUJUMBURA", "CAMBODIA_ASIA_PHNOM_PENH", "CAMEROON_AFRICA_DOUALA", "CANADA_AMERICA_ATIKOKAN", "CANADA_AMERICA_BLANC_SABLON", "CANADA_AMERICA_CAMBRIDGE_BAY", "CANADA_AMERICA_DAWSON_CREEK", "CANADA_AMERICA_DAWSON", "CANADA_AMERICA_EDMONTON", "CANADA_AMERICA_GLACE_BAY", "CANADA_AMERICA_GOOSE_BAY", "CANADA_AMERICA_HALIFAX", "CANADA_AMERICA_INUVIK", "CANADA_AMERICA_IQALUIT", "CANADA_AMERICA_MONCTON", "CANADA_AMERICA_MONTREAL", "CANADA_AMERICA_NIPIGON", "CANADA_AMERICA_PANGNIRTUNG", "CANADA_AMERICA_RAINY_RIVER", "CANADA_AMERICA_RANKIN_INLET", "CANADA_AMERICA_REGINA", "CANADA_AMERICA_RESOLUTE", "CANADA_AMERICA_ST_JOHNS", "CANADA_AMERICA_SWIFT_CURRENT", "CANADA_AMERICA_THUNDER_BAY", "CANADA_AMERICA_TORONTO", "CANADA_AMERICA_VANCOUVER", "CANADA_AMERICA_WHITEHORSE", "CANADA_AMERICA_WINNIPEG", "CANADA_AMERICA_YELLOWKNIFE", "CAPE_VERDE_ATLANTIC_CAPE_VERDE", "CAYMAN_ISLANDS_AMERICA_CAYMAN", "CENTRAL_AFRICAN_REPUBLIC_AFRICA_BANGUI", "CHAD_AFRICA_NDJAMENA", "CHILE_AMERICA_SANTIAGO", "CHILE_PACIFIC_EASTER", "CHINA_ASIA_CHONGQING", "CHINA_ASIA_HARBIN", "CHINA_ASIA_KASHGAR", "CHINA_ASIA_SHANGHAI", "CHINA_ASIA_URUMQI", "CHRISTMAS_ISLAND_INDIAN_CHRISTMAS", "COCOS_KEELING_ISLANDS_INDIAN_COCOS", "COLOMBIA_AMERICA_BOGOTA", "COMOROS_INDIAN_COMORO", "DEMOCRATIC_REPUBLIC_OF_CONGO_CONGO_KINSHASA_AFRICA_KINSHASA", "DEMOCRATIC_REPUBLIC_OF_CONGO_CONGO_KINSHASA_AFRICA_LUBUMBASHI", "CONGO_CONGO_BRAZZAVILLE_AFRICA_BRAZZAVILLE", "COOK_ISLANDS_PACIFIC_RAROTONGA", "COSTA_RICA_AMERICA_COSTA_RICA", "COTE_DIVOIRE_AFRICA_ABIDJAN", "CROATIA_EUROPE_ZAGREB", "CUBA_AMERICA_HAVANA", "CYPRUS_ASIA_NICOSIA", "CZECH_REPUBLIC_EUROPE_PRAGUE", "DENMARK_EUROPE_COPENHAGEN", "DJIBOUTI_AFRICA_DJIBOUTI", "DOMINICAN_REPUBLIC_AMERICA_SANTO_DOMINGO", "DOMINICA_AMERICA_DOMINICA", "ECUADOR_AMERICA_GUAYAQUIL", "ECUADOR_PACIFIC_GALAPAGOS", "EGYPT_AFRICA_CAIRO", "EL_SALVADOR_AMERICA_EL_SALVADOR", "EQUATORIAL_GUINEA_AFRICA_MALABO", "ERITREA_AFRICA_ASMARA", "ESTONIA_EUROPE_TALLINN", "ETHIOPIA_AFRICA_ADDIS_ABABA", "FALKLAND_ISLANDS_ATLANTIC_STANLEY", "FAROE_ISLANDS_ATLANTIC_FAROE", "FIJI_PACIFIC_FIJI", "FINLAND_EUROPE_HELSINKI", "FRANCE_EUROPE_PARIS", "FRENCH_GUIANA_AMERICA_CAYENNE", "FRENCH_POLYNESIA_PACIFIC_GAMBIER", "FRENCH_POLYNESIA_PACIFIC_MARQUESAS", "FRENCH_POLYNESIA_PACIFIC_TAHITI", "FRENCH_SOUTHERN_TERRITORIES_INDIAN_KERGUELEN", "GABON_AFRICA_LIBREVILLE", "GAMBIA_AFRICA_BANJUL", "GEORGIA_ASIA_TBILISI", "GERMANY_EUROPE_BERLIN", "GHANA_AFRICA_ACCRA", "GIBRALTAR_EUROPE_GIBRALTAR", "GREECE_EUROPE_ATHENS", "GREENLAND_AMERICA_DANMARKSHAVN", "GREENLAND_AMERICA_GODTHAB", "GREENLAND_AMERICA_SCORESBYSUND", "GREENLAND_AMERICA_THULE", "GRENADA_AMERICA_GRENADA", "GUADELOUPE_AMERICA_GUADELOUPE", "GUAM_PACIFIC_GUAM", "GUATEMALA_AMERICA_GUATEMALA", "GUERNSEY_EUROPE_GUERNSEY", "GUINEA_BISSAU_AFRICA_BISSAU", "GUINEA_AFRICA_CONAKRY", "GUYANA_AMERICA_GUYANA", "HAITI_AMERICA_PORT_AU_PRINCE", "HONDURAS_AMERICA_TEGUCIGALPA", "HONG_KONG_ASIA_HONG_KONG", "HUNGARY_EUROPE_BUDAPEST", "ICELAND_ATLANTIC_REYKJAVIK", "INDIA_ASIA_KOLKATA", "INDONESIA_ASIA_JAKARTA", "INDONESIA_ASIA_JAYAPURA", "INDONESIA_ASIA_MAKASSAR", "INDONESIA_ASIA_PONTIANAK", "IRAN_ASIA_TEHRAN", "IRAQ_ASIA_BAGHDAD", "IRELAND_EUROPE_DUBLIN", "ISLE_OF_MAN_EUROPE_ISLE_OF_MAN", "ISRAEL_ASIA_JERUSALEM", "ITALY_EUROPE_ROME", "JAMAICA_AMERICA_JAMAICA", "JAPAN_ASIA_TOKYO", "JERSEY_EUROPE_JERSEY", "JORDAN_ASIA_AMMAN", "KAZAKHSTAN_ASIA_ALMATY", "KAZAKHSTAN_ASIA_AQTAU", "KAZAKHSTAN_ASIA_AQTOBE", "KAZAKHSTAN_ASIA_ORAL", "KAZAKHSTAN_ASIA_QYZYLORDA", "KENYA_AFRICA_NAIROBI", "KIRIBATI_PACIFIC_ENDERBURY", "KIRIBATI_PACIFIC_KIRITIMATI", "KIRIBATI_PACIFIC_TARAWA", "KOREA_DEMOCRATIC_PEOPLES_REPUBLIC_OF_ASIA_PYONGYANG", "KOREA_REPUBLIC_OF_ASIA_SEOUL", "KUWAIT_ASIA_KUWAIT", "KYRGYZSTAN_ASIA_BISHKEK", "LAO_PEOPLES_DEMOCRATIC_REPUBLIC_ASIA_VIENTIANE", "LATVIA_EUROPE_RIGA", "LEBANON_ASIA_BEIRUT", "LESOTHO_AFRICA_MASERU", "LIBERIA_AFRICA_MONROVIA", "LIBYAN_ARAB_JAMAHIRIYA_AFRICA_TRIPOLI", "LIECHTENSTEIN_EUROPE_VADUZ", "LITHUANIA_EUROPE_VILNIUS", "LUXEMBOURG_EUROPE_LUXEMBOURG", "MACAO_ASIA_MACAU", "MACEDONIA_THE_FORMER_YUGOSLAV_REPUBLIC_OF_EUROPE_SKOPJE", "MADAGASCAR_INDIAN_ANTANANARIVO", "MALAWI_AFRICA_BLANTYRE", "MALAYSIA_ASIA_KUALA_LUMPUR", "MALAYSIA_ASIA_KUCHING", "MALDIVES_INDIAN_MALDIVES", "MALI_AFRICA_BAMAKO", "MALTA_EUROPE_MALTA", "MARSHALL_ISLANDS_PACIFIC_KWAJALEIN", "MARSHALL_ISLANDS_PACIFIC_MAJURO", "MARTINIQUE_AMERICA_MARTINIQUE", "MAURITANIA_AFRICA_NOUAKCHOTT", "MAURITIUS_INDIAN_MAURITIUS", "MAYOTTE_INDIAN_MAYOTTE", "MEXICO_AMERICA_CANCUN", "MEXICO_AMERICA_CHIHUAHUA", "MEXICO_AMERICA_HERMOSILLO", "MEXICO_AMERICA_MAZATLAN", "MEXICO_AMERICA_MERIDA", "MEXICO_AMERICA_MEXICO_CITY", "MEXICO_AMERICA_MONTERREY", "MEXICO_AMERICA_TIJUANA", "MICRONESIA_FEDERATED_STATES_OF_PACIFIC_KOSRAE", "MICRONESIA_FEDERATED_STATES_OF_PACIFIC_PONAPE", "MICRONESIA_FEDERATED_STATES_OF_PACIFIC_TRUK", "MOLDOVA_EUROPE_CHISINAU", "MONACO_EUROPE_MONACO", "MONGOLIA_ASIA_CHOIBALSAN", "MONGOLIA_ASIA_HOVD", "MONGOLIA_ASIA_ULAANBAATAR", "MONTENEGRO_EUROPE_PODGORICA", "MONTSERRAT_AMERICA_MONTSERRAT", "MOROCCO_AFRICA_CASABLANCA", "MOZAMBIQUE_AFRICA_MAPUTO", "MYANMAR_ASIA_RANGOON", "NAMIBIA_AFRICA_WINDHOEK", "NAURU_PACIFIC_NAURU", "NEPAL_ASIA_KATMANDU", "NETHERLANDS_ANTILLES_AMERICA_CURACAO", "NETHERLANDS_EUROPE_AMSTERDAM", "NEW_CALEDONIA_PACIFIC_NOUMEA", "NEW_ZEALAND_PACIFIC_AUCKLAND", "NEW_ZEALAND_PACIFIC_CHATHAM", "NICARAGUA_AMERICA_MANAGUA", "NIGERIA_AFRICA_LAGOS", "NIGER_AFRICA_NIAMEY", "NIUE_PACIFIC_NIUE", "NORFOLK_ISLAND_PACIFIC_NORFOLK", "NORTHERN_MARIANA_ISLANDS_PACIFIC_SAIPAN", "NORWAY_EUROPE_OSLO", "OCCUPIED_PALESTINIAN_TERRITORY_ASIA_GAZA", "OMAN_ASIA_MUSCAT", "PAKISTAN_ASIA_KARACHI", "PALAU_PACIFIC_PALAU", "PANAMA_AMERICA_PANAMA", "PAPUA_NEW_GUINEA_PACIFIC_PORT_MORESBY", "PARAGUAY_AMERICA_ASUNCION", "PERU_AMERICA_LIMA", "PHILIPPINES_ASIA_MANILA", "PITCAIRN_PACIFIC_PITCAIRN", "POLAND_EUROPE_WARSAW", "PORTUGAL_ATLANTIC_AZORES", "PORTUGAL_ATLANTIC_MADEIRA", "PORTUGAL_EUROPE_LISBON", "PUERTO_RICO_AMERICA_PUERTO_RICO", "QATAR_ASIA_QATAR", "REUNION_INDIAN_REUNION", "ROMANIA_EUROPE_BUCHAREST", "RUSSIAN_FEDERATION_ASIA_ANADYR", "RUSSIAN_FEDERATION_ASIA_IRKUTSK", "RUSSIAN_FEDERATION_ASIA_KAMCHATKA", "RUSSIAN_FEDERATION_ASIA_KRASNOYARSK", "RUSSIAN_FEDERATION_ASIA_MAGADAN", "RUSSIAN_FEDERATION_ASIA_NOVOSIBIRSK", "RUSSIAN_FEDERATION_ASIA_OMSK", "RUSSIAN_FEDERATION_ASIA_SAKHALIN", "RUSSIAN_FEDERATION_ASIA_VLADIVOSTOK", "RUSSIAN_FEDERATION_ASIA_YAKUTSK", "RUSSIAN_FEDERATION_ASIA_YEKATERINBURG", "RUSSIAN_FEDERATION_EUROPE_KALININGRAD", "RUSSIAN_FEDERATION_EUROPE_MOSCOW", "RUSSIAN_FEDERATION_EUROPE_SAMARA", "RUSSIAN_FEDERATION_EUROPE_VOLGOGRAD", "RWANDA_AFRICA_KIGALI", "SAINT_BARTHELEMY_AMERICA_ST_BARTHELEMY", "SAINT_KITTS_AND_NEVIS_AMERICA_ST_KITTS", "SAINT_LUCIA_AMERICA_ST_LUCIA", "SAINT_MARTIN_FRENCH_PART_AMERICA_MARIGOT", "SAINT_VINCENT_AND_THE_GRENADINES_AMERICA_ST_VINCENT", "SAMOA_PACIFIC_APIA", "SAN_MARINO_EUROPE_SAN_MARINO", "SAO_TOME_AND_PRINCIPE_AFRICA_SAO_TOME", "SAUDI_ARABIA_ASIA_RIYADH", "SENEGAL_AFRICA_DAKAR", "SERBIA_EUROPE_BELGRADE", "SEYCHELLES_INDIAN_MAHE", "SIERRA_LEONE_AFRICA_FREETOWN", "SINGAPORE_ASIA_SINGAPORE", "SLOVAKIA_EUROPE_BRATISLAVA", "SLOVENIA_EUROPE_LJUBLJANA", "SOLOMON_ISLANDS_PACIFIC_GUADALCANAL", "SOMALIA_AFRICA_MOGADISHU", "SOUTH_AFRICA_AFRICA_JOHANNESBURG", "SOUTH_GEORGIA_AND_THE_SOUTH_SANDWICH_ISLANDS_ATLANTIC_SOUTH_GEORGIA", "SPAIN_AFRICA_CEUTA", "SPAIN_ATLANTIC_CANARY", "SPAIN_EUROPE_MADRID", "SRI_LANKA_ASIA_COLOMBO", "ST_HELENA_ATLANTIC_ST_HELENA", "ST_PIERRE_AND_MIQUELON_AMERICA_MIQUELON", "SUDAN_AFRICA_KHARTOUM", "SURINAME_AMERICA_PARAMARIBO", "SVALBARD_AND_JAN_MAYEN_ISLANDS_ARCTIC_LONGYEARBYEN", "SWAZILAND_AFRICA_MBABANE", "SWEDEN_EUROPE_STOCKHOLM", "SWITZERLAND_EUROPE_ZURICH", "SYRIAN_ARAB_REPUBLIC_ASIA_DAMASCUS", "TAIWAN_ASIA_TAIPEI", "TAJIKISTAN_ASIA_DUSHANBE", "TANZANIA_AFRICA_DAR_ES_SALAAM", "THAILAND_ASIA_BANGKOK", "TIMOR_LESTE_ASIA_DILI", "TOGO_AFRICA_LOME", "TOKELAU_PACIFIC_FAKAOFO", "TONGA_PACIFIC_TONGATAPU", "TRINIDAD_AND_TOBAGO_AMERICA_PORT_OF_SPAIN", "TUNISIA_AFRICA_TUNIS", "TURKEY_EUROPE_ISTANBUL", "TURKMENISTAN_ASIA_ASHGABAT", "TURKS_AND_CAICOS_ISLANDS_AMERICA_GRAND_TURK", "TUVALU_PACIFIC_FUNAFUTI", "UGANDA_AFRICA_KAMPALA", "UKRAINE_EUROPE_KIEV", "UKRAINE_EUROPE_SIMFEROPOL", "UKRAINE_EUROPE_UZHGOROD", "UKRAINE_EUROPE_ZAPOROZHYE", "UNITED_ARAB_EMIRATES_ASIA_DUBAI", "UNITED_KINGDOM_EUROPE_LONDON", "UNITED_STATES_MINOR_OUTLYING_ISLANDS_PACIFIC_JOHNSTON", "UNITED_STATES_MINOR_OUTLYING_ISLANDS_PACIFIC_MIDWAY", "UNITED_STATES_MINOR_OUTLYING_ISLANDS_PACIFIC_WAKE", "UNITED_STATES_AMERICA_ADAK", "UNITED_STATES_AMERICA_ANCHORAGE", "UNITED_STATES_AMERICA_BOISE", "UNITED_STATES_AMERICA_CHICAGO", "UNITED_STATES_AMERICA_DENVER", "UNITED_STATES_AMERICA_DETROIT", "UNITED_STATES_AMERICA_INDIANA_INDIANAPOLIS", "UNITED_STATES_AMERICA_INDIANA_KNOX", "UNITED_STATES_AMERICA_INDIANA_MARENGO", "UNITED_STATES_AMERICA_INDIANA_PETERSBURG", "UNITED_STATES_AMERICA_INDIANA_TELL_CITY", "UNITED_STATES_AMERICA_INDIANA_VEVAY", "UNITED_STATES_AMERICA_INDIANA_VINCENNES", "UNITED_STATES_AMERICA_INDIANA_WINAMAC", "UNITED_STATES_AMERICA_JUNEAU", "UNITED_STATES_AMERICA_KENTUCKY_LOUISVILLE", "UNITED_STATES_AMERICA_KENTUCKY_MONTICELLO", "UNITED_STATES_AMERICA_LOS_ANGELES", "UNITED_STATES_AMERICA_MENOMINEE", "UNITED_STATES_AMERICA_NEW_YORK", "UNITED_STATES_AMERICA_NOME", "UNITED_STATES_AMERICA_NORTH_DAKOTA_CENTER", "UNITED_STATES_AMERICA_NORTH_DAKOTA_NEW_SALEM", "UNITED_STATES_AMERICA_PHOENIX", "UNITED_STATES_AMERICA_SHIPROCK", "UNITED_STATES_AMERICA_YAKUTAT", "UNITED_STATES_PACIFIC_HONOLULU", "URUGUAY_AMERICA_MONTEVIDEO", "UZBEKISTAN_ASIA_SAMARKAND", "UZBEKISTAN_ASIA_TASHKENT", "VANUATU_PACIFIC_EFATE", "VATICAN_CITY_STATE_EUROPE_VATICAN", "VENEZUELA_AMERICA_CARACAS", "VIET_NAM_ASIA_SAIGON", "VIRGIN_ISLANDS_BRITISH_AMERICA_TORTOLA", "VIRGIN_ISLANDS_US_AMERICA_ST_THOMAS", "WALLIS_AND_FUTUNA_ISLANDS_PACIFIC_WALLIS", "WESTERN_SAHARA_AFRICA_EL_AAIUN", "YEMEN_ASIA_ADEN", "ZAMBIA_AFRICA_LUSAKA", "ZIMBABWE_AFRICA_HARARE",
}

func validateLocationManagementTimeZones() schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		value, ok := i.(string)
		if !ok {
			return diag.Diagnostics{
				{
					Severity: diag.Error,
					Summary:  "Expected type to be string",
					Detail:   "Type assertion failed, expected string type for Location Management Timezones validation",
				},
			}
		}

		// Convert the cty.Path to a string representation
		pathStr := fmt.Sprintf("%+v", path)

		// Use StringInSlice from helper/validation package
		var diags diag.Diagnostics
		if _, errs := validation.StringInSlice(supportedLocationManagemeTimeZones, false)(value, pathStr); len(errs) > 0 {
			for _, err := range errs {
				diags = append(diags, diag.FromErr(err)...)
			}
		}

		return diags
	}
}

// Validate Cloud Firewall Network Service Applications

func validateISOCountryCodes(value interface{}, key string) ([]string, []error) {
	var warnings []string
	var errors []error

	code, ok := value.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", key))
		return warnings, errors
	}

	// Special values
	if code == "ANY" || code == "NONE" {
		return warnings, errors
	}

	// Check against ISO 3166 Alpha-2
	if !iso3166.ExistsIso3166ByAlpha2Code(code) {
		errors = append(errors, fmt.Errorf("'%s' is not a valid ISO-3166 Alpha-2 country code. Please visit the following site for reference: https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes", code))
		return warnings, errors
	}

	return warnings, errors
}

var supportedCloudFirewallNwApplications = []string{
	"NOT_AVAILABLE", "APNS", "APPSTORE", "DICT", "EPM", "GARP", "ICLOUD", "IOS_OTA_UPDATE", "LDAP", "LDAPS", "MQ", "NSPI", "PERFORCE", "PORTMAP", "SAMSUNG_APPS", "SRVLOC", "SSDP", "SYSLOG", "WINDOWS_MARKETPLACE", "XFS", "CHAP", "DIAMETER", "IDENT", "KRB5", "PAP", "RADIUS", "TACACS_PLUS", "YPPASSWD", "YPSERV", "SALESFORCE", "GOOGLEANALYTICS", "OFFICE365", "EVERNOTE", "TUMBLR", "ZOHOCRM", "SUGARCRM", "ZSCLOUD_HTTP_BYPASS", "DB2", "DRDA", "FILEMAKER_PRO", "INFORMIX", "MOBILINK", "MYSQL", "POSTGRES", "SQLI", "SYBASE", "TDS", "TNS", "GOTOMEETING", "MICROSOFTLIVEMEETING", "YAMMER", "ACROBATCONNECT", "NETVIEWER", "AMAZON_CLD_DRIVE", "WEBINAR", "AIMINI", "CFT", "FTP", "FTP_DATA", "FTPS", "GOOGLE_DRIVE", "HOTLINE", "IBACKUP", "MOUNT", "NETBIOS", "NFS", "NLOCKMGR", "RQUOTA", "RSTAT", "RSYNC", "RUSERS", "SMB", "SYNC", "TFTP", "YPUPDATE", "AIM_TRANSFER", "IRC_TRANSFER", "JABBER_TRANSFER", "PALTALK_TRANSFER", "QQ_TRANSFER", "YMSG_TRANSFER", "FTPS_DATA", "CSTRIKE", "EVE_ONLINE", "EVERQUEST", "HALFLIFE", "IMVU", "LINEAGE2", "POKER_STARS", "PSN", "QUAKE", "RUNESCAPE", "STEAM", "WFC", "WIICONNECT24", "WOW", "XBOXLIVE", "FLICKR", "NOTAPPLICABLE", "TELEGRAM", "ZOOM", "MIBBIT", "AIM", "AIM_EXPRESS", "AIMS", "AIRAIM", "BADOO", "EBUDDY", "FRIENDVOX", "GTALK", "GCM", "HOVRS", "ICQ2GO", "ILOVEIM", "IRC", "JABBER", "KAIXIN_CHAT", "KAKAOTALK", "KIK", "KOOLIM", "LINE", "LOTUS_SAMETIME", "MESSENGERFX", "MITALK", "MPLUS_MESSENGER", "MS_COMMUNICATOR", "MSN", "MSNMOBILE", "MXIT", "NETMEETING_ILS", "OICQ", "OOVOO", "PALTALK", "QQ", "RADIUSIM", "SKYPE", "TEAMSPEAK", "TEAMSPEAK_V3", "WECHAT", "WHATSAPP", "YMSG", "YMSG_CONF", "YOUNI", "ZOHO", "ZOHO_IM", "IOS_APPSTORE", "ANZHI", "GOOGLE_PLAY", "ALTIRIS", "CDP", "IPERF", "LCP", "NETFLOW", "RSVP", "SNMP", "SONMP", "ZABBIX_AGENT", "APOLLO", "BGP", "EIGRP", "HSRP", "IPXRIP", "LOOP", "MPLS", "NLSP", "OSPF", "PWE", "RIP1", "RIP2", "RIPNG1", "STP", "VRRP", "IP", "IP6", "TCP", "UDP", "ICMP", "DHCP", "NTP", "DNS", "TCP_OVER_DNS", "DNS_OVER_HTTPS", "OFW_HTTP_BYPASS", "OFW_HTTPS_BYPASS", "OFW_TCP_BYPASS", "OFW_UDP_BYPASS", "OFW_ICMP_BYPASS", "ADC", "APPLEJUICE", "ARES", "BITTORRENT", "BLUBSTER", "CLIP2NET", "DIINO", "DIRECTCONNECT", "EDONKEY", "FILETOPIA", "FOXY", "GNUNET", "GNUTELLA", "GOBOOGY", "IMESH", "KAZAA", "KUGOU", "LIVE_MESH", "MANOLITO", "MUTE", "OPENFT", "PANDO", "QQMUSIC", "QQSTREAM", "SHARE", "SLSK", "SORIBADA", "THUNDER", "UTP", "VAKAKA", "WINMX", "WINNY", "HTTPTUNNEL", "CLEARCASE", "GOTODEVICE", "GOTOMYPC", "ICA", "JEDI", "RADMIN", "RDP", "RFB", "SHOWMYPC", "TEAMVIEWER", "VMWARE", "X11", "XDMCP", "RLOGIN", "RSH", "TELNET", "TNVIP", "DCERPC", "DIOP", "GIOP", "GIOPS", "IIOP", "MSRPC", "RMI_IIOP", "RPC", "SOAP", "PINTEREST", "DIGG", "WORDPRESS", "GOOGLE_GROUPS", "IRCS", "LINKEDIN", "NNTP", "NNTPS", "ODNOKLASSNIKI", "VKONTAKTE", "YAHOO_GROUPS", "XING", "LIVEJOURNAL", "REDDIT", "GOOGLE_PLUS", "MIXI", "TWITTER", "BLOGSPOT", "HI5", "FRIENDSTER", "BEBO", "FACEBOOK", "MYSPACE", "ORKUT", "CRAIGSLIST", "GOOGLE_VIDEO", "WEBEX", "BOXNET", "APP_050PLUS", "ADOBE_CONNECT", "BABELGUM", "BAIDU_PLAYER", "BAOFENG", "BBC_PLAYER", "BLIP_TV", "BLOCKBUSTER", "BMFF", "CCTV_VOD", "CNET", "CNET_TV", "COMM", "FACETIME", "FLASH", "FRING", "GROOVESHARK", "H225", "H245", "H248_BINARY", "H248_TEXT", "HBO_GO", "HOWCAST", "HULU", "IAX", "ICALL", "ICECAST", "IHEARTRADIO", "ITUNES", "JAJAH", "MEETINGPLACE", "METACAFE", "MGCP", "MMS", "MOG", "MOGULUS", "MPEGTS", "MSN_VIDEO", "MSRP", "NETFLIX", "NICONICO_DOUGA", "PALTALK_AUDIO", "PALTALK_VIDEO", "PPLIVE", "PPSTREAM", "Q931", "QIK_VIDEO", "QQLIVE", "QVOD", "RDT", "RHAPSODY", "RTMP", "RTCP", "RTP", "RTSP", "SCCP", "SHOUTCAST", "SILVERLIGHT", "SIP", "SKY", "SKY_PLAYER", "SLACKER", "SLINGBOX", "SOPCAST", "SOUNDCLOUD", "SPOTIFY", "TANGO", "TU", "TUNEIN", "TVANTS", "TVUPLAYER", "UUSEE", "VEOHTV", "VEVO", "VIBER", "VODDLER", "YAHOO_SCREEN", "YMSG_VIDEO", "ZATTOO", "DAILYMOTION", "AOL_VIDEO", "ZIDDU", "LASTFM", "DEEZER", "MYVIDEODE", "PANDORA", "YOUTUBE", "VIMEO", "DROPBOX", "HIGHTAIL", "FILESTUBE", "RINGCENTRAL", "SHAREFILE", "THREEGPP_LI", "CAPWAP", "ETHERIP", "ETSI_LI", "GRE", "GTP", "GTPV2", "IPV6CP", "L2TP", "LQR", "MOBILE_IP", "OPENVPN", "PPP", "PPPOE", "PPTP", "SOCKS2HTTP", "SOCKS4", "SOCKS5", "TEREDO", "ULTRASURF", "VJC_COMP", "VJC_UNCOMP", "XOT", "IPSEC", "ISAKMP", "OCSP", "SSH", "HTTP_PROXY", "TOR", "WSTUNNEL", "SSL", "ABCNEWS", "ACCUWEATHER", "ACER", "ACROBAT", "ACTIVESYNC", "ADDICTINGGAMES", "ADNSTREAM", "ADOBE_MEETING_RC", "ADOBE_UPDATE", "ADRIVE", "ADULTADWORLD", "ADULTFRIENDFINDER", "ADVOGATO", "AGAME", "AIAIGAME", "AILI", "AIZHAN", "AKAMAI", "ALEXA", "ALIBABA", "ALIMAMA", "ALIPAY", "ALJAZEERA", "ALLOCINE", "AMAZON", "AMAZON_AWS", "AMAZON_MP3", "AMAZON_VIDEO", "AMEBA", "AMERICANEXPRESS", "AMIE_STREET", "ANOBII", "ANSWERS", "ZYNGA", "ZUM", "APPLE", "APPLE_AIRPORT", "APPLEDAILY", "APPLE_MAPS", "APPLE_SIRI", "APPLE_UPDATE", "APPSHOPPER", "ARCHIVE", "ARMORGAMES", "ASIAE", "ZSHARE", "ASMALLWORLD", "ATHLINKS", "ATT", "ATWIKI", "AUCTION", "AUFEMININ", "AUONE", "AVATARS_UNITED", "AVG_UPDATE", "AVIRA_UPDATE", "AVOIDR", "BABYCENTER", "BABYHOME", "BACKPACKERS", "BADONGO", "ZOO", "BAIKE", "BEANFUN", "BERNIAGA", "BIGADDA", "BIGLOBE_NE", "BIGTENT", "BIGUPLOAD", "BIIP", "ZOL", "BITDEFENDER_UPDATE", "BLACKBERRY", "BLACKPLANET", "BLOGDETIK", "BLOGGER", "BLOGIMG", "BLOGSTER", "BLOKUS", "BLOOMBERG", "BLUEJAYFILMS", "BOLT", "BONPOO", "BRIGHTTALK", "BUGS", "BUSINESSWEEK", "BUSINESSWEEKLY", "BUZZFEED", "BUZZNET", "BYPASSTHAT", "CAFEMOM", "CAM4", "CAMPFIRE", "CAMZAP", "CAPITALONE", "CARE2", "CARTOONNETWORK", "CDISCOUNT", "CELLUFUN", "CHANNEL4", "CHINA_AIRLINES", "CHINACOM", "CHINACOMCN", "CHINANEWS", "CHINATIMES", "CHINAZ", "CHOSUN", "CHOSUN_DAILY", "CHROME_UPDATE", "CITRIX_ONLINE", "CJ_MALL", "CK101", "CLASSMATES", "CLOOB", "CLOUDFLARE", "CLOUDME", "CLUBIC", "CNN", "CNTV", "CNYES", "CNZZ", "COCOLOG_NIFTY", "COLLEGE_BLENDER", "COMCAST", "CONCUR", "CONDUIT", "COUCH_SURFING", "COUPANG", "CROCKO", "CSDN", "CTRIP", "DAILYMAIL", "YUUGUU", "DAILY_STRENGH", "DANGDANG", "DAUM", "DAVIDOV", "DEBIAN_UPDATE", "DECAYENNE", "DELICIOUS", "DEPOSITFILES", "DETIK", "DETIKNEWS", "DEVIANT_ART", "DIGITALVERSE", "DIRECTDOWNLOADLINKS", "DIRECTV", "DISABOOM", "DIVSHARE", "DMM_CO", "DNSHOP", "DOL2DAY", "DONGA", "DONTSTAYIN", "DOORBLOG", "DOUBAN", "DOUBLECLICK_ADS", "DRAUGIEM", "DREAMWIZ", "DRUPAL", "DUOWAN", "DYNAMICINTRANET", "EARTHCAM", "EASTMONEY", "EASYTRAVEL", "EBAY", "ELFTOWN", "ELLE_TW", "EONS", "EPERNICUS", "EROOM_NET", "ESNIPS", "ESPN", "ETAO", "ETTODAY", "ZOHO_SHOW", "EVONY", "EXBLOG", "EXPEDIA", "EXPERIENCE_PROJECT", "EXPLOROO", "EYEJOT", "EYNY", "EZFLY", "EZTRAVEL", "FACEBOOK_APPS", "FACEPARTY", "FACES", "FASHIONGUIDE", "FC2", "FETLIFE", "FILE_DROPPER", "FILEFLYER", "FILE_HOST", "FILER_CX", "YOUSEEMORE", "FILLOS_DE_GALICIA", "FILMAFFINITY", "FIREFOX_UPDATE", "FLASHPLUGIN_UPDATE", "FLEDGEWING", "ZOHO_SHEET", "FLIXSTER", "FLUMOTION", "FLUXIOM", "FLY_PROXY", "FOGBUGZ", "FORTUNECHINA", "FOTKI", "FOTOLOG", "FOURSQUARE", "FOXMOVIES", "FOXNEWS", "FOXSPORTS", "FREEBSD_UPDATE", "FREEETV", "FRIENDS_REUNITED", "FRUHSTUCKSTREFF", "FSECURE_UPDATE", "FUBAR", "FUNSHION", "GAIAONLINE", "GAMEBASE_TW", "GAMERDNA", "GAMER_TW", "GAMES_CO", "GAMESMOMO", "GANJI", "GATHER", "GAYS", "GENI", "GFAN", "GIGAUP", "GLIDE", "GMARKET", "GOGOYOKO", "GOHAPPY", "GOMTV_VOD", "GOODREADS", "ZOHO_SHARE", "GOOGLE_ADS", "GOOGLE_APPENGINE", "GOOGLE_CACHE", "GOOGLE_DESKTOP", "GOOGLE_DOCS", "GOOGLE_EARTH", "GOOGLE_GEN", "GOOGLE_MAPS", "GOOGLE_PICASA", "GOOGLE_SKYMAP", "GOOGLE_TOOLBAR", "GOOGLE_TRANSLATE", "GOO_NE", "GOUGOU", "GRATISINDO", "GREE", "GRONO", "GROUPON", "GROUPWISE", "GSSHOP", "GSTATIC", "GUDANGLAGU", "GYAO", "HABBO", "HANGAME", "HANKOOKI", "HANKYUNG", "HAO123", "HARDSEXTUBE", "HATENA_NE", "HERALDM", "HERE", "HEXUN", "HGTV", "HINET_GAMES", "HOSPITALITY_CLUB", "HOTFILE", "HOWSTUFFWORKS", "HTTP", "HTTPS", "HTTP2", "HUDONG", "HYVES", "IAPP", "IBIBO", "ICAP", "IFENG", "IFENG_FINANCE", "IFILE_IT", "I_GAMER", "IKEA", "IMAGESHACK", "IMDB", "IMEEM", "IMEET", "IMGUR", "IMPRESS", "INDABA_MUSIC", "INDONETWORK", "INDOWEBSTER", "INILAH", "INSTAGRAM", "INTALKING", "INTERNATIONS", "INTERPARK", "INTUIT", "I_PART", "IQIYI", "IRC_GALLERIA", "ITALKI", "ITSMY", "IWIW", "JAIKU", "JAMMERDIRECT", "JANGO", "JAVA_UPDATE", "JINGDONG", "JNE", "JOBSTREET", "JOONGANG_DAILY", "JUBII", "JUSTIN_TV", "KAIOO", "KAKAKU", "KANKAN", "KAPANLAGI", "KAROSGAME", "KASKUS", "KASPERSKY", "KASPERSKY_UPDATE", "KBS", "KEEZMOVIES", "KEMENKUMHAM", "KHAN", "KIWIBOX", "KOMPAS", "KOMPASIANA", "KONAMINET", "KPROXY", "KU6", "KUXUN", "LADY8844", "LAREDOUTE", "YUGMA", "LATIV", "LDBLOG", "LEAPFILE", "LEBONCOIN", "LETV", "LEVEL3", "LG_ESHOP", "LIBERO_VIDEO", "LIBRARYTHING", "LIFEKNOT", "LINTASBERITA", "LIONAIR", "LIONTRAVEL", "LISTOGRAFY", "LIVEDOOR", "LIVEINTERNET", "LIVE_MEETING", "LIVEMOCHA", "LIVINGSOCIAL", "LOTOUR", "LOTTE", "LOTUS_LIVE", "LUNARSTORM", "LVPING", "SKYPE_FOR_BUSINESS", "MANDRIVA_UPDATE", "MANGOCITY", "MAPQUEST", "MASHABLE", "MASHARE", "MATCH", "MBC", "MBN", "MEDIAFIRE", "MEETIN", "MEETME", "MEETTHEBOSS", "MEETUP", "MEGA", "MK", "MOBAGE", "MOBILE01", "MOBILE_ME", "MOCOSPACE", "MOMOSHOP", "MONEX", "MONEY_163", "MONEYDJ", "MONSTER", "MOP", "MOUTHSHUT", "MOZILLA", "MPQUEST", "MSN_SEARCH", "MT", "MTV", "MULTIPLY", "MULTIUPLOAD", "MUSICA", "MYANIMELIST", "MYCHURCH", "MYHERITAGE", "MYLIFE", "MYVIDEO", "MYWEBSEARCH", "MY_YAHOO", "MYYEARBOOK", "NAPSTER", "NASA", "NASZA_KLASA", "NATECYWORLD", "NATIONALGEOGRAPHIC", "NATIONALLOTTERY", "NAVER", "NBA", "NBA_CHINA", "NDUOA", "NEND", "NETBSD_UPDATE", "NETLOAD", "NETLOG", "NETMARBLE", "NETTBY", "NEXIAN", "NEXON", "NEXOPIA", "NFL", "NGO_POST", "NIFTY", "NIKE", "NIKKEI", "NIMBUZZ_WEB", "NING", "NOD32_UPDATE", "NOKIA_OVI", "NORTON_UPDATE", "NOWNEWS", "NTV", "NYDAILYNEWS", "NYTIMES", "ZOHO_PLANNER", "OFFICEDEPOT", "OKEZONE", "OKWAVE", "ONLINEDOWN", "OOYALA", "OPENBSD_UPDATE", "OPEN_DIARY", "ORB", "OUTLOOK", "PAIPAI", "PANDA_UPDATE", "PANDORA_TV", "PARTNERUP", "PARTY_POKER", "PASSPORTSTAMP", "PAYEASY", "PCGAMES", "PCHOME", "PCLADY", "PCONLINE", "PEERCAST", "PENGYOU", "PEOPLE", "PERFSPOT", "PHOTOBUCKET", "PIMANG", "PINGSTA", "PIXIV", "PIXNET", "PLAXO", "PLAYSTATION", "PLURK", "POCO", "POGO", "PORNHUB", "PPFILM", "PPTV", "PRESENT", "PRICEMINISTER", "PRICERUNNER", "PRIVAX", "PROXEASY", "PSIPHON", "QQ_BLOG", "QQDOWNLOAD", "QQ_FINANCE", "QQ_GAMES", "QQ_LADY", "QQ_NEWS", "QQ_WEB", "QQ_WEIBO", "QUARTERLIFE", "QUNAR", "QZONE", "RADIKO", "RAKUTEN", "RAMBLER", "RAPIDSHARE", "RAVELRY", "REALTOR", "REDHAT_UPDATE", "REDTUBE", "RENREN", "REPUBLIKA", "RESEARCHGATE", "REUTERS", "REVERBNATION", "REVERSO", "RSS", "RTL", "RUTEN", "RYANAIR", "RYZE", "SABERINDO", "SAKURA_NE", "ZOHO_PEOPLE", "SAYCLUB", "SBS", "SCIENCESTAGE", "SCISPACE", "SCRIBD", "SDO", "SECONDLIFE", "SEEQPOD", "SEESAA", "SEESMIC", "SEGYE", "SENDSPACE", "SEOUL_NEWS", "SFR", "SHAREPOINT", "SHAREPOINT_ADMIN", "SHAREPOINT_BLOG", "SHAREPOINT_CALENDAR", "SHAREPOINT_DOCUMENT", "SHAREPOINT_ONLINE", "SHELFARI", "SHINHAN", "SHUTTERFLY", "SIEBEL_CRM", "SINA", "SINA_BLOG", "SINA_FINANCE", "SINA_NEWS", "SINA_VIDEO", "SINA_WEIBO", "SKYBLOG", "SKYCN", "ONEDRIVE", "SLIDESHARE", "SOCIALTV", "SOFT4FUN", "SOFTBANK", "SOGOU", "SOHU", "SOHU_BLOG", "SOKU", "SONET_NE", "SONICO", "SOSO", "SOUFUN", "SOUTHWEST", "SPDY", "SPEEDTEST", "SPIEGEL", "SPORTCHOSUN", "SPORTSILLUSTRATED", "SPORTSSEOUL", "SPRINT", "STAFABAND", "STAGEVU", "STAYFRIENDS", "STICKAM", "STOCKQ", "STREAMAUDIO", "STUDIVZ", "STUMBLEUPON", "SUGARSYNC", "SUNING", "SUPPERSOCCER", "SURROGAFIER", "SURVEYMONKEY", "SVTPLAY", "TABELOG", "TAGGED", "TAGOO", "TAIWANLOTTERY", "TAKU_FILE_BIN", "TALENTTROVE", "TALKBIZNOW", "TALTOPIA", "TAOBAO", "TARINGA", "TCHATCHE", "TEACHERTUBE", "TEACHSTREET", "TECHINLINE", "TEMPOINTERAKTIF", "TENCENT", "TF1", "TGBUS", "THEPIRATEBAY", "THREE", "TIANYA", "TICKETMONSTER", "TIDALTV", "TISTORY", "TMALL", "TOKBOX", "TOKOBAGUS", "TORRENTDOWNLOADS", "TORRENTZ", "TRAVBUDDY", "TRAVELLERSPOINT", "TRAVELOCITY", "TRAVIAN", "TRENDMICRO_UPDATE", "TRIBE", "TRIBUNNEWS", "TROMBI", "TUBE8", "TUDOU", "TUENTI", "ZOHO_NOTEBOOK", "TUNEWIKI", "TV", "TV4PLAY", "TWITPIC", "UBUNTU_ONE", "UDN", "UNIVISION", "UPLOADING", "USATODAY", "USEJUMP", "USTREAM", "VAMPIREFREAKS", "VEETLE", "VIADEO", "VIDEOBASH", "VIDEOSURF", "VIETBAO", "YOUTUBE_HD", "VIVANEWS", "VOX", "VTUNNEL", "VYEW", "WAKOOPA", "WALLSTREETJOURNAL_CHINA", "WANDOUJIA", "WASABI", "WASHINGTONPOST", "WAT", "WAYN", "WEATHER", "WEB_CRAWLER", "WEOURFAMILY", "WERKENNTWEN", "WIKIA", "WIKIPEDIA", "WINDOWS_AZURE", "WINDOWSLIVE", "WINDOWSLIVESPACE", "WINDOWSMEDIA", "WINDOWS_UPDATE", "WIXI", "WOORIBANK", "WRETCH", "XANGA", "XBOX", "XHAMSTER", "XIAMI", "XINHUANET", "XLNET", "XL_WAP", "XL_WEBMAIL", "XMLRPC", "XM_RADIO", "XNXX", "XREA", "XT3", "XUITE", "XVIDEOS", "XVIDEOSLIVE", "ZOHO_MEETING", "YAHOO360PLUSVIETNAM", "YAHOO_ANSWERS", "YAHOO_BIZ", "YAHOO_BUY", "YAHOO_DOUGA", "YAHOO_GAMES", "YAHOO_GEOCITIES", "YAHOO_KOREA", "YAHOO_MAPS", "YAHOO_REALESTATE", "YAHOO_SEARCH", "YAHOO_STOCK_TW", "YAHOO_TRAVEL", "ZOHO_DB", "YELP", "YESKY", "YIHAODIAN", "YOKA", "YOMIURI", "YOUDAO", "YOUKU", "YOUM7", "YOUMEO", "YOUPORN", "YOURFILEHOST", "BBC", "INDIATIMES", "DOCUSIGN", "PASTEBIN", "YANDEXDISK", "LASTPASS", "ZENDESK", "SLACK", "FOURSHARED", "KICKASSTORRENTS", "ADP", "THUGLAK", "REDIFF", "MEGAVIDEO", "MOTIONBOX", "QUIC", "REEBOK", "ADIDAS", "EGNYTE", "WETRANSFER", "INFOARMOR", "FASTMAIL", "YANDEX_MAIL", "GMAIL", "YAHOOMAIL", "HOTMAIL", "DIMP", "FACEBOOK_MAIL", "GMAIL_BASIC", "GMAIL_MOBILE", "GMX", "IMP", "LIVEMAIL_MOBILE", "MAIL_189", "MAIL2000", "MAILRU", "MAKTOOB", "MIMP", "OWA", "QQ_MAIL", "RAMBLER_WEBMAIL", "SQUIRRELMAIL", "YMAIL_CLASSIC", "YMAIL_MOBILE", "ZIMBRA", "ZIMBRA_STANDARD", "ORANGEMAIL", "FLIPKART", "GOOGLE", "YAHOO", "BAIDU", "ASK", "AOL", "BING", "YANDEX", "VMWARE_HORIZON_VIEW",
	"ADOBE_CREATIVE_CLOUD", "ZOOMINFO", "SERVICE_NOW", "MS_SSAS", "GOOGLE_DNS", "CLOUDFLARE_DNS", "ADGUARD", "QUAD9", "OPENDNS", "CLEANBROWSING",
	"COMCAST_DNS", "NEXTDNS", "POWERDNS", "BLAHDNS", "SECUREDNS", "RUBYFISH", "DOH_UNKNOWN", "GOOGLE_KEEP", "AMAZON_CHIME", "WORKDAY", "FIFA", "ROBLOX", "WANGWANG", "S7COMM_PLUS", "DOH", "AGORA_IO", "MS_DFSR", "WS_DISCOVERY", "STUN", "FOLDINGATHOME", "GE_PROCIFY", "MOXA_ASPP", "APP_CH", "GLASSDOOR", "TINDER", "BAIDU_TIEBA", "MIMEDIA", "FILESANYWHERE", "HOUSEPARTY", "GBRIDGE", "HAMACHI", "HEXATECH", "HOTSPOT_SHIELD", "MEGAPROXY", "OPERA_VPN", "SPOTFLUX", "TUNNELBEAR", "ZENMATE", "OPENGW", "VPNOVERDNS", "HOXX_VPN", "VPN1_COM", "SPRINGTECH_VPN", "BARRACUDA_VPN", "HIDEMAN_VPN", "WINDSCRIBE", "BROWSEC_VPN", "EPIC_BROWSER_VPN", "SKYVPN", "KPN_TUNNEL", "ERSPAN",
	"EVASIVE_PROTOCOL", "DOTDASH", "ADOBE_DOCUMENT_CLOUD", "FLIPKART_BOOKS",
}

func validateCloudFirewallNwApplications() schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		value, ok := i.(string)
		if !ok {
			return diag.Diagnostics{
				{
					Severity: diag.Error,
					Summary:  "Expected type to be string",
					Detail:   "Type assertion failed, expected string type for Cloud firewall NW Applications validation",
				},
			}
		}

		// Convert the cty.Path to a string representation
		pathStr := fmt.Sprintf("%+v", path)

		// Use StringInSlice from helper/validation package
		var diags diag.Diagnostics
		if _, errs := validation.StringInSlice(supportedCloudFirewallNwApplications, false)(value, pathStr); len(errs) > 0 {
			for _, err := range errs {
				diags = append(diags, diag.FromErr(err)...)
			}
		}

		return diags
	}
}

var supportedCloudFirewallNwServicesTag = []string{
	"ICMP_ANY", "UDP_ANY", "TCP_ANY", "OTHER_NETWORK_SERVICE", "DNS", "NETBIOS", "FTP", "GNUTELLA", "H_323", "HTTP", "HTTPS", "IKE", "IMAP", "ILS", "IKE_NAT", "IRC", "LDAP", "QUIC", "TDS", "NETMEETING", "NFS", "NTP", "SIP", "SNMP", "SMB", "SMTP", "SSH", "SYSLOG", "TELNET", "TRACEROUTE", "POP3", "PPTP", "RADIUS", "REAL_MEDIA", "RTSP", "VNC", "WHOIS", "KERBEROS_SEC", "TACACS", "SNMPTRAP", "NMAP", "RSYNC", "L2TP", "HTTP_PROXY", "PC_ANYWHERE", "MSN", "ECHO", "AIM", "IDENT", "YMSG", "SCCP", "MGCP_UA", "MGCP_CA", "VDO_LIVE", "OPENVPN", "TFTP", "FTPS_IMPLICIT", "ZSCALER_PROXY_NW_SERVICES", "GRE_PROTOCOL", "ESP_PROTOCOL", "DHCP",
}

func validateCloudFirewallNwServicesTag() schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		value, ok := i.(string)
		if !ok {
			return diag.Diagnostics{
				{
					Severity: diag.Error,
					Summary:  "Expected type to be string",
					Detail:   "Type assertion failed, expected string type for Cloud Firewall NW Services Tags validation",
				},
			}
		}

		// Convert the cty.Path to a string representation
		pathStr := fmt.Sprintf("%+v", path)

		// Use StringInSlice from helper/validation package
		var diags diag.Diagnostics
		if _, errs := validation.StringInSlice(supportedCloudFirewallNwServicesTag, false)(value, pathStr); len(errs) > 0 {
			for _, err := range errs {
				diags = append(diags, diag.FromErr(err)...)
			}
		}

		return diags
	}
}

func validateOCRDlpWebRules(dlp dlp_web_rules.WebDLPRules) error {
	// Define supported file types for OCR enabled scenarios
	supportedFileTypesWithOCREnabled := []string{"BITMAP", "PNG", "JPEG", "TIFF"}

	// Check if OCR is enabled
	if dlp.OcrEnabled {
		// Validate that dlp.FileTypes must be a subset of supportedFileTypesWithOCREnabled
		for _, fileType := range dlp.FileTypes {
			if !contains(supportedFileTypesWithOCREnabled, fileType) {
				return fmt.Errorf("web dlp rule file type '%s' is not supported when OCR is enabled. Supported types: %v", fileType, supportedFileTypesWithOCREnabled)
			}
		}
	}

	return nil
}

func validateDLPRuleFileTypes(dlp dlp_web_rules.WebDLPRules) error {
	// Define allowed file types for both true and false states of `withoutContentInspection`
	allowedFileTypesWithoutInspection := []string{
		"FORM_DATA_POST", "DB", "JAVASCRIPT", "FOR", "MS_POWERPOINT", "TMP", "MATLAB_FILES", "NATVIS", "PNG", "SC", "RUBY_FILES",
		"CAB", "PERL_FILES", "APPLE_DOCUMENTS", "CSX", "POSTSCRIPT", "ZIP", "CATALOG", "BITMAP", "SCZIP", "BORLAND_CPP_FILES",
		"RAR", "SQL", "APPX", "NETMON", "MS_RTF", "PARASOLID", "INF", "ACCDB", "IGS", "HIGH_EFFICIENCY_IMAGE_FILES", "RPY",
		"OAB", "CER", "FTCATEGORY_ENCRYPT", "ENCRYPT", "MM", "DSP", "YAML_FILES", "CHEMDRAW_FILES", "HBS", "SCT", "PS2", "INI", "CERT", "SLDPRT",
		"ICS", "MS_EXCEL", "MS_MSG", "QLIKVIEW_FILES", "MS_MDB", "VISUAL_BASIC_SCRIPT", "MAKE_FILES", "BCP", "MS_CPP_FILES",
		"AAC", "COMPILED_HTML_HELP", "DB2", "SDB", "MS_PST", "JAVA_APPLET", "ADE", "COBOL", "AUTOCAD", "VSDX", "MS_WORD", "CP",
		"BGI", "DAT", "DER", "ASM", "TAR", "BASH_SCRIPTS", "MUI", "PYTHON", "TLB", "HIVE", "KEY", "IMG", "GIF", "STL", "STUFFIT",
		"INCLUDE_FILES", "TABLEAU_FILES", "XZ", "AU3", "PCAP", "DELPHI", "P12", "PHOTOSHOP", "TIFF", "FLASH", "TLI", "VISUAL_CPP_FILES",
		"EML_FILES", "GREENSHOT", "C_FILES", "JAVA_FILES", "MANIFEST", "NFM", "IFC", "VIRTUAL_HARD_DISK", "ISO", "LOG_FILES", "GZIP",
		"EXP", "FCL", "BZIP2", "DMD", "P7Z", "PRT", "NCB", "X1B", "DRAWIO", "XAML", "CML", "ASHX", "PGP", "PS3", "ACIS", "VISUAL_BASIC_FILES",
		"TXT", "DRV", "NLS", "F_FILES", "P7B", "JPEG", "TLH", "CSV", "POD", "SAS", "WINDOWS_META_FORMAT", "RSP", "KDBX", "WINDOWS_SCRIPT_FILES",
		"SCALA", "ONENOTE", "CGR", "BASIC_SOURCE_CODE", "MSC", "POWERSHELL", "PEM", "INTEGRATED_CIRCUIT_FILES", "GO_FILES", "PDF_DOCUMENT",
		"DBF", "JKS", "VDA", "RES_FILES", "A_FILE", "SHELL_SCRAP", "ALL_OUTBOUND", "FTCATEGORY_ENCRYPT", "FTCATEGORY_P7Z", "FTCATEGORY_BZIP2",
		"FTCATEGORY_CAB", "FTCATEGORY_CPIO", "FTCATEGORY_FCL", "FTCATEGORY_GZIP", "FTCATEGORY_ISO", "FTCATEGORY_LZH", "FTCATEGORY_RAR", "FTCATEGORY_STUFFIT",
		"FTCATEGORY_TAR", "FTCATEGORY_XZ", "FTCATEGORY_ZIP", "FTCATEGORY_SCZIP", "FTCATEGORY_ZIPX",
	}

	allowedFileTypesWithInspection := []string{
		"BASH_SCRIPTS", "FORM_DATA_POST", "PYTHON", "INCLUDE_FILES", "TABLEAU_FILES", "JAVASCRIPT", "AU3", "DELPHI", "FOR", "TIFF",
		"MS_POWERPOINT", "TLI", "MATLAB_FILES", "NATVIS", "PNG", "SC", "RUBY_FILES", "VISUAL_CPP_FILES", "EML_FILES", "PERL_FILES",
		"APPLE_DOCUMENTS", "CSX", "C_FILES", "JAVA_FILES", "BITMAP", "IFC", "LOG_FILES", "SCZIP", "BORLAND_CPP_FILES", "SQL",
		"MS_RTF", "INF", "ACCDB", "X1B", "XAML", "RPY", "VISUAL_BASIC_FILES", "DSP", "TXT", "F_FILES", "YAML_FILES", "JPEG", "TLH",
		"CSV", "POD", "SCT", "SAS", "RSP", "WINDOWS_SCRIPT_FILES", "SCALA", "MS_EXCEL", "MS_MSG", "MS_MDB", "BASIC_SOURCE_CODE",
		"MSC", "VISUAL_BASIC_SCRIPT", "POWERSHELL", "GO_FILES", "MAKE_FILES", "BCP", "PDF_DOCUMENT", "MS_CPP_FILES", "RES_FILES",
		"SHELL_SCRAP", "JAVA_APPLET", "COBOL", "VSDX", "MS_WORD", "DAT", "ASM", "ALL_OUTBOUND",
		`FTCATEGORY_ACCDB`, `FTCATEGORY_APPLE_DOCUMENTS`, `FTCATEGORY_ASM`, `FTCATEGORY_AU3`, `FTCATEGORY_BASH_SCRIPTS`,
		`FTCATEGORY_BASIC_SOURCE_CODE`, `FTCATEGORY_BCP`, `FTCATEGORY_BITMAP`, `FTCATEGORY_BORLAND_CPP_FILES`, `FTCATEGORY_C_FILES`,
		`FTCATEGORY_COBOL`, `FTCATEGORY_CSV`, `FTCATEGORY_CSX`, `FTCATEGORY_DAT`, `FTCATEGORY_DCM`, `FTCATEGORY_DELPHI`,
		`FTCATEGORY_DSP`, `FTCATEGORY_EML_FILES`, `FTCATEGORY_F_FILES`, `FTCATEGORY_FOR`, `FTCATEGORY_FORM_DATA_POST`,
		`FTCATEGORY_GO_FILES`, `FTCATEGORY_HTTP`, `FTCATEGORY_IFC`, `FTCATEGORY_INCLUDE_FILES`, `FTCATEGORY_INF`,
		`FTCATEGORY_JAVA_FILES`, `FTCATEGORY_JPEG`, `FTCATEGORY_JSON`, `FTCATEGORY_LOG_FILES`, `FTCATEGORY_MAKE_FILES`,
		`FTCATEGORY_MATLAB_FILES`, `FTCATEGORY_MS_CPP_FILES`, `FTCATEGORY_MS_EXCEL`, `FTCATEGORY_MS_MDB`, `FTCATEGORY_MS_MSG`,
		`FTCATEGORY_MS_POWERPOINT`, `FTCATEGORY_MS_PUB`, `FTCATEGORY_MS_RTF`, `FTCATEGORY_MS_WORD`, `FTCATEGORY_MSC`, `FTCATEGORY_NATVIS`,
		`FTCATEGORY_OLM`, `FTCATEGORY_OPEN_OFFICE_PRESENTATIONS`, `FTCATEGORY_OPEN_OFFICE_SPREADSHEETS`, `FTCATEGORY_PDF_DOCUMENT`,
		`FTCATEGORY_PERL_FILES`, `FTCATEGORY_PNG`, `FTCATEGORY_POD`, `FTCATEGORY_POWERSHELL`, `FTCATEGORY_PYTHON`, `FTCATEGORY_RES_FILES`,
		`FTCATEGORY_RPY`, `FTCATEGORY_RSP`, `FTCATEGORY_RUBY_FILES`, `FTCATEGORY_SAS`, `FTCATEGORY_SC`, `FTCATEGORY_SCALA`,
		`FTCATEGORY_SCT`, `FTCATEGORY_SHELL_SCRAP`, `FTCATEGORY_SQL`, `FTCATEGORY_TABLEAU_FILES`, `FTCATEGORY_TIFF`, `FTCATEGORY_TLH`,
		`FTCATEGORY_TLI`, `FTCATEGORY_TXT`, `FTCATEGORY_UNK_TXT`, `FTCATEGORY_VISUAL_BASIC_FILES`, `FTCATEGORY_VISUAL_BASIC_SCRIPT`,
		`FTCATEGORY_VISUAL_CPP_FILES`, `FTCATEGORY_VSDX`, `FTCATEGORY_WINDOWS_SCRIPT_FILES`, `FTCATEGORY_X1B`, `FTCATEGORY_XAML`,
		`FTCATEGORY_XML`, `FTCATEGORY_YAML_FILES`, `FTCATEGORY_JAVA_APPLET`, `FTCATEGORY_JAVASCRIPT`,
	}

	// Check if `ALL_OUTBOUND` is selected and `withoutContentInspection` is false
	allOutboundSelected := contains(dlp.FileTypes, "ALL_OUTBOUND")
	if allOutboundSelected && !dlp.WithoutContentInspection {
		return fmt.Errorf("when file_type ALL_OUTBOUND is present, without_content_inspection must be true")
	}

	// If ALL_OUTBOUND is selected and no other file types are present, allow it
	if allOutboundSelected && len(dlp.FileTypes) > 1 {
		return fmt.Errorf("cannot have other file types when ALL_OUTBOUND is selected")
	}

	// Validate file types based on the `withoutContentInspection` flag
	var allowedFileTypes []string
	if dlp.WithoutContentInspection {
		allowedFileTypes = allowedFileTypesWithoutInspection
	} else {
		allowedFileTypes = allowedFileTypesWithInspection
	}

	// Ensure all selected file types are in the allowed list
	for _, fileType := range dlp.FileTypes {
		if !contains(allowedFileTypes, fileType) {
			return fmt.Errorf("the file_type '%s' is not accepted when without_content_inspection is %v", fileType, dlp.WithoutContentInspection)
		}
	}

	return nil
}

/*
func validateDLPRuleFileTypes(dlp dlp_web_rules.WebDLPRules) error {
	// New check: If FileTypes is not defined, WithoutContentInspection must be false
	if len(dlp.FileTypes) == 0 && dlp.WithoutContentInspection {
		return fmt.Errorf("without_content_inspection must be set to false when no file types are defined")
	}

	var allowedFileTypes []string

	allOutboundSelected := contains(dlp.FileTypes, "ALL_OUTBOUND")

	// If ALL_OUTBOUND is selected and withoutContentInspection is true, it should not trigger an error.
	if allOutboundSelected && len(dlp.FileTypes) == 1 && dlp.WithoutContentInspection {
		return nil
	}

	if allOutboundSelected && len(dlp.FileTypes) > 1 {
		return fmt.Errorf("cannot have other file types when ALL_OUTBOUND is selected")
	}

	if dlp.WithoutContentInspection {
		// Define allowed file types when without_content_inspection is true
		allowedFileTypes = []string{
			"ACCDB", "APPLE_DOCUMENTS", "APPX", "ASM", "AU3", "BASH_SCRIPTS", "BASIC_SOURCE_CODE", "BCP", "BORLAND_CPP_FILES", "C_FILES", "CHEMDRAW_FILES", "CML", "COBOL", "COMPILED_HTML_HELP", "CP", "CSV", "CSX", "DAT", "DELPHI", "DMD", "DSP", "EML_FILES", "F_FILES", "FOR", "FORM_DATA_POST", "GO_FILES", "IFC", "INCLUDE_FILES", "INF", "JAVA_FILES", "LOG_FILES", "MAKE_FILES", "MATLAB_FILES", "MM", "MS_CPP_FILES", "MS_EXCEL", "MS_MDB", "MS_MSG", "MS_POWERPOINT", "MS_RTF", "MS_WORD", "MSC", "NATVIS", "OAB", "PDF_DOCUMENT", "PERL_FILES", "POD", "POSTSCRIPT", "POWERSHELL", "PYTHON", "QLIKVIEW_FILES", "RES_FILES", "RPY", "RSP", "RUBY_FILES", "SAS", "SC", "SCALA", "SCT", "SCZIP", "SHELL_SCRAP", "SQL", "TABLEAU_FILES", "TLH", "TLI", "TXT", "VISUAL_BASIC_FILES", "VISUAL_BASIC_SCRIPT", "VISUAL_CPP_FILES", "VSDX", "WINDOWS_META_FORMAT", "WINDOWS_SCRIPT_FILES", "X1B", "XAML", "YAML_FILES", "JAVA_APPLET", "JAVASCRIPT",
		}
	} else {
		// Define allowed file types when without_content_inspection is false
		allowedFileTypes = []string{
			"A_FILE", "ACCDB", "ADE", "APPLE_DOCUMENTS", "APPX", "ASM", "AU3", "AUTOCAD", "BASH_SCRIPTS", "BASIC_SOURCE_CODE", "BCP", "BGI", "BITMAP", "BORLAND_CPP_FILES", "BZIP2", "C_FILES", "CAB", "CER", "CERT", "CHEMDRAW_FILES", "CML", "COBOL", "COMPILED_HTML_HELP", "CP", "CSV", "CSX", "DAT", "DB", "DB2", "DBF", "DELPHI", "DER", "DMD", "DRV", "DSP", "EML_FILES", "ENCRYPT", "F_FILES", "FOR", "FORM_DATA_POST", "GIF", "GO_FILES", "GZIP", "IFC", "INCLUDE_FILES", "INF", "INI", "INTEGRATED_CIRCUIT_FILES", "ISO", "JAVA_FILES", "JKS", "JPEG", "KEY", "LOG_FILES", "MAKE_FILES", "MANIFEST", "MATLAB_FILES", "MM", "MS_CPP_FILES", "MS_EXCEL", "MS_MDB", "MS_MSG", "MS_POWERPOINT", "MS_RTF", "MS_WORD", "MSC", "NATVIS", "NCB", "NFM", "NLS", "OAB", "ONENOTE", "P12", "P7B", "P7Z", "PCAP", "PDF_DOCUMENT", "PEM", "PERL_FILES", "PHOTOSHOP", "PNG", "POD", "POSTSCRIPT", "POWERSHELL", "PYTHON", "QLIKVIEW_FILES", "RAR", "RES_FILES", "RPY", "RSP", "RUBY_FILES", "SAS", "SC", "SCALA", "SCT", "SCZIP", "SHELL_SCRAP", "SQL", "STL", "STUFFIT", "TABLEAU_FILES", "TAR", "TIFF", "TLH", "TLI", "TXT", "VISUAL_BASIC_FILES", "VISUAL_BASIC_SCRIPT", "VISUAL_CPP_FILES", "VSDX", "WINDOWS_META_FORMAT", "WINDOWS_SCRIPT_FILES", "X1B", "XAML", "YAML_FILES", "ZIP", "FLASH", "JAVA_APPLET", "JAVASCRIPT",
		}
	}
	for _, fileType := range dlp.FileTypes {
		if !contains(allowedFileTypes, fileType) {
			return fmt.Errorf("the file_type '%s' is not accepted when without_content_inspection is false", fileType)
		}
	}

	return nil
}
*/

func validateDeviceTrustLevels() schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		value, ok := i.(string)
		if !ok {
			return diag.Diagnostics{
				{
					Severity: diag.Error,
					Summary:  "Expected type to be string",
					Detail:   "Type assertion failed, expected string type for Device Trust Levels validation",
				},
			}
		}

		// Convert the cty.Path to a string representation
		pathStr := fmt.Sprintf("%+v", path)

		// Use StringInSlice from helper/validation package
		var diags diag.Diagnostics
		if _, errs := validation.StringInSlice(supportedDeviceTrustLevels, false)(value, pathStr); len(errs) > 0 {
			for _, err := range errs {
				diags = append(diags, diag.FromErr(err)...)
			}
		}

		return diags
	}
}

var supportedDeviceTrustLevels = []string{
	"UNKNOWN_DEVICETRUSTLEVEL", "LOW_TRUST", "MEDIUM_TRUST", "HIGH_TRUST",
}

func validateURLFilteringActions(rule urlfilteringpolicies.URLFilteringRule) error {
	switch rule.Action {
	case "ISOLATE":
		// Validation 1: Check if any field in CBIProfile is set
		if rule.CBIProfile.ID == "" && rule.CBIProfile.Name == "" && rule.CBIProfile.URL == "" {
			return errors.New("cbi_profile attribute is required when action is ISOLATE")
		}

		// Validation 2: Check user_agent_types does not contain "OTHER"
		for _, userAgent := range rule.UserAgentTypes {
			if userAgent == "OTHER" {
				return errors.New("user_agent_types should not contain 'OTHER' when action is ISOLATE. Valid options are: FIREFOX, MSIE, MSEDGE, CHROME, SAFARI, MSCHREDGE")
			}
		}

		// Validation 3: Check Protocols should be HTTP or HTTPS
		validProtocols := map[string]bool{"HTTPS_RULE": true, "HTTP_RULE": true}
		for _, protocol := range rule.Protocols {
			if !validProtocols[strings.ToUpper(protocol)] {
				return errors.New("when action is ISOLATE, valid options for protocols are: HTTP and/or HTTPS")
			}
		}

	case "CAUTION":
		// Validation 4: Ensure request_methods only contain CONNECT, GET, HEAD
		validMethods := map[string]bool{"CONNECT": true, "GET": true, "HEAD": true}
		for _, method := range rule.RequestMethods { // Assuming RequestMethods is the correct field
			if !validMethods[strings.ToUpper(method)] {
				return errors.New("'CAUTION' action is allowed only for CONNECT/GET/HEAD request methods")
			}
		}
	}

	return nil
}

var predefinedIdentifiersMap = map[string][]string{
	"ASPP_LEAKAGE": {
		"ASPP_CN", "ASPP_JP", "ASPP_KR", "ASPP_MY", "ASPP_PH", "ASPP_SG", "ASPP_TW", "ASPP_TR",
	},
	"CRED_LEAKAGE": {
		"CRED_AMAZON_MWS_TOKEN", "CRED_GIT_TOKEN", "CRED_GITHUB_TOKEN", "CRED_GOOGLE_API", "CRED_GOOGLE_OAUTH_TOKEN",
		"CRED_GOOGLE_OAUTH_ID", "CRED_JWT_TOKEN", "CRED_PAYPAL_TOKEN", "CRED_PICATIC_API_KEY", "CRED_PRIVATE_KEY",
		"CRED_SENDGRID_API_KEY", "CRED_SLACK_TOKEN", "CRED_SLACK_WEBHOOK", "CRED_SQUARE_ACCESS_TOKEN", "CRED_SQUARE_OAUTH_SECRET",
		"CRED_STRIPE_API_KEY",
	},
	"EUIBAN_LEAKAGE": {
		"EUIBAN_AD", "EUIBAN_AT", "EUIBAN_BE", "EUIBAN_BA", "EUIBAN_BG", "EUIBAN_HR", "EUIBAN_CY",
		"EUIBAN_CZ", "EUIBAN_DK", "EUIBAN_EE", "EUIBAN_FO", "EUIBAN_FI", "EUIBAN_FR", "EUIBAN_DE",
		"EUIBAN_GI", "EUIBAN_GR", "EUIBAN_GL", "EUIBAN_HU", "EUIBAN_IS", "EUIBAN_IE", "EUIBAN_IL",
		"EUIBAN_IT", "EUIBAN_LV", "EUIBAN_LI", "EUIBAN_LT", "EUIBAN_LU", "EUIBAN_MK", "EUIBAN_MT",
		"EUIBAN_MC", "EUIBAN_ME", "EUIBAN_NL", "EUIBAN_NO", "EUIBAN_PL", "EUIBAN_PT", "EUIBAN_RO",
		"EUIBAN_SM", "EUIBAN_RS", "EUIBAN_SK", "EUIBAN_SI", "EUIBAN_ES", "EUIBAN_SE", "EUIBAN_CH",
		"EUIBAN_TN", "EUIBAN_TR", "EUIBAN_GB",
	},
	"PPEU_LEAKAGE": {
		"EUPP_AT", "EUPP_BE", "EUPP_BG", "EUPP_CZ", "EUPP_DK", "EUPP_EE", "EUPP_FL", "EUPP_FR", "EUPP_DE",
		"EUPP_GR", "EUPP_HU", "EUPP_IE", "EUPP_IT", "EUPP_LV", "EUPP_LU", "EUPP_NL", "EUPP_PL", "EUPP_PT",
		"EUPP_RO", "EUPP_SK", "EUPP_SI", "EUPP_ES", "EUPP_SE",
	},
	"USDL_LEAKAGE": {
		"USDL_AL", "USDL_AK", "USDL_AZ", "USDL_CA", "USDL_CO", "USDL_CT", "USDL_DE", "USDL_DC", "USDL_FL", "USDL_GA",
		"USDL_HI", "USDL_ID", "USDL_IL", "USDL_IN", "USDL_IA", "USDL_KS", "USDL_KY", "USDL_LA", "USDL_ME", "USDL_MD",
		"USDL_MA", "USDL_MI", "USDL_MN", "USDL_MS", "USDL_MO", "USDL_MT", "USDL_NE", "USDL_NV", "USDL_NH", "USDL_NJ",
		"USDL_NM", "USDL_NY", "USDL_NC", "USDL_ND", "USDL_OH", "USDL_OK", "USDL_OR", "USDL_PA", "USDL_RI", "USDL_SC",
		"USDL_SD", "USDL_TN", "USDL_TX", "USDL_UT", "USDL_VT", "USDL_VA", "USDL_WA", "USDL_WV", "USDL_WI", "USDL_WY",
	},
}

func validateDLPHierarchicalIdentifiersDiff(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
	identifiers := diff.Get("hierarchical_identifiers").(*schema.Set).List()

	var invalidIdentifiers []string
	for _, identifier := range identifiers {
		idStr := identifier.(string)
		if _, exists := predefinedIdentifiersMap[idStr]; !exists {
			valid := false
			for _, ids := range predefinedIdentifiersMap {
				for _, id := range ids {
					if id == idStr {
						valid = true
						break
					}
				}
				if valid {
					break
				}
			}
			if !valid {
				invalidIdentifiers = append(invalidIdentifiers, idStr)
			}
		}
	}

	if len(invalidIdentifiers) > 0 {
		return fmt.Errorf("invalid hierarchical identifiers: %v. Supported identifiers are: %v", invalidIdentifiers, getAllSupportedIdentifiers())
	}
	return nil
}

func getAllSupportedIdentifiers() []string {
	var identifiers []string
	for key, values := range predefinedIdentifiersMap {
		identifiers = append(identifiers, key)
		identifiers = append(identifiers, values...)
	}
	return identifiers
}

func validatePredefinedIdentifier(val interface{}, key string) (warns []string, errs []error) {
	supportedIdentifiers := []string{
		"EUIBAN_LEAKAGE", "ASPP_LEAKAGE", "REGON_LEAKAGE", "BRCNPJ_LEAKAGE", "MXCURP_LEAKAGE",
		"HRPIN_LEAKAGE", "JMBG_LEAKAGE", "TREATMENTS_LEAKAGE", "DRUGS_LEAKAGE", "DISEASES_LEAKAGE",
		"NDC_PROD_LEAKAGE", "NDC_PKG_LEAKAGE", "CRED_LEAKAGE", "PPEU_LEAKAGE", "USDL_LEAKAGE",
		"IEVAT_LEAKAGE", "FITIN_LEAKAGE", "IETIN_LEAKAGE", "NZNHIN_LEAKAGE", "GRTIN_LEAKAGE",
		"PLTIN_LEAKAGE", "USITIN_LEAKAGE", "AUPP_LEAKAGE", "CIF_LEAKAGE", "PERUC_LEAKAGE",
		"NPI_LEAKAGE", "ATSSN_LEAKAGE", "NLVAT_LEAKAGE", "LUVAT_LEAKAGE", "FRVAT_LEAKAGE",
		"DEVAT_LEAKAGE", "ATVAT_LEAKAGE", "BEVAT_LEAKAGE", "EGN_LEAKAGE", "LUTIN_LEAKAGE",
		"HUTIN_LEAKAGE", "ATTIN_LEAKAGE", "FRTIN_LEAKAGE", "SETIN_LEAKAGE", "PTTIN_LEAKAGE",
		"DKTIN_LEAKAGE", "BETIN_LEAKAGE", "DETIN_LEAKAGE", "NZTIN_LEAKAGE", "ML_IMMIGRATION_LEAKAGE",
		"ML_CORPORATE_LEGAL_LEAKAGE", "ML_COURT_LEAKAGE", "ML_LEGAL_LEAKAGE", "ML_TAX_LEAKAGE",
		"ML_INSURANCE_LEAKAGE", "ML_INVOICE_LEAKAGE", "ML_RESUME_LEAKAGE", "ML_REAL_ESTATE_LEAKAGE",
		"ML_MEDICAL_DOC_LEAKAGE", "ML_TECH_LEAKAGE", "ML_CORPORATE_FINANCE_LEAKAGE", "ML_DMV_LEAKAGE",
		"INSEE_LEAKAGE", "DNI_LEAKAGE", "AADHAR_LEAKAGE", "TICN_LEAKAGE", "MYKAD_LEAKAGE",
		"TIN_LEAKAGE", "PESEL_LEAKAGE", "AHV_LEAKAGE", "CYBER_BULLY", "JCN_LEAKAGE", "JMN_LEAKAGE",
		"SDS_LEAKAGE", "NAME_ES_LEAKAGE", "NAME_CA_LEAKAGE", "FISCAL_LEAKAGE", "HKID_LEAKAGE",
		"CREDIT_CARD", "FINANCIAL", "MEDICAL", "SOURCE_CODE", "NRIC_LEAKAGE", "CSIN_LEAKAGE",
		"SALESFORCE_REPORT_LEAKAGE", "ADULT_CONTENT", "ILLEGAL_DRUGS", "GAMBLING", "WEAPONS",
		"NINO_LEAKAGE", "NAME_LEAKAGE", "MCN_LEAKAGE", "TFN_LEAKAGE", "BSN_LEAKAGE", "ABA_LEAKAGE",
		"CLABE_LEAKAGE", "CPF_LEAKAGE", "NHS_LEAKAGE", "CNID_LEAKAGE", "KRRN_LEAKAGE", "TNID_LEAKAGE",
	}

	name := val.(string)
	for _, identifier := range supportedIdentifiers {
		if name == identifier {
			return
		}
	}

	errs = append(errs, fmt.Errorf("%q is not a valid predefined identifier. Supported identifiers are: %v", key, supportedIdentifiers))
	return
}

func validateDestAddress(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	// Regular expression to match FQDNs and wildcard FQDNs
	fqdnRegex := `^(\*\.)?([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`

	// Regular expression to match IPv4 ranges
	ipRangeRegex := `^(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})-(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})$`

	// Check if the value is a valid IPv4 CIDR
	if _, _, err := net.ParseCIDR(value); err == nil {
		ip, _, _ := net.ParseCIDR(value)
		if ip.To4() != nil {
			// It's a valid IPv4 CIDR
			return
		} else {
			errors = append(errors, fmt.Errorf("invalid IPv4 address: %s. IPv6 addresses are not allowed", value))
			return
		}
	}

	// Check if the value is a valid IPv4 address (without CIDR)
	if ip := net.ParseIP(value); ip != nil {
		if ip.To4() != nil {
			// It's a valid IPv4 address
			return
		} else {
			errors = append(errors, fmt.Errorf("invalid IPv4 address: %s. IPv6 addresses are not allowed", value))
			return
		}
	}

	// Check if the value is a valid IPv4 range
	if matched, _ := regexp.MatchString(ipRangeRegex, value); matched {
		parts := strings.Split(value, "-")
		if len(parts) == 2 {
			startIP := net.ParseIP(parts[0])
			endIP := net.ParseIP(parts[1])
			if startIP != nil && endIP != nil && startIP.To4() != nil && endIP.To4() != nil {
				// It's a valid IPv4 range
				return
			}
		}
		errors = append(errors, fmt.Errorf("invalid IPv4 range: %s. Must be a valid IPv4 range", value))
		return
	}

	// Check if the value is a valid FQDN or wildcard FQDN
	if matched, _ := regexp.MatchString(fqdnRegex, value); matched {
		// It's a valid FQDN or wildcard FQDN
		return
	}

	// If none of the above checks passed, it must be an invalid address
	errors = append(errors, fmt.Errorf("invalid address: %s. Must be a valid IPv4 address, IPv4 CIDR, IPv4 range, FQDN, or wildcard FQDN", value))
	return
}
