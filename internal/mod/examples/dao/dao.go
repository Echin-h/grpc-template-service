package dao

type user struct {
	ID   int    `gorm:"primaryKey"`
	Name string `gorm:"type:varchar(255);not null"`
	Age  int    `gorm:"type:int;not null"`
}

var (
	User = &user{}
)

//func Init(DB *gorm.DB, rds *redis.Client) error {
//	err = User.Init(DB)
//	if err != nil {
//		return err
//	}
//	return nil
//}
