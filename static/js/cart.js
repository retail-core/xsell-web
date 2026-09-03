let cart = [];
let paymentMethod = 'CASH';

function addToCart(id, name, price, imageUrl) {
    if (navigator.vibrate) navigator.vibrate(20);
    const existing = cart.find(item => item.id === id);
    if (existing) {
        existing.qty += 1;
    } else {
        cart.push({ id, name, price, imageUrl, qty: 1 });
    }
    renderCart();
    // showToast(name + ' added to cart', 'success');
}

function changeQty(id, delta) {
    const item = cart.find(i => i.id === id);
    if (!item) return;
    item.qty += delta;
    if (item.qty <= 0) {
        cart = cart.filter(i => i.id !== id);
    }
    renderCart();
}

function computeTotal() {
    return cart.reduce((sum, item) => sum + (item.price * item.qty), 0);
}

function formatNaira(n) {
    return '₦' + n.toLocaleString('en-NG');
}

function renderCart() {
    const container = document.getElementById('cart-items-container');
    const totalDisplay = document.getElementById('cart-total-display');
    const cartCountBadge = document.getElementById('cart-count-badge');

    if (cart.length === 0) {
        container.innerHTML = '<div class="flex items-center justify-center py-16 text-sm text-text-secondary">No item selected</div>';
    } else {
        container.innerHTML = cart.map(item => `
            <div class="flex items-center gap-3 py-3 border-b border-border">
                <div class="w-11 h-11 rounded-lg bg-surface-alt borde border-borde overflow-hidden shrink-0">
                    ${item.imageUrl ? `<img src="${item.imageUrl}" alt="${item.name}" class="w-full h-full object-cover" />` : ''}
                </div>
                <div class="flex-1 min-w-0">
                    <p class="text-sm text-text-primary truncate">${item.name}</p>
                    <p class="text-sm text-text-secondary">${formatNaira(item.price)}</p>
                </div>
                <div class="flex items-center gap-2 shrink-0">
                    <button onclick="tapFeedback(); changeQty('${item.id}', -1)" class="w-7 h-7 rounded-lg border border-border flex items-center justify-center text-text-secondary hover:border-accent hover:text-accent">−</button>
                    <span class="w-5 text-center text-sm font-medium text-text-primary">${item.qty}</span>
                    <button onclick="tapFeedback(); changeQty('${item.id}', 1)" class="w-7 h-7 rounded-lg border border-border flex items-center justify-center text-text-secondary hover:border-accent hover:text-accent">+</button>
                </div>
            </div>
        `).join('');
    }

    totalDisplay.textContent = formatNaira(computeTotal());

    if (cartCountBadge) {
        const count = cart.reduce((s, i) => s + i.qty, 0);
        cartCountBadge.textContent = count;
        cartCountBadge.classList.toggle('hidden', count === 0);
    }
}

function selectPaymentMethod(method) {
    paymentMethod = method;
    document.querySelectorAll('.payment-tab').forEach(tab => {
        const active = tab.dataset.method === method;
        tab.classList.toggle('bg-surface', active);
        tab.classList.toggle('text-accent', active);
        tab.classList.toggle('shadow-sm', active);
        tab.classList.toggle('text-text-secondary', !active);
    });
}

function openCartSheet() {
    renderCart();
    document.getElementById('cart-default-view').style.display = 'flex';
    document.getElementById('cart-default-view').classList.add('flex-col', 'flex-1');
    document.getElementById('cart-success-view').classList.add('hidden');
    document.getElementById('cart-sheet').classList.remove('translate-y-full');
    document.getElementById('cart-sheet-backdrop').classList.remove('hidden');
    lockBodyScroll(true);
}

function closeCartSheet() {
    document.getElementById('cart-sheet').classList.add('translate-y-full');
    document.getElementById('cart-sheet-backdrop').classList.add('hidden');
    lockBodyScroll(false);
}

async function submitOrder() {
    if (cart.length === 0) {
        showToast('Add at least one item', 'error');
        return;
    }

    const btn = document.getElementById('place-order-btn');
    const originalContent = btn.innerHTML;
    btn.disabled = true;
    btn.innerHTML = `
        <span class="flex items-center gap-2 mx-auto">
            <svg class="animate-spin" width="18" height="18" viewBox="0 0 24 24" fill="none">
                <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="3" opacity="0.25"/>
                <path d="M12 2a10 10 0 0 1 10 10" stroke="currentColor" stroke-width="3" stroke-linecap="round"/>
            </svg>
            Processing...
        </span>
    `;

    const payload = {
        payment_method: paymentMethod,
        inventory_items: cart.map(item => ({
            inventory_id: item.id,
            quantity: item.qty
        }))
    };

    try {
        const resp = await fetch('/pos/checkout', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });

        if (!resp.status) {
            throw new Error('Checkout failed');
        }

        document.getElementById('cart-default-view').style.display = 'none';
        document.getElementById('cart-success-view').classList.remove('hidden');
        document.getElementById('cart-success-view').classList.add('flex');

        cart = [];

        setTimeout(() => {
            closeCartSheet();
            htmx.ajax('GET', '/pos', { target: '#main-content', swap: 'innerHTML' });
        }, 1600)

    } catch (err) {
        showToast('Could not complete sale. Try again.', 'error');
        btn.disabled = false;
        btn.innerHTML = originalContent;
    }
}

function toggleSwitch(btn) {
    const isOn = btn.dataset.on === "true";
    btn.dataset.on = (!isOn).toString();
    btn.classList.toggle("bg-accent", !isOn);
    btn.classList.toggle("bg-border", isOn);
    btn.querySelector("span").classList.toggle("translate-x-5", !isOn);
    btn.querySelector("span").classList.toggle("translate-x-0.5", isOn);
}

let activeCategory = 'All';

function togglePOSSearch() {
    const bar = document.getElementById('pos-search-bar');
    bar.classList.toggle('hidden');
    if (!bar.classList.contains('hidden')) {
        document.getElementById('pos-search-input').focus();
    } else {
        document.getElementById('pos-search-input').value = '';
        applyFilters();
    }
}

function filterPOSItems() {
    applyFilters();
}

function filterByCategory(category, chipEl) {
    activeCategory = category;

    document.querySelectorAll('.category-chip').forEach(chip => {
        chip.classList.remove('border-accent', 'bg-accent-soft', 'text-accent');
        chip.classList.add('border-border', 'text-text-secondary');
    });
    chipEl.classList.remove('border-border', 'text-text-secondary');
    chipEl.classList.add('border-accent', 'bg-accent-soft', 'text-accent');

    applyFilters();
}

function applyFilters(query) {
    query = query || '';
    document.querySelectorAll('[data-category]').forEach(card => {
        const matchesCategory = activeCategory === 'All' || card.dataset.category === activeCategory;
        const matchesSearch = query === '' || card.dataset.name.toLowerCase().includes(query);
        card.style.display = (matchesCategory && matchesSearch) ? '' : 'none';
    });
}