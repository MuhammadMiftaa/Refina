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
  const numericAmount = typeof amount === "string" ? parseFloat(amount) : amount;
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
