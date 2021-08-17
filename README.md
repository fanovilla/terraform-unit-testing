# Terraform Unit Testing with Go and JMESPath



## FAQs

Why Go?
- Just nice symmetry with Terraform written in go

Why JMESPath?
- See https://gatling.io/2019/07/introducing-jmespath-support/
- AWS CLI familiarity - see https://docs.aws.amazon.com/cli/latest/userguide/cli-usage-filter.html

Why not jq?
- See Why JMESPath above
- However, you can use jq to process`PlanFixture.Json` which is exposed as a string


## Directory Layouts

Simple case, main.tf in project root.
```text
project_root
  main.tf
  tests
    suite
      sample_test.go
```

Multiple root modules, main.tf one folder down from project root
```text
project_root
  s3_module
    main.tf
    tests
      suite
        sample_test.go
  ec2_module
    main.tf
    tests
      suite
        sample_test.go
```

Multiple root modules using common submodule, main.tf one folder down from project root
```text
project_root
  modules
    common_module
      main.tf
  s3_module
    main.tf
    tests
      suite
        sample_test.go
  ec2_module
    main.tf
    tests
      suite
        sample_test.go
```
