package poohchanbot

import (
	"errors"
	"log"
	"math/rand"
	"os"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
	"github.com/line/line-bot-sdk-go/linebot"
)

const (
	ImageRootPath = "/images"
)

func replyRandomImage(string) *sendingMessageQueue {
	token := os.Getenv("DROPBOX_TOKEN")
	config := dropbox.Config{
		Token: token,
	}

	cli := files.New(config)
	lfarg := files.NewListFolderArg(ImageRootPath)

	res, err := cli.ListFolder(lfarg)
	if err != nil {
		log.Print(err)
		return nil
	}

	r := rand.Intn(len(res.Entries))
	file, ok := res.Entries[r].(*files.FileMetadata)
	if !ok {
		log.Print(errors.New("type assartion error"))
		return nil
	}

	path := file.PathLower

	gtlArg := files.NewGetTemporaryLinkArg(path)
	gtlResult, err := cli.GetTemporaryLink(gtlArg)
	if err != nil {
		log.Print(err)
		return nil
	}

	link := gtlResult.Link
	q := newSendingMessageQueue()
	q.enque(linebot.NewImageMessage(link, link))

	return q
}
