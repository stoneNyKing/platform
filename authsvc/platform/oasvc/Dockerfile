FROM alpine
ADD oasvc /

ENV HTTP_ADDRESS=:9060 HTTPS_ADDRESS=:8060 SD_TYPE=consul SD_ADDRESS=127.0.0.1:8500 RPCX_ADDRESS=:19060 \
	CONFIG_FILE=oasvc.conf

EXPOSE 8060 9060 19060
VOLUME  /data /etc/oasvc.conf
ENTRYPOINT  /oasvc --http=$HTTP_ADDRESS --https=$HTTPS_ADDRESS --sdt=$SD_TYPE --sda=$SD_ADDRESS --rpcx=$RPCX_ADDRESS \
		--file=$CONFIG_FILE \
		--container-http=:9060 \
		--container-https=:8060 \
		--container-rpcx=:19060
