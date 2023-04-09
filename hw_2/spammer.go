package main

import (
	"fmt"
	"os"
	"sort"
	"sync"
)

func RunPipeline(cmds ...cmd) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	chans := make([]chan any, 1+len(cmds))
	for i := range chans {
		chans[i] = make(chan any)
	}

	for i, f := range cmds {
		wg.Add(1)
		go func(f cmd, in, out chan any) {
			defer wg.Done()
			defer close(out)

			f(in, out)
		}(f, chans[i], chans[i+1])
	}
}

func SelectUsers(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	userIDs := &sync.Map{}

	for data := range in {
		wg.Add(1)
		go func(email string) {
			defer wg.Done()

			user := GetUser(email)
			if _, exists := userIDs.LoadOrStore(user.ID, true); !exists {
				out <- user
			}
		}(data.(string))
	}
}

func SelectMessages(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	sendBatch := func(batch []User) {
		defer wg.Done()

		results, err := GetMessages(batch...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, result := range results {
			out <- result
		}
	}

	batch := []User{}
	for data := range in {
		batch = append(batch, data.(User))
		if len(batch) == GetMessagesMaxUsersBatch {
			wg.Add(1)
			go sendBatch(batch)
			batch = []User{}
		}
	}
	if len(batch) > 0 { // если начали заполнять, но не заполнили полностью
		wg.Add(1)
		go sendBatch(batch)
	}
}

func CheckSpam(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	type none struct{}
	limiter := make(chan none, HasSpamMaxAsyncRequests)

	for data := range in {
		limiter <- none{}

		wg.Add(1)
		go func(msgID MsgID) {
			defer wg.Done()

			hasSpam, err := HasSpam(msgID)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			<-limiter

			out <- MsgData{
				ID:      msgID,
				HasSpam: hasSpam,
			}
		}(data.(MsgID))
	}
}

func CombineResults(in, out chan interface{}) {
	messages := []MsgData{}
	for data := range in {
		messages = append(messages, data.(MsgData))
	}
	sort.Slice(messages, func(i, j int) bool {
		if messages[i].HasSpam != messages[j].HasSpam {
			return messages[i].HasSpam
		}
		return messages[i].ID < messages[j].ID
	})

	for _, message := range messages {
		out <- fmt.Sprint(message.HasSpam, message.ID)
	}
}
