from django.shortcuts import render
from django.http import JsonResponse
from django.views.decorators.http import require_http_methods
from .latex_parser import LatexParser, MathExpression
from .slae_solver import SLAESolver
import json
import numpy as np

def index(request):
    """Главная страница"""
    return render(request, 'main/index.html')

def calculate(request):
    """Страница расчётов"""
    return render(request, 'main/index.html')

@require_http_methods(["POST"])
def parse_latex(request):
    """API endpoint для парсирования LaTeX"""
    try:
        data = json.loads(request.body)
        latex_expr = data.get('latex', '')
        
        result = LatexParser.parse(latex_expr)
        
        if result['valid']:
            request.session['current_expression'] = result['expression'].to_dict()
            response_data = {
                'original': result['original'],
                'parsed': result['parsed'],
                'valid': True,
                'error': None,
                'variables': result['expression'].variables
            }
        else:
            response_data = {
                'original': result['original'],
                'parsed': result['parsed'],
                'valid': False,
                'error': result['error']
            }
        
        return JsonResponse(response_data)
    
    except json.JSONDecodeError:
        return JsonResponse({
            'valid': False,
            'error': 'Невалидный JSON'
        }, status=400)
    except Exception as e:
        return JsonResponse({
            'valid': False,
            'error': str(e)
        }, status=500)

@require_http_methods(["POST"])
def solve_slae(request):
    """Решает СЛАУ методом Гаусса или простой итерации"""
    try:
        data = json.loads(request.body)
        
        # Парсим матрицу A
        A = np.array(data.get('matrix_A', []))
        b = np.array(data.get('vector_b', []))
        epsilon = float(data.get('epsilon', 0.01))
        method = data.get('method', 'gaussian')
        
        # Проверка размеров
        if A.shape[0] != len(b):
            return JsonResponse({
                'valid': False,
                'error': 'Размеры матрицы и вектора не совпадают'
            }, status=400)
        
        # Решаем систему
        if method == 'gaussian':
            result = SLAESolver.gaussian_method(A, b, epsilon)
        else:  # simple_iteration
            result = SLAESolver.simple_iteration_method(A, b, epsilon)
        
        if result['valid']:
            # Преобразуем в LaTeX
            solution_latex = SLAESolver.solution_to_latex(result['solution'])
            matrix_latex = SLAESolver.matrix_to_latex(A, b)
            
            response_data = {
                'valid': True,
                'solution': result['solution'].tolist(),
                'solution_latex': solution_latex,
                'matrix_latex': matrix_latex,
                'residual': float(result['residual']),
                'method': result['method'],
                'epsilon': result['epsilon'],
                'steps': result['steps'],
                'error': None
            }
            
            if 'num_iterations' in result:
                response_data['num_iterations'] = result['num_iterations']
        else:
            response_data = {
                'valid': False,
                'error': result['error'],
                'steps': result.get('steps', [])
            }
        
        return JsonResponse(response_data)
    
    except Exception as e:
        return JsonResponse({
            'valid': False,
            'error': f'Ошибка сервера: {str(e)}'
        }, status=500)

@require_http_methods(["GET"])
def get_latex_help(request):
    """Возвращает справку по поддерживаемым командам"""
    commands = LatexParser.get_supported_commands()
    return JsonResponse({'commands': commands})