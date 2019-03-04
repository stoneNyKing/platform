FROM alpine
ADD hrsvc /

ENV HTTP_ADDRESS=:10114 HTTPS_ADDRESS=:8114 SD_TYPE=consul SD_ADDRESS=127.0.0.1:8500 RPCX_ADDRESS=:20114 \
	CONFIG_FILE=hrsvc.conf

EXPOSE 8114 10114 20114
VOLUME  /data /etc/hrsvc.conf
ENTRYPOINT  /hrsvc --http=$HTTP_ADDRESS --https=$HTTPS_ADDRESS --sdt=$SD_TYPE --sda=$SD_ADDRESS --rpcx=$RPCX_ADDRESS  \
		--file=$CONFIG_FILE \
		--container-http=:10114 \
		--container-https=:8114 \
		--container-rpcx=:20114
