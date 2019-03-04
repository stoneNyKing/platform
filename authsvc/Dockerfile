FROM alpine
ADD authsvc /

ENV HTTP_ADDRESS=:10103 HTTPS_ADDRESS=:8103 SD_TYPE=consul SD_ADDRESS=127.0.0.1:8500 RPCX_ADDRESS=:20103 \
	CONFIG_FILE=authsvc.conf

EXPOSE 10103 20103 8103
VOLUME  /data /etc/authsvc.conf
ENTRYPOINT  /authsvc --http=$HTTP_ADDRESS --https=$HTTPS_ADDRESS --sdt=$SD_TYPE --sda=$SD_ADDRESS --rpcx=$RPCX_ADDRESS  \
		--file=$CONFIG_FILE \
		--container-http=:10103 \
		--container-https=:8103 \
		--container-rpcx=:20103
