module github.com/apache/apisix-dashboard

go 1.15

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.5

require (
	github.com/coreos/bbolt v1.3.2 // indirect
	github.com/coreos/etcd v3.3.25+incompatible // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/evanphx/json-patch/v5 v5.1.0
	github.com/getkin/kin-openapi v0.33.0
	github.com/gin-contrib/gzip v0.0.3
	github.com/gin-contrib/static v0.0.0-20200916080430-d45d9a37d28e
	github.com/gin-gonic/gin v1.6.3
	github.com/go-resty/resty/v2 v2.7.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/google/uuid v1.2.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2 // indirect
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/juliangruber/go-intersect v1.1.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.8.0 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/shiningrush/droplet v0.2.6-0.20210127040147-53817015cd1b
	github.com/shiningrush/droplet/wrapper/gin v0.2.1
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/sony/sonyflake v1.0.0
	github.com/spf13/cobra v0.0.3
	github.com/spf13/viper v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/tidwall/gjson v1.6.7
	github.com/tmc/grpc-websocket-proxy v0.0.0-20200427203606-3cfed13b9966 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0
	github.com/yuin/gopher-lua v0.0.0-20200816102855-ee81675732da
	go.etcd.io/etcd v3.3.25+incompatible
	go.uber.org/zap v1.17.0
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a // indirect
	golang.org/x/sys v0.0.0-20210514084401-e8d321eab015 // indirect
	golang.org/x/time v0.0.0-20201208040808-7e3f01d25324 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)
