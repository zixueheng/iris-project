/*
 * @Description: The program is written by the author, if modified at your own risk.
 * @Author: heyongliang
 * @Email: 356126067@qq.com
 * @Phone: 15215657185
 * @Date: 2023-02-25 10:34:06
 * @LastEditTime: 2023-02-25 10:43:18
 */
package formtoken

import (
	"context"
	"encoding/json"
	"errors"
	"iris-project/lib/cache"
	"iris-project/lib/lock"
	"log"
	"time"

	"github.com/google/uuid"
)

// FormToken 表单令牌，用于防止表单重复提交
type FormToken struct {
	Token    string `json:"token"`
	Deadline int64  `json:"dead_line"`
}

const (
	tokenPrefix = "form_token_"
	tokenTTL    = 5
)

func GetFormToken() string {
	var (
		token     = uuid.New().String()
		deadline  = time.Now().Add(tokenTTL * time.Minute)
		formToken = &FormToken{Token: token, Deadline: deadline.Unix()}
	)
	bytes, _ := json.Marshal(formToken)
	cache.Set(tokenPrefix+token, string(bytes), time.Duration(tokenTTL*time.Minute))
	return token
}

func CheckFormToken(token string) error {
	value, err := cache.Get(tokenPrefix + token)
	if err != nil {
		return err
	}
	defer cache.Del(tokenPrefix + token)
	formToken := &FormToken{}
	err = json.Unmarshal([]byte(value), formToken)
	if err != nil {
		return err
	}
	if formToken == nil || formToken.Token != token {
		return errors.New("formToken is null or invalid")
	}
	return nil
}

// Access 表单令牌 FormToken 验证
func Access(ctx context.Context, token string) bool {
	rl := lock.NewRedisLockWithCtx(ctx, token)
	if err := rl.Lock(); err != nil {
		log.Println(err.Error())
		return false
	}
	defer rl.Unlock()
	if err := CheckFormToken(token); err != nil {
		// log.Printf("TIME: %d，进程执行失败：%s", time.Now().UnixNano(), err.Error())
		return false
	} else {
		// log.Printf("TIME: %d，进程执行成功", time.Now().UnixNano())
		return true
	}
}
