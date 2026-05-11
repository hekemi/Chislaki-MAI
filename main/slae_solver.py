import numpy as np
from typing import Dict, Tuple, List, Any

class SLAESolver:
    """Решатель систем линейных алгебраических уравнений"""
    
    @staticmethod
    def gaussian_method(A: np.ndarray, b: np.ndarray, epsilon: float = 0.01) -> Dict[str, Any]:
        """
        Решает СЛАУ методом Гаусса с частичным выбором главного элемента
        
        Args:
            A: матрица коэффициентов (n x n)
            b: вектор правых частей (n,)
            epsilon: точность (для информации)
            
        Returns:
            Dict с решением и информацией о процессе
        """
        try:
            n = len(b)
            A = A.astype(float)
            b = b.astype(float)
            
            # Прямой ход (приведение к треугольному виду)
            steps = []
            
            for i in range(n):
                # Поиск максимального элемента в столбце (частичный выбор)
                max_row = i
                for k in range(i + 1, n):
                    if abs(A[k, i]) > abs(A[max_row, i]):
                        max_row = k
                
                # Перестановка строк
                if max_row != i:
                    A[[i, max_row]] = A[[max_row, i]]
                    b[i], b[max_row] = b[max_row], b[i]
                    steps.append(f"Перестановка строк {i+1} и {max_row+1}")
                
                # Проверка на вырожденность
                if abs(A[i, i]) < 1e-10:
                    return {
                        'valid': False,
                        'error': 'Матрица вырождена (det = 0)',
                        'solution': None,
                        'steps': steps
                    }
                
                # Исключение
                for k in range(i + 1, n):
                    factor = A[k, i] / A[i, i]
                    A[k, i:] -= factor * A[i, i:]
                    b[k] -= factor * b[i]
            
            steps.append("Прямой ход завершён")
            
            # Обратный ход
            x = np.zeros(n)
            for i in range(n - 1, -1, -1):
                x[i] = b[i]
                for j in range(i + 1, n):
                    x[i] -= A[i, j] * x[j]
                x[i] /= A[i, i]
            
            steps.append("Обратный ход завершён")
            
            # Проверка решения
            residual = np.linalg.norm(np.dot(np.array(A), x) - b)
            
            return {
                'valid': True,
                'solution': x,
                'residual': residual,
                'steps': steps,
                'method': 'Метод Гаусса',
                'epsilon': epsilon
            }
        
        except Exception as e:
            return {
                'valid': False,
                'error': f'Ошибка при решении методом Гаусса: {str(e)}',
                'solution': None,
                'steps': steps if 'steps' in locals() else []
            }
    
    @staticmethod
    def simple_iteration_method(A: np.ndarray, b: np.ndarray, 
                               epsilon: float = 0.01, max_iterations: int = 1000) -> Dict[str, Any]:
        """
        Решает СЛАУ методом простой итерации
        
        Args:
            A: матрица коэффициентов (n x n)
            b: вектор правых частей (n,)
            epsilon: точность сходимости
            max_iterations: максимум итераций
            
        Returns:
            Dict с решением и информацией о процессе
        """
        try:
            n = len(b)
            A = A.astype(float)
            b = b.astype(float)
            
            # Приведение к виду x = Bx + c
            # x = (E - A)x + b  или x = x - Ax + b
            # Более стабильно: используем диагональное доминирование
            
            steps = []
            
            # Проверка условия диагонального доминирования
            for i in range(n):
                diag = abs(A[i, i])
                sum_other = sum(abs(A[i, j]) for j in range(n) if j != i)
                if diag < sum_other:
                    steps.append(f"Внимание: нет диагонального доминирования в строке {i+1}")
            
            # Матрица B = E - A, вектор c = b
            B = np.eye(n) - A
            c = b
            
            # Начальное приближение
            x_prev = np.zeros(n)
            x_curr = c.copy()
            
            iterations = []
            for iteration in range(max_iterations):
                x_prev = x_curr.copy()
                x_curr = np.dot(B, x_prev) + c
                
                # Норма разности
                diff = np.linalg.norm(x_curr - x_prev)
                iterations.append({
                    'iteration': iteration + 1,
                    'norm': diff,
                    'solution': x_curr.copy()
                })
                
                steps.append(f"Итерация {iteration + 1}: ||x_k - x_{{k-1}}|| = {diff:.6e}")
                
                if diff < epsilon:
                    steps.append(f"Сходимость достигнута за {iteration + 1} итераций")
                    break
            
            if diff >= epsilon:
                steps.append(f"Внимание: не достигнута точность {epsilon} за {max_iterations} итераций")
            
            # Проверка решения
            residual = np.linalg.norm(np.dot(A, x_curr) - b)
            
            return {
                'valid': True,
                'solution': x_curr,
                'residual': residual,
                'steps': steps,
                'iterations': iterations,
                'method': 'Метод простой итерации',
                'epsilon': epsilon,
                'num_iterations': len(iterations)
            }
        
        except Exception as e:
            return {
                'valid': False,
                'error': f'Ошибка при решении методом простой итерации: {str(e)}',
                'solution': None,
                'steps': steps if 'steps' in locals() else []
            }
    
    @staticmethod
    def matrix_to_latex(A: np.ndarray, b: np.ndarray = None) -> str:
        """Преобразует матрицу в LaTeX формат"""
        n, m = A.shape
        latex = r"\begin{pmatrix}"
        
        for i in range(n):
            row = " & ".join(f"{A[i, j]:.4g}" for j in range(m))
            if b is not None:
                row += f" & {b[i]:.4g}"
            latex += row
            if i < n - 1:
                latex += r" \\ "
        
        latex += r"\end{pmatrix}"
        return latex
    
    @staticmethod
    def solution_to_latex(x: np.ndarray, var_names: List[str] = None) -> str:
        """Преобразует решение в LaTeX формат"""
        if var_names is None:
            var_names = [f"x_{i+1}" for i in range(len(x))]
        
        latex = r"\begin{cases}"
        for i, (name, val) in enumerate(zip(var_names, x)):
            latex += f"{name} = {val:.6g}"
            if i < len(x) - 1:
                latex += r" \\ "
        latex += r"\end{cases}"
        
        return latex