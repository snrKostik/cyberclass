from django.shortcuts import render, redirect, get_object_or_404
from .models import Note
from django.contrib.auth import authenticate, login, logout
from django.contrib.auth.decorators import login_required
from .forms import RegisterForm
from django.contrib import messages

# ----------------------
# Авторизация
# ----------------------
def register_view(request):
    if request.method == "POST":
        form = RegisterForm(request.POST)
        if form.is_valid():
            user = form.save()
            login(request, user)
            messages.success(request, "Регистрация успешна!")
            return redirect("note_list")
    else:
        form = RegisterForm()
    return render(request, "notes/register.html", {"form": form})

def login_view(request):
    if request.method == "POST":
        username = request.POST.get("username")
        password = request.POST.get("password")
        user = authenticate(request, username=username, password=password)
        if user:
            login(request, user)
            return redirect("note_list")
        else:
            messages.error(request, "Неверный логин или пароль")
    return render(request, "notes/login.html")

def logout_view(request):
    logout(request)
    return redirect("login")

# ----------------------
# Заметки
# ----------------------
@login_required
def note_list(request):
    query = request.GET.get('q', '')
    if query:
        notes = Note.objects.filter(user=request.user, title__icontains=query)
    else:
        notes = Note.objects.filter(user=request.user)
    return render(request, 'notes/note_list.html', {'notes': notes})

@login_required
def note_create(request):
    if request.method == "POST":
        title = request.POST['title']
        content = request.POST['content']
        Note.objects.create(user=request.user, title=title, content=content)
        return redirect('note_list')
    return render(request, 'notes/note_form.html')

@login_required
def note_detail(request, pk):
    note = get_object_or_404(Note, pk=pk, user=request.user)
    return render(request, 'notes/note_detail.html', {'note': note})

@login_required
def note_edit(request, pk):
    note = get_object_or_404(Note, pk=pk, user=request.user)
    if request.method == "POST":
        note.title = request.POST['title']
        note.content = request.POST['content']
        note.save()
        return redirect('note_detail', pk=pk)
    return render(request, 'notes/note_form.html', {'note': note})

@login_required
def note_delete(request, pk):
    note = get_object_or_404(Note, pk=pk, user=request.user)
    note.delete()
    return redirect('note_list')