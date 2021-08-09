.PHONY: docker docker_push
install:
	cp build/bin/${BINNAME}_${GOOS}_${GOARCH} /usr/bin/${BINNAME}
	cp usr/ /usr/share/goprod -r
docker:
	docker build . -f docker/Dockerfile.core -t mrcyjanek/goprod:core
	docker build . -f docker/Dockerfile.nodejs -t mrcyjanek/goprod:nodejs
	docker build . -f docker/Dockerfile.core.android -t mrcyjanek/goprod:core-android
	docker build . -f docker/Dockerfile.nodejs.android -t mrcyjanek/goprod:nodejs-android
	docker build . -f docker/Dockerfile.core.linux -t mrcyjanek/goprod:core-linux
	docker build . -f docker/Dockerfile.nodejs.linux -t mrcyjanek/goprod:nodejs-linux
	docker build . -f docker/Dockerfile.core.windows -t mrcyjanek/goprod:core-windows
	docker build . -f docker/Dockerfile.nodejs.windows -t mrcyjanek/goprod:nodejs-windows
docker_push:
	docker push mrcyjanek/goprod:core
	docker push mrcyjanek/goprod:nodejs
	docker push mrcyjanek/goprod:core-android
	docker push mrcyjanek/goprod:nodejs-android
	docker push mrcyjanek/goprod:core-linux
	docker push mrcyjanek/goprod:nodejs-linux
	docker push mrcyjanek/goprod:core-windows
	docker push mrcyjanek/goprod:nodejs-windows