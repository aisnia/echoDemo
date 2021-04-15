package email

import (
	"fmt"
	"github.com/jordan-wright/email"
	"io"
	"learn_together/initer"
	"net/smtp"
	"sync"
	"time"
)

var ch chan *email.Email

func InitEmial(config *initer.Config, logFile io.Writer) {
	//最多先发10个
	ch = make(chan *email.Email, 10)
	var num = config.Email.PoolSize
	p, err := email.NewPool(
		config.Email.Addr,
		num,
		smtp.PlainAuth("", config.Email.Username, config.Email.Password, config.Email.Host))
	if err != nil {
		panic("failed to create pool:" + err.Error())
	}

	//直接监听邮件chan 并且发送
	var wg sync.WaitGroup
	wg.Add(num)
	for i := 0; i < num; i++ {
		go func() {
			defer wg.Done()
			for e := range ch {
				err := p.Send(e, 10*time.Second)
				if err != nil {
					fmt.Fprintf(logFile, "email:%v send error:%v\n", e, err)
				}
			}
		}()
	}
}

func Send(e *email.Email) {
	ch <- e
}
