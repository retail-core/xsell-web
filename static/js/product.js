function toggleProductMenu(btn) {
    const popover = btn.nextElementSibling;
    const isOpen = !popover.classList.contains("hidden");

    // close any other open menus first
    document
        .querySelectorAll(".product-menu-popover")
        .forEach((p) => p.classList.add("hidden"));

    if (!isOpen) popover.classList.remove("hidden");
}

// close on outside tap
document.addEventListener("click", function (e) {
    const isMenuButton = e.target.closest(
        'button[onclick*="toggleProductMenu"]',
    );
    const isInsidePopover = e.target.closest(".product-menu-popover");

    if (!isMenuButton && !isInsidePopover) {
        const anyOpenMenu = document.querySelector(
            ".product-menu-popover:not(.hidden)",
        );
        if (anyOpenMenu) {
            e.preventDefault();
            e.stopPropagation();
        }
        document
            .querySelectorAll(".product-menu-popover")
            .forEach((p) => p.classList.add("hidden"));
    }
});

function toggleAddProductSheet() {
    const sheet = document.getElementById("add-product-sheet");
    const isOpening = sheet.classList.contains("translate-y-full");
    sheet.classList.toggle("translate-y-full");
    document.getElementById("add-product-backdrop").classList.toggle("hidden");
    lockBodyScroll(isOpening);
}

let activeField = "quantity";
const restockValues = { quantity: 0, selling: 0, cost: 0 };

function setActiveField(field) {
    activeField = field;
    document
        .querySelectorAll(".restock-field")
        .forEach((el) => el.classList.remove("active"));
    document.getElementById("field-" + field).classList.add("active");
}

function numpadInput(digit) {
    const current = restockValues[activeField].toString();
    const next = current === "0" ? digit : current + digit;
    restockValues[activeField] = parseInt(next, 10);
    document.getElementById("value-" + activeField).textContent =
        restockValues[activeField];
}

function numpadBackspace() {
    const current = restockValues[activeField].toString();
    const next = current.slice(0, -1) || "0";
    restockValues[activeField] = parseInt(next, 10);
    document.getElementById("value-" + activeField).textContent =
        restockValues[activeField];
}

function adjustQty(delta) {
    restockValues.quantity = Math.max(0, restockValues.quantity + delta);
    document.getElementById("value-quantity").textContent =
        restockValues.quantity;
}

function openRestockSheet(
    id,
    name,
    category,
    imageUrl,
    currentQty,
    sellingPrice,
    costPrice,
) {
    console.log(id, name, category, imageUrl)
    document.getElementById('restock-sheet').dataset.productId = id;
    document.getElementById("restock-product-name").textContent = name;
    document.getElementById("restock-product-category").textContent = category;
    document.getElementById("restock-product-image").innerHTML = imageUrl
        ? `<img src="${imageUrl}" alt="${name}" class="w-full h-full object-cover" />`
        : "";

    restockValues.quantity = currentQty;
    restockValues.selling = sellingPrice;
    restockValues.cost = costPrice;
    document.getElementById("value-quantity").textContent = currentQty;
    document.getElementById("value-selling").textContent = sellingPrice;
    document.getElementById("value-cost").textContent = costPrice;

    setActiveField("quantity");
    document
        .getElementById("restock-sheet")
        .classList.remove("translate-y-full");
    document.getElementById("restock-backdrop").classList.remove("hidden");
    lockBodyScroll(true);
}

function closeRestockSheet() {
    document.getElementById("restock-sheet").classList.add("translate-y-full");
    document.getElementById("restock-backdrop").classList.add("hidden");
    lockBodyScroll(false);
}

async function submitRestock() {
    const productId = document.getElementById('restock-sheet').dataset.productId;

    try {
        const resp = await fetch(`/products/${productId}/restock`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                quantity: restockValues.quantity,
                selling_price: restockValues.selling,
                cost_price: restockValues.cost
            })
        });

        const toastMsg = resp.headers.get('X-Toast-Message');
        closeRestockSheet();

        if (!resp.ok) {
            showToast(toastMsg || 'Could not update stock', 'error');
            return;
        }

        showToast(toastMsg || 'Stock updated', 'success');
        setTimeout(() => {
            htmx.ajax('GET', '/products', { target: '#main-content', swap: 'innerHTML' });
        }, 800);

    } catch (err) {
        closeRestockSheet();
        showToast('Network error. Try again.', 'error');
    }
}

function filterByCategory(category, chipEl) {
    document.querySelectorAll('[data-category]').forEach(card => {
        const show = category === 'All' || card.dataset.category === category;
        card.style.display = show ? '' : 'none';
    });

    document.querySelectorAll('.category-chip').forEach(chip => {
        chip.classList.remove('border-accent', 'bg-accent-soft', 'text-accent');
        chip.classList.add('border-border', 'text-text-secondary');
    });
    chipEl.classList.remove('border-border', 'text-text-secondary');
    chipEl.classList.add('border-accent', 'bg-accent-soft', 'text-accent');
}

function deleteProductQuick(id, btnEl) {
    showConfirmDialog(
        'Delete Product',
        'This cannot be undone.',
        async () => {
            try {
                const resp = await fetch(`/products/${id}`, { method: 'DELETE' });
                const toastMsg = resp.headers.get('X-Toast-Message');
                if (!resp.ok) {
                    showToast(toastMsg || 'Could not delete product', 'error');
                    return;
                }
                showToast(toastMsg || 'Product deleted', 'success');
                const row = btnEl.closest('a[data-category]');
                if (row) row.remove();
            } catch (err) {
                showToast('Network error. Try again.', 'error');
            }
        }
    );
}

async function toggleProductStatus(id, currentlyActive, btnEl) {
    const newStatus = !currentlyActive;

    // Close the popover immediately
    const popover = btnEl.closest('.product-menu-popover');
    if (popover) popover.classList.add('hidden');

    try {
        const resp = await fetch(`/products/${id}/status`, {
            method: 'PATCH',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ active: newStatus })
        });

        const toastMsg = resp.headers.get('X-Toast-Message');

        if (!resp.ok) {
            showToast(toastMsg || 'Could not update status', 'error');
            return;
        }

        showToast(toastMsg || 'Status updated', 'success');

        // Update the row in place: toggle badge, update button's own state/label
        const row = btnEl.closest('a[data-category]');
        if (row) {
            let badge = row.querySelector('.inactive-badge');
            if (!newStatus) {
                if (!badge) {
                    badge = document.createElement('span');
                    badge.className = 'inactive-badge text-[10px] px-2 py-0.5 rounded-full bg-danger-soft text-danger shrink-0 ml-2';
                    badge.textContent = 'Inactive';
                    row.querySelector('.flex-1').appendChild(badge);
                }
            } else if (badge) {
                badge.remove();
            }
        }

        btnEl.setAttribute('onclick', `toggleProductStatus('${id}', ${newStatus}, this)`);
        btnEl.querySelector('span').textContent = newStatus ? 'Deactivate' : 'Activate';

    } catch (err) {
        console.error(err);
        showToast('Network error. Try again.', 'error');
    }
}