package consumer

import (
	"context"
	"encoding/json"
	"log"
	"practice/internal/pkg/config"
	repoComp "practice/internal/repository/mongodb/computer"
	repoUser "practice/internal/repository/postgres/user"
	"practice/internal/service/computer"
	"practice/internal/service/user"
)

func ConsumeCreateUser(cfg *config.Config, svc user.ServiceUser) func(message []byte) {
	return func(message []byte) {
		log.Printf("Received message from topic %s: %s\n", cfg.KAFKA_TOPIC_USER_CREATED, string(message))

		var req repoUser.User
		if err := json.Unmarshal(message, &req); err != nil {
			log.Printf("error while unmarshalling user: %v\n", err)
		}

		log.Printf("Received user data for insertion: %v\n", &req)

		resp, err := svc.Create(context.Background(), &req)
		if err != nil {
			log.Printf("error while creating user: %v\n", err)
		}

		log.Printf("Created user: %v\n", resp)
	}
}

func ConsumeUpdateUser(cfg *config.Config, svc user.ServiceUser) func(message []byte) {
	return func(message []byte) {
		log.Printf("Received message from topic %s: %s\n", cfg.KAFKA_TOPIC_USER_UPDATED, string(message))

		var req repoUser.User
		if err := json.Unmarshal(message, &req); err != nil {
			log.Printf("error while unmarshalling user: %v\n", err)
		}

		log.Printf("Received user data for update: %v\n", &req)

		resp, err := svc.Update(context.Background(), &req)
		if err != nil {
			log.Printf("error while updating user: %v\n", err)
		}

		log.Printf("Updated user: %v\n", resp)
	}
}

func ConsumeDeleteUser(cfg *config.Config, svc user.ServiceUser) func(message []byte) {
	return func(message []byte) {
		log.Printf("Received message from topic %s: %s\n", cfg.KAFKA_TOPIC_USER_DELETED, string(message))

		resp, err := svc.Delete(context.Background(), string(message))
		if err != nil {
			log.Printf("error while deleting user: %v\n", err)
		}

		log.Printf("Deleted user id: %v\n", resp)
	}
}

func ConsumeCreateComputer(cfg *config.Config, svc computer.ServiceComputer) func(message []byte) {
	return func(message []byte) {
		log.Printf("Received message from topic %s: %s\n", cfg.KAFKA_TOPIC_COMPUTER_CREATED, string(message))

		var req repoComp.Computer
		if err := json.Unmarshal(message, &req); err != nil {
			log.Printf("error while unmarshalling computer: %v\n", err)
		}

		log.Printf("Received computer data for insertion: %v\n", &req)

		resp, err := svc.Create(context.Background(), &req)
		if err != nil {
			log.Printf("error while creating computer: %v\n", err)
		}

		log.Printf("Created computer: %v\n", *resp)
	}
}

func ConsumeUpdateComputer(cfg *config.Config, svc computer.ServiceComputer) func(message []byte) {
	return func(message []byte) {
		log.Printf("Received message from topic %s: %s\n", cfg.KAFKA_TOPIC_COMPUTER_UPDATED, string(message))

		var req repoComp.Computer
		if err := json.Unmarshal(message, &req); err != nil {
			log.Printf("error while unmarshalling computer: %v\n", err)
		}

		log.Printf("Received computer data for update: %v\n", &req)

		resp, err := svc.Update(context.Background(), &req)
		if err != nil {
			log.Printf("error while updating computer: %v\n", err)
		}

		log.Printf("Updated computer: %v\n", resp)
	}
}

func ConsumeDeleteComputer(cfg *config.Config, svc computer.ServiceComputer) func(message []byte) {
	return func(message []byte) {
		log.Printf("Received message from topic %s: %s\n", cfg.KAFKA_TOPIC_COMPUTER_DELETED, string(message))

		resp, err := svc.Delete(context.Background(), string(message))
		if err != nil {
			log.Printf("error while deleting computer: %v\n", err)
		}

		log.Printf("Deleted computer id: %v\n", resp)
	}
}
