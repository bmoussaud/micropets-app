package=$1
environment=$2

echo "
apiVersion: xl-deploy/v1
kind: Deployment
spec:
  package: ${package}
  environment: ${environment}
  onSuccessPolicy: ARCHIVE
" > /tmp/deploy.yaml

cat /tmp/deploy.yaml
./xlw preview --xl-deploy-url http://localhost:4516 -f /tmp/deploy.yaml

./xlw apply --xl-deploy-url http://localhost:4516 -f /tmp/deploy.yaml
