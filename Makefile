.PHONY: concourse/login concourse/p1 concourse/p2 concourse/trigger-p2 clean
.PHONY: test test/integration test/unit retool

.DEFAULT_GOAL := install

FLYBIN=.bin/fly
FLY=$(FLYBIN) -t local

.bin:
	@mkdir .bin

out:
	@mkdir out

$(FLYBIN): .bin
	@curl -L -s -o .bin/fly 'http://127.0.0.1:8080/api/v1/cli?arch=amd64&platform=darwin'
	@chmod u+x .bin/fly
	@$(FLY) sync

concourse/login: $(FLYBIN)
	@$(FLY) login -c 'http://127.0.0.1:8080/' -u test -p test

concourse/p1: $(FLYBIN)
	@$(FLY) validate-pipeline -c ci-test/pipeline1.yaml
	@$(FLY) set-pipeline --config ci-test/pipeline1.yaml --pipeline p1 -n
	@$(FLY) expose-pipeline --pipeline p1
	@$(FLY) unpause-pipeline --pipeline p1

concourse/p2: $(FLYBIN)
	@$(FLY) validate-pipeline -c ci-test/pipeline2.yaml
	@$(FLY) set-pipeline --config ci-test/pipeline2.yaml --pipeline p2 -n
	@$(FLY) expose-pipeline --pipeline p2
	@$(FLY) unpause-pipeline --pipeline p2

concourse/trigger-p2: concourse/p2
	@$(FLY) trigger-job --job p2/build-node

clean:
	@rm -rf .bin
	@rm -rf out

retool:
	@go get github.com/twitchtv/retool

gen:
	@retool do go generate ./...

test: test/unit


test/unit: retool
	@retool do go test ./...

test/integration: retool
	@retool do go test -tags integration ./...

install: retool
	@dep ensure

tunnel:
	@ngrok http --bind-tls false 8080

infra/up:
	@docker-compose up -d

infra/down:
	@docker-compose down

out/lambda: out aws/cmd/handler/lambda.go
	@GOOS=linux go build -o out/lambda ./aws/cmd/handler/

out/lambda.zip: out/lambda
	cd out/; zip lambda.zip lambda

target: out/lambda.zip
