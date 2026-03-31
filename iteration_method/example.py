"""
Пример использования метода простых итераций.

Задача: найти корень уравнения  x^3 - x - 2 = 0
Приведём к виду x = phi(x):
    x = (x + 2)^(1/3)   --  phi(x) = cbrt(x + 2)

Корень уравнения находится вблизи x ≈ 1.5214.
"""

import sys
import os

sys.path.insert(0, os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
from iteration_method.iteration import simple_iteration


def phi(x):
    """Итерационная функция: x = (x + 2)^(1/3)"""
    return (x + 2) ** (1 / 3)


if __name__ == "__main__":
    x0 = 1.0        # начальное приближение
    eps = 1e-6      # требуемая точность

    root = simple_iteration(phi, x0, epsilon=eps)

    print(f"Найденный корень: x ≈ {root:.8f}")
    print(f"Проверка: x^3 - x - 2 = {root**3 - root - 2:.2e}  (должно быть ≈ 0)")
