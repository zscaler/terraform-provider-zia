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
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/dlp/dlp_web_rules"
	"github.com/zscaler/zscaler-sdk-go/v3/zscaler/zia/services/urlfilteringpolicies"
)

func intBetween(min, max int) schema.SchemaValidateDiagFunc {
	return func(i interface{}, k cty.Path) diag.Diagnostics {
		v, ok := i.(int)
		if !ok {
			return diag.Errorf("expected type of %s to be integer", k)
		}
		if v < min || v > max {
			return diag.Errorf("expected %s to be in the range (%d - %d), got %d", k, min, max, v)
		}
		return nil
	}
}

func intAtMost(max int) schema.SchemaValidateDiagFunc {
	return func(i interface{}, k cty.Path) diag.Diagnostics {
		v, ok := i.(int)
		if !ok {
			return diag.Errorf("expected type of %s to be integer", k)
		}
		if v > max {
			return diag.Errorf("expected %s to be at most (%d), got %d", k, max, v)
		}
		return nil
	}
}

// Validate URL Filtering Category Options
var supportedURLCategories = []string{
	"ANY", "NONE",
	"OTHER_ADULT_MATERIAL", "ADULT_THEMES", "LINGERIE_BIKINI", "NUDITY", "PORNOGRAPHY", "SEXUALITY",
	"ADULT_SEX_EDUCATION", "K_12_SEX_EDUCATION", "SOCIAL_ADULT", "OTHER_BUSINESS_AND_ECONOMY",
	"CORPORATE_MARKETING", "FINANCE", "PROFESSIONAL_SERVICES", "CLASSIFIEDS", "TRADING_BROKARAGE_INSURANCE",
	"OTHER_DRUGS", "MARIJUANA", "OTHER_EDUCATION", "CONTINUING_EDUCATION_COLLEGES", "HISTORY",
	"K_12", "REFERENCE_SITES", "SCIENCE_AND_TECHNOLOGY", "OTHER_ENTERTAINMENT_AND_RECREATION",
	"ENTERTAINMENT", "TELEVISION_AND_MOVIES", "MUSIC", "STREAMING_MEDIA", "RADIO_STATIONS", "GAMBLING",
	"OTHER_GAMES", "SOCIAL_NETWORKING_GAMES", "OTHER_GOVERNMENT_AND_POLITICS", "GOVERNMENT", "POLITICS",
	"HEALTH", "OTHER_ILLEGAL_OR_QUESTIONABLE", "COPYRIGHT_INFRINGEMENT", "COMPUTER_HACKING", "QUESTIONABLE",
	"PROFANITY", "MATURE_HUMOR", "ANONYMIZER", "OTHER_INFORMATION_TECHNOLOGY", "TRANSLATORS", "IMAGE_HOST",
	"FILE_HOST", "SHAREWARE_DOWNLOAD", "WEB_BANNERS", "WEB_HOST", "WEB_SEARCH", "PORTALS", "SAFE_SEARCH_ENGINE",
	"CDN", "OSS_UPDATES", "DNS_OVER_HTTPS", "OTHER_INTERNET_COMMUNICATION", "INTERNET_SERVICES",
	"DISCUSSION_FORUMS", "ONLINE_CHAT", "EMAIL_HOST", "BLOG", "P2P_COMMUNICATION", "REMOTE_ACCESS",
	"WEB_CONFERENCING", "ZSPROXY_IPS", "JOB_SEARCH", "MILITANCY_HATE_AND_EXTREMISM", "OTHER_MISCELLANEOUS",
	"MISCELLANEOUS_OR_UNKNOWN", "NEWLY_REG_DOMAINS", "NON_CATEGORIZABLE", "NEWS_AND_MEDIA", "OTHER_RELIGION",
	"TRADITIONAL_RELIGION", "CULT", "ALT_NEW_AGE", "OTHER_SECURITY", "ADWARE_OR_SPYWARE", "ENCR_WEB_CONTENT",
	"MALICIOUS_TLD", "OTHER_SHOPPING_AND_AUCTIONS", "SPECIALIZED_SHOPPING", "REAL_ESTATE", "ONLINE_AUCTIONS",
	"OTHER_SOCIAL_AND_FAMILY_ISSUES", "SOCIAL_ISSUES", "FAMILY_ISSUES", "OTHER_SOCIETY_AND_LIFESTYLE",
	"ART_CULTURE", "ALTERNATE_LIFESTYLE", "HOBBIES_AND_LEISURE", "DINING_AND_RESTAURANT", "ALCOHOL_TOBACCO",
	"SOCIAL_NETWORKING", "SPECIAL_INTERESTS", "SPORTS", "TASTELESS", "TRAVEL", "USER_DEFINED", "VEHICLES",
	"VIOLENCE", "WEAPONS_AND_BOMBS", "DYNAMIC_DNS", "MILITARY", "NEWLY_REVIVED_DOMAINS", "AI_ML_APPS",
	"FILE_CONVERTORS", "GENERAL_AI_ML", "INSURANCE", "CUSTOM_00", "CUSTOM_01", "CUSTOM_02", "CUSTOM_03",
	"CUSTOM_04", "CUSTOM_05", "CUSTOM_06", "CUSTOM_07", "CUSTOM_08", "CUSTOM_09", "CUSTOM_10", "CUSTOM_11",
	"CUSTOM_12", "CUSTOM_13", "CUSTOM_14", "CUSTOM_15", "CUSTOM_16", "CUSTOM_17", "CUSTOM_18", "CUSTOM_19",
	"CUSTOM_20", "CUSTOM_21", "CUSTOM_22", "CUSTOM_23", "CUSTOM_24", "CUSTOM_25", "CUSTOM_26", "CUSTOM_27",
	"CUSTOM_28", "CUSTOM_29", "CUSTOM_30", "CUSTOM_31", "CUSTOM_32", "CUSTOM_33", "CUSTOM_34", "CUSTOM_35",
	"CUSTOM_36", "CUSTOM_37", "CUSTOM_38", "CUSTOM_39", "CUSTOM_40", "CUSTOM_41", "CUSTOM_42", "CUSTOM_43",
	"CUSTOM_44", "CUSTOM_45", "CUSTOM_46", "CUSTOM_47", "CUSTOM_48", "CUSTOM_49", "CUSTOM_50", "CUSTOM_51",
	"CUSTOM_52", "CUSTOM_53", "CUSTOM_54", "CUSTOM_55", "CUSTOM_56", "CUSTOM_57", "CUSTOM_58", "CUSTOM_59",
	"CUSTOM_60", "CUSTOM_61", "CUSTOM_62", "CUSTOM_63", "CUSTOM_64", "CUSTOM_65", "CUSTOM_66", "CUSTOM_67",
	"CUSTOM_68", "CUSTOM_69", "CUSTOM_70", "CUSTOM_71", "CUSTOM_72", "CUSTOM_73", "CUSTOM_74", "CUSTOM_75",
	"CUSTOM_76", "CUSTOM_77", "CUSTOM_78", "CUSTOM_79", "CUSTOM_80", "CUSTOM_81", "CUSTOM_82", "CUSTOM_83",
	"CUSTOM_84", "CUSTOM_85", "CUSTOM_86", "CUSTOM_87", "CUSTOM_88", "CUSTOM_89", "CUSTOM_90", "CUSTOM_91",
	"CUSTOM_92", "CUSTOM_93", "CUSTOM_94", "CUSTOM_95", "CUSTOM_96", "CUSTOM_97", "CUSTOM_98", "CUSTOM_99",
	"CUSTOM_100", "CUSTOM_101", "CUSTOM_102", "CUSTOM_103", "CUSTOM_104", "CUSTOM_105", "CUSTOM_106", "CUSTOM_107",
	"CUSTOM_108", "CUSTOM_109", "CUSTOM_110", "CUSTOM_111", "CUSTOM_112", "CUSTOM_113", "CUSTOM_114", "CUSTOM_115",
	"CUSTOM_116", "CUSTOM_117", "CUSTOM_118", "CUSTOM_119", "CUSTOM_120", "CUSTOM_121", "CUSTOM_122", "CUSTOM_123",
	"CUSTOM_124", "CUSTOM_125", "CUSTOM_126", "CUSTOM_127", "CUSTOM_128", "CUSTOM_129", "CUSTOM_130", "CUSTOM_131",
	"CUSTOM_132", "CUSTOM_133", "CUSTOM_134", "CUSTOM_135", "CUSTOM_136", "CUSTOM_137", "CUSTOM_138", "CUSTOM_139",
	"CUSTOM_140", "CUSTOM_141", "CUSTOM_142", "CUSTOM_143", "CUSTOM_144", "CUSTOM_145", "CUSTOM_146", "CUSTOM_147",
	"CUSTOM_148", "CUSTOM_149", "CUSTOM_150", "CUSTOM_151", "CUSTOM_152", "CUSTOM_153", "CUSTOM_154", "CUSTOM_155",
	"CUSTOM_156", "CUSTOM_157", "CUSTOM_158", "CUSTOM_159", "CUSTOM_160", "CUSTOM_161", "CUSTOM_162", "CUSTOM_163",
	"CUSTOM_164", "CUSTOM_165", "CUSTOM_166", "CUSTOM_167", "CUSTOM_168", "CUSTOM_169", "CUSTOM_170", "CUSTOM_171",
	"CUSTOM_172", "CUSTOM_173", "CUSTOM_174", "CUSTOM_175", "CUSTOM_176", "CUSTOM_177", "CUSTOM_178", "CUSTOM_179",
	"CUSTOM_180", "CUSTOM_181", "CUSTOM_182", "CUSTOM_183", "CUSTOM_184", "CUSTOM_185", "CUSTOM_186", "CUSTOM_187",
	"CUSTOM_188", "CUSTOM_189", "CUSTOM_190", "CUSTOM_191", "CUSTOM_192", "CUSTOM_193", "CUSTOM_194", "CUSTOM_195",
	"CUSTOM_196", "CUSTOM_197", "CUSTOM_198", "CUSTOM_199", "CUSTOM_200", "CUSTOM_201", "CUSTOM_202", "CUSTOM_203",
	"CUSTOM_204", "CUSTOM_205", "CUSTOM_206", "CUSTOM_207", "CUSTOM_208", "CUSTOM_209", "CUSTOM_210", "CUSTOM_211",
	"CUSTOM_212", "CUSTOM_213", "CUSTOM_214", "CUSTOM_215", "CUSTOM_216", "CUSTOM_217", "CUSTOM_218", "CUSTOM_219",
	"CUSTOM_220", "CUSTOM_221", "CUSTOM_222", "CUSTOM_223", "CUSTOM_224", "CUSTOM_225", "CUSTOM_226", "CUSTOM_227",
	"CUSTOM_228", "CUSTOM_229", "CUSTOM_230", "CUSTOM_231", "CUSTOM_232", "CUSTOM_233", "CUSTOM_234", "CUSTOM_235",
	"CUSTOM_236", "CUSTOM_237", "CUSTOM_238", "CUSTOM_239", "CUSTOM_240", "CUSTOM_241", "CUSTOM_242", "CUSTOM_243",
	"CUSTOM_244", "CUSTOM_245", "CUSTOM_246", "CUSTOM_247", "CUSTOM_248", "CUSTOM_249", "CUSTOM_250", "CUSTOM_251",
	"CUSTOM_252", "CUSTOM_253", "CUSTOM_254", "CUSTOM_255", "CUSTOM_256",
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

func validateFileTypeControlProtocols() schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		value, ok := i.(string)
		if !ok {
			return diag.Diagnostics{
				{
					Severity: diag.Error,
					Summary:  "Expected type to be string",
					Detail:   "Type assertion failed, expected string type for File Type Control Protocols validation",
				},
			}
		}

		// Convert the cty.Path to a string representation
		pathStr := fmt.Sprintf("%+v", path)

		// Use StringInSlice from helper/validation package
		var diags diag.Diagnostics
		if _, errs := validation.StringInSlice(supportedFileTypeProtocols, false)(value, pathStr); len(errs) > 0 {
			for _, err := range errs {
				diags = append(diags, diag.FromErr(err)...)
			}
		}

		return diags
	}
}

