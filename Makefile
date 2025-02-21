INPUT_FOLDER ?= protos
OUPUT_FOLDER ?= pb
PROTOVERSION ?= 4.22.0
OWNER ?= ociscloud
PREVERSION ?= 0.2.3
VERSION ?= $(shell cat VERSION)
PWD := $(shell pwd)

sed = sed
ifeq ("$(shell uname -s)", "Darwin")	# BSD sed, like MacOS
	sed += -i ''
else	# GNU sed, like LinuxOS
	sed += -i''
endif

.PHONY: update-submodules
update-submodules:
	@git submodule update --init --recursive
	@git submodule update --recursive --remote
	@git submodule foreach --recursive 'tag="$$(git config -f $$toplevel/.gitmodules submodule.$$name.tag)" && git reset --hard $$tag || echo "this module has no tag $$toplevel $$name $$tag"'

.PHONY: set-version
set-version:
	@echo "Set Version"
	@$(sed) -e'/$(PREVERSION)/{s//$(VERSION)/;:b' -e'n;bb' -e\} $(PWD)/constants/version.go
	@$(sed) -e'/$(PREVERSION)/{s//$(VERSION)/;:b' -e'n;bb' -e\} $(PWD)/Makefile

.PHONY: build-protofile
build-protofile:
	@export PATH=/root/go/bin:${PATH}  && protoc -I=$(INPUT_FOLDER) $(INPUT_FOLDER)/*.proto --go_out=$(OUPUT_FOLDER) --go_opt=paths=source_relative --go-grpc_out=$(OUPUT_FOLDER) --go-grpc_opt=paths=source_relative

.PHONY: container-build-protofile
container-build-protofile: 
	@docker run  --name build-env -w /home -v $(PWD):/home $(OWNER)/protobuf:$(PROTOVERSION)-$(OS) make build-protofile
	@docker rm build-env

FORMAT=""
.PHONY: verify-modules
verify-modules:
	@bash verify_modules.sh $(FORMAT)
