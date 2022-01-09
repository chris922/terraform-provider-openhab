resource "openhab_item" "example" {
  name = "test_item"

  type = "Number"

  label = "Test Number"

  category = "energy"

  tags        = ["tag1", "tag2"]
  group_names = ["group_1"]
}
