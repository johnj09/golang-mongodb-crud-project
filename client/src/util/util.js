function UnixToDate(unix) {
  const date = new Date(unix * 1000); // Convert seconds to milliseconds
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0"); // Months are zero-based
  const day = String(date.getDate()).padStart(2, "0");
  const hours = String(date.getHours()).padStart(2, "0");
  const minutes = String(date.getMinutes()).padStart(2, "0");
  const seconds = String(date.getSeconds()).padStart(2, "0");
  const timeCreated = `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
  return timeCreated;
}

export { UnixToDate };
