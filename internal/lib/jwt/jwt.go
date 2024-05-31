package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"vizapSSO/internal/entity"
)

func NewAccessToken(user entity.User, app entity.App, duration time.Duration) (accessToken string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.ID

	accessToken, err = token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func NewRefreshToken(lastChar string, app entity.App, duration time.Duration) (refreshToken string, err error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["last_char"] = lastChar

	refreshToken, err = token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

func ValidateToken(accessToken string) (isValid bool, uid int64, err error) {
	//TODO: исправить данную реализацию
	mySigningKey := "secret123"

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(mySigningKey), nil
	})
	if err != nil {
		return false, 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		uid, ok := claims["uid"].(float64)
		if !ok {
			return false, 0, err
		}
		return true, int64(uid), nil
	} else {
		return false, 0, err
	}
}

func CheckRefreshToken(refreshToken, accessToken string) error {
	// replace this
	mySigningKey := "secret123"

	lastChar := accessToken[len(accessToken)-6:]

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(mySigningKey), nil
	})
	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		forCheck, ok := claims["last_char"].(string)
		if !ok {
			return err
		}
		if lastChar != forCheck {
			return fmt.Errorf("invalid token pair")
		}
		return nil
	} else {
		return err
	}
}

func IdFromJWT(accessToken string) (uid int64, err error) {
	token, _, err := new(jwt.Parser).ParseUnverified(accessToken, jwt.MapClaims{})
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		floatUID, ok := claims["uid"].(float64) // Используем временную переменную для хранения float значения
		if !ok {
			return 0, fmt.Errorf("uid must be a float64, got %T", claims["uid"]) // Создаем новую ошибку, которая объясняет проблему
		}
		uid = int64(floatUID) // Приведение float64 к int64 безопасно после проверки типа
		return uid, nil
	}
	return 0, fmt.Errorf("invalid JWT claims, unable to assert to MapClaims")
}
