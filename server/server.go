package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/mmcdole/gofeed"

	"github.com/yseto/podcaster/client"
	"github.com/yseto/podcaster/ent"
	"github.com/yseto/podcaster/ent/entries"
	"github.com/yseto/podcaster/ent/feeds"
	"github.com/yseto/podcaster/ent/users"
)

type server struct {
	ent *ent.Client
}

func NewServer(ent *ent.Client) *server {
	return &server{ent: ent}
}

var _ StrictServerInterface = (*server)(nil)

func (s *server) Index(ctx context.Context, request IndexRequestObject) (IndexResponseObject, error) {
	// fmt.Sprint(userFromContext(ctx))

	f, err := os.Open("public/index.html")
	if err != nil {
		return nil, err
	}

	return Index200TexthtmlResponse{
		Body: f,
	}, nil
}

func (s *server) IndexFile(ctx context.Context, request IndexFileRequestObject) (IndexFileResponseObject, error) {
	f, err := os.Open("public/app.js")
	if err != nil {
		return nil, err
	}

	return IndexFile200TextjavascriptResponse{
		Body: f,
	}, nil
}

func (s *server) GetEntries(ctx context.Context, request GetEntriesRequestObject) (GetEntriesResponseObject, error) {
	feed, err := s.ent.Feeds.Get(ctx, request.Id)
	if err != nil {
		if ent.IsNotFound(err) {
			return GetEntries404Response{}, nil
		}
		return GetEntries400Response{}, err
	}

	entries, err := feed.QueryEntries().All(ctx)
	if err != nil {
		return nil, err
	}

	var resp []Entry
	for _, entry := range entries {
		resp = append(resp, Entry{
			ID:          uint64(entry.ID),
			Description: entry.Description,
			Title:       entry.Title,
			Url:         entry.URL,
			New:         entry.New,
			PublishedAt: entry.PublishedAt.Format(time.RFC3339),
		})
	}

	return GetEntries200JSONResponse(resp), nil
}

type entryItem struct {
	url         string
	description string
	title       string
	publishedAt time.Time
}

func (s *server) RegisterSubscription(ctx context.Context, request RegisterSubscriptionRequestObject) (RegisterSubscriptionResponseObject, error) {
	tx, err := s.ent.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	f, err := parseFeed(request.Body.Url)
	if err != nil {
		return nil, err
	}

	var items []entryItem
	for _, item := range f.Items {
		// if item.ITunesExt != nil {
		// 	fmt.Println(item.ITunesExt.Duration)
		// }

		for _, enclosure := range item.Enclosures {
			if !strings.HasPrefix(enclosure.Type, "audio/") {
				continue
			}
			publishedAt := time.Now()
			if item.PublishedParsed != nil {
				publishedAt = *item.PublishedParsed
			}
			items = append(items, entryItem{
				url:         enclosure.URL,
				description: item.Description,
				title:       item.Title,
				publishedAt: publishedAt,
			})
		}
	}

	// if f.ITunesExt != nil {
	// 	fmt.Println(f.ITunesExt.Image)
	// }

	feed, err := tx.Feeds.
		Create().
		SetTitle(f.Title).
		SetURL(request.Body.Url).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	u, err := tx.Users.Query().Where(users.ID(int(userFromContext(ctx)))).First(ctx)
	if err != nil {
		return nil, err
	}

	u.Update().AddFeeds(feed).Save(ctx)

	for _, item := range items {
		entry, err := tx.Entries.Create().
			SetURL(item.url).
			SetDescription(item.description).
			SetTitle(item.title).
			SetPublishedAt(item.publishedAt).
			Save(ctx)
		if err != nil {
			log.Fatal(err)
		}

		if _, err = feed.Update().AddEntries(entry).Save(ctx); err != nil {
			log.Fatal(err)
		}
	}

	defer tx.Commit()

	return RegisterSubscription200JSONResponse(Subscription{
		ID:    uint64(feed.ID),
		Title: feed.Title,
		Url:   feed.URL,
	}), nil
}

func parseFeed(url string) (*gofeed.Feed, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	feed, err := gofeed.NewParser().Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	return feed, nil
}

func (s *server) DeleteSubscription(ctx context.Context, request DeleteSubscriptionRequestObject) (DeleteSubscriptionResponseObject, error) {
	tx, err := s.ent.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	n, err := tx.Feeds.Delete().Where(feeds.IDEQ(request.Id)).Exec(ctx)
	if err != nil {
		return DeleteSubscription400Response{}, nil
	}

	if n == 0 {
		return DeleteSubscription404Response{}, nil
	}

	defer tx.Commit()
	return DeleteSubscription204Response{}, nil
}

