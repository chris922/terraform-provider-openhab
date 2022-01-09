resource "openhab_link" "example_link" {
  item_name   = "test_item"
  channel_uid = "modbus:data:smartenergymeter:L3:Current:number"

  configuration = {
    profile         = "modbus:gainOffset"
    gain            = "0.001 A"
    pre-gain-offset = "0"
  }
}
