SOURCE_IMAGE = os.getenv("SOURCE_X_IMAGE", default='akseutap4registry.azurecr.io/cats-source')
LOCAL_PATH = os.getenv("LOCAL_PATH", default='.')
NAMESPACE = os.getenv("NAMESPACE", default='micropets-dev')
OUTPUT_TO_NULL_COMMAND = os.getenv("OUTPUT_TO_NULL_COMMAND", default=' > /dev/null ')

compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/cats -buildmode pie -trimpath ./cmd/cats/main.go'

local_resource(
  'go-build',
  compile_cmd,
  deps=['./cmd', './service','./internal'],
  dir='.')

allow_k8s_contexts('aks-eu-tap-4')

k8s_custom_deploy(
    'cats',
    apply_cmd="tanzu apps workload apply -f config/workload.yaml --update-strategy replace --debug --live-update" +
              " --local-path " + LOCAL_PATH +
              " --source-image " + SOURCE_IMAGE +
              " --namespace " + NAMESPACE +
              " --yes " +
              OUTPUT_TO_NULL_COMMAND +
              " && kubectl get workload cats-golang --namespace " + NAMESPACE + " -o yaml",
    delete_cmd="tanzu apps workload delete -f config/workload.yaml --namespace " + NAMESPACE + " --yes",
    deps=['./build'],
    container_selector='workload',
    live_update=[      
      sync('./build', '/tmp/tilt')  ,      
      run('cp -rf /tmp/tilt/* /layers/tanzu-buildpacks_go-build/targets/bin', trigger=['./build']),
    ]
)

k8s_resource('cats', port_forwards=["8080:8080"],
            extra_pod_selectors=[{'carto.run/workload-name': 'cats-golang','app.kubernetes.io/component': 'run'}])
