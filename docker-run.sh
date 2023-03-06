# mysql
docker run -d --name mysql --cap-add sys_nice -e MYSQL_ROOT_PASSWORD=123 -p 3306:3306 -v mysql:/var/lib/mysql mysql:latest --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci --default-authentication-plugin=mysql_native_password --skip-name-resolve

# rabbitmq
docker run -d --name rabbitmq --restart=always -p 15672:15672 -p 5672:5672 -v rabbitmq:/var/lib/rabbitmq -e RABBITMQ_DEFAULT_USER=root -e RABBITMQ_DEFAULT_PASS=password rabbitmq:3.11.5-management

# redis
docker run -d --name redis -p 6379:6379 redis --requirepass 123 

# elasticsearch & kibana
docker network create localnetwork
docker run -d --name elasticsearch --net localnetwork -e "discovery.type=single-node" -e ES_JAVA_OPTS="-Xms64m -Xmx512m" -e "xpack.security.enabled=false" -e "xpack.security.transport.ssl.enabled=false" -e "http.cors.enabled=true" -e "http.cors.allow-origin='*'" -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" elasticsearch:8.5.2

docker run -d --name kibana --net localnetwork -p 5601:5601 kibana:8.5.2

# zookeeper
# docker run --name zookeeper --restart always -d -p 2181:2181 -p 2888:2888 -p 3888:3888 -p  8080:8080  zookeeper

# mongodb
docker run -d --name mongo -e MONGO_INITDB_ROOT_USERNAME=root -e MONGO_INITDB_ROOT_PASSWORD=123 -v mongo-volume:/data/db -p 27017:27017 mongo:6.0.4