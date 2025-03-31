from src.collector.prometheus import PrometheusClient
from src.analyzer.forecasting import ETS, AutoRegressor
from src.analyzer.optimizer import ScalingOptimizer
from src.orchestrator.k8s_manager import KubernetesManager
import yaml
import time

def main():
    # Инициализация компонентов
    prometheus = PrometheusClient()
    ets = ETS()
    ar = AutoRegressor()
    optimizer = ScalingOptimizer()
    k8s = KubernetesManager()
    
    with open("config/settings.yaml") as f:
        config = yaml.safe_load(f)
    
    while True:
        # Сбор данных
        metrics = prometheus.get_metrics()
        cpu_data = list(metrics['cpu_usage'].values())
        
        # Прогнозирование
        ets.fit(cpu_data[-10:])  # Используем последние 10 точек
        ets_forecast = ets.predict(config['forecasting']['horizon'])
        
        ar.fit(cpu_data)
        ar_forecast = ar.predict(cpu_data[-3:], config['forecasting']['horizon'])
        
        # Ансамбль
        combined = [(e + a)/2 for e, a in zip(ets_forecast, ar_forecast)]
        
        # Оптимизация
        current_nodes = k8s.get_current_nodes()
        desired = optimizer.calculate(
            current_nodes, 
            combined,
            config['scaling']
        )
        
        # Принятие решения
        if desired != current_nodes:
            k8s.scale(desired)
        
        time.sleep(config['scaling']['check_interval'])

if __name__ == "__main__":
    main()