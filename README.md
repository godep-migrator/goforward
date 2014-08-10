goforward
=========
[![Build Status](https://travis-ci.org/CapillarySoftware/goforward.svg?branch=master)](https://travis-ci.org/CapillarySoftware/goforward)

Log forwarding tool used to forward logs over zeromq with protobuffers

if you want to re generate proto files run this from within the syslogMessage folder.

<h3>install</h3>
<pre>
<code>
//nanomsg
http://nanomsg.org/download.html
./configure && make && make check && sudo make install
go get github.com/tools/godep
godep restore
go get github.com/op/go-nanomsg
go get github.com/CapillarySoftware/goforward
</code>
</pre>

<h3>Rebuild gogoprotobuffer messages</h3>
<pre>
<code>
go get -u code.google.com/p/gogoprotobuf/{proto,protoc-gen-gogo,gogoproto}
protoc --gogo_out=. -I=.:code.google.com/p/gogoprotobuf/protobuf -I=$GOPATH/src/ -I=$GOPATH/src/code.google.com/p/gogoprotobuf/protobuf *.proto
</code></pre>

<h3>Install statically linked version</h3>
<pre>
<code>
go install --ldflags '-extldflags "-static"'  github.com/CapillarySoftware/gomasticate
</code>
</pre>

