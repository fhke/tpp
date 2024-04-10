module "some_mod" {
  source = "somebucket/replaced"

  formatted = format("%s-asdf", local.name)
}

// a comment
