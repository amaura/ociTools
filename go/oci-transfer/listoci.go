package main

import (
        "context"
        "fmt"
        "errors"
        "github.com/oracle/oci-go-sdk/common"
        "github.com/oracle/oci-go-sdk/example/helpers"
        "github.com/oracle/oci-go-sdk/objectstorage"
        "os"
	"strconv"
)

func main() {
	bucketArg := ""
	prefixArg := ""
	limitArg := ""
        if len(os.Args) == 4 {
		bucketArg = os.Args[1]
		prefixArg = os.Args[2]
		limitArg = os.Args[3]
        } else {
                helpers.FatalIfError(errors.New("Usage : listoci <bucket name> <prefix> <limit>"))
        }

        c, clerr := objectstorage.NewObjectStorageClientWithConfigurationProvider(common.ConfigurationProviderEnvironmentVariables("oci", ""))
        helpers.FatalIfError(clerr)
	limitArgInt,err := strconv.Atoi(limitArg)
        helpers.FatalIfError(err)
	listFiles(c, bucketArg,prefixArg,limitArgInt)

}

func listFiles(client objectstorage.ObjectStorageClient, bucketname, prefix string, limit int) {

        ctx := context.Background()

        namespace := getNamespace(ctx, client)

	r,e := listObject(ctx, client, namespace,bucketname,prefix,limit)
        helpers.FatalIfError(e)
	for _, v := range r.ListObjects.Objects {
		fmt.Println(*v.Name)
	}
}

func getNamespace(ctx context.Context, c objectstorage.ObjectStorageClient) string {
        request := objectstorage.GetNamespaceRequest{}
        r, err := c.GetNamespace(ctx, request)
        helpers.FatalIfError(err)
        return *r.Value
}

func listObject(ctx context.Context, c objectstorage.ObjectStorageClient, namespace, bucketname,prefix string ,limit int) (objectstorage.ListObjectsResponse,error) {
        request := objectstorage.ListObjectsRequest{
                NamespaceName: &namespace,
                BucketName:    &bucketname,
		Prefix:        &prefix,
		Limit:         &limit}

        res, err := c.ListObjects(ctx, request)
        return res,err
}
