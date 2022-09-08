package model

type Place struct {
	BaseModel
	// 设置使用ID和CityCode作为联合的唯一索引
	ID        string `gorm:"uniqueIndex:idx_places_unique_identity"`
	Name      string
	CityName  string
	Addr      string
	MingE     string
	Condition string
	Method    string
	Tel       string
	OrderId   string
	YYTime    string
	Course    string
	CityCode  string `gorm:"uniqueIndex:idx_places_unique_identity"`
}
