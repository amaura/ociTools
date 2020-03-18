package main

import (
        "context"
//        "fmt"
        "errors"
        "github.com/oracle/oci-go-sdk/common"
        "github.com/oracle/oci-go-sdk/example/helpers"
        "github.com/oracle/oci-go-sdk/resourcesearch"
        "github.com/oracle/oci-go-sdk/database"
        "os"
	"strconv"
)

func main() {
	dbSystemNameArg := ""
	dbSystemCoreCountArg := ""
        if len(os.Args) == 3 {
		dbSystemNameArg = os.Args[1]
		dbSystemCoreCountArg  = os.Args[2]
        } else {
                helpers.FatalIfError(errors.New("Usage : set_db_system_cpu  <DB system name> <Number of CPU>"))
        }

	dbSystemCoreCountArgInt,err := strconv.Atoi(dbSystemCoreCountArg)
        helpers.FatalIfError(err)

        c, clerr := database.NewDatabaseClientWithConfigurationProvider(common.ConfigurationProviderEnvironmentVariables("oci", ""))
        helpers.FatalIfError(clerr)

	dbSystemOCID, e := dbsystemOCIDSearch(dbSystemNameArg)
        helpers.FatalIfError(e)

	updateCoreCountReq := database.UpdateDbSystemRequest {
		DbSystemId:	common.String(dbSystemOCID),
		UpdateDbSystemDetails: database.UpdateDbSystemDetails{
		CpuCoreCount: common.Int(dbSystemCoreCountArgInt),
		},
	}

	_, reqerr := c.UpdateDbSystem(context.Background(), updateCoreCountReq)
	helpers.FatalIfError(reqerr)
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