var supportedFileTypeProtocols = []string{
	"FOHTTP_RULE", "FTP_RULE", "HTTPS_RULE", "HTTP_RULE",
}

func validateSandboxRuleProtocols() schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		value, ok := i.(string)
		if !ok {
			return diag.Diagnostics{
				{
					Severity: diag.Error,
					Summary:  "Expected type to be string",
					Detail:   "Type assertion failed, expected string type for Sandbox Protocols validation",
				},
			}
		}

		// Convert the cty.Path to a string representation
		pathStr := fmt.Sprintf("%+v", path)

		// Use StringInSlice from helper/validation package
		var diags diag.Diagnostics
		if _, errs := validation.StringInSlice(supportedSandboxProtocols, false)(value, pathStr); len(errs) > 0 {
			for _, err := range errs {
				diags = append(diags, diag.FromErr(err)...)
			}
		}

		return diags
	}
}

var supportedSandboxProtocols = []string{
	"ANY_RULE", "FOHTTP_RULE", "SMRULEF_CASCADING_ALLOWED", "FTP_RULE", "HTTPS_RULE", "HTTP_RULE",
}

func validateDNSRuleProtocols() schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		value, ok := i.(string)
		if !ok {
			return diag.Diagnostics{
				{
					Severity: diag.Error,
					Summary:  "Expected type to be string",
					Detail:   "Type assertion failed, expected string type for Sandbox Protocols validation",
				},
			}
		}

		// Convert the cty.Path to a string representation
		pathStr := fmt.Sprintf("%+v", path)

		// Use StringInSlice from helper/validation package
		var diags diag.Diagnostics
		if _, errs := validation.StringInSlice(supportedDNSProtocols, false)(value, pathStr); len(errs) > 0 {
			for _, err := range errs {
				diags = append(diags, diag.FromErr(err)...)
			}
		}

		return diags
	}
}

var supportedDNSProtocols = []string{
	"ANY_RULE", "DOHTTPS_RULE", "TCP_RULE", "UDP_RULE",
}

func validateSandboxRuleFileTypes() schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		value, ok := i.(string)
		if !ok {
			return diag.Diagnostics{
				{
					Severity: diag.Error,
					Summary:  "Expected type to be string",
					Detail:   "Type assertion failed, expected string type for Sandbox Protocols validation",
				},
			}
		}

		// Convert the cty.Path to a string representation
		pathStr := fmt.Sprintf("%+v", path)

		// Use StringInSlice from helper/validation package
		var diags diag.Diagnostics
		if _, errs := validation.StringInSlice(supportedSandboxFileTypes, false)(value, pathStr); len(errs) > 0 {
			for _, err := range errs {
				diags = append(diags, diag.FromErr(err)...)
			}
		}

		return diags
	}
}

var supportedSandboxFileTypes = []string{
	"FTCATEGORY_HTA", "FTCATEGORY_RAR", "FTCATEGORY_FLASH", "FTCATEGORY_TAR", "FTCATEGORY_JAVA_APPLET",
	"FTCATEGORY_WINDOWS_EXECUTABLES", "FTCATEGORY_WINDOWS_LIBRARY", "FTCATEGORY_MS_EXCEL", "FTCATEGORY_MS_RTF",
	"FTCATEGORY_APK", "FTCATEGORY_VISUAL_BASIC_SCRIPT", "FTCATEGORY_MS_POWERPOINT", "FTCATEGORY_P7Z",
	"FTCATEGORY_SCZIP", "FTCATEGORY_BZIP2", "FTCATEGORY_POWERSHELL", "FTCATEGORY_MS_WORD",
	"FTCATEGORY_ZIP", "FTCATEGORY_PDF_DOCUMENT",
}

func validateSandboxPolicyCategories() schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		value, ok := i.(string)
		if !ok {
			return diag.Diagnostics{
				{
					Severity: diag.Error,
					Summary:  "Expected type to be string",
					Detail:   "Type assertion failed, expected string type for Sandbox Policy Categories validation",
				},
			}
		}

		// Convert the cty.Path to a string representation
		pathStr := fmt.Sprintf("%+v", path)

		// Use StringInSlice from helper/validation package
		var diags diag.Diagnostics
		if _, errs := validation.StringInSlice(supportedSandboxPolicyCategories, false)(value, pathStr); len(errs) > 0 {
			for _, err := range errs {
				diags = append(diags, diag.FromErr(err)...)
			}
		}

		return diags
	}
}

