module "namespaced_mod" {
  source  = "test::mymodule"
  version = "9.3.1"

  formatted = format("%s-asdf", local.name)
}

// a comment
