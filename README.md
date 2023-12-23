## Description

Imagine a client has a large set of potentially small files {F0, F1, …, Fn} and wants to upload them to a server and then delete its local copies. The client wants, however, to later download an arbitrary file from the server and be convinced that the file is correct and is not corrupted in any way (in transport, tampered with by the server, etc.).

You should implement the client, the server and a Merkle tree to support the above (we expect you to implement the Merkle tree rather than use a library, but you are free to use a library for the underlying hash functions).

The client must compute a single Merkle tree root hash and keep it on its disk after uploading the files to the server and deleting its local copies. The client can request the i-th file Fi and a Merkle proof Pi for it from the server. The client uses the proof and compares the resulting root hash with the one it persisted before deleting the files - if they match, file is correct.

## Usage
check tests and linter
```shell
make lint && make test 
```

### Local run
Start app and postgres in docker
```shell
make run 
```
Send file to app:
```shell
# this command will generate some sample files in data_folder
make upload

# or upload files from custom folder
go run cmd/uploader/main.go --path=/tmp/7bb9389dcfa58ad8965aa77f643dc60eeef881cb2fa59def965fddb1bec4cfeb/
```

## Implementation details
#### Merkle Tree
tree represented as slice of slices, where 0 index is bottom of tree and last index is root of tree. 

See `internal/service/merkle_tree/service.go`. For more abstract usage it build with generics, so it can be used with files of memory objects, only need provide correct hash function.

For quick rebuild and build proof tree have map, which allow by hash quickly find changed item and rebuild only changed branch. 

See `internal/service/merkle_tree/verify.go` for more details.

#### File Storage
File storage is simple file system storage, which store files in `/tmp/zm` folder. It is not production ready, but it can be replaced by some scalable storage.

As example, another test solution can be used for adopt as flexible storage https://github.com/abergasov/extendable_storage

so final flow may be:
* user upload file to app
* app verify files, calculate hashes and pass it in storage
* storage based on hash determine where to store file
* storage store file and return file id
* app store file id and hashes in db
* user can download file by id

### Improvements
* Better API Responses (now return status code only):
  * Return more detailed responses with status codes and meaningful messages. 
  * Include additional information in the response, like file metadata or IDs.
* Idempotent Upload (user for some reason may send same file several times or later send same files). Current implementation will serve 500, in case unique constraint violation. Depending on business logic we can:
  * Implement single-flight pattern to avoid processing the same file multiple times.
  * Consider versioning files, allowing multiple uploads without deletion.
  * Make entity `upload_session` which will be unique for each upload session and will allow to upload same files several times
* Merkle Tree:
  * Rebuild not implemented. For quick rebuild tree have map, which allow by hash quickly find changed item and rebuild only changed branch
  * Current implementation of tree probably may not be optimal, so it can be optimized for better performance. Need pprof and benchmarks to find bottlenecks
  * Current implementation is not concurrency safe
* File Storage (current file storage simply save data in `/tmp/zm` folder, so it should be improved by following ways):
  * Implement a garbage collector to remove old or unprocessed files.
  * For non-blocking upload, user can generate `session_id` which allow to upload files in several requests and then commit them. It will allow to upload files in parallel and then commit them in one request.
  * Add scalable storage which will allow to store files in case of usage growth
  * Sample implementation https://github.com/abergasov/extendable_storage
    * here is scale logic based on consistent hashing, so we can add new nodes and remove old nodes without full data migration
    * file can be chunked to files to allow better equal distribution of data between nodes and better disk utilization
* Add metrics for monitoring and alerting
* Security Measures:
  * Hide app behind proxy which will take care of initial checks, like:
    * limit number of connections to app
    * limit size of request body
    * limit number of requests per second
  * Add authentication to prevent that user will send files which he should not send
  * Add PoW to prevent spamming of app. For easy usage it can be raw TCP connection, so it will be easy to integrate with any language, and PoW will be done on client side from SDK