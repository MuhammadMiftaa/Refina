export const GetInitials = (username: string): string => {
  if (!username) return "";

  const words = username.trim().split(/\s+/); // Pisahkan berdasarkan spasi

  if (words.length > 1) {
    // Jika ada lebih dari 1 kata, ambil huruf pertama dari dua kata pertama
    return (words[0][0] + words[1][0]).toUpperCase();
  }

  // Jika hanya ada satu kata, ambil dua huruf pertama
  return words[0].slice(0, 2).toUpperCase();
};

export function shortenMoney(value: number): string {
  if (value >= 1_000_000_000) {
    return (value / 1_000_000_000).toFixed(1).replace(/\.0$/, "") + " B";
  } else if (value >= 1_000_000) {
    return (value / 1_000_000).toFixed(1).replace(/\.0$/, "") + " M";
  } else if (value >= 1_000) {
    return (value / 1_000).toFixed(1).replace(/\.0$/, "") + " K";
  } else {
    return value.toString();
  }
}

export function formatCurrency(amount: string | number): string {
  const numericAmount =
    typeof amount === "string" ? parseFloat(amount) : amount;
  if (isNaN(numericAmount)) {
    throw new Error("Invalid number format");
  }
  return numericAmount.toLocaleString("id-ID");
}

export async function handleCopy(textToCopy: string) {
  try {
    await navigator.clipboard.writeText(textToCopy);
  } catch (err) {
    console.error("Gagal menyalin teks:", err);
  }
}

export function formatRawDate(rawDate: string): [string, string, string] {
  const date = new Date(rawDate);

  // Konversi ke waktu lokal (WIB = UTC+7)
  // const wibOffset = 7 * 60; // 7 jam dalam menit
  // const localDate = new Date(date.getTime() + wibOffset * 60 * 1000);

  const days = [
    "Sunday",
    "Monday",
    "Tuesday",
    "Wednesday",
    "Thursday",
    "Friday",
    "Saturday",
  ];
  const months = [
    "Januari",
    "Februari",
    "Maret",
    "April",
    "Mei",
    "Juni",
    "Juli",
    "Agustus",
    "September",
    "Oktober",
    "November",
    "Desember",
  ];

  const dayName = days[date.getDay()];
  const hours = date.getHours().toString().padStart(2, "0");
  const minutes = date.getMinutes().toString().padStart(2, "0");
  const time = `${hours}:${minutes} WIB`;

  const dateNum = date.getDate();
  const month = months[date.getMonth()];
  const year = date.getFullYear();
  const fullDate = `${dateNum} ${month} ${year}`;

  return [dayName, time, fullDate];
}

export function generateColorByType(type: string): string {
  switch (type) {
    case "income":
      return "green-500";
    case "expense":
      return "red-500";
    case "fund_transfer":
      return "orange-500";
    default:
      return "gray-500"; // Warna default jika tipe tidak dikenali
  }
}

export function toLocalISOString(date: Date): string {
  const pad = (n: number) => n.toString().padStart(2, '0');

  const year = date.getFullYear();
  const month = pad(date.getMonth() + 1);
  const day = pad(date.getDate());
  const hours = pad(date.getHours());
  const minutes = pad(date.getMinutes());
  const seconds = pad(date.getSeconds());

  const offsetMinutes = date.getTimezoneOffset(); // dalam menit, negatif untuk GMT+
  const offsetSign = offsetMinutes <= 0 ? '+' : '-';
  const offsetAbs = Math.abs(offsetMinutes);
  const offsetHours = pad(Math.floor(offsetAbs / 60));
  const offsetMins = pad(offsetAbs % 60);

  return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}${offsetSign}${offsetHours}:${offsetMins}`;
}
