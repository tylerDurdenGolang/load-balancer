
from ..autoscaler.scaler import AutoScaler, ScalingConfig

def test_scaler_increases_when_overloaded():
    sc = AutoScaler(ScalingConfig(target_rps_per_pod=100, max_replicas=5))
    # predicted load 400 RPS with 1 replica should upscale
    desired = sc.optimise(predicted_rps=400, current_replicas=1)
    assert desired > 1
