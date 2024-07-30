module code.cloudfoundry.org/go-log-cache/v3

go 1.22

require (
	code.cloudfoundry.org/go-envstruct v1.7.0
	code.cloudfoundry.org/go-loggregator/v10 v10.0.0
	github.com/blang/semver/v4 v4.0.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.21.0
	github.com/onsi/ginkgo/v2 v2.19.1
	github.com/onsi/gomega v1.34.1
	google.golang.org/genproto/googleapis/api v0.0.0-20240723171418-e6d459c13d2a
	google.golang.org/grpc v1.65.0
	google.golang.org/protobuf v1.34.2
)

require (
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/pprof v0.0.0-20240424215950-a892ee059fd6 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	golang.org/x/exp v0.0.0-20240719175910-8a7402abbf56 // indirect
	golang.org/x/net v0.27.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	golang.org/x/tools v0.23.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240723171418-e6d459c13d2a // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract v3.0.0 // tagged an unreachable commit
