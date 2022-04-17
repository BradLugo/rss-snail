package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/JoeyDrew/rss-snail/internal"

	"github.com/mmcdole/gofeed"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.uber.org/zap"
)

type rssFeed struct {
	feed *gofeed.Feed
	To   string
	From string

	Logger *zap.Logger
}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize logger: %v", err)
	}
	if logger == nil {
		log.Fatal("null logger")
	}

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			logger.Error("couldn't close logger", zap.Error(err))
		}
	}(logger)

	s, err := internal.NewSqliteDao("identifier.sqlite")
	if err != nil {
		logger.Fatal("cant provision", zap.Error(err))
	}
	if s != nil {
		defer s.Close()
	}

	return

	var feed = NewFeed(logger)
	if feed == nil {
		logger.Fatal("couldn't initialize new feed")
	}
	feed.To = os.Getenv("SENDGRID_TO")
	feed.From = os.Getenv("SENDGRID_FROM")

	sendMail(feed)
}

func NewFeed(logger *zap.Logger) *rssFeed {
	fp := gofeed.NewParser()
	newfeed, err := fp.ParseURL("https://sanantonioreport.org/feed/")
	if err != nil {
		logger.Error("couldn't parse RSS URL", zap.Error(err))
		return nil
	}

	return &rssFeed{
		feed: newfeed,

		Logger: logger,
	}
}

func sendMail(rf *rssFeed) {
	from := mail.NewEmail("RSS Snail", rf.From)
	subject := "RSS Digest"
	to := mail.NewEmail("", rf.To)
	plainTextContent := rf.feed.Title

	var htmlContentBuilder strings.Builder
	htmlContentBuilder.WriteString("<h1>RSS Snail</h1>")
	htmlContentBuilder.WriteString(fmt.Sprintf("<h2>%s</h2>", rf.feed.Title))

	feedItems := rf.feed.Items
	if feedItems == nil {
		rf.Logger.Warn("feed doesn't have any items", zap.String("Feed title", rf.feed.Title))
	} else {
		if len(feedItems) > 5 {
			feedItems = feedItems[:5]
		}
		htmlContentBuilder.WriteString(fmt.Sprintf("<h3>Latest %d articles</h3>", len(feedItems)))
		for _, item := range feedItems {
			htmlContentBuilder.WriteString(fmt.Sprintf("<p><a href=\"%s\">%s</a></p>", item.Link, item.Title))
		}
	}

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContentBuilder.String())
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))

	response, err := client.Send(message)
	if err != nil {
		rf.Logger.Error("couldn't send email", zap.Error(err))
	} else {
		responseStatusCode := response.StatusCode

		// Abstract structure logging data, so we don't need to repeat it
		loggerFields := zap.Fields(
			zap.Int("Response status code", responseStatusCode),
			zap.String("Response status text", http.StatusText(responseStatusCode)),
			zap.String("Response body", response.Body),
			zap.Reflect("Response headers", response.Headers),
		)
		responseLogger := rf.Logger.WithOptions(loggerFields)

		if responseStatusCode >= 200 && responseStatusCode < 300 {
			responseLogger.Info("Successfully sent email")
		} else if responseStatusCode >= 300 && responseStatusCode < 400 {
			responseLogger.Warn("Request was redirected")
		} else if responseStatusCode >= 400 && responseStatusCode < 500 {
			responseLogger.Error("Client error, check response body for more info")
		} else if responseStatusCode >= 500 && responseStatusCode < 600 {
			responseLogger.Error("Client error, check response body for more info")
		} else {
			responseLogger.Warn("Unknown response code, review status code and response body for more info")
		}
	}
}
