package rule

import (
	"github.com/gookit/validate"
	"gorm.io/gorm"
)

func GlobalRules(db *gorm.DB) {
	validate.AddValidators(validate.M{
		"exists":    NewExists(db).Passes,
		"notExists": NewNotExists(db).Passes,
		"password":  NewPassword().Passes,
		"cron":      NewCron().Passes,
		"ipcidr":    NewIPCIDR().Passes,
	})
	validate.AddGlobalMessages(map[string]string{
		"exists":    "{field} does not exist",
		"notExists": "{field} already exists",
		"password":  "Password does not meet requirements (8-20 characters, containing at least two of: letters, numbers, special characters)",
		"cron":      "Invalid Cron expression",
		"ipcidr":    "Invalid IP or CIDR format",
	})
}
