#!/usr/bin/env bash
set -e

NAMESPACE=loadtest
RUN_TIME=30m

function run_test() {
  SUFFIX=$1         # hpa / custom
  SCALE_CMD=$2      # enable-hpa, enable-keda

  echo "[+] Reset environment"
  kubectl scale deploy mathcruncher -n $NAMESPACE --replicas=1
  kubectl delete hpa,scaledobject -n $NAMESPACE --all
  $SCALE_CMD

  echo "[+] Warm-up 60 s"
  sleep 60

  START=$(date +%s)
  locust -f locustfile.py,shape_step.py --headless \
         --csv $SUFFIX \
         --run-time $RUN_TIME \
         --host http://math-svc.loadtest \
         --only-summary
  END=$(date +%s)
  echo "[+] Load finished, $((END-START)) s"
}

run_test hpa    "kubectl apply -f hpa.yaml"
run_test custom "kubectl apply -f keda.yaml"
