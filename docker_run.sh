@echo off
docker run -d \
	--restart=always \
	-e TZ=Asia/Shanghai \
	--name chat_server \
	-p 9106:9106 \
	-v `pwd`/bin/:/work \
	-w /work \
	ubuntu \
	./chat_server