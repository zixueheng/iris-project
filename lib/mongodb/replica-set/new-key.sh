mkdir -p /etc/mongodb
openssl rand -base64 765 > /etc/mongodb/key.file
chmod 400 /etc/mongodb/key.file