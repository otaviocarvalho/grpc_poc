# Move to base dir 
cd /home/ocarval/

# Install pssh
sudo apt-get install pssh

# Install golang
mkdir /home/ocarval/go
time (while ps -opid= -C apt-get > /dev/null; do sleep 1; done); # wait for debian update to release dpkg lock
sudo apt-get install golang-go -y
echo "export GOPATH=/home/ocarval/go" >> /home/ocarval/.bash_aliases
echo "export PATH=$PATH:$GOROOT/bin:$GOPATH/bin" >> /home/ocarval/.bash_aliases
source /home/ocarval/.bash_aliases
sudo chown -R ocarval.ocarval *
sudo chown -R ocarval.ocarval .*

# Install itself as a dependency 
go get -u github.com/otaviocarvalho/grpc_poc

# Install golang and dependencies
go get google.golang.org/grpc
go get -u github.com/golang/protobuf/proto
go get -u github.com/golang/protobuf/protoc-gen-go
go get github.com/otaviocarvalho/hdrhistogram
go get github.com/VividCortex/ewma
go get go4.org/net/throttle
