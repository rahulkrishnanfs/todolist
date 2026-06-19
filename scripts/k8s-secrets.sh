kubectl create secret generic todolist-secrets \
  --from-file=keystore.p12=./secrets/keystore.p12 \
  --from-file=servercert.pem=./secrets/servercert.pem \
  --from-file=serverkey.pem=./secrets/serverkey.pem