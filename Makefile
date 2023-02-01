PROJECT:=iris-project-app

.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags="-w -s" -a -installsuffix -o iris-project-app .

# make build-linux
build-linux:
	@docker build -t iris-project-app:latest .
	@echo "build successful"

build-sqlite:
	go build -tags sqlite3 -ldflags="-w -s" -a -installsuffix -o iris-project-app .

# make run
run:
    # delete iris-project-app-api container
	@if [ $(shell docker ps -aq --filter name=iris-project-app --filter publish=8080) ]; then docker rm -f iris-project-app; fi

    # 启动方法一 run iris-project-app-api container  docker-compose 启动方式
    # 进入到项目根目录 执行 make run 命令
	@docker-compose up -d

	# 启动方式二 docker run  这里注意-v挂载的宿主机的地址改为部署时的实际决对路径
    #@docker run --name=iris-project-app -p 8080:8080 -v /home/code/go/src/iris-project-app/iris-project-app/config:/iris-project-app-api/config  -v /home/code/go/src/iris-project-app/iris-project-app-api/static:/iris-project-app/static -v /home/code/go/src/iris-project-app/iris-project-app/temp:/iris-project-app-api/temp -d --restart=always iris-project-app:latest

	@echo "iris-project-app service is running..."

	# delete Tag=<none> 的镜像
	@docker image prune -f
	@docker ps -a | grep "iris-project-app"

stop:
    # delete iris-project-app-api container
	@if [ $(shell docker ps -aq --filter name=iris-project-app --filter publish=8080) ]; then docker-compose down; fi
	#@if [ $(shell docker ps -aq --filter name=iris-project-app --filter publish=8080) ]; then docker rm -f iris-project-app; fi
	#@echo "iris-project-app stop success"


#.PHONY: test
#test:
#	go test -v ./... -cover

#.PHONY: docker
#docker:
#	docker build . -t iris-project-app:latest

# make deploy
deploy:

	#@git checkout master
	#@git pull origin master
	make build-linux
	make run