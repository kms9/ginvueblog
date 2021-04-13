module ginvueblog

go 1.16

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-contrib/multitemplate v0.0.0-20200916052041-666a7309d230
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/go-playground/validator/v10 v10.4.1
	github.com/gomodule/redigo v1.8.2
	github.com/jonboulle/clockwork v0.2.0 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/kms9/publicyc v0.0.0-00010101000000-000000000000
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v1.0.4 // indirect
	github.com/nacos-group/nacos-sdk-go v1.0.7
	github.com/qiniu/go-sdk/v7 v7.9.1
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5
	github.com/shima-park/agollo v1.2.10
	github.com/silenceper/log v0.0.0-20171204144354-e5ac7fa8a76a
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/viper v1.7.1
	github.com/ugorji/go v1.2.4 // indirect
	github.com/xordataexchange/crypt v0.0.3-0.20170626215501-b2862e3d0a77
	go.etcd.io/etcd v3.3.25+incompatible
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
	golang.org/x/sys v0.0.0-20210124154548-22da62e12c0c // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/ini.v1 v1.62.0
	gorm.io/gorm v1.20.12
)

replace github.com/kms9/publicyc => ../yc

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
