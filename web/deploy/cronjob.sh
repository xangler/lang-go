# 添加crontab -e
# * * * * * /bin/sh /root/robot/cronjob.sh >> /root/robot/log/$(/bin/date '+\%y\%m\%d') 2>&1 
# 0 */6 * * * /bin/ls -t /root/robot/log | /bin/sed -n '90,$p' | /bin/xargs -I {} /bin/rm -f /root/robot/log/{}

container=go-web-demo
image=registry.demo.com/demo/go-web-demo:main-af7c484
run=$(docker ps -a|grep ${container}|wc -l)
echo "===============$(date)==================="
if [ $run == 1 ]; then
    echo "${container} running..."
else
    echo "run ${container} start..."
    docker run --rm -d --network=host \
	--name=${container} \
	-e MYSQL_ADDRESS="127.0.0.1:3306" \
	-e MYSQL_DATABASE="demo" \
	-e MYSQL_USERNAME="root" \
	-e MYSQL_PASSWORD="root" \
	-v ${PWD}/config:/work/config \
	${image} sh -c '/work/bin/go-web-server.out -c /work/config/demo.yaml'
    echo "run ${container} succ..."
fi