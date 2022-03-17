package main

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"sync"
	"time"
)

var (
	rdb *redis.Client
)

func Init() (err error) {
	rdb = redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
	_, err = rdb.Ping().Result()
	if err != nil {
		errors.New("conn redis fail")
		fmt.Println("conn redis fail, err :", err)
		return
	}
	return nil
}

func main() {
	Init()
	if false {
		showPipeline()
		showPipelined()
		showTxPipeline()
		showTxPipelined()
		showWatch()
		showWatchDemo()
	}
}

func showWatchDemo() {
	const routineCount = 100

	// Increment 使用GET和SET命令以事务方式递增Key的值
	increment := func(key string) error {
		// 事务函数
		txf := func(tx *redis.Tx) error {
			// 获得key的当前值或零值
			n, err := tx.Get(key).Int()
			if err != nil && err != redis.Nil {
				return err
			}

			// 实际的操作代码（乐观锁定中的本地操作）
			n++

			// 操作仅在 Watch 的 Key 没发生变化的情况下提交
			_, err = tx.Pipelined(func(pipe redis.Pipeliner) error {
				// pipe handles the error case
				pipe.Set(key, n, 0)
				return nil
			})
			return err
		}

		//最多重试 maxRetries 次
		for retries := routineCount; retries > 0; retries-- {
			err := rdb.Watch(txf, key)
			if err != redis.TxFailedErr {
				return err
			}
			//优化的锁丢失
		}
		return errors.New("increment reached maximum number of retries")
	}

	//模拟 routineCount 个并发同时去修改 counter3 的值
	var wg sync.WaitGroup
	wg.Add(routineCount)
	for i := 0; i < routineCount; i++ {
		go func() {
			defer wg.Done()

			//并发安全地将redis值+1
			if err := increment("counter3"); err != nil {
				fmt.Println("increment error:", err)
			}
		}()
	}
	wg.Wait()

	n, err := rdb.Get("counter3").Int()
	fmt.Println("ended with", n, err)
}

//用来确保键的值不会被其他人修改,操作赋值
//如果被人修改删除替换,就会收到一个错误
func showWatch() (err error) {
	key := "watch_count"
	err = rdb.Watch(func(tx *redis.Tx) error {
		n, err := tx.Get(key).Int()
		if err != nil && err != redis.Nil {
			return err
		}
		tx.Pipelined(func(pipeliner redis.Pipeliner) error {
			pipeliner.Set(key, n+1, time.Minute*5)
			return nil
		})
		return err
	}, key)
	return
}

func showTxPipelined() {
	var incr *redis.IntCmd
	_, err := rdb.TxPipelined(func(pipeliner redis.Pipeliner) error {
		incr = pipeliner.Incr("tx_pipelined_counter")
		pipeliner.Expire("tx_pipelined_counter", time.Hour)
		return nil
	})
	fmt.Println(incr.Val(), err)
}

func showTxPipeline() {
	pipe := rdb.TxPipeline()

	incr := pipe.Incr("tx_pipeline_counter")
	pipe.Expire("tx_pipeline_counter", time.Hour)

	_, err := pipe.Exec()
	fmt.Println(incr.Val(), err)
}

func showPipelined() {
	var incr *redis.IntCmd
	_, err := rdb.Pipelined(func(pipeliner redis.Pipeliner) error {
		incr = pipeliner.Incr("pipelined_counter")
		pipeliner.Expire("pipelined_counter", time.Hour)
		return nil
	})
	fmt.Println(incr.Val(), err)
}

func showPipeline() {
	//不互相依赖的操作可以用pipeline 一起发送到redis,一次执行
	pipe := rdb.Pipeline()

	incr := pipe.Incr("pipeline_counter")      //原值+1
	pipe.Expire("pipeline_counter", time.Hour) //设置国企时间

	_, err := pipe.Exec()
	fmt.Println(incr.Val(), err)
}
