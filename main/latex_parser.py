import re
import numpy as np
from typing import Dict, List, Any
from scipy import integrate

class MathExpression:
    """Класс для представления математического выражения в удобном виде"""
    
    def __init__(self, latex_expr: str, parsed_expr: str):
        """
        Args:
            latex_expr: Оригинальное LaTeX выражение
            parsed_expr: Распаршеное выражение
        """
        self.latex = latex_expr
        self.parsed = parsed_expr
        self.variables = self._extract_variables()
    
    def _extract_variables(self) -> List[str]:
        """Извлекает переменные из выражения"""
        # Ищем одиночные буквы, которые не являются функциями
        vars_found = set(re.findall(r'(?<![a-z])([a-z])(?![a-z])', self.parsed))
        # Исключаем буквы, которые часть функций (sin, cos и т.д.)
        excluded = {'s', 'i', 'n', 'c', 'o', 't', 'a', 'l', 'g', 'e', 'x', 'p', 'd'}
        return sorted(list(vars_found - excluded))
    
    def evaluate(self, **kwargs) -> float:
        """
        Вычисляет значение выражения для заданных переменных
        
        Args:
            **kwargs: значения переменных (x=1, y=2 и т.д.)
            
        Returns:
            Вычисленное значение
        """
        expr = self.parsed
        
        # Подставляем значения переменных
        for var, value in kwargs.items():
            expr = re.sub(rf'\b{var}\b', f'({value})', expr)
        
        # Заменяем математические функции на numpy
        expr = expr.replace('sin', 'np.sin')
        expr = expr.replace('cos', 'np.cos')
        expr = expr.replace('tan', 'np.tan')
        expr = expr.replace('log', 'np.log10')
        expr = expr.replace('ln', 'np.log')
        expr = expr.replace('exp', 'np.exp')
        expr = expr.replace('√', 'np.sqrt')
        
        try:
            return float(eval(expr))
        except Exception as e:
            raise ValueError(f"Ошибка при вычислении: {str(e)}")
    
    def vectorize(self, var: str, values: np.ndarray) -> np.ndarray:
        """
        Вычисляет функцию для массива значений переменной
        
        Args:
            var: имя переменной
            values: массив значений
            
        Returns:
            Массив результатов
        """
        results = []
        for val in values:
            try:
                result = self.evaluate(**{var: val})
                results.append(result)
            except:
                results.append(np.nan)
        return np.array(results)
    
    def derivative_numeric(self, var: str, point: float, h: float = 1e-5) -> float:
        """
        Численная производная в точке
        
        Args:
            var: переменная
            point: точка, в которой считать производную
            h: шаг
            
        Returns:
            Значение производной
        """
        f_plus = self.evaluate(**{var: point + h})
        f_minus = self.evaluate(**{var: point - h})
        return (f_plus - f_minus) / (2 * h)
    
    def integral_numeric(self, var: str, a: float, b: float, n: int = 1000) -> float:
        """
        Численное интегрирование методом трапеций
        
        Args:
            var: переменная интегрирования
            a: нижний предел
            b: верхний предел
            n: количество разбиений
            
        Returns:
            Значение интеграла
        """
        x = np.linspace(a, b, n)
        y = self.vectorize(var, x)
        # Используем numpy.trapezoid вместо trapz
        return float(np.trapezoid(y, x))
    
    def find_root(self, var: str, a: float, b: float, tol: float = 1e-6) -> float:
        """
        Поиск корня методом бисекции
        
        Args:
            var: переменная
            a, b: границы интервала
            tol: точность
            
        Returns:
            Найденный корень
        """
        while abs(b - a) > tol:
            mid = (a + b) / 2
            f_a = self.evaluate(**{var: a})
            f_mid = self.evaluate(**{var: mid})
            
            if f_a * f_mid < 0:
                b = mid
            else:
                a = mid
        
        return (a + b) / 2
    
    def __str__(self) -> str:
        return f"LaTeX: {self.latex}\nПаршено: {self.parsed}\nПеременные: {', '.join(self.variables)}"
    
    def to_dict(self) -> Dict[str, Any]:
        """Преобразует в словарь для передачи в JSON"""
        return {
            'latex': self.latex,
            'parsed': self.parsed,
            'variables': self.variables
        }


class LatexParser:
    """Парсер для преобразования LaTeX выражений в читаемый математический формат"""
    
    # Словарь замен LaTeX команд
    LATEX_REPLACEMENTS = [
        (r'\\frac\{([^}]+)\}\{([^}]+)\}', r'(\1)/(\2)'),
        (r'\\sqrt\{([^}]+)\}', r'√(\1)'),
        (r'\\sin\b', 'sin'),
        (r'\\cos\b', 'cos'),
        (r'\\tan\b', 'tan'),
        (r'\\log\b', 'log'),
        (r'\\ln\b', 'ln'),
        (r'\\exp\b', 'exp'),
        (r'\\pi\b', 'π'),
        (r'\\infty\b', '∞'),
        (r'\\alpha\b', 'α'),
        (r'\\beta\b', 'β'),
        (r'\\gamma\b', 'γ'),
        (r'\\delta\b', 'δ'),
    ]
    
    @staticmethod
    def parse(latex_expr: str) -> Dict[str, Any]:
        """
        Парсит LaTeX выражение и возвращает структурированные данные
        
        Args:
            latex_expr: LaTeX строка
            
        Returns:
            Dict с ключами:
            - 'original': оригинальное выражение
            - 'parsed': распаршеное выражение
            - 'valid': валидно ли выражение
            - 'error': ошибка парсирования (если есть)
            - 'expression': объект MathExpression (если валидно)
        """
        try:
            if not latex_expr.strip():
                return {
                    'original': latex_expr,
                    'parsed': '',
                    'valid': False,
                    'error': 'Поле не может быть пустым'
                }
            
            parsed = latex_expr
            
            # Проверка на сбалансированные скобки
            if not LatexParser._check_brackets(parsed):
                return {
                    'original': latex_expr,
                    'parsed': '',
                    'valid': False,
                    'error': 'Несбалансированные скобки'
                }
            
            # Применяем замены
            for latex_pattern, replacement in LatexParser.LATEX_REPLACEMENTS:
                parsed = re.sub(latex_pattern, replacement, parsed)
            
            # Удаляем лишние пробелы
            parsed = re.sub(r'\s+', '', parsed)
            
            # Создаём объект MathExpression
            expr_obj = MathExpression(latex_expr, parsed)
            
            return {
                'original': latex_expr,
                'parsed': parsed,
                'valid': True,
                'error': None,
                'expression': expr_obj
            }
        
        except Exception as e:
            return {
                'original': latex_expr,
                'parsed': '',
                'valid': False,
                'error': f'Ошибка парсирования: {str(e)}'
            }
    
    @staticmethod
    def _check_brackets(expr: str) -> bool:
        """Проверяет сбалансированность скобок"""
        brackets = {'(': ')', '{': '}', '[': ']'}
        stack = []
        
        for char in expr:
            if char in brackets:
                stack.append(char)
            elif char in brackets.values():
                if not stack or brackets[stack.pop()] != char:
                    return False
        
        return len(stack) == 0
    
    @staticmethod
    def get_supported_commands() -> List[str]:
        """Возвращает список поддерживаемых LaTeX команд"""
        return [pattern for pattern, _ in LatexParser.LATEX_REPLACEMENTS]