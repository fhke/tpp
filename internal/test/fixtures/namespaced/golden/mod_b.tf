locals {
  name = "test"
}

module "some_other_mod" {
  source  = "s3::otherbucket/a/sdf"
  version = "3.0.9"
}