package converter

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type CurrencyConverter struct {
	redis *redis.Client
}

func NewCurrencyConverter(redis *redis.Client) *CurrencyConverter {
	return &CurrencyConverter{redis: redis}
}

func (cc *CurrencyConverter) ConvertFromRub(currency string, amount float64) (float64, error) {
	koef, err := cc.redis.Get(currency).Result()
	if err != nil {
		return -1, err
	}
	koefFloat, err := strconv.ParseFloat(koef, 64)
	if err != nil {
		return -1, err
	}

	return amount * koefFloat, nil
}

func (cc *CurrencyConverter) ConvertToRub(currency string, amount float64) (float64, error) {
	koef, err := cc.redis.Get(currency).Result()
	if err != nil {
		return -1, err
	}
	koefFloat, err := strconv.ParseFloat(koef, 64)
	if err != nil {
		return -1, err
	}

	return amount / koefFloat, nil
}

func (cc *CurrencyConverter) UpdateInfo() {
	resp, err := http.Get(viper.GetString("currency.transferRateUrl"))
	if err != nil {
		logrus.Errorf("error in accesing currency transfer rate site error=%s", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("eror while reading currency transfer rate site, error=%s", err)
	}
	var generalStruct ModelAnswer
	err = json.Unmarshal(body, &generalStruct)
	if err != nil {
		logrus.Errorf("error in unmarshaling defalut structure, error=%s", err)
	}
	var currencyJson map[string]float64
	raw, err := json.Marshal(generalStruct.Rates)
	if err != nil {
		logrus.Errorf("error in unmarshaling currency values, error=%s", err)
	}
	json.Unmarshal(raw, &currencyJson)

	for key, value := range currencyJson {
		err := cc.redis.Set(key, value, time.Minute).Err()
		if err != nil {
			logrus.Errorf("error while setting key err=%s", err)
		}
	}
	r, _ := cc.redis.Get("AED").Result()
	logrus.Info(r)

}
