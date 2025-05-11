
import numpy as np
from ..autoscaler.predictor import ARPredictor

def test_predictor_converges():
    p = ARPredictor(p=2)
    # generate AR(2) data: y_t = 0.6 y_{t-1} + 0.2 y_{t-2} + noise
    y1, y2 = 1.0, 1.0
    for _ in range(100):
        y = 0.6 * y1 + 0.2 * y2 + np.random.normal(0, 0.01)
        p.update(y)
        y2, y1 = y1, y
    assert np.allclose(p.phi, [0.6, 0.2], atol=0.1)
