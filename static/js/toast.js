function showToast(message, type = 'error') {
    const isError = type === 'error';

    Toastify({
        text: message,
        duration: 4000,
        gravity: "top",
        position: "center",
        stopOnFocus: true,
        style: {
            background: 'var(--color-surface)',
            color: isError ? 'var(--color-danger)' : 'var(--color-accent)',
            border: `1px solid ${isError ? 'var(--color-danger)' : 'var(--color-accent)'}`,
            borderRadius: '10px',
            boxShadow: '0 10px 25px -5px rgba(0,0,0,0.1)',
            fontSize: '14px',
            fontWeight: '500',
            fontFamily: 'Inter, sans-serif',
            padding: '14px 18px',
        },
    }).showToast();
}

document.body.addEventListener('htmx:responseError', function (evt) {
    const msg = evt.detail.xhr.getResponseHeader('X-Toast-Message');
    if (msg) showToast(msg, 'error');
});

document.body.addEventListener('htmx:afterRequest', function (evt) {
    if (evt.detail.xhr.status >= 200 && evt.detail.xhr.status < 300) {
        const msg = evt.detail.xhr.getResponseHeader('X-Toast-Message');
        if (msg) showToast(msg, 'success');
    }
});