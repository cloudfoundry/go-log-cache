module code.cloudfoundry.org/go-log-cache/v3

go 1.22.0

toolchain go1.22.8

require (
	code.cloudfoundry.org/go-envstruct v1.7.0
	code.cloudfoundry.org/go-loggregator/v10 v10.0.1
	github.com/blang/semver/v4 v4.0.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.25.1
	github.com/onsi/ginkgo/v2 v2.22.2
	github.com/onsi/gomega v1.36.2
	google.golang.org/genproto/googleapis/api v0.0.0-20241219192143-6b3ec007d9bb
	google.golang.org/grpc v1.70.0
	google.golang.org/protobuf v1.36.3
)

require (
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/pprof v0.0.0-20241210010833-40e02aabc2ad // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	golang.org/x/tools v0.28.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241219192143-6b3ec007d9bb // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract v3.0.0 // tagged an unreachable commit
