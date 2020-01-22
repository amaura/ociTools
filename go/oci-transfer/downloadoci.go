package main

import (
        "context"
//        "fmt"
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
        fileOutArg := ""

        if len(os.Args) == 4 {
                fileArg = os.Args[1]
                bucketArg = os.Args[2]
                fileOutArg = os.Args[3]
        } else {
                helpers.FatalIfError(errors.New("Usage : downloadoci <file name> <bucket name> <destination file>"))
        }

        c, clerr := objectstorage.NewObjectStorageClientWithConfigurationProvider(common.ConfigurationProviderEnvironmentVariables("oci", ""))
        helpers.FatalIfError(clerr)
        downloadFile(c, fileArg, bucketArg, fileOutArg)

}

func downloadFile(client objectstorage.ObjectStorageClient, filename string, bucketname string, outFilename string) {

        ctx := context.Background()
        namespace := getNamespace(ctx, client)
        r, e := getObject(ctx, client, namespace, bucketname, filename)
        helpers.FatalIfError(e)

        outFile, e := os.Create(outFilename)
        helpers.FatalIfError(e)
	defer outFile.Close()

	_, e = io.Copy(outFile, r.Content)
        helpers.FatalIfError(e)


}

func getNamespace(ctx context.Context, c objectstorage.ObjectStorageClient) string {
        request := objectstorage.GetNamespaceRequest{}
        r, err := c.GetNamespace(ctx, request)
        helpers.FatalIfError(err)
        //log.Println("Object Namespace is", *r.Value)
        return *r.Value
}

func getObject(ctx context.Context, c objectstorage.ObjectStorageClient, namespace, bucketname, objectname string) (objectstorage.GetObjectResponse, error) {
        request := objectstorage.GetObjectRequest{
                NamespaceName: &namespace,
                BucketName:    &bucketname,
                ObjectName:    &objectname,
        }
        o, err := c.GetObject(ctx, request)
        log.Println("Getting object", objectname, "from bucket", bucketname)
        return o, err
}
