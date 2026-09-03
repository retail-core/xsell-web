function tapFeedback(intensity = 10) {
    if (navigator.vibrate) navigator.vibrate(intensity);
}

function openPageSearch() {
    document.getElementById('topbar-default').classList.add('hidden');
    document.getElementById('topbar-search').classList.remove('hidden');
    document.getElementById('topbar-search').classList.add('flex');
    document.getElementById('page-search-input').focus();
}

function closePageSearch() {
    document.getElementById('topbar-search').classList.add('hidden');
    document.getElementById('topbar-search').classList.remove('flex');
    document.getElementById('topbar-default').classList.remove('hidden');
    document.getElementById('page-search-input').value = '';
    onPageSearch();
}

function onPageSearch() {
    const query = document.getElementById('page-search-input').value.trim().toLowerCase();
    // Delegate to whichever page's filter function is active
    if (typeof applyFilters === 'function') {
        applyFilters(query);
    }
}

let confirmCallback = null;

function showConfirmDialog(title, message, onConfirm, actionLabel = 'Delete') {
    document.getElementById('confirm-title').textContent = title;
    document.getElementById('confirm-message').textContent = message;
    document.getElementById('confirm-action-btn').textContent = actionLabel;

    confirmCallback = onConfirm;

    document.getElementById('confirm-backdrop').classList.remove('hidden');
    document.getElementById('confirm-dialog').classList.remove('hidden');
}

function closeConfirmDialog() {
    document.getElementById('confirm-backdrop').classList.add('hidden');
    document.getElementById('confirm-dialog').classList.add('hidden');
    confirmCallback = null;
}

document.getElementById('confirm-action-btn')?.addEventListener('click', () => {
    if (confirmCallback) confirmCallback();
    closeConfirmDialog();
});

function showFieldError(name, message) {
    const errorEl = document.getElementById('error-' + name);
    const inputEl = document.getElementById('field-' + name);
    if (errorEl) {
        errorEl.textContent = message;
        errorEl.classList.remove('hidden');
    }
    if (inputEl) {
        inputEl.classList.add('border-danger');
        inputEl.classList.remove('border-border');
    }
}

function clearFieldError(name) {
    const errorEl = document.getElementById('error-' + name);
    const inputEl = document.getElementById('field-' + name);
    if (errorEl) errorEl.classList.add('hidden');
    if (inputEl) {
        inputEl.classList.remove('border-danger');
        inputEl.classList.add('border-border');
    }
}

function clearAllFieldErrors(formId) {
    document.querySelectorAll(`#${formId} [id^="error-"]`).forEach(el => el.classList.add('hidden'));
    document.querySelectorAll(`#${formId} input`).forEach(el => {
        el.classList.remove('border-danger');
        el.classList.add('border-border');
    });
}