# --go_outbandera es el directorio donde desea que el compilador escriba su salida de Go
# --go_out=. generate files in the same folder of *.proto
protoc --go_out=. --go_opt=paths=source_relative \
     --go-grpc_out=. --go-grpc_opt=paths=source_relative \
   greet/greetpb/greet.proto

# if you need different outputs, such as java or python
#protoc --python_out=/your-fullpath/python greet/greetpb/*.proto

#protoc --java_out=/your-fullpath/java greet/greetpb/*.proto

 