data "zia_devices" "device_model"{
    device_model = "VMware Virtual Platform"
}

output "zia_devices_device_model"{
    value = data.zia_devices.device_model
}

data "zia_devices" "owner_name"{
    owner_name = "jonh.doe@acme.com"
}

output "zia_devices_owner_name"{
    value = data.zia_devices.owner_name
}

data "zia_devices" "os_type"{
    os_type = "WINDOWS_OS"
}

output "zia_devices_os_type"{
    value = data.zia_devices.os_type
}

data "zia_devices" "os_version"{
    os_version = "Microsoft Windows 10 Pro;64 bit"
}

output "zia_devices_os_version"{
    value = data.zia_devices.os_version
}