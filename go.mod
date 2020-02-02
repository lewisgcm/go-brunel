module go-brunel

require (
	github.com/Microsoft/go-winio v0.4.11 // indirect
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/Sirupsen/logrus v1.3.0
	github.com/boltdb/bolt v1.3.1 // indirect
	github.com/buildkite/terminal-to-html v3.2.0+incompatible
	github.com/casbin/casbin v1.8.1
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/docker/containerd v0.2.8 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v0.0.0-20170502054910-90d35abf7b35
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.3.3
	github.com/docker/libkv v0.2.1 // indirect
	github.com/docker/libnetwork v0.5.6 // indirect
	github.com/docker/libtrust v0.0.0-20160708172513-aabc10ec26b7 // indirect
	github.com/docker/swarmkit v1.12.0 // indirect
	github.com/emicklei/go-restful v2.9.0+incompatible // indirect
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-chi/chi v4.0.2+incompatible
	github.com/go-openapi/jsonpointer v0.18.0 // indirect
	github.com/go-openapi/jsonreference v0.18.0 // indirect
	github.com/go-openapi/spec v0.18.0 // indirect
	github.com/go-openapi/swag v0.18.0 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/godbus/dbus v4.1.0+incompatible // indirect
	github.com/gogo/protobuf v1.2.1 // indirect
	github.com/golang/mock v1.3.0
	github.com/golang/snappy v0.0.1 // indirect
	github.com/google/btree v1.0.0 // indirect
	github.com/google/go-cmp v0.3.0 // indirect
	github.com/google/go-jsonnet v0.12.1
	github.com/google/gofuzz v0.0.0-20170612174753-24818f796faf // indirect
	github.com/google/uuid v1.1.0
	github.com/googleapis/gnostic v0.2.0 // indirect
	github.com/gregjones/httpcache v0.0.0-20190203031600-7a902570cb17 // indirect
	github.com/hashicorp/consul/api v1.3.0 // indirect
	github.com/howeyc/gopass v0.0.0-20170109162249-bf9dde6d0d2c // indirect
	github.com/imdario/mergo v0.3.6 // indirect
	github.com/json-iterator/go v1.1.5 // indirect
	github.com/juju/ratelimit v1.0.1 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/markbates/goth v1.50.0
	github.com/miekg/dns v1.1.27 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742 // indirect
	github.com/mongodb/mongo-go-driver v0.3.0
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/opencontainers/runc v0.1.1 // indirect
	github.com/opencontainers/runtime-spec v1.0.1 // indirect
	github.com/pelletier/go-toml v1.4.0 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/pkg/errors v0.8.1
	github.com/samuel/go-zookeeper v0.0.0-20190923202752-2cc03de413da // indirect
	github.com/sirupsen/logrus v1.4.1 // indirect
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.3.2
	github.com/tidwall/pretty v0.0.0-20190325153808-1166b9ac2b65 // indirect
	github.com/vbatts/tar-split v0.11.1 // indirect
	github.com/xanzy/ssh-agent v0.2.1 // indirect
	github.com/xdg/scram v0.0.0-20180814205039-7eeb5667e42c // indirect
	github.com/xdg/stringprep v0.0.0-20180714160509-73f8eece6fdc // indirect
	golang.org/x/text v0.3.2 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	google.golang.org/grpc v1.27.0 // indirect
	gopkg.in/go-playground/webhooks.v5 v5.9.0
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/src-d/go-billy.v4 v4.3.0 // indirect
	gopkg.in/src-d/go-git.v4 v4.10.0
	gopkg.in/yaml.v2 v2.2.2
	k8s.io/api v0.0.0-20171214033149-af4bc157c3a2
	k8s.io/apimachinery v0.0.0-20171207040834-180eddb345a5
	k8s.io/client-go v6.0.0+incompatible
	k8s.io/kube-openapi v0.0.0-20190208205540-d7c86cdc46e3 // indirect
)

replace github.com/Sirupsen/logrus v1.3.0 => github.com/sirupsen/logrus v1.3.0

replace sourcegraph.com/sourcegraph/go-diff v0.5.1 => github.com/sourcegraph/go-diff v0.5.1

go 1.13
