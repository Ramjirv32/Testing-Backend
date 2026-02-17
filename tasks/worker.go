package tasks

import (
	"backend/utils"
	"log"
)

type EmailTask struct {
	From    utils.EmailType
	To      string
	Subject string
	Body    string
}

var EmailQueue = make(chan EmailTask, 1000)

func StartEmailWorker() {
	go func() {
		for task := range EmailQueue {
			err := utils.SendEmail(task.From, task.To, task.Subject, task.Body)
			if err != nil {
				log.Printf("Failed to send async email to %s: %v", task.To, err)
			} else {
				log.Printf("Async email sent to %s", task.To)
			}
		}
	}()
}

func EnqueueEmail(from utils.EmailType, to, subject, body string) {
	select {
	case EmailQueue <- EmailTask{From: from, To: to, Subject: subject, Body: body}:
		// Enqueued
	default:
		log.Printf("Email queue full, falling back to sync sending for %s", to)
		_ = utils.SendEmail(from, to, subject, body)
	}
}
