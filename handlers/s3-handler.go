package handlers

import (
	"awsx-api/log"
	"encoding/json"
	"fmt"
	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-s3/controller"
	"net/http"
)

func GetS3(w http.ResponseWriter, r *http.Request) {
	log.Info("Starting /awsx/s3 api")
	w.Header().Set("Content-Type", "application/json")

	region := r.URL.Query().Get("zone")
	vaultUrl := r.URL.Query().Get("vaultUrl")

	if vaultUrl != "" {
		accountId := r.URL.Query().Get("accountId")
		vaultToken := r.URL.Query().Get("vaultToken")
		authFlag, clientAuth, err := authenticate.AuthenticateData(vaultUrl, vaultToken, accountId, region, "", "", "", "")
		if err != nil || !authFlag {
			log.Error(err.Error())
			http.Error(w, fmt.Sprintf("Exception: "+err.Error()), http.StatusInternalServerError)
			return
		}
		result, respErr := controller.GetBucketList(clientAuth)
		if respErr != nil {
			log.Error(respErr.Error())
			http.Error(w, fmt.Sprintf("Exception: "+respErr.Error()), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(result)
	} else {
		accessKey := r.URL.Query().Get("accessKey")
		secretKey := r.URL.Query().Get("secretKey")
		crossAccountRoleArn := r.URL.Query().Get("crossAccountRoleArn")
		externalId := r.URL.Query().Get("externalId")
		authFlag, clientAuth, err := authenticate.AuthenticateData("", "", "", region, accessKey, secretKey, crossAccountRoleArn, externalId)
		if err != nil || !authFlag {
			log.Error(err.Error())
			http.Error(w, fmt.Sprintf("Exception: "+err.Error()), http.StatusInternalServerError)
			return
		}
		result, respErr := controller.GetBucketList(clientAuth)
		if respErr != nil {
			log.Error(respErr.Error())
			http.Error(w, fmt.Sprintf("Exception: "+respErr.Error()), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(result)
	}

	log.Info("/awsx/s3 completed")

}
