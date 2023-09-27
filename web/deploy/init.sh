container=go-web-demo
image=registry.demo.com/demo/go-web-demo:main-af7c484

docker run --rm --network=host \
--name=${container} \
-e MYSQL_ADDRESS="127.0.0.1:3306" \
-e MYSQL_DATABASE="demo" \
-e MYSQL_USERNAME="root" \
-e MYSQL_PASSWORD="root" \
${image} sh -c '/migrations.sh migrate'
