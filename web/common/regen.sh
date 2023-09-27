
PROJECT=$1
ENV_GOPATH=${GOPATH}/src
PB_PACKAGE=${PROJECT}/pb
GO_PACKAGE=${PROJECT}/api
SWAGGER_PACKAGE=${PROJECT}/swagger

rm -rf $GO_PACKAGE
rm -rf $SWAGGER_PACKAGE
cd ${ENV_GOPATH}
for j in $(ls $PB_PACKAGE); do
	for i in $(ls $ENV_GOPATH/$PB_PACKAGE/$j/*.proto); do
		echo $i
		/usr/local/protoc/bin/protoc \
		-I$GOPATH/src \
		-I$GOPATH/pkg/mod/github.com/googleapis/api-common-protos \
		-I$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway \
		-I. \
		--go_out=. --go-grpc_out=. \
		--grpc-gateway_out=logtostderr=true:. \
		--openapiv2_out=logtostderr=true,json_names_for_fields=false:. \
		"$i";
	done
	mkdir -p $SWAGGER_PACKAGE/$j
	mv $PB_PACKAGE/$j/*.swagger.json $SWAGGER_PACKAGE/$j
done
