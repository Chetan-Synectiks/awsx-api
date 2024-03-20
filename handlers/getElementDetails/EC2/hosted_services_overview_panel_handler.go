package EC2

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getelementdetails/handler/EC2"
	"github.com/spf13/cobra"
)

func HostedServicesOverviewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract parameters from the URL query
	queries := r.URL.Query()
	region := queries.Get("zone")
	elementId := queries.Get("elementId")
	elementApiUrl := queries.Get("cmdbApiUrl")
	crossAccountRoleArn := queries.Get("crossAccountRoleArn")
	externalId := queries.Get("externalId")
	responseType := queries.Get("responseType")
	instanceId := queries.Get("instanceId")
	startTime := queries.Get("startTime")
	endTime := queries.Get("endTime")
	logGroupName := queries.Get("logGroupName")

	fmt.Printf("Received request with parameters: region=%s, elementId=%s, elementApiUrl=%s, crossAccountRoleArn=%s, externalId=%s, responseType=%s, instanceId=%s, startTime=%s, endTime=%s, logGroupName=%s\n",
		region, elementId, elementApiUrl, crossAccountRoleArn, externalId, responseType, instanceId, startTime, endTime, logGroupName)

	// Prepare command parameters
	commandParam := model.CommandParam{}
	if elementId != "" {
		commandParam.CloudElementId = elementId
		commandParam.CloudElementApiUrl = elementApiUrl
		commandParam.Region = region
	} else {
		commandParam.CrossAccountRoleArn = crossAccountRoleArn
		commandParam.ExternalId = externalId
		commandParam.Region = region
	}

	// Authenticate and get client credentials
	// clientAuth, err := authenticateAndCacheTracking(commandParam)
	// if err != nil {
	// 	sendErrorResponse(w, fmt.Sprintf("Authentication failed: %s", err), http.StatusInternalServerError)
	// 	return
	// }

	fmt.Println("Authentication successful")

	// Create Cobra command for passing flags
	cmd := &cobra.Command{}
	cmd.PersistentFlags().String("elementId", elementId, "Description of the elementId flag")
	cmd.PersistentFlags().String("instanceId", instanceId, "Description of the instanceId flag")
	cmd.PersistentFlags().String("elementType", queries.Get("elementType"), "Description of the elementType flag")
	cmd.PersistentFlags().String("startTime", startTime, "Description of the startTime flag")
	cmd.PersistentFlags().String("endTime", endTime, "Description of the endTime flag")
	cmd.PersistentFlags().String("responseType", responseType, "responseType flag - json/frame")
	cmd.PersistentFlags().String("logGroupName", logGroupName, "logGroupName flag - json/frame")

	// Parse flags
	if err := cmd.ParseFlags(nil); err != nil {
		sendErrorResponse(w, fmt.Sprintf("Failed to parse flags: %s", err), http.StatusInternalServerError)
		return
	}

	fmt.Println("Flags parsed successfully")
	// EC2.ListErrorEvents()
	// Call the function to get hosted service data
	hostedServiceData, err := EC2.GetHostedServicesData(cmd)
	if err != nil {
		esendErrorResponse(w, fmt.Sprintf("Failed to get hosted service data: %s", err), http.StatusInternalServerError)
		return
	}

	// Marshal the hosted service data into JSON
	hostedServiceJSON, err := json.Marshal(hostedServiceData)
	if err != nil {
		sendErrorResponse(w, fmt.Sprintf("Failed to marshal hosted service data to JSON: %s", err), http.StatusInternalServerError)
		return
	}

	// Write the JSON response
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(hostedServiceJSON); err != nil {
		sendErrorResponse(w, fmt.Sprintf("Failed to write response: %s", err), http.StatusInternalServerError)
		return
	}
}

func iauthenticateAndCacheTracking(commandParam model.CommandParam) (*model.Auth, error) {
	cacheKey := commandParam.CloudElementId

	authCacheLockTracking.Lock()
	defer authCacheLockTracking.Unlock()

	if _, ok := authCacheTracking.Load(cacheKey); ok {
		return nil, nil
	}

	_, _, err := authenticate.DoAuthenticate(commandParam)
	if err != nil {
		return nil, err
	}

	authCacheTracking.Store(cacheKey, nil)

	return nil, nil
}

func esendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := map[string]string{"error": message}
	json.NewEncoder(w).Encode(response)
}
