FROM alpine
ADD filesvc /

ENV HTTP_ADDRESS=:10107 HTTPS_ADDRESS=:8107 SD_TYPE=consul SD_ADDRESS=127.0.0.1:8500 PREFIX=download FILEPATH=/var/www/attach \
	CONFIG_FILE=filesvc.conf

EXPOSE 8107 10107 
VOLUME  /data /etc/filesvc.conf
ENTRYPOINT  /filesvc --http=$HTTP_ADDRESS --https=$HTTPS_ADDRESS --sdt=$SD_TYPE --sda=$SD_ADDRESS --rpcx=$RPCX_ADDRESS  \
		--file=$CONFIG_FILE --prefix=$PREFIX --filepath=$FILEPATH \
		--container-http=:10107 \
		--container-https=:8107
