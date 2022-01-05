package env

import "github.com/spf13/viper"

func GetWithDefault[T any](key string, value T) T {
	viper.SetDefault(key, value)
	if v, ok := viper.Get(key).(T); ok {
		return v
	}

	return value
}
