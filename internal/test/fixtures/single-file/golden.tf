locals {
  name = "test"
}

resource "aws_vpc" "this" {
  for_each = toset(["a", "b", "c"])

  cidr_block = "10.0.0.0/16"
  tags = {
    Name = local.name
  }
}

moved {
  from = aws_vpc.this[0]
  to   = aws_vpc.this
}

module "some_mod" {
  source = "somebucket/replaced"

  formatted = format("%s-asdf", local.name)
}

module "some_other_mod" {
  source  = "somebucket/replaced/2"
  version = "1.0.0"
}