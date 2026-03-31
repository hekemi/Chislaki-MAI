"""
Метод простых итераций (Simple Iteration Method)

Решает уравнение вида f(x) = 0, предварительно приведя его к виду x = phi(x).
Итерационная формула: x_{n+1} = phi(x_n)
Итерации продолжаются до тех пор, пока |x_{n+1} - x_n| < epsilon.
"""


def simple_iteration(phi, x0, epsilon=1e-6, max_iter=10000):
    """
    Метод простых итераций для нахождения корня уравнения x = phi(x).

    Параметры:
        phi      (callable) -- итерационная функция phi(x), полученная из f(x) = 0
                               приведением к виду x = phi(x)
        x0       (float)    -- начальное приближение
        epsilon  (float)    -- точность (критерий остановки |x_new - x_old| < epsilon)
        max_iter (int)      -- максимальное число итераций

    Возвращает:
        float -- найденное значение корня

    Исключения:
        ValueError -- если метод не сошёлся за max_iter итераций
    """
    x_old = x0
    for iteration in range(1, max_iter + 1):
        x_new = phi(x_old)
        if abs(x_new - x_old) < epsilon:
            return x_new
        x_old = x_new

    raise ValueError(
        f"Метод не сошёлся за {max_iter} итераций. "
        "Проверьте начальное приближение или условие сходимости |phi'(x)| < 1."
    )
