/**
 * formatBytes converts a numeric byte count into a human readable unit (`KB`, `MB`, etc)
 * @param {Number} bytes - The bytes to format
 * @param {Number} decimals - Decimal count to use
 * @returns {String} Human readable string format
 */
export function formatBytes(bytes, decimals = 2) {
  if (bytes === undefined || bytes === null || Number.isNaN(Number(bytes))) {
    return "—";
  }

  if (bytes === 0) return "0 Bytes";

  const k = 1024;
  const dm = decimals < 0 ? 0 : decimals;
  const sizes = ["Bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"];

  const i = Math.floor(Math.log(bytes) / Math.log(k));

  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + " " + sizes[i];
}

export function formatInteger(value) {
  if (value === undefined || value === null || Number.isNaN(Number(value))) {
    return "—";
  }

  return Number(value).toLocaleString();
}

export function formatDateTime(value) {
  if (!value) return "—";

  return new Date(value).toLocaleString(undefined, {
    dateStyle: "medium",
    timeStyle: "medium",
  });
}

export function formatDuration(start, end) {
  if (!start || !end) return "—";

  const ms = new Date(end).getTime() - new Date(start).getTime();
  if (!Number.isFinite(ms) || ms < 0) return "—";

  const seconds = Math.round(ms / 1000);
  if (seconds < 60) return `${seconds}s`;

  const minutes = Math.floor(seconds / 60);
  const remainingSeconds = seconds % 60;
  if (minutes < 60) return `${minutes}m ${remainingSeconds}s`;

  const hours = Math.floor(minutes / 60);
  const remainingMinutes = minutes % 60;
  return `${hours}h ${remainingMinutes}m`;
}
