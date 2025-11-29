const SharePage = (() => {
  const state = {
    shareLink: "",
    transferId: "",
    token: "",
    currentDeviceId: null,
  };

  function renderQR() {
    const el = document.getElementById("qr");
    if (!el) return;
    el.innerHTML = "";
    new QRCode(el, {
      text: state.shareLink,
      width: 220,
      height: 220,
      colorDark: "#38bdf8",
      colorLight: "#0f172a",
    });
  }

  function copyLink() {
    if (!navigator.clipboard) {
      fallbackCopy(state.shareLink);
      return;
    }
    navigator.clipboard.writeText(state.shareLink).then(showCopied, () => fallbackCopy(state.shareLink));
  }

  function fallbackCopy(text) {
    const textarea = document.createElement("textarea");
    textarea.value = text;
    textarea.style.position = "fixed";
    textarea.style.opacity = "0";
    document.body.appendChild(textarea);
    textarea.select();
    document.execCommand("copy");
    document.body.removeChild(textarea);
    showCopied();
  }

  function showCopied() {
    const button = document.getElementById("copy-link");
    if (!button) return;
    const defaultText = button.textContent;
    button.textContent = "Copied!";
    setTimeout(() => (button.textContent = defaultText), 1600);
  }

  async function fetchDevices() {
    try {
      const res = await fetch("/api/devices");
      if (!res.ok) throw new Error("Failed to load devices");
      const devices = await res.json();
      renderDevices(devices);
    } catch (err) {
      renderDevices([]);
    }
  }

  function renderDevices(devices) {
    const container = document.getElementById("devices");
    if (!container) return;
    container.innerHTML = "";
    const filtered = devices.filter((device) => !(state.currentDeviceId && device.id === state.currentDeviceId));
    if (!filtered.length) {
      container.innerHTML = '<p class="device-meta">No devices registered yet. Open the device listener to register.</p>';
      return;
    }
    filtered.forEach((device) => {
      const item = document.createElement("div");
      item.className = "device-item";
      const details = document.createElement("div");
      details.innerHTML = `<div class="device-name">${device.name}</div><div class="device-meta">ID: ${device.id}</div>`;
      const actions = document.createElement("div");
      actions.className = "device-actions";
      const button = document.createElement("button");
      button.textContent = "Send";
      button.addEventListener("click", () => notifyDevice(device.id, button));
      actions.appendChild(button);
      item.append(details, actions);
      container.appendChild(item);
    });
  }

  async function notifyDevice(deviceId, button) {
    button.disabled = true;
    const original = button.textContent;
    button.textContent = "Sending...";
    try {
      const res = await fetch("/api/devices/notify", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          deviceId,
          transferId: state.transferId,
          token: state.token,
        }),
      });
      if (!res.ok) throw new Error("failed");
      button.textContent = "Sent!";
    } catch (err) {
      button.textContent = "Try again";
    } finally {
      setTimeout(() => {
        button.textContent = original;
        button.disabled = false;
      }, 1800);
    }
  }

  function bindEvents() {
    const copyBtn = document.getElementById("copy-link");
    copyBtn && copyBtn.addEventListener("click", copyLink);
    const openScanner = document.getElementById("open-scanner");
    openScanner &&
      openScanner.addEventListener("click", () =>
        QRScanner.open({
          message: "Scan a QR code to open its link instantly.",
          onResult: (data) => {
            if (data) {
              window.location.href = data;
            }
          },
        }),
      );
  }

  return {
    init(config) {
      state.shareLink = config.shareLink;
      state.transferId = config.transferId;
      state.token = config.token;
      renderQR();
      bindEvents();
      fetchDevices();
      DeviceIdentity.onReady((id) => {
        state.currentDeviceId = id;
        fetchDevices();
      });
      setInterval(fetchDevices, 15000);
    },
  };
})();