func (s *server) FetchSubscription(ctx context.Context, request FetchSubscriptionRequestObject) (FetchSubscriptionResponseObject, error) {
	tx, err := s.ent.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	feed, err := tx.Feeds.Get(ctx, request.Id)
	if ent.IsNotFound(err) {
		return FetchSubscription404Response{}, nil
	}
	if err != nil {
		return nil, err
	}

	latestEntry, err := feed.QueryEntries().Order(entries.ByPublishedAt(sql.OrderDesc())).First(ctx)
	if err != nil {
		return nil, err
	}

	parsedFeed, err := parseFeed(feed.URL)
	if err != nil {
		return nil, err
	}

	var items []entryItem
	for _, item := range parsedFeed.Items {
		// if item.ITunesExt != nil {
		// 	fmt.Println(item.ITunesExt.Duration)
		// }

		if item.PublishedParsed != nil && item.PublishedParsed.Before(latestEntry.PublishedAt) {
			continue
		}

		for _, enclosure := range item.Enclosures {
			if !strings.HasPrefix(enclosure.Type, "audio/") {
				continue
			}
			publishedAt := time.Now()
			if item.PublishedParsed != nil {
				publishedAt = *item.PublishedParsed
			}
			items = append(items, entryItem{
				url:         enclosure.URL,
				description: item.Description,
				title:       item.Title,
				publishedAt: publishedAt,
			})
		}
	}

	for _, item := range items {
		_, err := tx.Entries.Create().
			SetURL(item.url).
			SetDescription(item.description).
			SetTitle(item.title).
			SetPublishedAt(item.publishedAt).
			SetFeeds(feed).
			Save(ctx)

		if err != nil && ent.IsConstraintError(err) {
			continue
		}
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return FetchSubscription202Response{}, nil
}

func (s *server) Subscriptions(ctx context.Context, request SubscriptionsRequestObject) (SubscriptionsResponseObject, error) {
	u, err := s.ent.Users.Query().Where(users.ID(int(userFromContext(ctx)))).First(ctx)
	if err != nil {
		return nil, err
	}

	// https://entgo.io/ja/docs/feature-flags/#modify-example-4
	/*
		SELECT `feeds`.`id`, `feeds`.`title`, `feeds`.`url`, COUNT(`t1`.`new`OR NULL) AS `new_count`
		FROM `feeds`
		LEFT JOIN `entries` AS `t1` ON `feeds`.`id` = `t1`.`feeds_entries`
		WHERE `users_feeds` = ? GROUP BY `feeds`.`id` args=[?]
	*/

	var values []struct {
		ent.Feeds
		NewCount int `sql:"new_count"`
	}

	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("recovered from: %v", rec)
		}
	}()

	u.QueryFeeds().
		Modify(func(s *sql.Selector) {
			t := sql.Table(entries.Table)
			s.LeftJoin(t).On(
				s.C(feeds.FieldID),
				t.C(entries.FeedsColumn),
			).
				AppendSelect(
					sql.As(sql.Count(t.C(entries.FieldNew)+" OR NULL"), "new_count"),
				).
				GroupBy(s.C(feeds.FieldID))
		}).ScanX(ctx, &values)

	var respFeeds []Subscription
	for _, feed := range values {
		respFeeds = append(respFeeds, Subscription{
			ID:            uint64(feed.ID),
			Title:         feed.Title,
			Url:           feed.URL,
			NewEntryCount: feed.NewCount,
		})
	}

	return Subscriptions200JSONResponse(respFeeds), nil
}

func (s *server) OpenedEntry(ctx context.Context, request OpenedEntryRequestObject) (OpenedEntryResponseObject, error) {
	tx, err := s.ent.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	feed, err := tx.Feeds.Get(ctx, request.Id)
	if ent.IsNotFound(err) {
		return OpenedEntry404Response{}, nil
	}
	if err != nil {
		return nil, err
	}

	entry, err := feed.QueryEntries().Where(entries.IDEQ(request.EntryId)).First(ctx)
	if ent.IsNotFound(err) {
		return OpenedEntry404Response{}, nil
	}
	if err != nil {
		return nil, err
	}

	_, err = entry.Update().SetNew(false).Save(ctx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return OpenedEntry202Response{}, nil
}

func (s *server) DeleteEntry(ctx context.Context, request DeleteEntryRequestObject) (DeleteEntryResponseObject, error) {
	tx, err := s.ent.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	feed, err := tx.Feeds.Get(ctx, request.Id)
	if ent.IsNotFound(err) {
		return DeleteEntry404Response{}, nil
	}
	if err != nil {
		return nil, err
	}

	entry, err := feed.QueryEntries().Where(entries.IDEQ(request.EntryId)).First(ctx)
	if ent.IsNotFound(err) {
		return DeleteEntry404Response{}, nil
	}
	if err != nil {
		return nil, err
	}

	_, err = tx.Entries.Delete().Where(entries.IDEQ(entry.ID)).Exec(ctx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return DeleteEntry202Response{}, nil
}
