FROM alpine
ADD ousvc /

ENV HTTP_ADDRESS=:9058 HTTPS_ADDRESS=:8058 SD_TYPE=consul SD_ADDRESS=127.0.0.1:8500 RPCX_ADDRESS=:19058 ADMIN_ADDRESS=:9060 \
	CONFIG_FILE=ousvc.conf

EXPOSE 8058 9058 19058
VOLUME  /data /etc/ousvc.conf
ENTRYPOINT  /ousvc --http=$HTTP_ADDRESS --https=$HTTPS_ADDRESS --sdt=$SD_TYPE --sda=$SD_ADDRESS --rpcx=$RPCX_ADDRESS --adminaddr=$ADMIN_ADDRESS \
		--file=$CONFIG_FILE \
		--container-http=:9058 \
		--container-https=:8058 \
		--container-rpcx=:19058
