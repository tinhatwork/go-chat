APP_NAME=pbill-chat
version=$(shell bash ./version.sh build)
SET_VERSION=$(eval APP_VERSION=$(version))
ENV_PROD="PROD"

run:
	env bee run

build:
	GOOS=linux CGO_ENABLED=0 go build -o ${APP_NAME} ./main.go

docs:
	#bee run -downdoc=true -gendoc=true
	bee generate docs

sbsync:
	rsync -aurv player.html supporter.html go.mod go.sum  ${APP_NAME} pbilling_sb:/home/deploy/pbilling_chat

sandbox: build sbsync
	ssh pbilling_sb docker restart pbilling.chat.sb

docker: build
	$(SET_VERSION)
	# repos at docker.io
	docker build . -t kenkinsai/${APP_NAME}:v${APP_VERSION}
	#docker push kenkinsai/${APP_NAME}:v${APP_VERSION}
	#docker image rm kenkinsai/${APP_NAME}:v${APP_VERSION}
