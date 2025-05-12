import os, time, requests, streamlit as st
import pandas as pd
import numpy as np
from ..autoscaler.predictor import ARPredictor
from ..autoscaler.scaler import AutoScaler, ScalingConfig

# ------------------ Настройка Prometheus live ------------------
PROM_URL = os.getenv("PROMETHEUS_URL", "http://localhost:9090")


def query(expr):
    try:
        resp = requests.get(f"{PROM_URL}/api/v1/query",
                            params={"query": expr}, timeout=5)
        resp.raise_for_status()
        data = resp.json()["data"]["result"]
        return float(data[0]["value"][1]) if data else None
    except Exception:
        return None


# ------------------ UI ------------------
st.title("Intelligent Autoscaler dashboard")

mode = st.sidebar.radio("Источник данных", ("Live из Prometheus", "CSV‑файл"))

# ---------- LIVE‑РЕЖИМ ----------
if mode == "Live из Prometheus":
    deploy = st.sidebar.text_input("Deployment", "my‑service")
    ns = st.sidebar.text_input("Namespace", "default")
    refresh = st.sidebar.number_input("Интервал (сек.)", 1, 30, 5)

    expr_rps = f'rate(http_requests_total{{deployment="{deploy}"}}[1m])'
    expr_rep = f'kube_deployment_status_replicas{{deployment="{deploy}",namespace="{ns}"}}'

    hist = st.empty()
    if "live_hist" not in st.session_state:
        st.session_state.live_hist = []

    while True:
        rps = query(expr_rps) or 0
        rep = query(expr_rep) or 0
        st.session_state.live_hist.append({"time": pd.Timestamp.utcnow(),
                                           "rps": rps, "replicas": rep})
        df = pd.DataFrame(st.session_state.live_hist).set_index("time")
        hist.line_chart(df)
        time.sleep(refresh)
        st.experimental_rerun()  # перезапустить цикл

# ---------- CSV‑РЕЖИМ ----------
else:
    up = st.file_uploader("CSV с колонками time,rps", type="csv")
    if up:
        df = pd.read_csv(up)
        if not {"time", "rps"}.issubset(df.columns):
            st.error("Нет колонок time,rps")
        else:
            df["time"] = pd.to_datetime(df["time"])
            df = df.sort_values("time")

            predictor = ARPredictor(p=2)
            scaler = AutoScaler(ScalingConfig())
            preds, reps = [], []

            for r in df["rps"]:
                predictor.update(r)
                pred = predictor.predict() or r
                preds.append(pred)
                reps.append(scaler.optimise(pred, scaler.prev_replicas))

            df["pred_rps"] = preds
            df["replicas"] = reps

            st.subheader("RPS и прогноз")
            st.line_chart(df.set_index("time")[["rps", "pred_rps"]])

            st.subheader("Рекомендуемые реплики")
            st.line_chart(df.set_index("time")[["replicas"]])
