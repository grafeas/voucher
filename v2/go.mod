module github.com/grafeas/voucher/v2

go 1.16

require (
	// if updating containeranalysis or grafeas, ensure the options in containeranalysis/client.go are still valid
	cloud.google.com/go/containeranalysis v0.1.0
	cloud.google.com/go/grafeas v0.1.0
	cloud.google.com/go/kms v1.0.0
	cloud.google.com/go/pubsub v1.3.1
	github.com/CycloneDX/cyclonedx-go v0.5.1
	github.com/DataDog/datadog-api-client-go v1.3.0
	github.com/DataDog/datadog-go v3.4.0+incompatible
	github.com/Shopify/ejson v1.2.0
	github.com/antihax/optional v1.0.0
	github.com/bradleyfalzon/ghinstallation v1.1.1
	github.com/coreos/clair v2.0.5+incompatible
	github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/docker v20.10.14+incompatible
	github.com/docker/libtrust v0.0.0-20160708172513-aabc10ec26b7
	github.com/dustin/gojson v0.0.0-20160307161227-2e71ec9dd5ad // indirect
	github.com/fernet/fernet-go v0.0.0-20211208181803-9f70042a33ee // indirect
	github.com/golang/mock v1.6.0
	github.com/google/go-containerregistry v0.8.0
	github.com/googleapis/gax-go/v2 v2.1.1
	github.com/gorilla/mux v1.8.0
	github.com/mennanov/fieldmask-utils v0.0.0-20190703161732-eca3212cf9f3
	github.com/mitchellh/go-homedir v1.1.0
	github.com/opencontainers/go-digest v1.0.0
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/pborman/uuid v0.0.0-20180906182336-adf5a7427709 // indirect
	github.com/shurcooL/githubv4 v0.0.0-20190718010115-4ba037080260
	github.com/shurcooL/graphql v0.0.0-20181231061246-d48a9a75455f // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/smartystreets/goconvey v0.0.0-20190731233626-505e41936337 // indirect
	github.com/spf13/cobra v1.3.0
	github.com/spf13/viper v1.10.0
	github.com/stretchr/testify v1.7.1
	go.mozilla.org/sops/v3 v3.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.0.0-20211215153901-e495a2d5b3d3
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8
	google.golang.org/api v0.63.0
	google.golang.org/genproto v0.0.0-20211208223120-3a66f561d7aa
	google.golang.org/grpc v1.43.0
	gotest.tools/v3 v3.1.0 // indirect
)

replace go.mozilla.org/sops/v3 => github.com/Shopify/sops/v3 v3.7.2-0.20220208200849-289e67dba11f
