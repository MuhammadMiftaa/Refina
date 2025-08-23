import { PieChartType } from "@/types/Chart";

export function createCookiesOpts(): Cookies.CookieAttributes {
  const mode = import.meta.env.VITE_MODE;
  let options: Cookies.CookieAttributes = { expires: 7 };
  switch (mode) {
    case "production":
      options = {
        expires: 7,
        secure: true,
        sameSite: "None",
        domain: ".miftech.web.id",
      };
      break;
    case "development":
      options = {
        expires: 7,
      };
      break;
  }

  return options;
}

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

  const dayName = days[date.getUTCDay()];
  const hours = date.getUTCHours().toString().padStart(2, "0");
  const minutes = date.getUTCMinutes().toString().padStart(2, "0");
  const time = `${hours}:${minutes} WIB`;

  const dateNum = date.getUTCDate();
  const month = months[date.getUTCMonth()];
  const year = date.getUTCFullYear();
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
  const pad = (n: number) => n.toString().padStart(2, "0");

  const year = date.getFullYear();
  const month = pad(date.getMonth() + 1);
  const day = pad(date.getDate());
  const hours = pad(date.getHours());
  const minutes = pad(date.getMinutes());
  const seconds = pad(date.getSeconds());

  const offsetMinutes = date.getTimezoneOffset(); // dalam menit, negatif untuk GMT+
  const offsetSign = offsetMinutes <= 0 ? "+" : "-";
  const offsetAbs = Math.abs(offsetMinutes);
  const offsetHours = pad(Math.floor(offsetAbs / 60));
  const offsetMins = pad(offsetAbs % 60);

  return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}${offsetSign}${offsetHours}:${offsetMins}`;
}

function fileToBase64(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();

    reader.onload = () => {
      if (typeof reader.result === "string") {
        resolve(reader.result);
      } else {
        reject("Failed to convert file to base64");
      }
    };

    reader.onerror = reject;
    reader.readAsDataURL(file);
  });
}

export async function convertFilesToBase64(files: File[]): Promise<string[]> {
  const base64Array = await Promise.all(
    files.map((file) => fileToBase64(file)),
  );
  return base64Array;
}

export function shortenFilename(filename: string): string {
  const ext = filename.split(".").pop() || "";
  const base = filename.slice(0, -ext.length - 1); // hapus ".ext"

  return `${base.slice(0, 35)}.....`;
  // return `${base.slice(0, 30)}........${base.slice(-1)}.${ext}`;
}

export function bytesToMegabytes(bytes: number): string {
  const mb = bytes / (1024 * 1024);
  return `${mb.toFixed(2)} MB`;
}

const generateSoftBlue = (input: number, total: number): string => {
  if (input < 1 || input > total) throw new Error("Input out of range");

  const hue = 216; // tetap di sekitar warna #2A7CFA
  const saturationBase = 96;
  const saturationRange = 10; // ±10% variasi
  const lightnessBase = 57;
  const lightnessRange = 15; // ±15% variasi

  const step = (input - 1) / (total - 1);

  const saturation = Math.max(
    0,
    Math.min(
      100,
      saturationBase - saturationRange / 2 + step * saturationRange,
    ),
  );
  const lightness = Math.max(
    0,
    Math.min(100, lightnessBase - lightnessRange / 2 + step * lightnessRange),
  );

  return `hsl(${hue}, ${saturation.toFixed(0)}%, ${lightness.toFixed(0)}%)`;
};

type ChartConfig = Record<
  string,
  {
    label: string;
    color?: string;
  }
>;

export const buildPieChartConfig = (
  categories: PieChartType[],
): ChartConfig => {
  const config: ChartConfig = {
    value: {
      label: "Value",
    },
  };

  categories.forEach((item, index) => {
    config[item.parent_category_name] = {
      label: item.parent_category_name,
      color: generateSoftBlue(index + 1, categories.length),
    };
  });

  return config;
};

export const getLast6MonthsRange = (): string => {
  const today = new Date();
  const end = new Intl.DateTimeFormat("en-US", { month: "long" }).format(today);
  const endYear = today.getFullYear();

  const past = new Date(today);
  past.setMonth(today.getMonth() - 5); // karena bulan dimulai dari 0
  const start = new Intl.DateTimeFormat("en-US", { month: "long" }).format(
    past,
  );
  const startYear = past.getFullYear();

  if (startYear !== endYear) {
    return `${start} ${startYear} – ${end} ${endYear}`;
  }

  return `${start} – ${end} ${endYear}`;
};

export const generateAvatarFromName = (
  name: string,
): { initials: string; textColor: string; backgroundColor: string } => {
  const words = name.trim().split(/\s+/).slice(0, 2);
  const initials = words.map((w) => w[0].toUpperCase()).join("");

  // Generate random background color
  const hex = Math.floor(Math.random() * 0xffffff)
    .toString(16)
    .padStart(6, "0");
  const backgroundColor = `#${hex}`;

  // Hitung brightness untuk menentukan warna teks
  const r = parseInt(hex.slice(0, 2), 16);
  const g = parseInt(hex.slice(2, 4), 16);
  const b = parseInt(hex.slice(4, 6), 16);
  const yiq = (r * 299 + g * 587 + b * 114) / 1000;
  const textColor = yiq >= 128 ? "#000000" : "#FFFFFF";

  return {
    initials,
    textColor: "text-[" + textColor + "]",
    backgroundColor: "bg-[" + backgroundColor + "]",
  };
};
