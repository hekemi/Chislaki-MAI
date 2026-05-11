from django.urls import path
from . import views

urlpatterns = [
    path('', views.index, name='index'),
    path('api/parse-latex/', views.parse_latex, name='parse_latex'),
    path('api/solve-slae/', views.solve_slae, name='solve_slae'),
    path('api/latex-help/', views.get_latex_help, name='latex_help'),
]