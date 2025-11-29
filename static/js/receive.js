const ReceivePage = (() => {
  function bindScanner() {
    const btn = document.getElementById("scan-receive");
    if (!btn) return;
    btn.addEventListener("click", () => {
      QRScanner.open({
        message: "Scan the sender's QR code to open the transfer.",
        onResult: (data) => {
          if (data) {
            window.location.href = data;
          }
        },
      });
    });
  }

  return {
    init() {
      bindScanner();
    },
  };
})();

document.addEventListener("DOMContentLoaded", () => ReceivePage.init());
