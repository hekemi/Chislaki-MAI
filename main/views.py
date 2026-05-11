from django.shortcuts import render
from django.http import JsonResponse
from django.views.decorators.http import require_http_methods
from .latex_parser import LatexParser
import json

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
        
        return JsonResponse(result)
    
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

@require_http_methods(["GET"])
def get_latex_help(request):
    """Возвращает справку по поддерживаемым командам"""
    commands = LatexParser.get_supported_commands()
    return JsonResponse({'commands': commands})