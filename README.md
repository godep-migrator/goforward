goforward
=========

Log forwarding tool used to forward logs over zeromq with protobuffers

if you want to re generate proto files run this from within the syslogMessage folder.
<pre>
<code>
go get -u code.google.com/p/gogoprotobuf/{proto,protoc-gen-gogo,gogoproto}
protoc --gogo_out=. -I=.:code.google.com/p/gogoprotobuf/protobuf -I=$GOPATH/src/ -I=$GOPATH/src/code.google.com/p/gogoprotobuf/protobuf *.proto
</code>
</pre>
