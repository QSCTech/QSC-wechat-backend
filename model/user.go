package model

import (
	"github.com/jinzhu/gorm"
)

type UserModel struct {
	gorm.Model
	ZJUid string `json:"ZJUid" gorm:"not null;unique_index;column:ZJUid"`
	INTLid string `json:"INTLid"`
	Password string `json:"password"`	// Warning
	PasswordINTL string `json:"password_intl"`	// Warning
	WechatOpenID string `json:"wechat_open_id"`
	WechatSessionID string `json:"wechat_session_id"`
}

func (user *UserModel) Create() error {
	return DB.Local.Create(&user).Error
}
func (user *UserModel) Save() error {
	return DB.Local.Save(&user).Error
}
func Delete(ZJUid string) error {
	user := UserModel{
		ZJUid: ZJUid,
	}
	return DB.Local.Delete(&user).Error
}
func GetUserByZJUid(ZJUid string) (*UserModel, error) {
	user := &UserModel{}
	d := DB.Local.Where("ZJUid = ?", ZJUid).First(&user)
	return user, d.Error
}