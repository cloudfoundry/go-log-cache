module code.cloudfoundry.org/go-log-cache/v3

go 1.22.0

toolchain go1.22.8

require (
	code.cloudfoundry.org/go-envstruct v1.7.0
	code.cloudfoundry.org/go-loggregator/v10 v10.0.1
	github.com/blang/semver/v4 v4.0.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.23.0
	github.com/onsi/ginkgo/v2 v2.22.0
	github.com/onsi/gomega v1.36.0
	google.golang.org/genproto/googleapis/api v0.0.0-20241021214115-324edc3d5d38
	google.golang.org/grpc v1.68.1
	google.golang.org/protobuf v1.35.2
)

require (
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/pprof v0.0.0-20241029153458-d1b30febd7db // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	golang.org/x/tools v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241021214115-324edc3d5d38 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract v3.0.0 // tagged an unreachable commit