var supportedSandboxPolicyCategories = []string{
	"ADWARE_BLOCK", "BOTMAL_BLOCK", "ANONYP2P_BLOCK", "RANSOMWARE_BLOCK", "OFFSEC_TOOLS_BLOCK", "SUSPICIOUS_BLOCK",
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
	"AI_ML": {"DENY_AI_ML_WEB_USE", "ALLOW_AI_ML_WEB_USE", "ISOLATE_AI_ML_WEB_USE", "CAUTION_AI_ML_WEB_USE", "DENY_AI_ML_UPLOAD",
		"ALLOW_AI_ML_UPLOAD", "DENY_AI_ML_SHARE", "ALLOW_AI_ML_SHARE", "DENY_AI_ML_DOWNLOAD", "ALLOW_AI_ML_DOWNLOAD", "DENY_AI_ML_DELETE",
		"ALLOW_AI_ML_DELETE", "DENY_AI_ML_INVITE", "ALLOW_AI_ML_INVITE", "DENY_AI_ML_CHAT", "ALLOW_AI_ML_CHAT", "DENY_AI_ML_CREATE",
		"ALLOW_AI_ML_CREATE", "DENY_AI_ML_RENAME", "ALLOW_AI_ML_RENAME"},

	"BUSINESS_PRODUCTIVITY": {"ALLOW_BUSINESS_PRODUCTIVITY_APPS", "BLOCK_BUSINESS_PRODUCTIVITY_APPS", "CAUTION_BUSINESS_PRODUCTIVITY_APPS", "ISOLATE_BUSINESS_PRODUCTIVITY_APPS"},

	"CONSUMER": {"ALLOW_CONSUMER_APPS", "BLOCK_CONSUMER_APPS", "CAUTION_CONSUMER_APPS", "ISOLATE_CONSUMER_APPS"},

	"CUSTOM_CAPP": {"BLOCK_CUSTOM_CAPP_USE", "ALLOW_CUSTOM_CAPP_USE", "ISOLATE_CUSTOM_CAPP_USE", "CAUTION_CUSTOM_CAPP_USE"},

	"DNS_OVER_HTTPS": {"ALLOW_DNS_OVER_HTTPS_USE", "DENY_DNS_OVER_HTTPS_USE"},

	"ENTERPRISE_COLLABORATION": {"ALLOW_ENTERPRISE_COLLABORATION_APPS", "ALLOW_ENTERPRISE_COLLABORATION_CHAT",
		"ALLOW_ENTERPRISE_COLLABORATION_UPLOAD", "ALLOW_ENTERPRISE_COLLABORATION_SHARE",
		"BLOCK_ENTERPRISE_COLLABORATION_APPS", "ALLOW_ENTERPRISE_COLLABORATION_EDIT",
		"ALLOW_ENTERPRISE_COLLABORATION_RENAME", "ALLOW_ENTERPRISE_COLLABORATION_CREATE",
		"ALLOW_ENTERPRISE_COLLABORATION_DOWNLOAD", "ALLOW_ENTERPRISE_COLLABORATION_HUDDLE",
		"ALLOW_ENTERPRISE_COLLABORATION_INVITE", "ALLOW_ENTERPRISE_COLLABORATION_MEETING",
		"ALLOW_ENTERPRISE_COLLABORATION_DELETE", "ALLOW_ENTERPRISE_COLLABORATION_SCREEN_SHARE",
		"BLOCK_ENTERPRISE_COLLABORATION_CHAT", "BLOCK_ENTERPRISE_COLLABORATION_UPLOAD",
		"BLOCK_ENTERPRISE_COLLABORATION_SHARE", "BLOCK_ENTERPRISE_COLLABORATION_EDIT",
		"BLOCK_ENTERPRISE_COLLABORATION_RENAME", "BLOCK_ENTERPRISE_COLLABORATION_CREATE",
		"BLOCK_ENTERPRISE_COLLABORATION_DOWNLOAD", "BLOCK_ENTERPRISE_COLLABORATION_DELETE",
		"BLOCK_ENTERPRISE_COLLABORATION_HUDDLE", "BLOCK_ENTERPRISE_COLLABORATION_INVITE",
		"BLOCK_ENTERPRISE_COLLABORATION_MEETING", "BLOCK_ENTERPRISE_COLLABORATION_SCREEN_SHARE",
		"ISOLATE_ENTERPRISE_COLLABORATION_APPS", "CAUTION_ENTERPRISE_COLLABORATION_APPS"},

	"FILE_SHARE": {"DENY_FILE_SHARE_VIEW", "ALLOW_FILE_SHARE_VIEW", "CAUTION_FILE_SHARE_VIEW", "DENY_FILE_SHARE_UPLOAD", "ALLOW_FILE_SHARE_UPLOAD", "ISOLATE_FILE_SHARE_VIEW",
		"DENY_FILE_SHARE_SHARE", "ALLOW_FILE_SHARE_SHARE", "DENY_FILE_SHARE_EDIT", "ALLOW_FILE_SHARE_EDIT", "DENY_FILE_SHARE_RENAME", "ALLOW_FILE_SHARE_RENAME",
		"DENY_FILE_SHARE_CREATE", "ALLOW_FILE_SHARE_CREATE", "DENY_FILE_SHARE_DOWNLOAD", "ALLOW_FILE_SHARE_DOWNLOAD", "DENY_FILE_SHARE_DELETE", "ALLOW_FILE_SHARE_DELETE",
		"DENY_FILE_SHARE_FORM_SHARE", "ALLOW_FILE_SHARE_FORM_SHARE", "DENY_FILE_SHARE_INVITE", "ALLOW_FILE_SHARE_INVITE"},

	"FINANCE": {"ALLOW_FINANCE_USE", "CAUTION_FINANCE_USE", "DENY_FINANCE_USE", "ISOLATE_FINANCE_USE"},

	"HEALTH_CARE": {"ALLOW_HEALTH_CARE_USE", "CAUTION_HEALTH_CARE_USE", "DENY_HEALTH_CARE_USE", "ISOLATE_HEALTH_CARE_USE"},

	"HOSTING_PROVIDER": {"DENY_HOSTING_PROVIDER_USE", "ALLOW_HOSTING_PROVIDER_USE", "ISOLATE_HOSTING_PROVIDER_USE", "CAUTION_HOSTING_PROVIDER_USE",
		"DENY_HOSTING_PROVIDER_CREATE", "ALLOW_HOSTING_PROVIDER_CREATE", "DENY_HOSTING_PROVIDER_EDIT", "ALLOW_HOSTING_PROVIDER_EDIT",
		"DENY_HOSTING_PROVIDER_DOWNLOAD", "ALLOW_HOSTING_PROVIDER_DOWNLOAD", "DENY_HOSTING_PROVIDER_DELETE", "ALLOW_HOSTING_PROVIDER_DELETE",
		"DENY_HOSTING_PROVIDER_MOVE", "ALLOW_HOSTING_PROVIDER_MOVE"},

	"HUMAN_RESOURCES": {"ALLOW_HUMAN_RESOURCES_USE", "CAUTION_HUMAN_RESOURCES_USE", "DENY_HUMAN_RESOURCES_USE", "ISOLATE_HUMAN_RESOURCES_USE"},

	"INSTANT_MESSAGING": {"ALLOW_CHAT", "ALLOW_FILE_TRANSFER_IN_CHAT", "BLOCK_CHAT", "BLOCK_FILE_TRANSFER_IN_CHAT", "CAUTION_CHAT", "CAUTION_FILE_TRANSFER_IN_CHAT", "ISOLATE_CHAT"},

	"IT_SERVICES": {"ALLOW_IT_SERVICES_USE", "CAUTION_LEGAL_USE", "DENY_IT_SERVICES_USE", "ISOLATE_IT_SERVICES_USE"},

	"LEGAL": {"ALLOW_LEGAL_USE", "CAUTION_LEGAL_USE", "DENY_LEGAL_USE", "ISOLATE_LEGAL_USE"},

	"SALES_AND_MARKETING": {"ALLOW_SALES_MARKETING_APPS", "BLOCK_SALES_MARKETING_APPS", "CAUTION_SALES_MARKETING_APPS", "ISOLATE_SALES_MARKETING_APPS"},

	"STREAMING_MEDIA": {"BLOCK_STREAMING_VIEW_LISTEN", "ALLOW_STREAMING_VIEW_LISTEN", "CAUTION_STREAMING_VIEW_LISTEN", "BLOCK_STREAMING_UPLOAD", "ALLOW_STREAMING_UPLOAD", "ISOLATE_STREAMING_VIEW_LISTEN"},

	"SOCIAL_NETWORKING": {"ALLOW_SOCIAL_NETWORKING_CHAT", "ALLOW_SOCIAL_NETWORKING_COMMENT", "ALLOW_SOCIAL_NETWORKING_CREATE", "ALLOW_SOCIAL_NETWORKING_EDIT", "ALLOW_SOCIAL_NETWORKING_POST", "ALLOW_SOCIAL_NETWORKING_SHARE", "ALLOW_SOCIAL_NETWORKING_UPLOAD",
		"ALLOW_SOCIAL_NETWORKING_VIEW", "BLOCK_SOCIAL_NETWORKING_CHAT", "BLOCK_SOCIAL_NETWORKING_COMMENT", "BLOCK_SOCIAL_NETWORKING_CREATE", "BLOCK_SOCIAL_NETWORKING_EDIT", "BLOCK_SOCIAL_NETWORKING_POST",
		"BLOCK_SOCIAL_NETWORKING_SHARE", "BLOCK_SOCIAL_NETWORKING_UPLOAD", "BLOCK_SOCIAL_NETWORKING_VIEW", "CAUTION_SOCIAL_NETWORKING_POST", "CAUTION_SOCIAL_NETWORKING_VIEW", "ISOLATE_SOCIAL_NETWORKING_VIEW"},

	"SYSTEM_AND_DEVELOPMENT": {"BLOCK_SYSTEM_DEVELOPMENT_APPS", "ALLOW_SYSTEM_DEVELOPMENT_APPS", "ISOLATE_SYSTEM_DEVELOPMENT_APPS", "BLOCK_SYSTEM_DEVELOPMENT_UPLOAD", "ALLOW_SYSTEM_DEVELOPMENT_UPLOAD",
		"CAUTION_SYSTEM_DEVELOPMENT_APPS", "BLOCK_SYSTEM_DEVELOPMENT_CREATE", "ALLOW_SYSTEM_DEVELOPMENT_CREATE", "BLOCK_SYSTEM_DEVELOPMENT_EDIT", "ALLOW_SYSTEM_DEVELOPMENT_EDIT",
		"BLOCK_SYSTEM_DEVELOPMENT_SHARE", "ALLOW_SYSTEM_DEVELOPMENT_SHARE", "BLOCK_SYSTEM_DEVELOPMENT_COMMENT", "ALLOW_SYSTEM_DEVELOPMENT_COMMENT", "BLOCK_SYSTEM_DEVELOPMENT_REACTION",
		"ALLOW_SYSTEM_DEVELOPMENT_REACTION"},

	"WEBMAIL": {"BLOCK_WEBMAIL_VIEW", "ALLOW_WEBMAIL_VIEW", "CAUTION_WEBMAIL_VIEW", "BLOCK_WEBMAIL_ATTACHMENT_SEND",
		"ALLOW_WEBMAIL_ATTACHMENT_SEND", "CAUTION_WEBMAIL_ATTACHMENT_SEND", "ALLOW_WEBMAIL_SEND", "BLOCK_WEBMAIL_SEND", "ISOLATE_WEBMAIL_VIEW"},
}

func validateDnsRuleRequestTypes() schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		value, ok := i.(string)
		if !ok {
			return diag.Diagnostics{
				{
					Severity: diag.Error,
					Summary:  "Expected type to be string",
					Detail:   "Type assertion failed, expected string type for DNS Request Types validation",
				},
			}
		}

		// Convert the cty.Path to a string representation
		pathStr := fmt.Sprintf("%+v", path)

		// Use StringInSlice from helper/validation package
		var diags diag.Diagnostics
		if _, errs := validation.StringInSlice(supportedDnsRuleRequestTypes, false)(value, pathStr); len(errs) > 0 {
			for _, err := range errs {
				diags = append(diags, diag.FromErr(err)...)
			}
		}

		return diags
	}
}

