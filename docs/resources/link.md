---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "openhab_link Resource - terraform-provider-openhab"
subcategory: ""
description: |-
  OpenHAB Link between an Item and a Channel.
---

# openhab_link (Resource)

OpenHAB Link between an Item and a Channel.

## Example Usage

```terraform
resource "openhab_link" "example_link" {
  item_name   = "test_item"
  channel_uid = "modbus:data:smartenergymeter:L3:Current:number"

  configuration = {
    profile         = "modbus:gainOffset"
    gain            = "0.001 A"
    pre-gain-offset = "0"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **channel_uid** (String) Channel UID
- **item_name** (String) Item name

### Optional

- **configuration** (Map of String) Link configuration

### Read-Only

- **id** (String) Resource ID


