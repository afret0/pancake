package user

import (
	"backend/notice"
	"backend/source"
	"backend/user/token"
	verificationCode2 "backend/user/verificationCode"
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var man *Manager

func init() {
	man = new(Manager)
	man.logger = source.Logger
	man.verifier = verificationCode2.GetVerifier()
}

type Manager struct {
	logger   *logrus.Logger
	ctx      context.Context
	verifier *verificationCode2.Verifier
}

func (m *Manager) RegisterUser(WXName, phone string) error {
	//user := new(User)
	//user.ID = man.newUserId()
	//user.Phone = phone
	//user.Name = WXName
	//user.WXName = WXName
	//user.Role = role.Customer

	//filter := bson.M{"phone": phone}
	//upt := bson.M{"$set": bson.M{"WXName": WXName, "id": m.newUserId(), "userName": WXName, "role": role.Customer}}
	//opt := new(options.UpdateOptions)
	//T := true
	//opt.Upsert = &T
	//r, err := operator.UpdateOne(m.ctx, filter, upt, opt)
	//if err != nil {
	//	m.logger.Errorln(phone, err)
	//	return err
	//}
	return nil
}

func (m *Manager) Login(phone, verificationCode string) (string, error) {
	pjt := bson.M{"_id": 1, "name": 1, "token": 1}
	user, err := m.FindOneUserByPhone(phone, pjt)
	if err != nil {
		return "", err
	}
	if !m.verifier.CheckVCode(phone, verificationCode) {
		return "", errors.New("verificationCode error")
	}
	t, err := token.GenerateToken(user.Name, user.Name)
	if t == user.Token {
		return t, nil
	}
	filter := bson.M{"phone": phone}
	upt := bson.M{"$set": bson.M{"token": t}}
	_, err = operator.UpdateOne(filter, upt)
	if err != nil {
		m.logger.Errorln(phone, err)
	}
	return t, err
}

func (m *Manager) SendVerificationCode(senderName, phone string) error {
	vCode := m.verifier.GenVerifyCode()
	m.verifier.SetVerifyCode(phone, vCode, 10)
	sender := notice.GetSender(senderName)
	return sender.SendVerificationCode(phone, vCode)
}

func (m *Manager) FindOneUserByPhone(phone string, pjt interface{}) (*User, error) {
	filter := bson.M{"phone": phone}
	opt := new(options.FindOneOptions)
	opt.Projection = pjt
	user, err := operator.FindOne(filter, opt)
	if err != nil {
		m.logger.Errorln(phone, err)
	}
	return user, err
}
