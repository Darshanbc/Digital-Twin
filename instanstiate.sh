jq --version > /dev/null 2>&1
if [ $? -ne 0 ]; then
	echo "Please Install 'jq' https://stedolan.github.io/jq/ to execute this script"
	echo
	exit 1
fi

starttime=$(date +%s)

# Print the usage message
function printHelp () {
  echo "Usage: "
  echo "  ./testAPIs.sh -l golang|node"
  echo "    -l <language> - chaincode language (defaults to \"golang\")"
}
# Language defaults to "golang"
LANGUAGE="golang"

# Parse commandline args
while getopts "h?l:" opt; do
  case "$opt" in
    h|\?)
      printHelp
      exit 0
    ;;
    l)  LANGUAGE=$OPTARG
    ;;
  esac
done

##set chaincode path
function setChaincodePath(){
	LANGUAGE=`echo "$LANGUAGE" | tr '[:upper:]' '[:lower:]'`
	case "$LANGUAGE" in
		"golang")
		CC_SRC_PATH="github.com/example_cc/go"
		;;
		"node")
		CC_SRC_PATH="$PWD/artifacts/src/github.com/example_cc/node"
		;;
		*) printf "\n ------ Language $LANGUAGE is not supported yet ------\n"$
		exit 1
	esac
}

setChaincodePath

echo "POST request Enroll on Org1  ..."
echo
ORG1_TOKEN=$(curl -s -X POST \
  http://localhost:4000/users \
  -H "content-type: application/x-www-form-urlencoded" \
  -d 'username=DarshanBC6&orgName=Org1&role="Farmer"')
echo $ORG1_TOKEN
ORG1_TOKEN=$(echo $ORG1_TOKEN | jq ".token" | sed "s/\"//g")



# echo "POST invoke KYC on peers of Org1 and Org2"
# echo
# TRX_ID=$(curl -s -X POST \
#   http://localhost:4000/channels/mychannel/chaincodes/dtwin \
#   -H "authorization: Bearer $ORG1_TOKEN" \
#   -H "content-type: application/json" \
#   -d '{
#   		"fcn":"newuser",
#   		"args":["Darshan","BC","25","bangalore","darshan@gmail.com","8904374405"] 	
# }')
# echo "Transacton ID is $TRX_ID"




# echo "POST invoke PlotRegisteration on peers of Org1 and Org2"

# echo
# TRX_ID=$(curl -s -X POST \
#   http://localhost:4000/channels/mychannel/chaincodes/dtwin \
#   -H "authorization: Bearer $ORG1_TOKEN" \
#   -H "content-type: application/json" \
#   -d '{

# 	"fcn":"PlotRegisteration",
# 	"args":["{\"survey_no\":\"Eg456-456-5677\",\"soil_type\":\"sasa\",\"co_ordinates\":[{\"lattitude\":45.15445,\"longitude\":32.45455},{\"lattitude\":45.15445,\"longitude\":32.45455},{\"lattitude\":45.15445,\"longitude\":32.45455},{\"lattitude\":45.15445,\"longitude\":32.45455}]}"]
 	
# }')
# echo "Transacton ID is $TRX_ID"


# echo "POST invoke insert crop history on peers of Org1 and Org2"
# echo
# TRX_ID=$(curl -s -X POST \
#   http://localhost:4000/channels/mychannel/chaincodes/dtwin \
#   -H "authorization: Bearer $ORG1_TOKEN" \
#   -H "content-type: application/json" \
#   -d '{

# 	"fcn":"CropDetails",
# 	"args":["Plot0","{\"crop_name\":\"rice\",\"crop_cycle\":[{\"from_month\":\"june\",\"to_month\":\"september\"}],\"fertilzer_used\":[{\"fertlizer_name\":\"urea\",\"fertlizer_id\":\"U35\",\"quantity\":2.5}]}","1"]
 	
# }')

# echo "Transacton ID is $TRX_ID"
# echo
# echo

# echo "POST invoke certRegistration on peers of Org1 and Org2"
# echo
# TRX_ID=$(curl -s -X POST \
#   http://localhost:4000/channels/mychannel/chaincodes/dtwin \
#   -H "authorization: Bearer $ORG1_TOKEN" \
#   -H "content-type: application/json" \
#   -d '{

# 	"fcn":"CropDetails",
# 	"args":["Plot0","{\"crop_name\":\"rice\",\"crop_cycle\":[{\"from_month\":\"june\",\"to_month\":\"september\"}],\"fertilzer_used\":[]}","certRegistration"]
 	
# }')

# echo "Transacton ID is $TRX_ID"
# echo
# echo


# echo "POST invoke addFertilizerToStore on peers of Org1 and Org2"
# echo
# TRX_ID=$(curl -s -X POST \
#   http://localhost:4000/channels/mychannel/chaincodes/dtwin \
#   -H "authorization: Bearer $ORG1_TOKEN" \
#   -H "content-type: application/json" \
#   -d '{

# 	"fcn":"addFertilizerToStore",
# 	"args":["{\"fertlizer_name\":\"urea\",\"fertlizer_id\":\"U35\",\"quantity\":2.5}"]
 	
# }')

# echo "Transacton ID is $TRX_ID"
# echo
# echo


# echo "POST invoke addFertilizerToCrop on peers of Org1 and Org2"
# echo
# TRX_ID=$(curl -s -X POST \
#   http://localhost:4000/channels/mychannel/chaincodes/dtwin \
#   -H "authorization: Bearer $ORG1_TOKEN" \
#   -H "content-type: application/json" \
#   -d '{

# 	"fcn":"addFertilizerToCrop",
# 	"args":["Plot0","{\"fertlizer_name\":\"urea\",\"fertlizer_id\":\"U35\",\"quantity\":2.5}"]
 	
# }')

# echo "Transacton ID is $TRX_ID"
# echo
# echo


# echo "GET query chaincode on peer1 of Org1"
# echo
# key="Criyagen"
# curl -s -X GET \
#   "http://localhost:4000/channels/mychannel/chaincodes/dtwin?peer=peer0.org1.example.com&fcn=query&args=%5B%22$key%22%5D" \
#   -H "authorization: Bearer $ORG1_TOKEN" \
#   -H "content-type: application/json"
# echo
# echo

TRX_ID=42f75f75d2ebfdbcf53173256c7a6ab627b97769514604e0fe116b476667dc10
echo "GET query Transaction by TransactionID"
echo
curl -s -X GET http://localhost:4000/channels/mychannel/transactions/$TRX_ID?peer=peer0.org1.example.com&fcn="ApproveOrDeny" \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json"
echo
echo




