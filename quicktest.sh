#!/bin/bash
source .env
echo "Login"
curl -d "username=admin&password=$ADMIN_PASSWORD" -c .cookie http://127.0.0.1:8000/api/login 
echo "New record"
curl -b .cookie -d "user_id=101234567&pass=true&temperature=0" -X POST http://127.0.0.1:8000/api/records
echo "List record"
curl -b .cookie http://127.0.0.1:8000/api/records
echo "Check record"
curl http://127.0.0.1:8000/api/check?user_id=101234567