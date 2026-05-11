import re
from typing import Dict, List

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
    def parse(latex_expr: str) -> Dict[str, any]:
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
            
            return {
                'original': latex_expr,
                'parsed': parsed,
                'valid': True,
                'error': None
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