import os
import time
import requests
import streamlit as st
import pandas as pd

# ------------------ Настройка Prometheus ------------------
PROM_URL = os.getenv("PROMETHEUS_URL", "http://localhost:9090")

def query(expr: str) -> float | None:
    """Выполняет запрос к Prometheus и возвращает числовое значение либо None."""
    try:
        resp = requests.get(
            f"{PROM_URL}/api/v1/query",
            params={"query": expr},
            timeout=5,
        )
        resp.raise_for_status()
        data = resp.json()["data"]["result"]
        return float(data[0]["value"][1]) if data else None
    except Exception:
        return None

# ------------------ UI ------------------
st.set_page_config(
    page_title="Intelligent Autoscaler — Live",
    page_icon="⚡",
    layout="wide",
)

st.title("⚡ Intelligent Autoscaler — Live dashboard")

# ----------- Сайдбар -------------
with st.sidebar:
    st.header("Параметры запроса")
    deploy = st.text_input("Deployment", "mathcruncher")
    ns = st.text_input("Namespace", "mathcruncher")
    refresh = st.number_input("Интервал обновления (сек.)", 1, 30, 5)
    col_run, col_stop = st.columns(2)
    run_clicked = col_run.button("▶️ Запустить")
    stop_clicked = col_stop.button("⏹️ Остановить")

# ----------- Управление состоянием -------------
if run_clicked:
    st.session_state["running"] = True
if stop_clicked:
    st.session_state["running"] = False
running: bool = st.session_state.get("running", False)

# ----------- Формулы запросов -------------
expr_rps = 'rate(mc_requests_total[1m])'

expr_rep = (
    f'kube_deployment_status_replicas{{deployment="{deploy}",namespace="{ns}"}}'
)

# ----------- Главная логика -------------
if running:
    # История значений сохраняется в сессии (максимум 1000 точек)
    history: pd.DataFrame = st.session_state.get(
        "history", pd.DataFrame(columns=["rps", "replicas"])
    )

    # Получаем очередные метрики
    rps = query(expr_rps) or 0
    rep = query(expr_rep) or 0
    now = pd.Timestamp.utcnow()

    history.loc[now, ["rps", "replicas"]] = [rps, rep]
    history = history.tail(1000).copy()
    st.session_state["history"] = history

    # Метрики текущего состояния
    col1, col2 = st.columns(2)
    with col1:
        st.metric("Текущий RPS", f"{rps:.0f}")
    with col2:
        st.metric("Количество реплик", f"{int(rep)}")

    # Графики
    line_col1, line_col2 = st.columns(2)
    with line_col1:
        st.subheader("RPS")
        st.line_chart(history["rps"])
    with line_col2:
        st.subheader("Replicas")
        st.line_chart(history["replicas"])

    # Пауза и перезапуск цикла
    time.sleep(refresh)
    st.rerun()
else:
    st.info("Нажмите **Запустить**, чтобы начать мониторинг.")
