const QRScanner = (() => {
  let overlay;
  let video;
  let messageEl;
  let closeBtn;
  let stream = null;
  let raf = null;
  let onResult = null;
  let canvas = null;
  let ctx = null;

  function ensureOverlay() {
    if (overlay) return;
    overlay = document.createElement("div");
    overlay.className = "popup hidden";

    const content = document.createElement("div");
    content.className = "popup-content scan-overlay";

    const title = document.createElement("h3");
    title.textContent = "Scan QR Code";

    messageEl = document.createElement("p");
    messageEl.className = "device-meta";

    video = document.createElement("video");
    video.setAttribute("playsinline", "true");
    video.autoplay = true;

    const actions = document.createElement("div");
    actions.className = "actions";
    actions.style.marginTop = "1rem";
    closeBtn = document.createElement("button");
    closeBtn.className = "btn-secondary";
    closeBtn.textContent = "Close";
    closeBtn.addEventListener("click", close);

    actions.appendChild(closeBtn);
    content.appendChild(title);
    content.appendChild(messageEl);
    content.appendChild(video);
    content.appendChild(actions);
    overlay.appendChild(content);
    document.body.appendChild(overlay);
  }

  function open(opts) {
    ensureOverlay();
    if (!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
      alert("Camera access is not supported in this browser.");
      return;
    }
    messageEl.textContent = opts?.message || "Point your camera at the QR code.";
    onResult = opts?.onResult || null;
    overlay.classList.remove("hidden");
    navigator.mediaDevices
      .getUserMedia({ video: { facingMode: "environment" } })
      .then((mediaStream) => {
        stream = mediaStream;
        video.srcObject = mediaStream;
        video.onloadedmetadata = () => {
          video.play();
          scanFrame();
        };
      })
      .catch((err) => {
        console.error("Camera error", err);
        close();
      });
  }

  function scanFrame() {
    if (!video || video.readyState !== video.HAVE_ENOUGH_DATA) {
      raf = requestAnimationFrame(scanFrame);
      return;
    }
    if (!canvas) {
      canvas = document.createElement("canvas");
      ctx = canvas.getContext("2d");
    }
    canvas.width = video.videoWidth;
    canvas.height = video.videoHeight;
    ctx.drawImage(video, 0, 0, canvas.width, canvas.height);
    const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height);
    const code = jsQR(imageData.data, canvas.width, canvas.height);
    if (code && code.data) {
      const result = code.data;
      close();
      if (onResult) {
        onResult(result);
      }
      return;
    }
    raf = requestAnimationFrame(scanFrame);
  }

  function close() {
    overlay && overlay.classList.add("hidden");
    if (raf) {
      cancelAnimationFrame(raf);
      raf = null;
    }
    if (stream) {
      stream.getTracks().forEach((track) => track.stop());
      stream = null;
    }
    onResult = null;
  }

  return { open, close };
})();
