document.addEventListener('DOMContentLoaded', () => {
    // DOM Elements
    const todoList = document.getElementById('todo-list');
    const newTodoInput = document.getElementById('new-todo');
    const addTodoButton = document.getElementById('add-todo');

    // Load todos when the page loads
    loadTodos();

    // Add event listener for adding a new todo
    addTodoButton.addEventListener('click', addTodo);
    newTodoInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            addTodo();
        }
    });

    // Function to load todos from the API
    async function loadTodos() {
        try {
            const response = await fetch('/api/todos');
            if (!response.ok) {
                throw new Error('Failed to load todos');
            }
            
            const todos = await response.json();
            renderTodos(todos);
        } catch (error) {
            console.error('Error loading todos:', error);
            showError('Failed to load todos. Please try again later.');
        }
    }

    // Function to add a new todo
    async function addTodo() {
        const title = newTodoInput.value.trim();
        if (!title) return;

        try {
            const response = await fetch('/api/todos', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ title })
            });

            if (!response.ok) {
                throw new Error('Failed to add todo');
            }

            const newTodo = await response.json();
            renderTodo(newTodo);
            newTodoInput.value = '';
        } catch (error) {
            console.error('Error adding todo:', error);
            showError('Failed to add todo. Please try again.');
        }
    }

    // Function to update a todo's completion status
    async function updateTodoStatus(id, completed) {
        try {
            const response = await fetch(`/api/todos/${id}`, {
                method: 'PATCH',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({ completed })
            });

            if (!response.ok) {
                throw new Error('Failed to update todo');
            }
        } catch (error) {
            console.error('Error updating todo:', error);
            showError('Failed to update todo. Please try again.');
            // Reset the checkbox to its previous state
            loadTodos();
        }
    }

    // Function to delete a todo
    async function deleteTodo(id) {
        try {
            const response = await fetch(`/api/todos/${id}`, {
                method: 'DELETE'
            });

            if (!response.ok) {
                throw new Error('Failed to delete todo');
            }

            // Remove the todo from the DOM
            const todoElement = document.getElementById(`todo-${id}`);
            if (todoElement) {
                todoElement.remove();
            }
        } catch (error) {
            console.error('Error deleting todo:', error);
            showError('Failed to delete todo. Please try again.');
        }
    }

    // Function to render todos in the DOM
    function renderTodos(todos) {
        todoList.innerHTML = '';
        todos.forEach(renderTodo);
    }

    // Function to render a single todo in the DOM
    function renderTodo(todo) {
        const todoItem = document.createElement('li');
        todoItem.id = `todo-${todo.id}`;
        todoItem.className = `todo-item ${todo.completed ? 'completed' : ''}`;

        // Format the date
        const created = new Date(todo.created_at);
        const formattedDate = created.toLocaleDateString('en-US', {
            month: 'short',
            day: 'numeric',
            year: 'numeric'
        });

        todoItem.innerHTML = `
            <input type="checkbox" class="todo-checkbox" ${todo.completed ? 'checked' : ''}>
            <span class="todo-text">${escapeHtml(todo.title)}</span>
            <span class="todo-date">${formattedDate}</span>
            <button class="todo-delete">Delete</button>
        `;

        // Add event listener for checkbox
        const checkbox = todoItem.querySelector('.todo-checkbox');
        checkbox.addEventListener('change', () => {
            updateTodoStatus(todo.id, checkbox.checked);
            todoItem.classList.toggle('completed', checkbox.checked);
        });

        // Add event listener for delete button
        const deleteButton = todoItem.querySelector('.todo-delete');
        deleteButton.addEventListener('click', () => {
            deleteTodo(todo.id);
        });

        todoList.appendChild(todoItem);
    }

    // Function to show error message
    function showError(message) {
        alert(message);
    }

    // Helper function to escape HTML to prevent XSS
    function escapeHtml(unsafe) {
        return unsafe
            .replace(/&/g, "&amp;")
            .replace(/</g, "&lt;")
            .replace(/>/g, "&gt;")
            .replace(/"/g, "&quot;")
            .replace(/'/g, "&#039;");
    }
}); 