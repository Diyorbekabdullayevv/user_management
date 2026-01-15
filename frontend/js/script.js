// Page Navigation
const navLinks = document.querySelectorAll('.nav-link');
const pages = document.querySelectorAll('.page');

function navigateTo(pageName) {
    // Remove active class from all
    navLinks.forEach(l => l.classList.remove('active'));
    pages.forEach(p => p.classList.remove('active'));
    
    // Add active class to corresponding page and nav link
    const navLink = document.querySelector(`[data-page="${pageName}"]`);
    if (navLink) {
        navLink.classList.add('active');
    }
    document.getElementById(pageName).classList.add('active');
    
    // Scroll to top
    window.scrollTo(0, 0);
}

navLinks.forEach(link => {
    link.addEventListener('click', (e) => {
        e.preventDefault();
        const pageName = link.getAttribute('data-page');
        navigateTo(pageName);
    });
});

// Set home as active by default
document.querySelector('[data-page="home"]').classList.add('active');

// Tab switching functionality (for Users page)
const tabButtons = document.querySelectorAll('.tab-button');
const tabContents = document.querySelectorAll('.tab-content');

tabButtons.forEach(button => {
    button.addEventListener('click', () => {
        const tabName = button.getAttribute('data-tab');
        
        // Remove active class from all
        tabButtons.forEach(btn => btn.classList.remove('active'));
        tabContents.forEach(content => content.classList.remove('active'));
        
        // Add active class to clicked tab
        button.classList.add('active');
        document.getElementById(tabName).classList.add('active');
    });
});

// Set first tab as active by default
if (tabButtons.length > 0) {
    tabButtons[0].classList.add('active');
    tabContents[0].classList.add('active');
}

// User Registration Form
const registerForm = document.getElementById('userForm');
const registerMessage = document.getElementById('registerMessage');
const registerLoading = document.getElementById('registerLoading');

registerForm.addEventListener('submit', async (e) => {
    e.preventDefault();

    const name = document.getElementById('name').value.trim();
    const surname = document.getElementById('surname').value.trim();
    const age = parseInt(document.getElementById('age').value);

    if (!name || !surname || !age) {
        showMessage(registerMessage, 'Please fill in all fields', 'error');
        return;
    }

    if (age < 0 || age > 150) {
        showMessage(registerMessage, 'Please enter a valid age', 'error');
        return;
    }

    registerLoading.style.display = 'block';
    registerMessage.style.display = 'none';

    try {
        const response = await fetch('/api/users', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                name: name,
                surname: surname,
                age: age
            })
        });

        registerLoading.style.display = 'none';

        if (response.ok) {
            const data = await response.json();
            showMessage(registerMessage, `✅ User registered successfully!`, 'success');
            registerForm.reset();
            
            // Switch to filter tab after successful registration
            setTimeout(() => {
                document.querySelector('[data-tab="filter"]').click();
            }, 1500);
        } else {
            const error = await response.json();
            showMessage(registerMessage, `❌ Error: ${error.error || 'Failed to register user'}`, 'error');
        }
    } catch (error) {
        registerLoading.style.display = 'none';
        showMessage(registerMessage, `❌ Error: ${error.message}`, 'error');
    }
});

// Filter Users Form
const filterForm = document.getElementById('filterForm');
const filterMessage = document.getElementById('filterMessage');
const filterLoading = document.getElementById('filterLoading');
const usersContainer = document.getElementById('usersContainer');

filterForm.addEventListener('submit', async (e) => {
    e.preventDefault();

    const minAge = parseInt(document.getElementById('minAge').value);
    const maxAge = parseInt(document.getElementById('maxAge').value);

    if (minAge < 0 || maxAge > 150) {
        showMessage(filterMessage, 'Age must be between 0 and 150', 'error');
        return;
    }

    if (minAge > maxAge) {
        showMessage(filterMessage, 'Minimum age cannot be greater than maximum age', 'error');
        return;
    }

    filterLoading.style.display = 'block';
    filterMessage.style.display = 'none';
    usersContainer.innerHTML = '';

    try {
        const response = await fetch(`/api/users/filter?min_age=${minAge}&max_age=${maxAge}`);

        filterLoading.style.display = 'none';

        if (response.ok) {
            const data = await response.json();
            
            if (data.users && data.users.length > 0) {
                showMessage(filterMessage, `✅ Found ${data.count} user(s)`, 'success');
                displayUsers(data.users);
            } else {
                usersContainer.innerHTML = `
                    <div class="no-results">
                        <div style="font-size: 48px; margin-bottom: 15px;">🔍</div>
                        <p>No users found in the age range ${minAge}-${maxAge}</p>
                    </div>
                `;
                showMessage(filterMessage, 'No users found in this age range', 'error');
            }
        } else {
            const error = await response.json();
            showMessage(filterMessage, `❌ Error: ${error.error || 'Failed to fetch users'}`, 'error');
        }
    } catch (error) {
        filterLoading.style.display = 'none';
        showMessage(filterMessage, `❌ Error: ${error.message}`, 'error');
    }
});

function displayUsers(users) {
    usersContainer.innerHTML = '<div class="users-grid">';
    
    users.forEach(user => {
        const userCard = document.createElement('div');
        userCard.className = 'user-card';
        userCard.innerHTML = `
            <h3>${user.name} ${user.surname}</h3>
            <div class="user-info">
                <div class="user-info-item">
                    <span class="user-info-label">Age</span>
                    <span class="user-info-value">${user.age} years</span>
                </div>
                <div class="user-info-item">
                    <span class="user-info-label">ID</span>
                    <span class="user-info-value" style="font-size: 12px; word-break: break-all;">${user.id}</span>
                </div>
            </div>
        `;
        usersContainer.querySelector('.users-grid').appendChild(userCard);
    });

    usersContainer.innerHTML += '</div>';
}

function showMessage(element, text, type) {
    element.textContent = text;
    element.className = `message ${type}`;
    element.style.display = 'block';
}
