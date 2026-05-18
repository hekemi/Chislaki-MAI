"""
Метод дихотомии (бисекции) для нахождения корня функции на отрезке [a, b].

Условие применимости: функция f должна быть непрерывна на [a, b],
и f(a) * f(b) < 0 (функция меняет знак на концах отрезка).
"""


def dichotomy(f, a: float, b: float, eps: float = 1e-6, max_iter: int = 1000) -> float:
    """
    Находит корень уравнения f(x) = 0 методом дихотомии (бисекции).

    Параметры
    ----------
    f       : callable — исследуемая функция одной переменной
    a, b    : float    — границы отрезка (f(a) * f(b) < 0)
    eps     : float    — допустимая погрешность (по умолчанию 1e-6)
    max_iter: int      — максимальное число итераций (по умолчанию 1000)

    Возвращает
    ----------
    float — приближённое значение корня

    Исключения
    ----------
    ValueError — если f(a) * f(b) >= 0 (условие смены знака не выполнено)
    """
    fa = f(a)
    fb = f(b)

    # Если один из концов уже является корнем
    if fa == 0:
        return a
    if fb == 0:
        return b

    if fa * fb >= 0:
        raise ValueError(
            f"Условие f(a)*f(b)<0 не выполнено: f({a})={fa}, f({b})={fb}"
        )

    for _ in range(max_iter):
        mid = a + (b - a) / 2.0
        fmid = f(mid)

        if fmid == 0 or (b - a) / 2.0 < eps:
            return mid

        if fa * fmid < 0:
            b = mid
            fb = fmid
        else:
            a = mid
            fa = fmid

    return a + (b - a) / 2.0


if __name__ == "__main__":
    import math

    # Пример 1: корень уравнения x^3 - x - 2 = 0 на [1, 2]
    def f1(x):
        return x**3 - x - 2

    root1 = dichotomy(f1, 1, 2)
    print(f"Пример 1: корень x^3 - x - 2 = 0  →  x ≈ {root1:.8f}")
    print(f"  Проверка: f({root1:.8f}) = {f1(root1):.2e}")

    # Пример 2: корень уравнения cos(x) - x = 0 на [0, π/2]
    def f2(x):
        return math.cos(x) - x

    root2 = dichotomy(f2, 0, math.pi / 2)
    print(f"\nПример 2: корень cos(x) - x = 0  →  x ≈ {root2:.8f}")
    print(f"  Проверка: f({root2:.8f}) = {f2(root2):.2e}")
