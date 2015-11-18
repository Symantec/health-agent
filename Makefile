all:
	@cd $(GOPATH)/src; go install github.com/Symantec/health-agent/cmd/*