var supportedDnsRuleRequestTypes = []string{
	"A", "NS", "MD", "MF", "CNAME", "SOA", "MB", "MG", "MR", "NULL",
	"WKS", "PTR", "HINFO", "MINFO", "MX", "TXT", "RP", "AFSDB",
	"X25", "ISDN", "RT", "NSAP", "NSAP_PTR", "SIG", "KEY", "PX",
	"GPOS", "AAAA", "LOC", "NXT", "EID", "NIMLOC", "SRV", "ATMA",
	"NAPTR", "KX", "CERT", "A6", "DNAME", "SINK", "OPT", "APL",
	"DS", "SSHFP", "PSECKEF", "RRSIG", "NSEC", "DNSKEY",
	"DHCID", "NSEC3", "NSEC3PARAM", "TLSA", "HIP", "NINFO",
	"RKEY", "TALINK", "CDS", "CDNSKEY", "OPENPGPKEY", "CSYNC",
	"ZONEMD", "SVCB", "HTTPS",
}

func validateBlockResponseCode(val interface{}, key string) ([]string, []error) {
	supportedBlockResponseCode := []string{
		"ANY", "NONE", "FORMERR", "SERVFAIL", "NXDOMAIN", "NOTIMP", "REFUSED",
		"YXDOMAIN", "YXRRSET", "NXRRSET", "NOTAUTH", "NOTZONE", "BADVERS",
		"BADKEY", "BADTIME", "BADMODE", "BADNAME", "BADALG", "BADTRUNC",
		"UNSUPPORTED", "BYPASS", "INT_ERROR", "SRV_TIMEOUT", "EMPTY_RESP",
		"REQ_BLOCKED", "ADMIN_DROP", "WCDN_TIMEOUT", "IPS_BLOCK", "FQDN_RESOLV_FAIL",
	}

	// Convert the value to a string
	v := val.(string)

	// Check if the value is in the supported list
	for _, supported := range supportedBlockResponseCode {
		if v == supported {
			return nil, nil // Valid value, no errors
		}
	}

	// If not found, return an error
	return nil, []error{fmt.Errorf("%q is not a valid block_response_code. Supported values are: %v", v, supportedBlockResponseCode)}
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

var supportedCloudApplications = []string{
	"ABCNEWS", "ACCUWEATHER", "ACER", "ACROBAT", "ACROBATCONNECT", "ACTIVESYNC", "ADC", "ADDICTINGGAMES", "ADGUARD", "ADIDAS",
	"ADNSTREAM", "ADOBE_CONNECT", "ADOBE_CREATIVE_CLOUD", "ADOBE_DOCUMENT_CLOUD", "ADOBE_MEETING_RC", "ADOBE_UPDATE", "ADP",
	"ADRIVE", "ADULTADWORLD", "ADULTFRIENDFINDER", "ADVOGATO", "AGAME", "AGORA_IO", "AH_CDN", "AIAIGAME", "AILI", "AIM", "AIMINI",
	"AIMS", "AIM_EXPRESS", "AIM_TRANSFER", "AIRAIM", "AIZHAN", "AKAMAI", "ALEXA", "ALIBABA", "ALIMAMA", "ALIPAY", "ALJAZEERA",
	"ALLOCINE", "ALTIRIS", "AMAZON", "AMAZONVIDEO", "AMAZON_AWS", "AMAZON_CHIME", "AMAZON_CLD_DRIVE", "AMAZON_MP3", "AMAZON_VIDEO",
	"AMEBA", "AMERICANEXPRESS", "AMIE_STREET", "ANOBII", "ANSWERS", "ANZHI", "AOL", "AOL_VIDEO", "APNS", "APOLLO", "APPLE",
	"APPLEDAILY", "APPLEJUICE", "APPLE_AIRPORT", "APPLE_MAPS", "APPLE_SIRI", "APPLE_UPDATE", "APPSHOPPER", "APPSTORE",
	"APP_050PLUS", "APP_CH", "ARCHIVE", "ARENADAEMON", "ARES", "ARMORGAMES", "ASIAE", "ASK", "ASMALLWORLD", "ATHLINKS", "ATT",
	"ATWIKI", "AUCTION", "AUFEMININ", "AUONE", "AVATARS_UNITED", "AVG_UPDATE", "AVIRA_UPDATE", "AVOIDR", "BABELGUM", "BABYCENTER",
	"BABYHOME", "BACKPACKERS", "BADONGO", "BADOO", "BAIDU", "BAIDUYUNDNS", "BAIDU_PLAYER", "BAIDU_TIEBA", "BAIKE", "BAOFENG",
	"BARRACUDA_VPN", "BBC", "BBC_PLAYER", "BEANFUN", "BEBO", "BERNIAGA", "BGP", "BIGADDA", "BIGLOBE_NE", "BIGTENT", "BIGUPLOAD",
	"BIIP", "BING", "BITDEFENDER_UPDATE", "BITTORRENT", "BLACKBERRY", "BLACKPLANET", "BLAHDNS", "BLIP_TV", "BLOCKBUSTER", "BLOGDETIK",
	"BLOGGER", "BLOGIMG", "BLOGSPOT", "BLOGSTER", "BLOKUS", "BLOOMBERG", "BLUBSTER", "BLUEJAYFILMS", "BMFF", "BOLT", "BONPOO",
	"BOSTON", "BOXNET", "BRIGHTSPCACE", "BRIGHTTALK", "BROWSEC_VPN", "BUGS", "BUSINESSWEEK", "BUSINESSWEEKLY", "BUZZFEED",
	"BUZZNET", "BYPASSTHAT", "CAFEMOM", "CAM4", "CAMPFIRE", "CAMPINA", "CAMZAP", "CAPITALONE", "CAPWAP", "CARE2", "CARTOONNETWORK",
	"CCM", "CCTV_VOD", "CDISCOUNT", "CDP", "CELLUFUN", "CFT", "CHANNEL4", "CHAP", "CHINACOM", "CHINACOMCN", "CHINANEWS",
	"CHINATIMES", "CHINAZ", "CHINA_AIRLINES", "CHOSUN", "CHOSUN_DAILY", "CHROME_UPDATE", "CITRIX_ONLINE", "CJ_MALL", "CK101",
	"CLASSMATES", "CLEANBROWSING", "CLEARCASE", "CLIP2NET", "CLOOB", "CLOUDFLARE", "CLOUDFLARE_DNS", "CLOUDME", "CLUBIC", "CNET",
	"CNET_TV", "CNN", "CNTV", "CNYES", "CNZZ", "COCOLOG_NIFTY", "CODE42", "COLLEGE_BLENDER", "COMCAST", "COMCAST_DNS", "COMM",
	"CONCUR", "CONDUIT", "COSTAIN", "COUCH_SURFING", "COUPANG", "CRAIGSLIST", "CRASHPLAN", "CROCKO", "CSDN", "CSTRIKE", "CSWG",
	"CTRIP", "CYMRU", "DAILYMAIL", "DAILYMOTION", "DAILY_STRENGH", "DANGDANG", "DAUM", "DAVIDOV", "DB2", "DCERPC", "DDSTUDIOS",
	"DEBIAN_UPDATE", "DECAYENNE", "DEEZER", "DELICIOUS", "DEPOSITFILES", "DETIK", "DETIKNEWS", "DEVIANT_ART", "DEVICESCAPE", "DHCP",
	"DIAMETER", "DICT", "DIGG", "DIGITALVERSE", "DIINO", "DIMP", "DIOP", "DIRECTCONNECT", "DIRECTDOWNLOADLINKS", "DIRECTV",
	"DISABOOM", "DISNEY", "DIVSHARE", "DMM_CO", "DNS", "DNSHOP", "DNSTUN_ADS", "DNSTUN_AUTH", "DNSTUN_BUSINESS", "DNSTUN_CONFERENCE",
	"DNSTUN_DATABASE", "DNSTUN_ENTERPRISE", "DNSTUN_FILESERVER_TRANSFER", "DNSTUN_GAMING", "DNSTUN_GOOD", "DNSTUN_IM",
	"DNSTUN_IMAGEHOST", "DNSTUN_MALICIOUS", "DNSTUN_MALWARE", "DNSTUN_MAPPSTORE", "DNSTUN_MOBILE", "DNSTUN_NETMGMT", "DNSTUN_P2P",
	"DNSTUN_REMOTE", "DNSTUN_SOCIAL", "DNSTUN_STREAMING", "DNSTUN_TUNNELING", "DNSTUN_UNKNOWN", "DNSTUN_WEBSEARCH", "DNS_OVER_HTTPS",
	"DOCUSIGN", "DOH", "DOH_UNKNOWN", "DOL2DAY", "DONGA", "DONTSTAYIN", "DOORBLOG", "DOTDASH", "DOUBAN", "DOUBLECLICK_ADS",
	"DRAUGIEM", "DRDA", "DREAMWIZ", "DROPBOX", "DRUPAL", "DUOWAN", "DXC", "DYNAMICINTRANET", "EARTHCAM", "EASTMONEY", "EASYTRAVEL",
	"EBAY", "EBUDDY", "EDONKEY", "EGNYTE", "EIGRP", "ELFTOWN", "ELLE_TW", "EONS", "EPERNICUS", "EPIC_BROWSER_UPDATE",
	"EPIC_BROWSER_VPN", "EPM", "EROOM_NET", "ERSPAN", "ESET", "ESNIPS", "ESPN", "ETAO", "ETHERIP", "ETSI_LI", "ETTODAY",
	"EVASIVE_PROTOCOL", "EVERNOTE", "EVERQUEST", "EVE_ONLINE", "EVONY", "EXBLOG", "EXPEDIA", "EXPERIENCE_PROJECT", "EXPLOROO",
	"EYEJOT", "EYNY", "EZFLY", "EZTRAVEL", "FACEBOOK", "FACEBOOK_APPS", "FACEBOOK_GAMES", "FACEBOOK_MAIL", "FACEPARTY", "FACES",
	"FACETIME", "FASHIONGUIDE", "FASTMAIL", "FC2", "FETLIFE", "FIFA", "FILEFLYER", "FILEMAKER_PRO", "FILER_CX", "FILESANYWHERE",
	"FILESTUBE", "FILETOPIA", "FILE_DROPPER", "FILE_HOST", "FILLOS_DE_GALICIA", "FILMAFFINITY", "FIREFOX_UPDATE", "FLASH",
	"FLASHPLUGIN_UPDATE", "FLEDGEWING", "FLICKR", "FLIPKART", "FLIPKART_BOOKS", "FLIXSTER", "FLUMOTION", "FLUXIOM", "FLYINGMAG",
	"FLY_PROXY", "FOGBUGZ", "FOLDINGATHOME", "FORTUNECHINA", "FOTKI", "FOTOLOG", "FOURSHARED", "FOURSQUARE", "FOXMOVIES", "FOXNEWS",
	"FOXSPORTS", "FOXY", "FRANCETELECOM", "FREEBSD_UPDATE", "FREEETV", "FRIENDSTER", "FRIENDS_REUNITED", "FRIENDVOX", "FRING",
	"FRUHSTUCKSTREFF", "FSECURE_UPDATE", "FTP", "FTPS", "FTPS_DATA", "FTP_DATA", "FUBAR", "FUNSHION", "GAIAONLINE", "GAMEBASE_TW",
	"GAMERDNA", "GAMER_TW", "GAMESMOMO", "GAMES_CO", "GANJI", "GARP", "GATHER", "GAYS", "GBRIDGE", "GCM", "GENESISMISSIONARYBAPTISTCHURCH",
	"GENI", "GE_PROCIFY", "GFAN", "GIGAUP", "GIOP", "GIOPS", "GLASSDOOR", "GLIDE", "GMAIL", "GMAIL_BASIC", "GMAIL_MOBILE", "GMARKET",
	"GMX", "GMX_DNSTUN", "GNUNET", "GNUTELLA", "GOBOOGY", "GOGOYOKO", "GOHAPPY", "GOMTV_VOD", "GOODREADS", "GOOGLE",
	"GOOGLEANALYTICS", "GOOGLE_ADS", "GOOGLE_APPENGINE", "GOOGLE_CACHE", "GOOGLE_DESKTOP", "GOOGLE_DNS", "GOOGLE_DNSTUN",
	"GOOGLE_DOCS", "GOOGLE_DRIVE", "GOOGLE_EARTH", "GOOGLE_GEN", "GOOGLE_GROUPS", "GOOGLE_KEEP", "GOOGLE_MAPS", "GOOGLE_PICASA",
	"GOOGLE_PLAY", "GOOGLE_PLUS", "GOOGLE_SKYMAP", "GOOGLE_STADIA", "GOOGLE_TOOLBAR", "GOOGLE_TRANSLATE", "GOOGLE_VIDEO", "GOO_NE",
	"GOTODEVICE", "GOTOMEETING", "GOTOMYPC", "GOUGOU", "GRATISINDO", "GRE", "GREE", "GRONO", "GROOVESHARK", "GROUPON",
	"GROUPSERVICES", "GROUPWISE", "GSSHOP", "GSTATIC", "GTALK", "GTP", "GTPV2", "GUDANGLAGU", "GYAO", "H225", "H245",
	"H248_BINARY", "H248_TEXT", "HABBO", "HALFLIFE", "HAMACHI", "HANGAME", "HANKOOKI", "HANKYUNG", "HAO123", "HARDSEXTUBE",
	"HATENA_NE", "HBO_GO", "HERALDM", "HERE", "HEXATECH", "HEXUN", "HGTV", "HI5", "HIDEMAN_VPN", "HIGHTAIL", "HINET_GAMES",
	"HITACHI", "HOEFLIGER", "HOFF", "HONDA", "HOSPITALITY_CLUB", "HOTFILE", "HOTLINE", "HOTMAIL", "HOTSPOT_SHIELD", "HOUSEPARTY",
	"HOVRS", "HOWCAST", "HOWSTUFFWORKS", "HOXX_VPN", "HSRP", "HTTP", "HTTP2", "HTTPS", "HTTPTUNNEL", "HTTP_PROXY", "HUDONG",
	"HULU", "HYVES", "IAPP", "IAX", "IBACKUP", "IBIBO", "ICA", "ICALL", "ICAP", "ICECAST", "ICLOUD", "ICMP", "ICQ2GO", "IDENT",
	"IDEXXI", "IFENG", "IFENG_FINANCE", "IFILE_IT", "IGMP", "IHEARTRADIO", "IIOP", "IKEA", "ILOVEIM", "IMAGESHACK", "IMDB",
	"IMEEM", "IMEET", "IMESH", "IMGUR", "IMP", "IMPRESS", "IMRWORLDWIDE", "IMVU", "INDABA_MUSIC", "INDIATIMES", "INDONETWORK",
	"INDOWEBSTER", "INFOARMOR", "INFORMIX", "INFOSYSBPM", "INILAH", "INSAGS", "INSTAGRAM", "INTALKING", "INTECH", "INTERNATIONS",
	"INTERPARK", "INTUIT", "IOS_APPSTORE", "IOS_OTA_UPDATE", "IP", "IP6", "IPASS", "IPERF", "IPSEC", "IPV6CP", "IPV6TEST", "IPXRIP",
	"IQIYI", "IRC", "IRCO", "IRCS", "IRC_GALLERIA", "IRC_TRANSFER", "ISAKMP", "ITALKI", "ITSMY", "ITUNES", "IWIW", "I_GAMER",
	"I_PART", "JABBER", "JABBER_TRANSFER", "JAIKU", "JAJAH", "JAMMERDIRECT", "JANGO", "JAVA_UPDATE", "JEDI", "JINGDONG", "JNE",
	"JOBSTREET", "JOONGANG_DAILY", "JUBII", "JUSTIN_TV", "KAIOO", "KAIXIN_CHAT", "KAKAKU", "KAKAOTALK", "KANKAN", "KAPANLAGI",
	"KAROSGAME", "KASKUS", "KASPERSKY", "KASPERSKY_UPDATE", "KAZAA", "KBS", "KCTEST", "KEEZMOVIES", "KELLYSERVICES", "KEMENKUMHAM",
	"KHAN", "KICKASSTORRENTS", "KIK", "KIWIBOX", "KLATENCOR", "KOMPAS", "KOMPASIANA", "KONAMINET", "KOOLIM", "KPN_TUNNEL",
	"KPROXY", "KR0", "KRB5", "KU6", "KUGOU", "KUXUN", "L2TP", "LADY8844", "LAREDOUTE", "LASTFM", "LASTPASS", "LATIV", "LCP",
	"LDAP", "LDAPS", "LDBLOG", "LEAPFILE", "LEARNZOLASUITE", "LEBONCOIN", "LETV", "LEVEL3", "LG_ESHOP", "LIBERO_VIDEO", "LIBRARYTHING",
	"LIFEKNOT", "LINE", "LINEAGE2", "LINKEDIN", "LINTASBERITA", "LIONAIR", "LIONTRAVEL", "LISTOGRAFY", "LIVEDOOR", "LIVEINTERNET",
	"LIVEJOURNAL", "LIVEMAIL_MOBILE", "LIVEMOCHA", "LIVE_GROUPS", "LIVE_MEETING", "LIVE_MESH", "LIVINGSOCIAL", "LOCATION_IQ", "LOOP",
	"LOTOUR", "LOTTE", "LOTUS_LIVE", "LOTUS_SAMETIME", "LQR", "LUNARSTORM", "LVPING", "MAIL2000", "MAILRU", "MAILSHELL", "MAIL_189",
	"MAKTOOB", "MANDRIVA_UPDATE", "MANGOCITY", "MANOLITO", "MAPQUEST", "MARLABS", "MASHABLE", "MASHARE", "MATCH", "MATTEL", "MBC",
	"MBN", "MCAFEE", "MEDIAFIRE", "MEETIN", "MEETINGPLACE", "MEETME", "MEETTHEBOSS", "MEETUP", "MEGA", "MEGAPROXY", "MEGAVIDEO",
	"MESSENGERFX", "METACAFE", "MGCP", "MIBBIT", "MICRON", "MICROSOFTLIVEMEETING", "MICROSOFT_DNSTUN", "MIMEDIA", "MIMP", "MITALK",
	"MIXI", "MK", "MMS", "MOBAGE", "MOBILE01", "MOBILE_IP", "MOBILE_ME", "MOBILINK", "MOCOSPACE", "MODSECURITY", "MOG", "MOGULUS",
	"MOMOSHOP", "MONEX", "MONEYDJ", "MONEY_163", "MONSTER", "MOP", "MOTIONBOX", "MOUNT", "MOUTHSHUT", "MOXA", "MOXA_ASPP", "MOZILLA",
	"MPEGTS", "MPLS", "MPLUS_MESSENGER", "MPQUEST", "MQ", "MRSHMC", "MSN", "MSNMOBILE", "MSN_GROUPS", "MSN_SEARCH", "MSN_VIDEO",
	"MSRP", "MSRPC", "MS_COMMUNICATOR", "MS_DFSR", "MS_SSAS", "MT", "MTV", "MULTIPLY", "MULTIUPLOAD", "MUSICA", "MUTE", "MXIT",
	"MYANIMELIST", "MYCHURCH", "MYHERITAGE", "MYL", "MYLIFE", "MYSPACE", "MYSQL", "MYVIDEO", "MYVIDEODE", "MYWEBSEARCH", "MYYEARBOOK",
	"MY_YAHOO", "NAPSTER", "NASA", "NASZA_KLASA", "NATECYWORLD", "NATIONALGEOGRAPHIC", "NATIONALLOTTERY", "NAVER", "NBA",
	"NBA_CHINA", "NDUOA", "NEND", "NESSUS", "NETBIOS", "NETBSD_UPDATE", "NETFLIX", "NETFLOW", "NETLOAD", "NETLOG", "NETMARBLE",
	"NETMEETING_ILS", "NETTBY", "NETVIEWER", "NEXIAN", "NEXON", "NEXOPIA", "NEXTDNS", "NFL", "NFS", "NGO_POST", "NICONICO_DOUGA",
	"NIFTY", "NIKE", "NIKKEI", "NIMBUZZ_WEB", "NING", "NLOCKMGR", "NLSP", "NNTP", "NNTPS", "NOD32_UPDATE", "NOKIA_OVI", "NORTON_UPDATE",
	"NOTAPPLICABLE", "NOT_AVAILABLE", "NOWNEWS", "NSPI", "NTP", "NTTDATA", "NTV", "NYDAILYNEWS", "NYTIMES", "OBJECTIVEFS", "OCSP",
	"ODNOKLASSNIKI", "OFFICE365", "OFFICEDEPOT", "OFW_HTTPS_BYPASS", "OFW_HTTP_BYPASS", "OFW_ICMP_BYPASS", "OFW_TCP_BYPASS",
	"OFW_UDP_BYPASS", "OICQ", "OKEZONE", "OKWAVE", "ONEDRIVE", "ONLINEDOWN", "OOVOO", "OOYALA", "OPENBSD_UPDATE", "OPENDNS",
	"OPENDNS_DNSTUN", "OPENFT", "OPENGW", "OPENVPN", "OPEN_DIARY", "OPERA_VPN", "ORANGEMAIL", "ORB", "ORKUT", "OSPF", "OTSUKA",
	"OUTLOOK", "OWA", "PAIPAI", "PALTALK", "PALTALK_AUDIO", "PALTALK_TRANSFER", "PALTALK_VIDEO", "PANDA_UPDATE", "PANDO", "PANDORA",
	"PANDORA_TV", "PAP", "PARADIGMGEO", "PARTNERUP", "PARTY_POKER", "PASSPORTSTAMP", "PASTEBIN", "PAYEASY", "PCGAMES", "PCHOME",
	"PCLADY", "PCONLINE", "PEERCAST", "PENGYOU", "PEOPLE", "PERFORCE", "PERFSPOT", "PHOTOBUCKET", "PIMANG", "PINGSTA", "PINTEREST",
	"PIXIV", "PIXNET", "PLAXO", "PLAYSTATION", "PLURK", "POCO", "POGO", "POKER_STARS", "POP3", "POP3S", "PORNHUB", "PORTMAP",
	"POSTGRES", "POWERDNS", "PPFILM", "PPLIVE", "PPP", "PPPOE", "PPSTREAM", "PPTP", "PPTV", "PRESENT", "PRICEMINISTER", "PRICERUNNER",
	"PRIVAX", "PROXEASY", "PSIPHON", "PSN", "PWE", "Q931", "QIK_VIDEO", "QQ", "QQDOWNLOAD", "QQLIVE", "QQMUSIC", "QQSTREAM",
	"QQ_BLOG", "QQ_FINANCE", "QQ_GAMES", "QQ_LADY", "QQ_MAIL", "QQ_NEWS", "QQ_TRANSFER", "QQ_WEB", "QQ_WEIBO", "QUAD9", "QUAKE",
	"QUARTERLIFE", "QUIC", "QUNAR", "QVOD", "QY", "QZONE", "RACKSPACE", "RADIKO", "RADIUS", "RADIUSIM", "RADIX", "RADMIN", "RAKUTEN",
	"RAMBLER", "RAMBLER_WEBMAIL", "RANDSTAD", "RAPIDSHARE", "RAVELRY", "RDP", "RDT", "REALTOR", "REDIFF", "REDTUBE", "REEBOK",
	"RENREN", "REPUBLIKA", "RESEARCHGATE", "REUTERS", "REVERBNATION", "REVERSO", "RFB", "RHAPSODY", "RINGCENTRAL", "RIP1", "RIP2",
	"RIPNG1", "RLOGIN", "RMI_IIOP", "ROBLOX", "RPC", "RQUOTA", "RSH", "RSS", "RSTAT", "RSVP", "RSYNC", "RTCP", "RTL", "RTMP", "RTP",
	"RTSP", "RUBYFISH", "RUNESCAPE", "RUSERS", "RUTEN", "RVBDCC", "RYANAIR", "RYDER", "RYZE", "S7COMM_PLUS", "SABERINDO", "SAKURA_NE",
	"SALESFORCE", "SALESFORCE_DNSTUN", "SAMSUNG_APPS", "SAYCLUB", "SBS", "SCCP", "SCIENCESTAGE", "SCISPACE", "SCRIBD", "SCTP", "SDO",
	"SECONDLIFE", "SECUREDNS", "SECURESERVER", "SEEQPOD", "SEESAA", "SEESMIC", "SEGYE", "SENDERBASE", "SENDSPACE", "SENSIC",
	"SEOUL_NEWS", "SERVICE_NOW", "SFR", "SHARE", "SHAREFILE", "SHAREPOINT", "SHAREPOINT_ADMIN", "SHAREPOINT_BLOG", "SHAREPOINT_CALENDAR",
	"SHAREPOINT_DOCUMENT", "SHAREPOINT_ONLINE", "SHELFARI", "SHINHAN", "SHOPIFY", "SHOUTCAST", "SHOWMYPC", "SHUTTERFLY", "SIEBEL_CRM",
	"SILVERLIGHT", "SINA", "SINA_BLOG", "SINA_FINANCE", "SINA_NEWS", "SINA_VIDEO", "SINA_WEIBO", "SIP", "SKY", "SKYBLOG", "SKYCN",
	"SKYPE", "SKYPE_FOR_BUSINESS", "SKYVPN", "SKY_PLAYER", "SLACK", "SLACKER", "SLIDESHARE", "SLINGBOX", "SLSK", "SMB", "SMTP",
	"SNMP", "SOAP", "SOCIALTV", "SOCKS2HTTP", "SOCKS4", "SOCKS5", "SODEXO", "SOFT4FUN", "SOFTBANK", "SOGOU", "SOHU", "SOHU_BLOG",
	"SOKU", "SONET_NE", "SONGMOUNTAINFINEART", "SONICO", "SONICWALL", "SONMP", "SOPCAST", "SOPHOS", "SORIBADA", "SOSO", "SOUFUN",
	"SOUNDCLOUD", "SOUTHWEST", "SPDY", "SPEEDTEST", "SPIEGEL", "SPORTCHOSUN", "SPORTSILLUSTRATED", "SPORTSSEOUL", "SPOTFLUX",
	"SPOTIFY", "SPOTIFYWLAN", "SPRINGTECH_VPN", "SPRINT", "SQLI", "SQUIRRELMAIL", "SRVLOC", "SSDP", "SSH", "SSL", "STACKPATH",
	"STAFABAND", "STAGEVU", "STAYFRIENDS", "STEAM", "STICKAM", "STOCKQ", "STP", "STREAMAUDIO", "STUDIVZ", "STUMBLEUPON", "STUN",
	"SUGARCRM", "SUGARSYNC", "SUNING", "SUPPERSOCCER", "SURROGAFIER", "SURVEYMONKEY", "SVTPLAY", "SYBASE", "SYNC", "SYSLOG", "TABELOG",
	"TACACS_PLUS", "TAGGED", "TAGOO", "TAIWANLOTTERY", "TAKEDAPHARMA", "TAKU_FILE_BIN", "TALENTTROVE", "TALKBIZNOW", "TALTOPIA",
	"TANGO", "TAOBAO", "TARINGA", "TCHATCHE", "TCP", "TCP_OVER_DNS", "TDS", "TEACHERTUBE", "TEACHSTREET", "TEAMSPEAK", "TEAMSPEAK_V3",
	"TEAMVIEWER", "TECHINLINE", "TEL", "TELEGRAM", "TELNET", "TEMPOINTERAKTIF", "TENCENT", "TEREDO", "TF1", "TFTP", "TGBUS", "TGIN",
	"THEPIRATEBAY", "THOUGHTCATALOG", "THREE", "THREEGPP_LI", "THREEMINUTEWEBSITE", "THUGLAK", "THUNDER", "TIANYA", "TICKETMONSTER",
	"TIDALTV", "TINDER", "TISTORY", "TMALL", "TNS", "TNVIP", "TOADTEXTURE", "TOKBOX", "TOKOBAGUS", "TOR", "TORRENTDOWNLOADS", "TORRENTZ",
	"TOUTIAO", "TQ", "TRAVBUDDY", "TRAVELLERSPOINT", "TRAVELOCITY", "TRAVIAN", "TRENDMICRO_UPDATE", "TRIBE", "TRIBUNNEWS", "TROMBI",
	"TRUCKINSURANCE", "TU", "TUBE8", "TUCHONG", "TUDOU", "TUENTI", "TUMBLR", "TUNEIN", "TUNEWIKI", "TUNNELBEAR", "TV", "TV4PLAY",
	"TVANTS", "TVUPLAYER", "TWITPIC", "TWITTER", "UBUNTU_ONE", "UDN", "UDP", "ULTRASURF", "UNISYS", "UNIVISION", "UPLOADING",
	"USATODAY", "USEJUMP", "USL", "USTREAM", "UTP", "UUSEE", "VAKAKA", "VAMPIREFREAKS", "VEETLE", "VEOHTV", "VEVO", "VIADEO", "VIBER",
	"VIDEOBASH", "VIDEOSURF", "VIETBAO", "VIMEO", "VIRTUSA", "VIVANEWS", "VJC_COMP", "VJC_UNCOMP", "VKONTAKTE", "VMWARE",
	"VMWARE_HORIZON_VIEW", "VODDLER", "VOX", "VPN1_COM", "VPNOVERDNS", "VRRP", "VTUNNEL", "VYEW", "WAKOOPA", "WALLSTREETJOURNAL_CHINA",
	"WANDOUJIA", "WANGWANG", "WASABI", "WASHINGTONPOST", "WAT", "WAYN", "WEATHER", "WEAVERPUBLISHING", "WEBEX", "WEBINAR", "WEBSOCKET",
	"WEB_CRAWLER", "WECHAT", "WEOURFAMILY", "WERKENNTWEN", "WETRANSFER", "WFC", "WHATSAPP", "WIICONNECT24", "WIKIA", "WIKIPEDIA",
	"WINDOWSCENTRAL", "WINDOWSLIVE", "WINDOWSLIVESPACE", "WINDOWSMEDIA", "WINDOWS_AZURE", "WINDOWS_LIVE_SPACES", "WINDOWS_MARKETPLACE",
	"WINDOWS_UPDATE", "WINDSCRIBE", "WINMX", "WINNY", "WIXI", "WOORIBANK", "WORDPRESS", "WORKDAY", "WOW", "WRETCH", "WSTUNNEL",
	"WS_DISCOVERY", "X11", "XANGA", "XBOX", "XBOXLIVE", "XDADEV", "XDMCP", "XFS", "XHAMSTER", "XIAMI", "XING", "XINHUANET", "XLNET",
	"XL_WAP", "XL_WEBMAIL", "XMLRPC", "XM_RADIO", "XNXX", "XOT", "XREA", "XT3", "XUITE", "XVIDEOS", "XVIDEOSLIVE", "YAHOO",
	"YAHOO360PLUSVIETNAM", "YAHOOMAIL", "YAHOO_ANSWERS", "YAHOO_BIZ", "YAHOO_BUY", "YAHOO_DOUGA", "YAHOO_GAMES", "YAHOO_GEOCITIES",
	"YAHOO_GROUPS", "YAHOO_KOREA", "YAHOO_MAPS", "YAHOO_REALESTATE", "YAHOO_SCREEN", "YAHOO_SEARCH", "YAHOO_STOCK_TW", "YAHOO_TRAVEL",
	"YAMMER", "YANDEX", "YANDEXDISK", "YANDEX_MAIL", "YELP", "YESKY", "YIHAODIAN", "YMAIL_CLASSIC", "YMAIL_MOBILE", "YMSG",
	"YMSG_CONF", "YMSG_TRANSFER", "YMSG_VIDEO", "YOKA", "YOMIURI", "YOUDAO", "YOUKU", "YOUM7", "YOUMEO", "YOUNI", "YOUPORN",
	"YOURFILEHOST", "YOUSEEMORE", "YOUTUBE", "YOUTUBE_HD", "YPPASSWD", "YPSERV", "YPUPDATE", "YUGMA", "YUUGUU", "ZABBIX_AGENT",
	"ZATTOO", "ZENDESK", "ZENMATE", "ZIDDU", "ZIMBRA", "ZIMBRA_STANDARD", "ZOHO", "ZOHOCRM", "ZOHO_DB", "ZOHO_IM", "ZOHO_MEETING",
	"ZOHO_NOTEBOOK", "ZOHO_PEOPLE", "ZOHO_PLANNER", "ZOHO_SHARE", "ZOHO_SHEET", "ZOHO_SHOW", "ZOL", "ZOO", "ZOOM", "ZOOMINFO",
	"ZSCLOUD_HTTP_BYPASS", "ZSHARE", "ZUM", "ZVELO", "ZYNGA",
}

func validateCloudApplications() schema.SchemaValidateDiagFunc {
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
		if _, errs := validation.StringInSlice(supportedCloudApplications, false)(value, pathStr); len(errs) > 0 {
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

func validateFileTypes() schema.SchemaValidateDiagFunc {
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
		if _, errs := validation.StringInSlice(supportedFileTypes, false)(value, pathStr); len(errs) > 0 {
			for _, err := range errs {
				diags = append(diags, diag.FromErr(err)...)
			}
		}

		return diags
	}
}

var supportedFileTypes = []string{
	"ANY", "NONE", "FTCATEGORY_JAVASCRIPT", "FTCATEGORY_FLASH", "FTCATEGORY_JAVA_APPLET", "FTCATEGORY_HTA",
	"FTCATEGORY_HAR", "FTCATEGORY_ZIP", "FTCATEGORY_GZIP", "FTCATEGORY_TAR", "FTCATEGORY_BZIP2", "FTCATEGORY_RAR",
	"FTCATEGORY_STUFFIT", "FTCATEGORY_ISO", "FTCATEGORY_CAB", "FTCATEGORY_P7Z", "FTCATEGORY_SCZIP", "FTCATEGORY_DMG",
	"FTCATEGORY_PKG", "FTCATEGORY_NUPKG", "FTCATEGORY_MF", "FTCATEGORY_EGG", "FTCATEGORY_ALZ", "FTCATEGORY_LZ4",
	"FTCATEGORY_LZOP", "FTCATEGORY_ZST", "FTCATEGORY_RZIP", "FTCATEGORY_LZIP", "FTCATEGORY_LRZIP", "FTCATEGORY_DACT",
	"FTCATEGORY_ZPAQ", "FTCATEGORY_BH", "FTCATEGORY_B64", "FTCATEGORY_LZMA", "FTCATEGORY_XZ", "FTCATEGORY_FCL",
	"FTCATEGORY_ZIPX", "FTCATEGORY_CPIO", "FTCATEGORY_LZH", "FTCATEGORY_MP3", "FTCATEGORY_WAV", "FTCATEGORY_OGG_VORBIS",
	"FTCATEGORY_M3U", "FTCATEGORY_VPR", "FTCATEGORY_AAC", "FTCATEGORY_ADE", "FTCATEGORY_DB2", "FTCATEGORY_SQL",
	"FTCATEGORY_EDMX", "FTCATEGORY_FRM", "FTCATEGORY_ACCDB", "FTCATEGORY_DBF", "FTCATEGORY_VIRTUAL_HARD_DISK",
	"FTCATEGORY_DB", "FTCATEGORY_SDB", "FTCATEGORY_KDBX", "FTCATEGORY_DXL", "FTCATEGORY_WINDOWS_EXECUTABLES",
	"FTCATEGORY_MICROSOFT_INSTALLER", "FTCATEGORY_WINDOWS_LIBRARY", "FTCATEGORY_WINDOWS_LNK", "FTCATEGORY_PYTHON",
	"FTCATEGORY_POWERSHELL", "FTCATEGORY_VISUAL_BASIC_SCRIPT", "FTCATEGORY_MSP", "FTCATEGORY_REG", "FTCATEGORY_BAT",
	"FTCATEGORY_BASH_SCRIPTS", "FTCATEGORY_SHELL_SCRAP", "FTCATEGORY_DEB", "FTCATEGORY_APPX", "FTCATEGORY_MSC",
	"FTCATEGORY_ELF", "FTCATEGORY_MACH", "FTCATEGORY_DRV", "FTCATEGORY_GBA", "FTCATEGORY_SMD", "FTCATEGORY_XBEH",
	"FTCATEGORY_PSX", "FTCATEGORY_THREETWOX", "FTCATEGORY_NDS", "FTCATEGORY_BITMAP", "FTCATEGORY_PHOTOSHOP",
	"FTCATEGORY_WINDOWS_META_FORMAT", "FTCATEGORY_GIF", "FTCATEGORY_JPEG", "FTCATEGORY_PNG", "FTCATEGORY_WEBP",
	"FTCATEGORY_TIFF", "FTCATEGORY_DCM", "FTCATEGORY_THREEDM", "FTCATEGORY_KML", "FTCATEGORY_JPD", "FTCATEGORY_DNG",
	"FTCATEGORY_RWZ", "FTCATEGORY_GREENSHOT", "FTCATEGORY_IMG", "FTCATEGORY_HIGH_EFFICIENCY_IMAGE_FILES",
	"FTCATEGORY_AAF", "FTCATEGORY_OMFI", "FTCATEGORY_PLS", "FTCATEGORY_HLP", "FTCATEGORY_MDZ", "FTCATEGORY_MST",
	"FTCATEGORY_WINDOWS_SCRIPT_FILES", "FTCATEGORY_GRP", "FTCATEGORY_PIF", "FTCATEGORY_JOB", "FTCATEGORY_PSW",
	"FTCATEGORY_ONENOTE", "FTCATEGORY_CATALOG", "FTCATEGORY_NETMON", "FTCATEGORY_HIVE", "FTCATEGORY_APK",
	"FTCATEGORY_IPA", "FTCATEGORY_MOBILECONFIG", "FTCATEGORY_MS_POWERPOINT", "FTCATEGORY_MS_WORD",
	"FTCATEGORY_MS_EXCEL", "FTCATEGORY_MS_RTF", "FTCATEGORY_MS_MDB", "FTCATEGORY_MS_MSG", "FTCATEGORY_MS_PST",
	"FTCATEGORY_MS_VSIX", "FTCATEGORY_VSDX", "FTCATEGORY_OAB", "FTCATEGORY_OLM", "FTCATEGORY_MS_PUB",
	"FTCATEGORY_TNEF", "FTCATEGORY_ENCROFF", "FTCATEGORY_OPEN_OFFICE_DOC", "FTCATEGORY_OPEN_OFFICE_DRAWINGS",
	"FTCATEGORY_OPEN_OFFICE_PRESENTATIONS", "FTCATEGORY_OPEN_OFFICE_SPREADSHEETS", "FTCATEGORY_ENCRYPT",
	"FTCATEGORY_PDF_DOCUMENT", "FTCATEGORY_POSTSCRIPT", "FTCATEGORY_COMPILED_HTML_HELP", "FTCATEGORY_DWG",
	"FTCATEGORY_CGR", "FTCATEGORY_SLDPRT", "FTCATEGORY_TXT", "FTCATEGORY_UNK", "FTCATEGORY_IPT", "FTCATEGORY_XPS",
	"FTCATEGORY_CSV", "FTCATEGORY_STL", "FTCATEGORY_IQY", "FTCATEGORY_CERT", "FTCATEGORY_INTERNET_SIGNUP",
	"FTCATEGORY_PCAP", "FTCATEGORY_TTF", "FTCATEGORY_CRX", "FTCATEGORY_CER", "FTCATEGORY_DER", "FTCATEGORY_P7B",
	"FTCATEGORY_PEM", "FTCATEGORY_JKS", "FTCATEGORY_KEY", "FTCATEGORY_P12", "FTCATEGORY_CHEMDRAW_FILES",
	"FTCATEGORY_CML", "FTCATEGORY_BPL", "FTCATEGORY_CCC", "FTCATEGORY_CP", "FTCATEGORY_DEVFILE", "FTCATEGORY_MM",
	"FTCATEGORY_AES", "FTCATEGORY_WOFF2", "FTCATEGORY_STEP_FILES", "FTCATEGORY_RVT", "FTCATEGORY_EMF",
	"FTCATEGORY_PCD", "FTCATEGORY_INF", "FTCATEGORY_SAM", "FTCATEGORY_PMD", "FTCATEGORY_EOT", "FTCATEGORY_OPENXML",
	"FTCATEGORY_FODT", "FTCATEGORY_JOBOPTIONS", "FTCATEGORY_IDML", "FTCATEGORY_CXP", "FTCATEGORY_ENEX",
	"FTCATEGORY_OTF", "FTCATEGORY_LGX", "FTCATEGORY_CBZ", "FTCATEGORY_DPB", "FTCATEGORY_GLB", "FTCATEGORY_PM3",
	"FTCATEGORY_CD3", "FTCATEGORY_FLN", "FTCATEGORY_IVR", "FTCATEGORY_VU3", "FTCATEGORY_PFB", "FTCATEGORY_WIM",
	"FTCATEGORY_APPLE_DOCUMENTS", "FTCATEGORY_TABLEAU_FILES", "FTCATEGORY_AUTOCAD",
	"FTCATEGORY_INTEGRATED_CIRCUIT_FILES", "FTCATEGORY_LOG_FILES", "FTCATEGORY_EML_FILES", "FTCATEGORY_DAT",
	"FTCATEGORY_INI", "FTCATEGORY_THREED", "FTCATEGORY_THREEDA", "FTCATEGORY_THREEDFA", "FTCATEGORY_THREEDL",
	"FTCATEGORY_THREEDZ", "FTCATEGORY_APR", "FTCATEGORY_REALFLOW", "FTCATEGORY_COMP", "FTCATEGORY_DDF",
	"FTCATEGORY_DEM", "FTCATEGORY_THREEDS_MAX", "FTCATEGORY_GSP", "FTCATEGORY_HCL", "FTCATEGORY_MOTION_ANALYSIS",
	"FTCATEGORY_IGS", "FTCATEGORY_K3D", "FTCATEGORY_LIGHTSCAPE", "FTCATEGORY_AUTODESK_MAYA", "FTCATEGORY_MXS",
	"FTCATEGORY_OBJ", "FTCATEGORY_SHP", "FTCATEGORY_SPB", "FTCATEGORY_WRL", "FTCATEGORY_TMP", "FTCATEGORY_MUI",
	"FTCATEGORY_HBS", "FTCATEGORY_ICS", "FTCATEGORY_PUB", "FTCATEGORY_DRAWIO", "FTCATEGORY_PRT", "FTCATEGORY_PS2",
	"FTCATEGORY_PS3", "FTCATEGORY_ACIS", "FTCATEGORY_VDA", "FTCATEGORY_PARASOLID", "FTCATEGORY_PGP",
	"FTCATEGORY_BIN", "FTCATEGORY_JSON", "FTCATEGORY_XML", "FTCATEGORY_BINHEX", "FTCATEGORY_QUARKXPRESS",
	"FTCATEGORY_GO_FILES", "FTCATEGORY_SWIFT_FILES", "FTCATEGORY_RUBY_FILES", "FTCATEGORY_PERL_FILES",
	"FTCATEGORY_MATLAB_FILES", "FTCATEGORY_INCLUDE_FILES", "FTCATEGORY_JAVA_FILES", "FTCATEGORY_MAKE_FILES",
	"FTCATEGORY_YAML_FILES", "FTCATEGORY_VISUAL_BASIC_FILES", "FTCATEGORY_C_FILES", "FTCATEGORY_XAML",
	"FTCATEGORY_BASIC_SOURCE_CODE", "FTCATEGORY_SCT", "FTCATEGORY_A_FILE", "FTCATEGORY_MS_CPP_FILES",
	"FTCATEGORY_ASM", "FTCATEGORY_BORLAND_CPP_FILES", "FTCATEGORY_CLW", "FTCATEGORY_COBOL", "FTCATEGORY_CSX",
	"FTCATEGORY_DELPHI", "FTCATEGORY_DMD", "FTCATEGORY_DSP", "FTCATEGORY_F_FILES", "FTCATEGORY_NATVIS",
	"FTCATEGORY_NCB", "FTCATEGORY_NFM", "FTCATEGORY_POD", "FTCATEGORY_QLIKVIEW_FILES", "FTCATEGORY_RES_FILES",
	"FTCATEGORY_RPY", "FTCATEGORY_RSP", "FTCATEGORY_SAS", "FTCATEGORY_SC", "FTCATEGORY_SCALA", "FTCATEGORY_SWC",
	"FTCATEGORY_TCC", "FTCATEGORY_TLH", "FTCATEGORY_TLI", "FTCATEGORY_VISUAL_CPP_FILES", "FTCATEGORY_X1B",
	"FTCATEGORY_IFC", "FTCATEGORY_BCP", "FTCATEGORY_FOR", "FTCATEGORY_NCI", "FTCATEGORY_AU3", "FTCATEGORY_BGI",
	"FTCATEGORY_MANIFEST", "FTCATEGORY_NLS", "FTCATEGORY_TLB", "FTCATEGORY_ASHX", "FTCATEGORY_EXP",
	"FTCATEGORY_FLASH_VIDEO", "FTCATEGORY_AVI", "FTCATEGORY_MPEG", "FTCATEGORY_MP4", "FTCATEGORY_3GPP",
	"FTCATEGORY_QUICKTIME_VIDEO", "FTCATEGORY_WINDOWS_MEDIA_MOVIE", "FTCATEGORY_MKV", "FTCATEGORY_WEBM",
	"FTCATEGORY_VS4", "FTCATEGORY_TS",
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

func validateSSLInspectionPlatforms() schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		value, ok := i.(string)
		if !ok {
			return diag.Diagnostics{
				{
					Severity: diag.Error,
					Summary:  "Expected type to be string",
					Detail:   "Zscaler Client Connector device platforms for which the rule must be applied.",
				},
			}
		}

		// Convert the cty.Path to a string representation
		pathStr := fmt.Sprintf("%+v", path)

		// Use StringInSlice from helper/validation package
		var diags diag.Diagnostics
		if _, errs := validation.StringInSlice(supportedSSLInspectionPlatforms, false)(value, pathStr); len(errs) > 0 {
			for _, err := range errs {
				diags = append(diags, diag.FromErr(err)...)
			}
		}

		return diags
	}
}

var supportedSSLInspectionPlatforms = []string{
	"SCAN_IOS", "SCAN_ANDROID", "SCAN_MACOS", "SCAN_WINDOWS", "NO_CLIENT_CONNECTOR", "SCAN_LINUX",
}

/*
func stringIsJSON(i interface{}, k cty.Path) diag.Diagnostics {
	v, ok := i.(string)
	if !ok {
		return diag.Errorf("expected type of %s to be string", k)
	}
	if v == "" {
		return diag.Errorf("expected %q JSON to not be empty, got %v", k, i)
	}
	if _, err := structure.NormalizeJsonString(v); err != nil {
		return diag.Errorf("%q contains an invalid JSON: %s", k, err)
	}
	return nil
}
*/

func stringIsMultiLine(i interface{}, k cty.Path) diag.Diagnostics {
	v, ok := i.(string)
	if !ok {
		return diag.Errorf("expected type of %s to be string", k)
	}
	if v == "" {
		return diag.Errorf("expected %q text to not be empty, got %v", k, i)
	}
	return nil
}

// Ensures consistent formatting for multi-line text (aligns properly)
func normalizeMultiLineString(val interface{}) string {
	str, ok := val.(string)
	if !ok || str == "" {
		return ""
	}

	// Trim leading/trailing whitespace for consistency
	str = strings.TrimSpace(str)

	// Ensure uniform indentation by trimming each line
	lines := strings.Split(str, "\n")
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}

	// Escape Terraform variable interpolation (`$` → `$$`)
	escapedStr := strings.Join(lines, "\n")
	escapedStr = strings.ReplaceAll(escapedStr, "$", "$$")

	// Ensure the final newline to match Terraform formatting
	return escapedStr + "\n"
}

// Suppresses differences in multi-line text by ignoring whitespace discrepancies
func noChangeInMultiLineText(k, oldText, newText string, d *schema.ResourceData) bool {
	if newText == "" {
		return true
	}

	// Normalize both values and compare
	oldTextNormalized := normalizeMultiLineString(oldText)
	newTextNormalized := normalizeMultiLineString(newText)

	return oldTextNormalized == newTextNormalized
}

func unescapeTerraformVariables(val string) string {
	// Convert `$$` back to `$` for the API
	return strings.ReplaceAll(val, "$$", "$")
}
