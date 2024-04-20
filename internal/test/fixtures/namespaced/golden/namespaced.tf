module "namespaced_mod" {
  source  = "s3::somebucket/replaced"
  version = "9.7.0"

  formatted = format("%s-asdf", local.name)
}

// a comment
