package prepopulate

import (
	"sync"

	"server/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CategoryTypeSeeder(db *gorm.DB) {
	const (
		Income  entity.CategoryType = "income"
		Expense entity.CategoryType = "expense"
	)

	// Data kategori utama dan child
	categoriesData := map[string]struct {
		Type  entity.CategoryType
		Items []string
	}{
		"Gaji/Pendapatan Tetap": {Income, []string{"Gaji Bulanan", "Gaji Freelance", "Honor Projek", "Bonus Tahunan"}},
		"Pendapatan Pasif":      {Income, []string{"Bunga Deposito", "Dividen Saham", "Sewa Properti", "Royalti"}},
		"Pendapatan Lainnya":    {Income, []string{"Hadiah/Uang Saku", "Penjualan Barang Bekas", "Refund/Pengembalian Dana", "Pendapatan Investasi"}},
		"Pendapatan Sampingan":  {Income, []string{"Jualan Online", "Pendapatan dari Hobi", "Pendapatan dari Afiliasi"}},

		"Kebutuhan Pokok":          {Expense, []string{"Makanan & Minuman", "Belanja Bulanan", "Transportasi Harian", "Listrik, Air, & Gas", "Internet & Telepon", "Perlengkapan Kebersihan", "Rokok & Alkohol", "Perawatan Hewan Peliharaan"}},
		"Perumahan":                {Expense, []string{"Sewa Rumah/Kontrakan", "Cicilan Rumah/KPR", "Perbaikan & Pemeliharaan Rumah", "Perlengkapan Rumah Tangga", "Pajak Properti", "Asuransi Rumah", "Dekorasi & Kebun", "Biaya Keamanan"}},
		"Kesehatan":                {Expense, []string{"Asuransi Kesehatan", "Obat-obatan & Resep Dokter", "Medical Check-up", "Donor Darah/Organ", "Layanan Kesehatan", "Peralatan Kesehatan", "Suplemen & Vitamin", "Kesehatan Mental"}},
		"Pendidikan":               {Expense, []string{"Biaya Sekolah/Kuliah", "Buku & Alat Tulis", "Kursus/Pelatihan", "Transportasi Sekolah/Kampus", "Seragam & Perlengkapan Khusus", "Kegiatan Ekstrakurikuler", "Biaya Ujian & Sertifikasi"}},
		"Hiburan & Rekreasi":       {Expense, []string{"Nongkrong di Kafe/Restoran", "Nonton Bioskop/Streaming", "Liburan", "Hobi", "Langganan Digital", "Koleksi & Kesenangan Pribadi", "Event & Konser"}},
		"Fashion & Perawatan Diri": {Expense, []string{"Belanja Pakaian/Aksesoris", "Kosmetik & Skincare", "Salon & Barbershop", "Perawatan Tubuh", "Laundry & Dry Cleaning", "Sewa Kostum/Pakaian", "Perhiasan & Jam Tangan"}},
		"Transportasi & Kendaraan": {Expense, []string{"Cicilan Kendaraan", "Servis & Perawatan Rutin", "Parkir & Tol", "Pajak Kendaraan", "Asuransi Kendaraan", "Bahan Bakar & Pengisian Daya", "Biaya SIM & Surat-surat Lainnya"}},
		"Tagihan & Utilitas":       {Expense, []string{"Tagihan Kartu Kredit", "Tagihan PLN/PDAM/Gas", "Tagihan TV Kabel/Streaming", "Iuran Lingkungan", "Layanan Digital", "Biaya Administrasi"}},
		"Zonasi & Amal":            {Expense, []string{"Zakat", "Sedekah", "Donasi Sosial", "Wakaf", "Bantuan Bencana Alam", "Donasi Pendidikan"}},
		"Teknologi & Gadget":       {Expense, []string{"Pembelian Gadget", "Perbaikan & Aksesori Gadget", "Langganan Software", "Upgrade Perangkat", "Biaya E-commerce"}},
		"Lainnya":                  {Expense, []string{"Hadiah & Ucapan", "Biaya Tak Terduga", "Biaya Hukum & Notaris", "Denda & Tilang", "Biaya Pernikahan", "Keanggotaan", "Biaya Perjalanan Dinas"}},
		"Investasi":                {Expense, []string{"Emas", "Saham", "Reksadana", "Obligasi", "Deposito", "Properti", "Cryptocurrency", "Peer-to-Peer Lending"}},
		"Utang":                    {Expense, []string{"Kartu Kredit", "Pinjaman Pribadi", "Cicilan Kendaraan", "Cicilan Rumah", "Pinjaman Online", "Utang Teman/Orang Tua"}},
		"Bisnis":                   {Expense, []string{"Modal Usaha", "Biaya Operasional", "Pemasaran & Iklan", "Gaji Karyawan", "Sewa Tempat Usaha", "Peralatan & Inventaris", "Biaya Transportasi Bisnis"}},
	}

		var wg sync.WaitGroup
	parentIDs := make(map[string]uuid.UUID)

	// Insert kategori utama dengan pengecekan duplikasi
	for parentName, data := range categoriesData {
		var existingParent entity.Categories

		// Cek apakah parent category sudah ada
		if err := db.Where("name = ?", parentName).First(&existingParent).Error; err == nil {
			// fmt.Printf("Parent category %s already exists, skipping...\n", parentName)
			parentIDs[parentName] = existingParent.ID
			continue
		} else if err != gorm.ErrRecordNotFound {
			// log.Printf("Error checking parent category %s: %v\n", parentName, err)
			continue
		}

		// Jika tidak ada, insert parent
		parentID := uuid.New()
		parentIDs[parentName] = parentID // Simpan UUID untuk parent
		category := entity.Categories{
			Base: entity.Base{ID: parentID},
			Name: parentName,
			Type: data.Type,
		}

		if err := db.Create(&category).Error; err != nil {
			// log.Printf("Error inserting parent category %s: %v\n", parentName, err)
			continue
		}
		// fmt.Printf("Inserted parent category: %s with ID: %s\n", parentName, parentID)

		// Insert child categories secara paralel menggunakan goroutine
		for _, childName := range data.Items {
			wg.Add(1)
			go func(parentID uuid.UUID, childName string, categoryType entity.CategoryType) {
				defer wg.Done()

				var existingChild entity.Categories

				// Cek apakah child category sudah ada
				if err := db.Where("name = ? AND parent_id = ?", childName, parentID).First(&existingChild).Error; err == nil {
					// fmt.Printf("Child category %s (Parent: %s) already exists, skipping...\n", childName, parentID)
					return
				} else if err != gorm.ErrRecordNotFound {
					// log.Printf("Error checking child category %s: %v\n", childName, err)
					return
				}

				// Jika tidak ada, insert child
				child := entity.Categories{
					Base:     entity.Base{ID: uuid.New()},
					ParentID: &parentID,
					Name:     childName,
					Type:     categoryType,
				}

				if err := db.Create(&child).Error; err != nil {
					// log.Printf("Error inserting child category %s: %v\n", childName, err)
				} else {
					// fmt.Printf("Inserted child category: %s (Parent: %s)\n", childName, parentID)
				}
			}(parentID, childName, data.Type)
		}
	}

	// Tunggu semua goroutine selesai
	wg.Wait()
}
