all:
	GOPATH=`pwd` GOBIN=`pwd`/bin ${GO} install webapp


release:
	make clean && make all && mkdir target
	cp -r bin target


clean:
	rm -rf bin target
