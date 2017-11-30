# Install docker
sudo apt-get update
sudo apt-get install docker

# Install git
sudo apt-get install git

# Clone repo
mkdir go && mkdir go/src
git clone https://github.com/otaviocarvalho/grpc_poc.git ./go/src/grpc_poc

# Install golang
curl -sSLO https://storage.googleapis.com/golang/go1.8.1.linux-armv6l.tar.gz
sudo mkdir -p /usr/local/go  
sudo tar -xvf go1.8.1.linux-armv6l.tar.gz -C /usr/local/go --strip-components=1

# Export env variables
export GOPATH=$HOME/go && export GOBIN=$GOPATH/bin && export PATH=$PATH:$GOBIN

# Install golang and dependencies
/usr/local/go/bin/go get google.golang.org/grpc
/usr/local/go/bin/go get -u github.com/golang/protobuf/proto
/usr/local/go/bin/go get -u github.com/golang/protobuf/protoc-gen-go
/usr/local/go/bin/go get github.com/otaviocarvalho/hdrhistogram
/usr/local/go/bin/go get github.com/VividCortex/ewma
/usr/local/go/bin/go get go4.org/net/throttle
