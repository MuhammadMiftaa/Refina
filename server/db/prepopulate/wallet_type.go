package prepopulate

import (
	"sync"

	"gorm.io/gorm"

	"server/internal/entity"
)

func WalletTypesSeeder(db *gorm.DB) {
	const (
		Bank         entity.WalletType = "bank"
		EWallet      entity.WalletType = "e-wallet"
		Physical     entity.WalletType = "physical"
		OthersWallet entity.WalletType = "others"
	)
	
	// Data prepopulate
	walletTypes := []entity.WalletTypes{
		// Bank
		{Name: "BCA (Bank Central Asia)", Type: Bank},
		{Name: "Mandiri", Type: Bank},
		{Name: "BNI (Bank Negara Indonesia)", Type: Bank},
		{Name: "BRI (Bank Rakyat Indonesia)", Type: Bank},
		{Name: "CIMB Niaga", Type: Bank},
		{Name: "Bank Danamon", Type: Bank},
		{Name: "Bank Mega", Type: Bank},
		{Name: "Bank Permata", Type: Bank},
		{Name: "Bank BTPN (Jenius)", Type: Bank},
		{Name: "Bank Jago", Type: Bank},
		{Name: "Bank Syariah Indonesia (BSI)", Type: Bank},
		{Name: "OCBC NISP", Type: Bank},
		{Name: "Maybank Indonesia", Type: Bank},
		{Name: "DBS Indonesia", Type: Bank},
		{Name: "Bank Panin", Type: Bank},
		{Name: "Bank DKI", Type: Bank},
		{Name: "Bank Bukopin", Type: Bank},
		{Name: "Bank BCA Syariah", Type: Bank},
		{Name: "Bank Muamalat", Type: Bank},
		{Name: "Bank Neo Commerce", Type: Bank},
		{Name: "Bank Tabungan Negara (BTN)", Type: Bank},
		{Name: "Bank Sinarmas", Type: Bank},
		{Name: "Bank Sahabat Sampoerna", Type: Bank},
		{Name: "Bank Victoria International", Type: Bank},
		{Name: "Bank Mayora", Type: Bank},
		{Name: "Bank BJB (Bank Jabar Banten)", Type: Bank},
		{Name: "Bank Aceh Syariah", Type: Bank},
		{Name: "Bank NTB Syariah", Type: Bank},
		{Name: "Bank Nagari", Type: Bank},
		{Name: "Seabank (Bank Digital SeaGroup)", Type: Bank},
		{Name: "Blu by BCA Digital", Type: Bank},

		// E-Wallet
		{Name: "GoPay (Gojek)", Type: EWallet},
		{Name: "OVO (Grab)", Type: EWallet},
		{Name: "DANA", Type: EWallet},
		{Name: "LinkAja", Type: EWallet},
		{Name: "ShopeePay (Shopee)", Type: EWallet},
		{Name: "Jenius (BTPN)", Type: EWallet},
		{Name: "Sakuku (BCA)", Type: EWallet},
		{Name: "Doku", Type: EWallet},
		{Name: "iSaku", Type: EWallet},
		{Name: "Paytren", Type: EWallet},
		{Name: "Flip", Type: EWallet},
		{Name: "Akulaku", Type: EWallet},
		{Name: "Kredivo PayLater", Type: EWallet},
		{Name: "OY! Indonesia", Type: EWallet},
		{Name: "Octo Clicks (CIMB Niaga)", Type: EWallet},
		{Name: "Nobu (Bank Nationalnobu)", Type: EWallet},
		{Name: "Livin' by Mandiri", Type: EWallet},
		{Name: "QRIS (Sistem Pembayaran Nasional)", Type: EWallet},
		{Name: "Jago Pocket (Bank Jago)", Type: EWallet},
		{Name: "Blu by BCA Digital", Type: EWallet},
		{Name: "DBS PayLah! (DBS Indonesia)", Type: EWallet},
		{Name: "LinkAja Syariah", Type: EWallet},
		{Name: "Dana Syariah", Type: EWallet},

		// Physical
		{Name: "Dompet Fisik (Cash)", Type: Physical},
		{Name: "Celengan (Piggy Bank)", Type: Physical},
		{Name: "Safe Deposit Box", Type: Physical},

		// Others
		{Name: "Koperasi Pegawai Negeri (KPN)", Type: OthersWallet},
		{Name: "Koperasi Syariah", Type: OthersWallet},
	}

	var wg sync.WaitGroup

	for _, walletType := range walletTypes {
		wg.Add(1)

		// Jalankan goroutine untuk setiap insert
		go func(w entity.WalletTypes) {
			defer wg.Done()
			var existing entity.WalletTypes
			result := db.Where("name = ?", w.Name).First(&existing)
			if result.Error != nil {
				db.Create(&w)
				// fmt.Println("Inserted:", w.Name)
			} else {
				// fmt.Println("Already exists:", w.Name)
			}
		}(walletType)
	}

	// Tunggu semua goroutine selesai
	wg.Wait()
	// fmt.Println("Prepopulate selesai!")
}
