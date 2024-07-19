module github.com/EasyDarwin/EasyDarwin

go 1.12

require (
	github.com/StackExchange/wmi v0.0.0-20180412205111-cdffdb33acae // indirect
	github.com/common-nighthawk/go-figure v0.0.0-20180619031829-18b2b544842c
	github.com/denisenkom/go-mssqldb v0.12.3 // indirect
	github.com/eiannone/keyboard v0.0.0-20190314115158-7169d0afeb4f // indirect
	github.com/erikstmartin/go-testdb v0.0.0-20160219214506-8d10e4a1bae5 // indirect
	github.com/gin-contrib/pprof v0.0.0-20180514151456-0ed7ffb6a189
	github.com/gin-contrib/sse v0.0.0-20170109093832-22d885f9ecc7 // indirect
	github.com/gin-contrib/static v0.0.0-20180301030858-73da7037e716
	github.com/gin-gonic/gin v1.3.0
	github.com/go-ini/ini v1.49.0 // indirect
	github.com/go-ole/go-ole v0.0.0-20180213002836-a1ec82a652eb // indirect
	github.com/go-redis/redis v6.15.6+incompatible // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/gorilla/context v1.1.1 // indirect
	github.com/gorilla/sessions v1.2.0 // indirect
	github.com/jinzhu/gorm v0.0.0-20180512062900-82eb9f8a5bbb
	github.com/jinzhu/inflection v0.0.0-20180308033659-04140366298a // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jonboulle/clockwork v0.4.0 // indirect
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible // indirect
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/mattn/go-isatty v0.0.4 // indirect
	github.com/mattn/go-sqlite3 v1.11.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.33.1 // indirect
	github.com/penggy/EasyGoLib v0.0.0-00010101000000-000000000000
	github.com/penggy/cors v0.0.0-20180918145040-d08bb28f7e48
	github.com/penggy/service v0.0.0-00010101000000-000000000000
	github.com/penggy/sessions v0.0.0-00010101000000-000000000000
	github.com/pixelbender/go-sdp v0.0.0-20181123094152-100bc9371a0c
	github.com/shirou/gopsutil v0.0.0-20180417021151-63047d776e07
	github.com/smartystreets/goconvey v1.8.1 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	github.com/teris-io/shortid v0.0.0-20171029131806-771a37caa5cf
	github.com/ugorji/go v0.0.0-20180307152341-02537d3a3e32 // indirect
	golang.org/x/net v0.27.0 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v8 v8.18.2
	gopkg.in/ini.v1 v1.67.0 // indirect
)

replace (
	github.com/penggy/EasyGoLib => ./vendor/github.com/penggy/EasyGoLib
	github.com/penggy/EasyGoLib/util => ./vendor/github.com/penggy/EasyGoLib/util
	github.com/penggy/db => ./vendor/github.com/penggy/EasyGoLib/db
	github.com/penggy/service => ./vendor/github.com/penggy/service
	github.com/penggy/sessions => ./vendor/github.com/penggy/sessions
)
