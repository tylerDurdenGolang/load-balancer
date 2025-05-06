
from pathlib import Path
import yaml
from pydantic import BaseModel
from typing import Optional

class Settings(BaseModel):
    prometheus_url: str = "http://prometheus:9090"
    deployment: str
    namespace: str = "default"
    interval: int = 30      # seconds
    scaling: dict = {}

def load(path: str | Path = "config.yaml") -> Settings:
    data = yaml.safe_load(Path(path).read_text())
    return Settings(**data)
