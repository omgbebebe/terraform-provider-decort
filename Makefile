TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=digitalenergy.online
NAMESPACE=decort
NAME=terraform-provider-decort
BINDIR = ./bin
ZIPDIR = ./zip
#BINARY=terraform-provider-${NAME}
BINARY=${NAME}.exe
WORKPATH= ./examples/terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAMESPACE}/${VERSION}/${OS_ARCH} 
MAINPATH = ./cmd/decort/
VERSION=3.6.0
#OS_ARCH=darwin_amd64
OS_ARCH=windows_amd64
#OS_ARCH=linux_amd64

FILES= 	${BINARY}_${VERSION}_darwin_amd64\
		${BINARY}_${VERSION}_freebsd_386\
		${BINARY}_${VERSION}_freebsd_amd64\
		${BINARY}_${VERSION}_freebsd_arm\
		${BINARY}_${VERSION}_linux_386\
		${BINARY}_${VERSION}_linux_amd64\
		${BINARY}_${VERSION}_linux_arm\
		${BINARY}_${VERSION}_openbsd_386\
		${BINARY}_${VERSION}_openbsd_amd64\
		${BINARY}_${VERSION}_solaris_amd64\
		${BINARY}_${VERSION}_windows_386  \
		${BINARY}_${VERSION}_windows_amd64\

BINS = $(addprefix bin/, $(FILES))

default: install

image:
	GOOS=linux GOARCH=amd64 go build -o terraform-provider-decort ./cmd/decort/
	docker build . -t rudecs/tf:3.2.2
	rm terraform-provider-decort

lint:
	golangci-lint run --timeout 600s

st:
	go build -o ${BINARY} ${MAINPATH}
	cp ${BINARY} ${WORKPATH}
	rm ${BINARY}

build:
	go build -o ${BINARY} ${MAINPATH}

release: $(FILES)

$(FILES) : $(BINDIR) $(ZIPDIR) $(BINS)
	zip -r $(ZIPDIR)/$@.zip $(BINDIR)/$@

$(BINDIR):
	mkdir $@

$(ZIPDIR):
	mkdir $@

$(BINS):
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_darwin_amd64 $(MAINPATH)
	GOOS=freebsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_freebsd_386 $(MAINPATH)
	GOOS=freebsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_freebsd_amd64 $(MAINPATH)
	GOOS=freebsd GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_freebsd_arm $(MAINPATH)
	GOOS=linux GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_linux_386 $(MAINPATH)
	GOOS=linux GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_linux_amd64 $(MAINPATH)
	GOOS=linux GOARCH=arm go build -o ./bin/${BINARY}_${VERSION}_linux_arm $(MAINPATH)
	GOOS=openbsd GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_openbsd_386 $(MAINPATH)
	GOOS=openbsd GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_openbsd_amd64 $(MAINPATH)
	GOOS=solaris GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_solaris_amd64 $(MAINPATH)
	GOOS=windows GOARCH=386 go build -o ./bin/${BINARY}_${VERSION}_windows_386 $(MAINPATH)
	GOOS=windows GOARCH=amd64 go build -o ./bin/${BINARY}_${VERSION}_windows_amd64 $(MAINPATH)

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test: 
	go test -i $(TEST) || exit 1                                                   
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4                    

testacc: 
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m
