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
    name,
    category,
    imageUrl,
    currentQty,
    sellingPrice,
    costPrice,
) {
    document.getElementById("restock-product-name").textContent = name;
    document.getElementById("restock-product-category").textContent = category;
    document.getElementById("restock-product-image").innerHTML = imageUrl
        ? `<img src="${imageUrl}" class="w-full h-full object-cover" />`
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

function submitRestock() {
    console.log("Restock submitted:", restockValues);
    // TODO: hx-post to inventory-service
    closeRestockSheet();
}