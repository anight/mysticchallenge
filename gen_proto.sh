#! /bin/bash

ALL_PROTO="proto/challenge.proto"

docker run -it --rm -v $(pwd):/app -w /app anight/pref-protoc \
	protoc \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		${ALL_PROTO}

docker run -it --rm -v $(pwd):/app -w /app anight/pref-protoc \
	python -m grpc_tools.protoc -Iproto \
		--python_out=proto \
		--grpc_python_out=proto \
		${ALL_PROTO}
