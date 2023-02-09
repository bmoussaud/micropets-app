SOURCE_IMAGE = os.getenv("SOURCE_X_IMAGE", default='akseutap4registry.azurecr.io/pets-source')
LOCAL_PATH = os.getenv("LOCAL_PATH", default='.')
NAMESPACE = os.getenv("NAMESPACE", default='dev-tap')
OUTPUT_TO_NULL_COMMAND = os.getenv("OUTPUT_TO_NULL_COMMAND", default=' > /dev/null ')

compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/pets -buildmode pie -trimpath ./cmd/pets/main.go'

local_resource(
  'pets-build',
  compile_cmd,
  deps=['./cmd', './service','./internal'],
  dir='.')

allow_k8s_contexts('aks-eu-tap-4')

k8s_custom_deploy(
    'pets',
    apply_cmd="tanzu apps workload apply -f config/workload.yaml --update-strategy replace --debug --live-update" +
              " --local-path " + LOCAL_PATH +
              " --source-image " + SOURCE_IMAGE +
              " --namespace " + NAMESPACE +
              " --yes " +
              OUTPUT_TO_NULL_COMMAND +
              " && kubectl get workload pets-golang --namespace " + NAMESPACE + " -o yaml",
    delete_cmd="tanzu apps workload delete -f config/workload.yaml --namespace " + NAMESPACE + " --yes",
    deps=['./build'],
    container_selector='workload',
    live_update=[      
      sync('./build', '/tmp/tilt')  ,      
      run('cp -rf /tmp/tilt/* /layers/tanzu-buildpacks_go-build/targets/bin', trigger=['./build']),
    ]
)

k8s_resource('pets', port_forwards=["8080:8080"],
            extra_pod_selectors=[{'carto.run/workload-name': 'pets-golang','app.kubernetes.io/component': 'run'}])
