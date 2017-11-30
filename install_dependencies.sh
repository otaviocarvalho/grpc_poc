# Install docker
sudo apt-get update
sudo apt-get install docker

# Install golang and dependencies
go get google.golang.org/grpc
go get -u github.com/golang/protobuf/proto
go get -u github.com/golang/protobuf/protoc-gen-go
go get github.com/otaviocarvalho/hdrhistogram
go get github.com/VividCortex/ewma
go get go4.org/net/throttle
