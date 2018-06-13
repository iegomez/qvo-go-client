requirements:
	@echo "Installing development tools"
	@go get -u github.com/pkg/errors
	@go get -u github.com/smartystreets/goconvey/convey
	@go get -u github.com/sirupsen/logrus

test:
	go test . -v -bench=none

benchmark:
	go test . -v -bench=. -run=^a