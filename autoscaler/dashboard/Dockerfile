# Dockerfile.dashboard
FROM python:3.11-slim

# ⬇︎ инструменты для сборки колёс и некоторые runtime-библиотеки
RUN apt-get update && apt-get install -y --no-install-recommends \
        build-essential gcc g++ make \
        libffi-dev libbz2-dev liblzma-dev libjpeg-dev zlib1g-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY requirements.txt .
RUN pip install --upgrade pip \
 && pip install --no-cache-dir -r requirements.txt

COPY . .

EXPOSE 8501
CMD ["streamlit", "run", "app.py","--server.port", "8501", "--server.address", "0.0.0.0"]
