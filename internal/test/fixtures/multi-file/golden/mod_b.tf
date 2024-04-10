locals {
  name = "test"
}

module "some_other_mod" {
  source  = "somebucket/replaced/2"
  version = "1.0.0"
}