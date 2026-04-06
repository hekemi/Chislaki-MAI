from django.shortcuts import render
from django.http import HttpResponse

def index(request):
    return render(request, 'main/index.html')

def calculate(request):
    if request.method == 'POST':
        # Получаем данные из полей формы (атрибуты 'name' из HTML)
        method = request.POST.get('method')
        function_text = request.POST.get('function')
        
        # Здесь будет ваша логика вычислений (бисекция, Ньютон и т.д.)
        count = eval(function_text)
        result = f"Вы выбрали метод {method} для функции {function_text}. \n Ответ: {count}"
        
        return HttpResponse(result) # Или рендерите новый HTML с результатом
    
    # Если зашли на страницу просто так (GET), возвращаем обратно на главную
    return render(request, 'main/index.html')