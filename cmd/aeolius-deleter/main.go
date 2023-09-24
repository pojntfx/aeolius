package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/pojntfx/aeolius/pkg/bluesky"
)

func main() {
	pdsURL := flag.String("pds-url", "https://bsky.social", "PDS URL")
	username := flag.String("username", "example.bsky.social", "Bluesky username")
	password := flag.String("password", "", "Bluesky password, preferably an app password (get one from https://bsky.app/settings/app-passwords)")
	postTTL := flag.Int("post-ttl", 3, "Maximum post age before considering it for deletion")
	cursorFlag := flag.String("cursor", "", "Cursor from which point forwards posts should be considered for deletion")
	batchSize := flag.Int("batch-size", 100, "How many posts to read at a time")
	limit := flag.Int("limit", 10, "Maximum amount of batches of posts to read/delete")

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	auth := &xrpc.AuthInfo{}

	client := &xrpc.Client{
		Client: http.DefaultClient,
		Host:   *pdsURL,
		Auth:   auth,
	}

	session, err := atproto.ServerCreateSession(ctx, client, &atproto.ServerCreateSession_Input{
		Identifier: *username,
		Password:   *password,
	})
	if err != nil {
		panic(err)
	}

	auth.AccessJwt = session.AccessJwt
	auth.RefreshJwt = session.RefreshJwt
	auth.Handle = session.Handle
	auth.Did = session.Did

	recordsToDelete, cursor, err := bluesky.GetPostsToDelete(client, *postTTL, *cursorFlag, *batchSize, *limit)
	if err != nil {
		panic(err)
	}

	log.Println("Deleting", recordsToDelete)

	log.Println("Setting refresh JWT to <redacted> and cursor to", cursor, "in database")
}
