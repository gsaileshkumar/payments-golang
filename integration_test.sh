#!/bin/bash

echo "Ensure DB is clean, before executing this test"

# API endpoint
URL="http://localhost:8080/accounts"

# Data
DATA='{"account_id":12345, "initial_balance":"1000.00"}'

# Perform the POST request
RESPONSE=$(curl -s -X POST -H "Content-Type: application/json" -d "$DATA" "$URL")

# Extract values using jq
ACCOUNT_ID=$(echo $RESPONSE | jq '.account_id')
BALANCE=$(echo $RESPONSE | jq '.balance')

# Check if the values are correct
if [[ $ACCOUNT_ID == 12345 && $BALANCE == "\"1000.000000"\" ]]; then
  echo "Create account test passed"
else
  echo "Create account test failed"
fi

# API endpoint
URL="http://localhost:8080/accounts/12345"

# Perform the GET request
RESPONSE=$(curl -s -X GET "$URL")

# Extract values using jq
ACCOUNT_ID=$(echo $RESPONSE | jq '.account_id')
BALANCE=$(echo $RESPONSE | jq '.balance')

# Check if the values are correct
if [[ $ACCOUNT_ID == 12345 && "$BALANCE" == "\"1000.000000"\"  ]]; then
  echo "Get account details test passed"
else
  echo "Get account details test failed"
fi

# Create second account
URL="http://localhost:8080/accounts"
DATA='{"account_id":67890, "initial_balance":"100.55555"}'
RESPONSE=$(curl -s -X POST -H "Content-Type: application/json" -d "$DATA" "$URL")
ACCOUNT_ID=$(echo $RESPONSE | jq '.account_id')
BALANCE=$(echo $RESPONSE | jq '.balance')

# Check if the values are correct
if [[ $ACCOUNT_ID == 67890 && $BALANCE == "\"100.555550"\" ]]; then
  echo "Create second account test passed"
else
  echo "Create second account test failed"
fi

# API endpoint
URL="http://localhost:8080/transactions"
# Data
DATA='{"source_account_id":12345, "destination_account_id":67890, "amount":"10.000002"}'
# Perform the POST request
RESPONSE=$(curl -s -X POST -H "Content-Type: application/json" -d "$DATA" "$URL" -w "\n%{http_code}")
HTTP_STATUS=$(echo "$RESPONSE" | tail -n 1)
# Check for success indication, assuming API sends a specific success message or status
if [[ $HTTP_STATUS -eq 200 ]]; then
  echo "Transaction test passed"
else
  echo "Transaction test failed"
fi

URL="http://localhost:8080/accounts/12345"
RESPONSE=$(curl -s -X GET "$URL")
ACCOUNT_ID=$(echo $RESPONSE | jq '.account_id')
BALANCE=$(echo $RESPONSE | jq '.balance')

# Check if the values are correct
if [[ $ACCOUNT_ID == 12345 && "$BALANCE" == "\"989.999998"\"  ]]; then
  echo "Source account updated test passed"
else
  echo "Source account updated test failed"
fi

URL="http://localhost:8080/accounts/67890"
RESPONSE=$(curl -s -X GET "$URL")
ACCOUNT_ID=$(echo $RESPONSE | jq '.account_id')
BALANCE=$(echo $RESPONSE | jq '.balance')

# Check if the values are correct
if [[ $ACCOUNT_ID == 67890 && "$BALANCE" == "\"110.555552"\"  ]]; then
  echo "Destination account updated test passed"
else
  echo "Destination account updated test failed"
fi