.PHONY: login p1 p2 trigger-p2

FLYBIN=.bin/fly
FLY=$(FLYBIN) -t local

.bin:
	mkdir .bin

$(FLYBIN): .bin
	@curl -L -s -o .bin/fly 'http://127.0.0.1:8080/api/v1/cli?arch=amd64&platform=darwin'
	@chmod u+x .bin/fly
	@$(FLY) sync

login: $(FLYBIN)
	@$(FLY) login -c 'http://127.0.0.1:8080/' -u test -p test

p1: $(FLYBIN)
	@$(FLY) validate-pipeline -c ci-test/pipeline1.yaml
	@$(FLY) set-pipeline --config ci-test/pipeline1.yaml --pipeline p1 -n
	@$(FLY) expose-pipeline --pipeline p1
	@$(FLY) unpause-pipeline --pipeline p1

p2: $(FLYBIN)
	@$(FLY) validate-pipeline -c ci-test/pipeline2.yaml
	@$(FLY) set-pipeline --config ci-test/pipeline2.yaml --pipeline p2 -n
	@$(FLY) expose-pipeline --pipeline p2
	@$(FLY) unpause-pipeline --pipeline p2

trigger-p2: p2
	@$(FLY) trigger-job --job p2/build-node

clean:
	@rm -rf .bin
