code_root=github.com/learn-go/web/common

docker run -it --rm \
	-u `id -u` \
	-v $PWD:/go/src/${code_root} \
	-e PROTOC_INSTALL=/go \
	-w /go/src/$code_root \
	protoc:v1.0.0 bash ./regen.sh $code_root
