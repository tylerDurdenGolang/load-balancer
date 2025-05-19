
# Intelligent Autoscaler

Proactive autoscaling microservice for Kubernetes.  
Uses an adaptive AR(*p*) predictor (Recursive Least Squares) to forecast load and decide the optimal replica count.

```bash
# local
pip install -r requirements.txt
python -m autoscaler.main --config config.yaml
```

See `tests/` for unit tests and `k8s/` for example manifests.


## Streamlit dashboard

Build & run locally:

```bash
export PROMETHEUS_URL=http://localhost:9090
streamlit run dashboard/app.py
```

Or with Docker:

```bash
docker build -f Dockerfile -t autoscaler-dashboard .
docker run -p 8501:8501 -e PROMETHEUS_URL=http://prometheus:9090 autoscaler-dashboard
```
