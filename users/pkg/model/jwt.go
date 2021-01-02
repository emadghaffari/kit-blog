package model

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	jjwt "github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"

	"github.com/emadghaffari/kit-blog/users/config"
	"github.com/emadghaffari/kit-blog/users/pkg/redis"
)

// jwt struct
type jwt struct {
	AccessToken  string `json:"at"`
	RefreshToken string `json:"rt"`
	AccessUUID   string `json:"uuid"`
	RefreshUUID  string `json:"rau"`
	AtExpires    int64  `json:"exp"`
	RtExpires    int64  `json:"rexp"`
}

var (
	// Conf variable instance of intef
	Conf intef = &wt{}
)

type intef interface {
	Generate(data User) (*jwt, error)
}
type wt struct{}

func (j *wt) Generate(user User) (*jwt, error) {

	td, err := j.genJWT()
	if err != nil {
		return nil, err
	}

	if err := j.genRefJWT(td); err != nil {
		return nil, err
	}

	if err := j.redis(user, td); err != nil {
		return nil, err
	}

	return td, nil
}

func (j *wt) genJWT() (*jwt, error) {
	// create new jwt
	td := &jwt{}
	td.AtExpires = time.Now().Add(time.Duration(time.Minute * viper.GetDuration("jwt.expire"))).Unix()
	td.RtExpires = time.Now().Add(time.Duration(time.Minute * viper.GetDuration("jwt.RTexpire"))).Unix()
	td.AccessUUID = hasher(30)
	td.RefreshUUID = hasher(60)

	// New MapClaims for access token
	atClaims := jjwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["uuid"] = td.AccessUUID
	atClaims["exp"] = td.AtExpires
	at := jjwt.NewWithClaims(jjwt.SigningMethodHS256, atClaims)

	var err error
	td.AccessToken, err = at.SignedString([]byte(config.Confs.JWT.Secret))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func (j *wt) genRefJWT(td *jwt) error {
	// New MapClaims for refresh access token
	rtClaims := jjwt.MapClaims{}
	rtClaims["uuid"] = td.RefreshUUID
	rtClaims["exp"] = td.RtExpires
	rt := jjwt.NewWithClaims(jjwt.SigningMethodHS256, rtClaims)

	var err error
	td.RefreshToken, err = rt.SignedString([]byte(config.Confs.JWT.RSecret))
	if err != nil {
		return err
	}
	return nil
}

func (j *wt) redis(user User, td *jwt) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	// make map for store in redis
	us := make(map[string]interface{}, 4)
	us["id"] = user.ID
	us["username"] = user.Username
	us["email"] = user.Email
	us["phone"] = user.Phone

	bt, err := json.Marshal(us)
	if err != nil {
		return err
	}
	if err := redis.DB.GetDB().Set(context.Background(), td.AccessUUID, string(bt), at.Sub(now)).Err(); err != nil {
		return err
	}

	if err := redis.DB.GetDB().Set(context.Background(), td.RefreshUUID, string(bt), rt.Sub(now)).Err(); err != nil {
		return err
	}
	return nil
}

// Generate hash key
func hasher(lenght int) string {
	letters := []int32("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789-&()_")
	rand.Seed(time.Now().UnixNano())
	b := make([]int32, lenght)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
