resource "zia_activation_status" "example1"{
    status = "ACTIVE"
}

output "zia_activation_status_example1"{
    value = zia_activation_status.example1
}

data "zia_activation_status" "example2"{
}

output "zia_activation_status_example2"{
    value = data.zia_activation_status.example2
}



