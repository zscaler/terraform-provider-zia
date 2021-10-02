package zia

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func SetToStringSlice(d *schema.Set) []string {
	list := d.List()
	return ListToStringSlice(list)
}

func ListToStringSlice(v []interface{}) []string {
	if len(v) == 0 {
		return []string{}
	}

	ans := make([]string, len(v))
	for i := range v {
		switch x := v[i].(type) {
		case nil:
			ans[i] = ""
		case string:
			ans[i] = x
		}
	}

	return ans
}

/*
// validateIpv4CIDRNetworkAddress ensures that the string value is a valid IPv4 CIDR that
// represents a network address - it adds an error otherwise
func validateIpv4CIDRNetworkAddress(v interface{}, k string) (ws []string, errors []error) {
	if err := validateIpv4CIDRBlock(v.(string)); err != nil {
		errors = append(errors, err)
		return
	}

	return
}


// validateIpv4CIDRBlock validates that the specified CIDR block is valid:
// - The CIDR block parses to an IP address and network
// - The IP address is an IPv4 address
// - The CIDR block is the CIDR block for the network
func validateIpv4CIDRBlock(cidr string) error {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return fmt.Errorf("%q is not a valid CIDR block: %w", cidr, err)
	}

	ipv4 := ip.To4()
	if ipv4 == nil {
		return fmt.Errorf("%q is not a valid IPv4 CIDR block", cidr)
	}

	if !tfnet.CIDRBlocksEqual(cidr, ipnet.String()) {
		return fmt.Errorf("%q is not a valid IPv4 CIDR block; did you mean %q?", cidr, ipnet)
	}

	return nil
}
*/

func getIntFromResourceData(d *schema.ResourceData, key string) (int, bool) {
	obj, isSet := d.GetOk(key)
	val, isInt := obj.(int)
	return val, isSet && isInt && val > 0
}
func getStringFromResourceData(d *schema.ResourceData, key string) (string, bool) {
	obj, isSet := d.GetOk(key)
	val, isStr := obj.(string)
	return val, isSet && isStr && val != ""
}
