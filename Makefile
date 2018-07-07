.PHONY: default build builder-image binary-image clean-images clean
BUILDER = k8swatchcrd-builder
BINARY = k8swatchcrd

VERSION=
BUILD=

GOCMD = go
GOFLAGS ?= $(GOFLAGS:)
LDFLAGS =
DOCKER_IMAGE ?= herekar/k8swatchcrd
# Default value "1.0"
DOCKER_TAG ?= 1.0
REPOSITORY = ${DOCKER_IMAGE}:${DOCKER_TAG}
default: build test

build:
	$(shell cd opt; go build -o k8swatchcrd) 

builder-image:
	@docker build --network host -t "${BUILDER}" -f opt/build/Dockerfile.build .

binary-image: builder-image
	@docker run --network host --rm "${BUILDER}" | docker build --network host -t "${REPOSITORY}" -f Dockerfile.run -

clean-images: stop
	@docker rmi "${BUILDER}" "${BINARY}"

test:
	$(shell cd opt/controller; go test -v)

clean:
	"$(GOCMD)" clean -i
