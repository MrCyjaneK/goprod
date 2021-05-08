install:
	cp build/bin/${BINNAME}_${GOOS}_${GOARCH} /usr/bin/${BINNAME}
	cp usr/ /usr/share/goprod -r