install:
	cp build/bin/${BINNAME}_${GOOS}_${GOARCH} /usr/bin/${BINNAME}
	cp usr/ /usr/share/goprod -r
docker:
	docker build . -f Dockerfile -t mrcyjanek/goprod:latest
	docker build . -f Dockerfile.nodejs -t mrcyjanek/goprod:latest-nodejs
docker_push:
	docker push mrcyjanek/goprod:latest
	docker push mrcyjanek/goprod:latest-nodejs