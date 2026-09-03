async function openOrderDetail(orderId) {
    // Open sheet immediately, show subtle loading state in existing elements
    document.getElementById('order-id').textContent = 'Loading...';
    document.getElementById('order-sold-by').textContent = '';
    document.getElementById('order-datetime').textContent = '';
    document.getElementById('order-items').innerHTML = '<p class="text-sm text-text-secondary text-center py-6">Loading order...</p>';
    document.getElementById('order-subtotal').textContent = '';
    document.getElementById('order-payment-badge').innerHTML = '';
    document.getElementById('order-status-badge').innerHTML = '';
    document.getElementById('order-credit-section').classList.add('hidden');

    document.getElementById('order-sheet').classList.remove('translate-y-full');
    document.getElementById('order-backdrop').classList.remove('hidden');
    lockBodyScroll(true);

    try {
        const resp = await fetch(`/orders/${orderId}`);
        if (!resp.ok) throw new Error('Failed to load order');
        const order = await resp.json();
        renderOrderDetail(order);
    } catch (err) {
        document.getElementById('order-id').textContent = 'Order';
        document.getElementById('order-items').innerHTML = '<p class="text-sm text-danger text-center py-6">Could not load order details</p>';
    }
}

function renderOrderDetail(order) {
    document.getElementById('order-id').textContent = 'Order #' + (order.id.slice(0, 8));
    document.getElementById('order-sold-by').textContent = 'Sold by ' + (order.sold_by || 'Unknown');
    document.getElementById('order-datetime').textContent = new Date(order.created_at).toLocaleString('en-NG', {
        month: 'short', day: 'numeric', year: 'numeric', hour: '2-digit', minute: '2-digit'
    });
    document.getElementById('order-subtotal').textContent = formatNaira(order.total_amount);

    const isCash = order.payment_method === 'CASH';
    document.getElementById('order-payment-badge').innerHTML =
        `<span class="text-xs px-2.5 py-1 rounded-full border ${isCash ? 'border-accent text-accent' : 'border-border text-text-secondary'}">${titleCase(order.payment_method)}</span>`;

    const isCredit = order.status === 'CREDIT';
    document.getElementById('order-status-badge').innerHTML =
        `<span class="text-xs px-2.5 py-1 rounded-full ${isCredit ? 'bg-danger-soft text-danger' : 'bg-accent-soft text-accent'}">${titleCase(order.status)}</span>`;

    document.getElementById('order-items').innerHTML = order.items.map(item => `
        <div class="flex items-center gap-3">
            <div class="w-11 h-11 rounded-lg bg-surface-alt borde border-border shrink-0 overflow-hidden">
                ${item.image_url ? `<img src="${item.image_url}" alt="${item.name}" class="w-full h-full object-cover" />` : ''}
            </div>
            <div class="flex-1 min-w-0">
                <p class="text-xs text-text-primary">${item.product_name}</p>
                <p class="text-xs text-text-secondary">Qty ${item.quantity} · ${formatNaira(item.unit_price)} each</p>
            </div>
            <span class="text-sm font-display font-semibold text-text-primary">${formatNaira(item.total_price)}</span>
        </div>
    `).join('');

    const creditSection = document.getElementById('order-credit-section');
    if (isCredit && order.customer_name) {
        creditSection.classList.remove('hidden');
        document.getElementById('order-credit-name').textContent = order.customer_name;
        document.getElementById('order-credit-avatar').textContent = order.customer_name.split(' ').map(n => n[0]).join('').slice(0, 2);
    } else {
        creditSection.classList.add('hidden');
    }
}

function titleCase(s) {
    if (!s) return '';
    return s.charAt(0) + s.slice(1).toLowerCase();
}

function closeOrderDetail() {
    document.getElementById('order-sheet').classList.add('translate-y-full');
    document.getElementById('order-backdrop').classList.add('hidden');
    lockBodyScroll(false);
}

function toggleSalesGroup(idx) {
    const body = document.getElementById('group-body-' + idx);
    const chevron = document.getElementById('group-chevron-' + idx);
    body.classList.toggle('hidden');
    chevron.classList.toggle('rotate-180');
}

// document.getElementById('order-items').innerHTML = order.items.map(item => `
//       <div class="flex items-center gap-3">
//         <div class="w-11 h-11 rounded-lg bg-surface-alt border border-border shrink-0"></div>
//         <div class="flex-1 min-w-0">
//           <p class="text-sm text-text-primary">${item.name}</p>
//           <p class="text-xs text-text-secondary">Qty ${item.qty} · ${item.unitPrice} each</p>
//         </div>
//         <span class="text-sm font-display font-semibold text-text-primary">${item.lineTotal}</span>
//       </div>
//     `).join('');