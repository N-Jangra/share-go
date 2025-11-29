const DevicePage = (() => {
  const state = {
    deviceId: null,
    pollTimer: null,
    pending: null,
  };

  function submitRegistration(event) {
    event.preventDefault();
    const input = document.getElementById("device-name");
    if (!input.value.trim() || !state.deviceId) return;
    fetch("/api/devices/register", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ id: state.deviceId, name: input.value.trim() }),
    })
      .then((res) => {
        if (!res.ok) throw new Error("Failed to register");
        return res.json();
      })
      .then((device) => {
        state.deviceId = device.id;
        input.value = "";
        showListener(device.id);
      })
      .catch(() => alert("Unable to register device. Please try another name."));
  }

  function showListener(id) {
    const listener = document.getElementById("listener");
    listener.classList.remove("hidden");
    const idLabel = document.getElementById("device-id");
    idLabel.textContent = id;
  }

  function copyDeviceId() {
    if (!state.deviceId) return;
    if (navigator.clipboard) {
      navigator.clipboard.writeText(state.deviceId);
      return;
    }
    const textarea = document.createElement("textarea");
    textarea.value = state.deviceId;
    textarea.style.position = "fixed";
    textarea.style.opacity = "0";
    document.body.appendChild(textarea);
    textarea.select();
    document.execCommand("copy");
    document.body.removeChild(textarea);
  }

  function startPolling() {
    stopPolling();
    if (!state.deviceId) return;
    const poll = async () => {
      try {
        const res = await fetch(`/api/devices/pending?id=${encodeURIComponent(state.deviceId)}`);
        if (res.status === 204) return;
        if (!res.ok) throw new Error("bad response");
        const data = await res.json();
        if (data && (!state.pending || state.pending.transferId !== data.transferId)) {
          state.pending = data;
          showPopup(data);
        }
      } catch (err) {
        // ignore errors, will try again
      }
    };
    state.pollTimer = setInterval(poll, 3500);
    poll();
  }

  function stopPolling() {
    if (state.pollTimer) {
      clearInterval(state.pollTimer);
      state.pollTimer = null;
    }
  }

  function showPopup(data) {
    const popup = document.getElementById("incoming-popup");
    popup.classList.remove("hidden");
    document.getElementById("incoming-name").textContent = data.fileName;
    const size = (data.fileSize / (1024 * 1024)).toFixed(2);
    document.getElementById("incoming-size").textContent = `${size} MB · ${data.mime}`;
  }

  function hidePopup() {
    document.getElementById("incoming-popup").classList.add("hidden");
    state.pending = null;
  }

  function openTransfer() {
    if (!state.pending) return;
    const target = `${window.location.origin}/incoming?id=${encodeURIComponent(state.pending.transferId)}&token=${encodeURIComponent(state.pending.token)}`;
    fetch("/api/devices/clear", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ deviceId: state.deviceId, transferId: state.pending.transferId }),
    }).finally(() => {
      hidePopup();
      window.location.href = target;
    });
  }

  function dismissTransfer() {
    if (!state.pending) {
      hidePopup();
      return;
    }
    fetch("/api/devices/clear", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ deviceId: state.deviceId, transferId: state.pending.transferId }),
    }).finally(() => hidePopup());
  }

  function init() {
    const form = document.getElementById("register-device");
    form.addEventListener("submit", submitRegistration);
    document.getElementById("copy-device-id").addEventListener("click", copyDeviceId);
    document.getElementById("open-transfer").addEventListener("click", openTransfer);
    document.getElementById("dismiss-transfer").addEventListener("click", dismissTransfer);

    DeviceIdentity.onReady((id) => {
      if (!id) return;
      state.deviceId = id;
      showListener(id);
      startPolling();
    });
  }

  return { init };
})();

document.addEventListener("DOMContentLoaded", () => DevicePage.init());
