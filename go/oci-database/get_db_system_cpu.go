package main

import (
        "context"
        "fmt"
        "errors"
        "github.com/oracle/oci-go-sdk/common"
        "github.com/oracle/oci-go-sdk/example/helpers"
        "github.com/oracle/oci-go-sdk/resourcesearch"
        "github.com/oracle/oci-go-sdk/database"
        "os"
)

func main() {
	dbSystemNameArg := ""
        if len(os.Args) == 2 {
		dbSystemNameArg = os.Args[1]
        } else {
                helpers.FatalIfError(errors.New("Usage : get_db_system_cpu <DB system name>"))
        }

        c, clerr := database.NewDatabaseClientWithConfigurationProvider(common.ConfigurationProviderEnvironmentVariables("oci", ""))
        helpers.FatalIfError(clerr)

	dbSystemOCID, e := dbsystemOCIDSearch(dbSystemNameArg)
        helpers.FatalIfError(e)

	getDbSystemReq := database.GetDbSystemRequest {
		DbSystemId:	common.String(dbSystemOCID),
	}

	res, err := c.GetDbSystem(context.Background(), getDbSystemReq)
	helpers.FatalIfError(err)

	fmt.Println(*res.DbSystem.CpuCoreCount)
}

func dbsystemOCIDSearch(name string) (string, error){
	client, err := resourcesearch.NewResourceSearchClientWithConfigurationProvider(common.ConfigurationProviderEnvironmentVariables("oci", ""))
	helpers.FatalIfError(err)

	q := "query dbsystem resources where displayName = '" + name + "'"
	searchReq := resourcesearch.SearchResourcesRequest{
		SearchDetails: resourcesearch.StructuredSearchDetails{
		MatchingContextType: resourcesearch.SearchDetailsMatchingContextTypeHighlights,
		Query:               common.String(q),
		},
	}

	structureSearchResp, err := client.SearchResources(context.Background(), searchReq)
	helpers.FatalIfError(err)
	if len(structureSearchResp.Items) == 1 {
		return *structureSearchResp.Items[0].Identifier,nil
	} else {
		return "", errors.New("Provided name did not return any db system")
	}
}
