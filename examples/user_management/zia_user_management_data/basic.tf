terraform {
  required_providers {
    zia = {
      version = "2.0.1"
      source  = "zscaler.com/zia/zia"
    }
  }
}

provider "zia" {}

resource "zia_user_management" "john_ashcroft" {
  name     = "John Ashcroft"
  email    = "john.ashcroft@securitygeek.io"
  password = "P@ssw0rd123*"

  groups {
    id = data.zia_group_management.normal_internet.id
    #idp_id = ""
  }
  groups {
    id = data.zia_group_management.devops.id
    #idp_id = ""
  }

  department {
    id = data.zia_department_management.engineering.id
  }

}

data "zia_group_management" "normal_internet" {
  name = "Normal_Internet"
}

data "zia_group_management" "devops" {
  name = "DevOps"
}

// data "zia_group_management" "engineering" {
//  name = "Engineering"
// }

data "zia_department_management" "engineering" {
  name = "Engineering"
}

// data.zia_group_management.engineering.id
