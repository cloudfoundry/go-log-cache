module code.cloudfoundry.org/go-log-cache/v3

go 1.23.0

toolchain go1.23.5

require (
	code.cloudfoundry.org/go-envstruct v1.7.0
	code.cloudfoundry.org/go-loggregator/v10 v10.2.0
	github.com/blang/semver/v4 v4.0.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.27.0
	github.com/onsi/ginkgo/v2 v2.23.4
	github.com/onsi/gomega v1.37.0
	google.golang.org/genproto/googleapis/api v0.0.0-20250603155806-513f23925822
	google.golang.org/grpc v1.73.0
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/pprof v0.0.0-20250403155104-27863c87afa6 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	go.uber.org/automaxprocs v1.6.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	golang.org/x/tools v0.33.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250603155806-513f23925822 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract v3.0.0 // tagged an unreachable commit
