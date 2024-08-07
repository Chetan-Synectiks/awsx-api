package getLandingZoneDetails

import (
	"encoding/json"
	"fmt"
	"github.com/Appkube-awsx/awsx-common/authenticate"
	"github.com/Appkube-awsx/awsx-common/model"
	"github.com/Appkube-awsx/awsx-getlandingzonedetails/handler/API_GW"
	"github.com/Appkube-awsx/awsx-getlandingzonedetails/handler/CDN"
	"github.com/Appkube-awsx/awsx-getlandingzonedetails/handler/CLOUDWATCH"
	"github.com/Appkube-awsx/awsx-getlandingzonedetails/handler/CONFIG_SERVICE"
	"github.com/Appkube-awsx/awsx-getlandingzonedetails/handler/DYNAMODB"
	"github.com/Appkube-awsx/awsx-getlandingzonedetails/handler/EC2"
	"github.com/Appkube-awsx/awsx-getlandingzonedetails/handler/ECS"
	"github.com/Appkube-awsx/awsx-getlandingzonedetails/handler/EKS"
	"github.com/Appkube-awsx/awsx-getlandingzonedetails/handler/KINESIS"
	"github.com/Appkube-awsx/awsx-getlandingzonedetails/handler/KMS"
	"github.com/Appkube-awsx/awsx-getlandingzonedetails/handler/LAMBDA"
	"github.com/Appkube-awsx/awsx-getlandingzonedetails/handler/LB"
	"github.com/Appkube-awsx/awsx-getlandingzonedetails/handler/RDS"
	"github.com/Appkube-awsx/awsx-getlandingzonedetails/handler/S3"
	"github.com/Appkube-awsx/awsx-getlandingzonedetails/handler/SSL"
	"github.com/Appkube-awsx/awsx-getlandingzonedetails/handler/VPC"
	"github.com/Appkube-awsx/awsx-getlandingzonedetails/handler/WAF"
	"net/http"
)

func ExecuteLandingzoneQueries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("query")
	landingZoneId := r.URL.Query().Get("landingZoneId")
	logGroupName := r.URL.Query().Get("logGroupName")
	commandParam := model.CommandParam{
		LandingZoneId: landingZoneId,
	}
	_, clientAuth, err := authenticate.DoAuthenticate(commandParam)
	if err != nil {
		http.Error(w, fmt.Sprintf("error in getting aws credentials: %s", err), http.StatusInternalServerError)
		return
	}
	var instances interface{}
	if query == "getEc2List" {
		instances, err = EC2.ListEc2Instances(clientAuth, nil)
	}
	if query == "getCdnList" {
		instances, err = CDN.CdnDistributionConfigWithTagList(clientAuth, nil)
	}
	if query == "getCdnFunctionList" {
		instances, err = CDN.ListCdnFunctionInstances(clientAuth, nil)
	}
	if query == "getApiGwList" {
		instances, err = API_GW.ListApiGwInstances(clientAuth, nil)
	}
	if query == "getLbList" {
		instances, err = LB.ListLbInstances(clientAuth, nil)
	}
	if query == "getDynamoDbList" {
		instances, err = DYNAMODB.ListDynamoDbInstance(clientAuth, nil)
	}
	if query == "getEcsList" {
		instances, err = ECS.ListEcsInstances(clientAuth, nil)
	}
	if query == "getEksList" {
		instances, err = EKS.ListEksInstances(clientAuth, nil)
	}
	if query == "getKinesisList" {
		instances, err = KINESIS.ListKinesisInstances(clientAuth, nil)
	}
	if query == "getKinesisRecordList" {
		instances, err = KINESIS.ListKinesisRecordInstances(clientAuth, nil)
	}
	if query == "getKmsList" {
		instances, err = KMS.ListKmsInstances(clientAuth, nil)
	}
	if query == "getLambdaList" {
		instances, err = LAMBDA.ListLambdaInstances(clientAuth, nil)
	}
	if query == "getRdsList" {
		instances, err = RDS.ListRdsInstances(clientAuth, nil)
	}
	if query == "" {
		instances, err = S3.ListS3Instances(clientAuth, nil)
	}
	if query == "getVpcList" {
		instances, err = VPC.ListVpcInstances(clientAuth, nil)
	}
	if query == "getWafList" {
		instances, err = WAF.ListWafInstances(clientAuth, nil)
	}
	if query == "getDiscoveredResourceCounts" {
		instances, err = CONFIG_SERVICE.DiscoveredResourceCounts(clientAuth, nil)
	}
	if query == "getSslList" {
		instances, err = SSL.ListSslInstances(clientAuth, nil)
	}
	if query == "getCwAlarmList" {
		instanceId := r.URL.Query().Get("instanceId")
		if instanceId == "" {
			http.Error(w, fmt.Sprintf("instance id missing"), http.StatusBadRequest)
			return
		}
		instances, err = CLOUDWATCH.ListCwAlarms(instanceId, clientAuth, nil)
	}
	if query == "getCwLogsStreamList" {
		logGroupName = r.URL.Query().Get("logGroupName")
		if logGroupName == "" {
			http.Error(w, fmt.Sprintf("logGroup name  missing"), http.StatusBadRequest)
			return
		}
		instances, err = CLOUDWATCH.ListCwLogsStream(logGroupName, clientAuth, nil)
	}
	if query == "getCwLogsGorupList" {

		instances, err = CLOUDWATCH.ListCwLogsGorup(clientAuth, nil)
	}
	if query == "getCwEventList" {
		instanceId := r.URL.Query().Get("instanceId")
		if instanceId == "" {
			http.Error(w, fmt.Sprintf("instance id missing"), http.StatusBadRequest)
			return
		}
		instances, err = CLOUDWATCH.ListCwEvent(instanceId, clientAuth, nil)
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("api error: %s", err), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(instances)
	if err != nil {
		http.Error(w, fmt.Sprintf("errror in json encoding %s ", err), http.StatusInternalServerError)
		return
	}
}
