package utils

import (
    "context"
    uuid "github.com/satori/go.uuid"
    "google.golang.org/grpc/metadata"
)

func NewUUID() string {
    return uuid.NewV4().String()
}

func NewContext(uuid string, userName string) context.Context {
    if uuid == "" {
        uuid = NewUUID()
    }

    md := metadata.Pairs("session-id", uuid, "user-name", userName)
    return metadata.NewOutgoingContext(context.Background(), md)
}
