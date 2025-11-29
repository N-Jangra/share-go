const DeviceIdentity = (() => {
  const STORAGE_KEY = "fs_device_id";
  let deviceId = null;
  const readyCallbacks = [];

  const storageSupported = (() => {
    try {
      const key = "__fs_probe__";
      window.localStorage.setItem(key, "1");
      window.localStorage.removeItem(key);
      return true;
    } catch (err) {
      return false;
    }
  })();

  function defaultName() {
    const nav = navigator;
    const platform = (nav.userAgentData && nav.userAgentData.platform) || nav.platform || "Device";
    const ua = (nav.userAgent || "").toLowerCase();
    let browser = "Browser";
    if (ua.includes("edg")) browser = "Edge";
    else if (ua.includes("chrome")) browser = "Chrome";
    else if (ua.includes("safari") && !ua.includes("chrome")) browser = "Safari";
    else if (ua.includes("firefox")) browser = "Firefox";
    else if (ua.includes("android")) browser = "Android";
    return `${platform} • ${browser}`;
  }

  function notifyReady(id) {
    deviceId = id;
    if (id) {
      document.documentElement.dataset.deviceId = id;
    }
    while (readyCallbacks.length) {
      const cb = readyCallbacks.shift();
      cb(id);
    }
  }

  async function register(existingId) {
    if (!window.fetch) {
      return null;
    }
    const payload = { name: defaultName() };
    if (existingId) payload.id = existingId;
    const res = await fetch("/api/devices/register", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });
    if (!res.ok) {
      throw new Error("device registration failed");
    }
    return res.json();
  }

  function getStoredId() {
    if (!storageSupported) return null;
    try {
      return window.localStorage.getItem(STORAGE_KEY);
    } catch {
      return null;
    }
  }

  function setStoredId(id) {
    if (!storageSupported || !id) return;
    try {
      window.localStorage.setItem(STORAGE_KEY, id);
    } catch (err) {
      console.warn("Unable to persist device id", err);
    }
  }

  async function init() {
    const stored = getStoredId();
    try {
      const device = await register(stored);
      if (device && device.id) {
        setStoredId(device.id);
        notifyReady(device.id);
        return;
      }
    } catch (err) {
      console.warn("Device auto-registration failed", err);
    }
    if (stored) {
      notifyReady(stored);
    } else {
      notifyReady(null);
    }
  }

  function onReady(cb) {
    if (deviceId !== null) {
      cb(deviceId);
    } else {
      readyCallbacks.push(cb);
    }
  }

  return {
    init,
    onReady,
    getId: () => deviceId,
  };
})();

document.addEventListener("DOMContentLoaded", () => {
  DeviceIdentity.init();
});
