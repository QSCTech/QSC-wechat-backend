package model

import (
	"time"
)

type UserModel struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	ZJUid string `json:"ZJUid" gorm:"not null;unique_index;column:ZJUid"`
	Password string `json:"password"`	// Warning

	INTLid string `json:"INTLid" gorm:"column:INTLid"`
	PasswordINTL string `json:"password_intl"`	// Warning

	PRINTid string `json:"PRINTid" gorm:"column:PRINTid"`
	PasswordPRINT string `json:"password_print"`

	WechatOpenID string `json:"wechat_open_id" gorm:"unique_index"`
	WechatSessionID string `json:"wechat_session_id"`
}

func (user *UserModel) Create() error {
	return DB.Local.Create(&user).Error
}
func (user *UserModel) Save() error {
	return DB.Local.Save(&user).Error
}
func DeleteZJU(ZJUid string) error {
	_, err := GetUserByZJUid(ZJUid)
	if err != nil {
		// Don't has user
		return nil
	}
	return DB.Local.Where("ZJUid = ?", ZJUid).Delete(&UserModel{}).Error
}
func DeleteWechat(OpenID string) error {
	_, err := GetUserByWechatID(OpenID)
	if err != nil {
		// User not exist
		return nil
	}
	return DB.Local.Where("wechat_open_id = ?", OpenID).Delete(&UserModel{}).Error
}
func GetUserByZJUid(ZJUid string) (*UserModel, error) {
	user := &UserModel{}
	d := DB.Local.Where("ZJUid = ?", ZJUid).First(&user)
	return user, d.Error
}
func GetUserByWechatID(OpenID string) (*UserModel, error) {
	user := &UserModel{}
	d := DB.Local.Where("wechat_open_id = ?", OpenID).First(&user)
	return user, d.Error
}