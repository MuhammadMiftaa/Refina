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
