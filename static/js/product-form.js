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