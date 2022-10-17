module github.com/apache/apisix-dashboard/api/test/e2e

go 1.15

require (
	github.com/apache/apisix-dashboard v0.0.0-00010101000000-000000000000
	github.com/gavv/httpexpect/v2 v2.3.1
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.16.0
	github.com/savsgio/gotils v0.0.0-20210617111740-97865ed5a873
	github.com/stretchr/testify v1.7.1
	github.com/tidwall/gjson v1.11.0
)

replace github.com/apache/apisix-dashboard => ../../../
