RUN PROTOC


protoc --proto_path=. --go_out=plugins=grpc:. .\common\model\*.proto
