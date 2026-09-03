function filterByCategory(category, chipEl) {
    if (navigator.vibrate) navigator.vibrate(10);
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