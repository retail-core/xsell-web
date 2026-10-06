function toggleFormSwitch(btn) {
  const isOn = btn.dataset.on === "true";
  btn.dataset.on = (!isOn).toString();
  btn.classList.toggle("bg-accent", !isOn);
  btn.classList.toggle("bg-border", isOn);
  btn.querySelector("span").classList.toggle("translate-x-5", !isOn);
  btn.querySelector("span").classList.toggle("translate-x-0.5", isOn);
  document.getElementById("is_active_input").value = (!isOn).toString();
}

function openBarcodeScanner() {
  // TODO: wire real scanner (zxing/html5-qrcode) later
  console.log("open barcode scanner");
}

function deleteProduct(id) {
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
                setTimeout(() => { window.location.href = '/products'; }, 800);

            } catch (err) {
                showToast('Network error. Try again.', 'error');
            }
        }
    );
}

async function submitProductForm() {
    if (imageUploadPromise) {
        showToast('Finishing image upload...', 'success');
        await imageUploadPromise;
    }
    
    const form = document.getElementById('product-form');
    clearAllFieldErrors('product-form');

    const formData = new FormData(form);
    const data = Object.fromEntries(formData.entries());
    let hasError = false;

    if (!data.name || data.name.trim() === '') {
        showFieldError('name', 'Product name is required');
        hasError = true;
    }
    if (!data.category || data.category.trim() === '') {
        showFieldError('category', 'Category is required');
        hasError = true;
    }
    if (!data.selling_price || parseFloat(data.selling_price) < 0) {
        showFieldError('selling_price', 'Enter a valid selling price');
        hasError = true;
    }

    if (hasError) return;

    const isEdit = form.dataset.isEdit === 'true';
    const productId = form.dataset.productId;
    const url = isEdit ? `/products/${productId}/edit` : '/products/new';

    try {
        const resp = await fetch(url, {
            method: 'POST',
            body: new URLSearchParams(data)
        });

        const toastMsg = resp.headers.get('X-Toast-Message');

        if (!resp.ok) {
            showToast(toastMsg || 'Could not save product', 'error');
            return;
        }

        showToast(toastMsg || 'Product saved', 'success');
        setTimeout(() => { window.location.href = '/products'; }, 800);

    } catch (err) {
        showToast('Network error. Try again.', 'error');
    }
}

let imageUploadPromise = null;

async function handleImageSelect(input) {
    const file = input.files[0];
    if (!file) return;

    // Instant local preview, before upload completes
    const reader = new FileReader();
    reader.onload = (e) => {
        const preview = document.getElementById('product-image-preview');
        preview.src = e.target.result;
        preview.classList.remove('hidden');
        document.getElementById('product-image-placeholder').classList.add('hidden');
    };
    reader.readAsDataURL(file);

    const formData = new FormData();
    formData.append('image', file);

     imageUploadPromise = fetch('/products/upload-image', { method: 'POST', body: formData })
        .then(resp => resp.json())
        .then(result => {
            document.getElementById('product-image-url').value = result.image_url || '';
            if (!result.image_url) {
                showToast('Image upload failed, continuing without image', 'error');
            }
        })
        .catch(() => {
            document.getElementById('product-image-url').value = '';
            showToast('Image upload failed, continuing without image', 'error');
        })
        .finally(() => {
            imageUploadPromise = null;
        });
}