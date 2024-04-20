module "some_mod" {
  source  = "s3::somebucket/a/sdf"
  version = "3.0.1"

  formatted = format("%s-asdf", local.name)
}

// a comment
