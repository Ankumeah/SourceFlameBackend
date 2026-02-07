# Run go tests on ./backend/
test *options:
  @ cd ./backend/ && go test ./... {{options}}
