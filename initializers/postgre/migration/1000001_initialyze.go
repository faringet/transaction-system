package migration

import (
	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10/orm"
	"transaction-system/internal/domain"
)

func init() {
	models := []interface{}{
		&domain.Clients{},
		&domain.Transactions{},
		&domain.Currencies{},
	}

	migrations.MustRegister(func(db migrations.DB) error {
		for _, model := range models {
			err := db.Model(model).CreateTable(&orm.CreateTableOptions{
				IfNotExists: true,
			})
			if err != nil {
				return err
			}
		}

		_, err := db.Exec(`
			ALTER TABLE transactions
			ADD CONSTRAINT fk_transactions_clients
			FOREIGN KEY (client_id)
			REFERENCES clients(id)
			ON DELETE CASCADE;

			ALTER TABLE transactions
    		ADD CONSTRAINT fk_transactions_currencies
       		FOREIGN KEY (currency_id)
            REFERENCES currencies(id)
            ON DELETE CASCADE;
		`)
		if err != nil {
			return err
		}

		currencies := []domain.Currencies{
			{ID: 1, CurrencyCode: 840, CurrencyName: "USD"},
			{ID: 2, CurrencyCode: 978, CurrencyName: "EUR"},
			{ID: 3, CurrencyCode: 643, CurrencyName: "RUB"},
			{ID: 4, CurrencyCode: 156, CurrencyName: "CNY"},
			{ID: 5, CurrencyCode: 826, CurrencyName: "GBP"},
			{ID: 6, CurrencyCode: 392, CurrencyName: "JPY"},
			{ID: 7, CurrencyCode: 124, CurrencyName: "CAD"},
			{ID: 8, CurrencyCode: 756, CurrencyName: "CHF"},
			{ID: 9, CurrencyCode: 360, CurrencyName: "IDR"},
			{ID: 10, CurrencyCode: 356, CurrencyName: "INR"},
			{ID: 11, CurrencyCode: 410, CurrencyName: "KRW"},
			{ID: 12, CurrencyCode: 458, CurrencyName: "MYR"},
			{ID: 13, CurrencyCode: 554, CurrencyName: "NZD"},
			{ID: 14, CurrencyCode: 578, CurrencyName: "NOK"},
			{ID: 15, CurrencyCode: 352, CurrencyName: "ISK"},
			{ID: 16, CurrencyCode: 752, CurrencyName: "SEK"},
			{ID: 17, CurrencyCode: 702, CurrencyName: "SGD"},
			{ID: 18, CurrencyCode: 380, CurrencyName: "CUP"},
		}

		_, err = db.Model(&currencies).Insert()
		if err != nil {
			return err
		}

		clients := []domain.Clients{
			{ID: 3456, WalletNumber: 123456789, CardNumber: 5321300240335856},
			{ID: 2567, WalletNumber: 234567890, CardNumber: 5478396041568712},
			{ID: 1254, WalletNumber: 345678901, CardNumber: 5123876098751234},
			{ID: 8745, WalletNumber: 456789012, CardNumber: 5256789012457890},
			{ID: 3400, WalletNumber: 567890123, CardNumber: 5409123456789012},
			{ID: 1111, WalletNumber: 678901234, CardNumber: 5312467890123456},
			{ID: 1234, WalletNumber: 789012345, CardNumber: 5267890123456789},
			{ID: 9845, WalletNumber: 890123456, CardNumber: 5456347890123456},
			{ID: 8709, WalletNumber: 901234567, CardNumber: 5334789012345678},
			{ID: 1267, WalletNumber: 101234567, CardNumber: 5146901234567890},
			{ID: 1231, WalletNumber: 112345678, CardNumber: 5489012345678901},
			{ID: 9832, WalletNumber: 126456789, CardNumber: 5367901234567890},
			{ID: 7099, WalletNumber: 134567890, CardNumber: 5101234567890123},
			{ID: 6078, WalletNumber: 145678901, CardNumber: 5278345612345678},
			{ID: 5033, WalletNumber: 156789012, CardNumber: 5309123456789012},
			{ID: 9124, WalletNumber: 167890123, CardNumber: 5412789012345678},
			{ID: 8766, WalletNumber: 178901234, CardNumber: 5389012345678901},
			{ID: 9333, WalletNumber: 189012345, CardNumber: 5437890123456789},
			{ID: 7477, WalletNumber: 190123456, CardNumber: 5490123456789012},
			{ID: 3466, WalletNumber: 201234567, CardNumber: 5356789012345678},
		}

		_, err = db.Model(&clients).Insert()
		if err != nil {
			return err
		}

		return nil
	}, func(db migrations.DB) error {
		for _, model := range models {
			err := db.Model(model).DropTable(&orm.DropTableOptions{
				IfExists: true,
				Cascade:  true,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}
