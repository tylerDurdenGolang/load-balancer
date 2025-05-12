# ⚖️ Интеллектуальный оркестратор в Kubernetes

Проект представляет собой реализацию интеллектуального балансировщика нагрузки и компонента динамического масштабирования для Kubernetes-кластера. Система адаптируется к изменениям пользовательской нагрузки и принимает решения на основе метрик и прогноза.

---

## 🧠 Основные компоненты

- **Load Balancer (Go)** — балансировщик с поддержкой стратегий:
    - Round Robin
    - Weighted Random
    - Exponential
    - Adaptive Fuzzy
    - Queuing
- **Orchestrator (Python)** — анализирует метрики (CPU, RPS), прогнозирует нагрузку и масштабирует реплики.
- **Prometheus** — сбор и агрегация метрик.
- **Kubernetes** — оркестрация подов на основе решений от оркестратора.

---

## 🏗️ Архитектура

```text
[ Клиенты ] 
     │
     ▼
[ Load Balancer (Go) ]
     │
     ▼
[ Backend-инстансы (симуляция) ]
     │
     ├──→ [ Метрики в Prometheus ]
     │
     └──→ [ Orchestrator (Python) ] → [Kubernetes API]
```

---

## 🚀 Быстрый старт

```bash
# Запуск балансировщика (Docker)
cd load-balancer
docker build -t load-balancer .
docker run -p 8080:8080 load-balancer

# Запуск оркестратора
cd orchestrator
docker build -t orchestrator .
docker run orchestrator

# Применение манифестов
kubectl apply -f load-balancer/k8s.yml
```

---

## 📊 Стратегии балансировки

Находятся в `internal/balancer/strategies/`:

- `round_robin.go`
- `weighted_random.go`
- `adaptive_fuzzy.go`
- `queuing.go`
- `exponential.go`

Каждая стратегия реализует интерфейс выбора следующего backend'а.

---

## 📈 Оркестратор (AutoScaler)

- Реализован на Python
- Прогнозирует нагрузку с помощью скользящего окна
- Оценивает "стоимость" масштабирования с учётом:
    - Латентности (нагрузка на одну реплику)
    - Стоимости (число реплик)
    - Скачков (`churn`)

---

## 📁 Структура проекта

```text
├── load-balancer/
│   ├── cmd/lb/main.go
│   ├── internal/balancer/strategies/
│   ├── config/config.yaml
│   ├── k8s.yml
│   └── Dockerfile
└── orchestrator/
   ├── main.py
   ├── predictor.py
   ├── autoscaler.py
   └── Dockerfile

```

---

## 📌 Зависимости

- Go 1.20+
- Python 3.10+
- Docker + Kubernetes (Minikube или K3s)
- Prometheus (установлен в кластер)

---

## 📎 Лицензия

MIT License
