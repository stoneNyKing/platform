FROM alpine
ADD confsvc /

ENV HTTP_ADDRESS=:10104 HTTPS_ADDRESS=:8104 SD_TYPE=consul SD_ADDRESS=127.0.0.1:8500 RPCX_ADDRESS=:20104 \
	CONFIG_FILE=confsvc.conf

EXPOSE 10104 20104
VOLUME  /data /etc/confsvc.conf
ENTRYPOINT  /confsvc --http=$HTTP_ADDRESS --https=$HTTPS_ADDRESS --sdt=$SD_TYPE --sda=$SD_ADDRESS --rpcx=$RPCX_ADDRESS  \
		--file=$CONFIG_FILE \
		--container-http=:10104 \
		--container-https=:8104 \
		--container-rpcx=:20104
