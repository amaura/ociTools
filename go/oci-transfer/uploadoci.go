package main

import (
        "context"
        //      "fmt"
        "errors"
        "github.com/oracle/oci-go-sdk/common"
        "github.com/oracle/oci-go-sdk/example/helpers"
        "github.com/oracle/oci-go-sdk/objectstorage"
        "io"
        "log"
        "os"
        //      "time"
        //      "github.com/oracle/oci-go-sdk/objectstorage/transfer"
)

func main() {
        bucketArg := ""
        fileArg := ""

        if len(os.Args) == 3 {
                fileArg = os.Args[1]
                bucketArg = os.Args[2]
        } else {
                helpers.FatalIfError(errors.New("Usage : upload <file name> <bucket name>"))
        }

        c, clerr := objectstorage.NewObjectStorageClientWithConfigurationProvider(common.ConfigurationProviderEnvironmentVariables("oci", ""))
        helpers.FatalIfError(clerr)
        uploadFile(c, fileArg, fileArg, bucketArg)
        //uploadFile(fileArg+"."+nowFormatted(), fileArg, bucketArg)

}

func uploadFile(client objectstorage.ObjectStorageClient, objectname string, filename string, bucketname string) {

        ctx := context.Background()

        namespace := getNamespace(ctx, client)

        file, e := os.Open(filename)
        defer file.Close()
        helpers.FatalIfError(e)

        stat, e := file.Stat()
        helpers.FatalIfError(e)

        e = putObject(ctx, client, namespace, bucketname, filename, stat.Size(), file, nil)
        helpers.FatalIfError(e)
}

func getNamespace(ctx context.Context, c objectstorage.ObjectStorageClient) string {
        request := objectstorage.GetNamespaceRequest{}
        r, err := c.GetNamespace(ctx, request)
        helpers.FatalIfError(err)
        //log.Println("Object Namespace is", *r.Value)
        return *r.Value
}

func putObject(ctx context.Context, c objectstorage.ObjectStorageClient, namespace, bucketname, objectname string, contentLen int64, content io.ReadCloser, metadata map[string]string) error {
        request := objectstorage.PutObjectRequest{
                NamespaceName: &namespace,
                BucketName:    &bucketname,
                ObjectName:    &objectname,
                ContentLength: &contentLen,
                PutObjectBody: content,
                OpcMeta:       metadata,
        }
        _, err := c.PutObject(ctx, request)
        log.Println("Putting object", objectname, "in bucket", bucketname)
        return err
}
