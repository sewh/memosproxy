package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type MemosWebhook struct {
	URL          string    `json:"url"`
	ActivityType string    `json:"activityType"`
	CreatorID    int       `json:"creatorId"`
	CreateTime   time.Time `json:"createTime"`
	Memo         struct {
		Name        string    `json:"name"`
		UID         string    `json:"uid"`
		RowStatus   string    `json:"rowStatus"`
		Creator     string    `json:"creator"`
		CreateTime  time.Time `json:"createTime"`
		UpdateTime  time.Time `json:"updateTime"`
		DisplayTime time.Time `json:"displayTime"`
		Content     string    `json:"content"`
		Nodes       []struct {
			Type          string `json:"type"`
			ParagraphNode struct {
				Children []struct {
					Type     string `json:"type"`
					TextNode struct {
						Content string `json:"content"`
					} `json:"textNode"`
				} `json:"children"`
			} `json:"paragraphNode"`
		} `json:"nodes"`
		Visibility string `json:"visibility"`
		Property   struct {
		} `json:"property"`
	} `json:"memo"`
}

func (mwh MemosWebhook) String() string {
	return mwh.Memo.Content
}

func HandleMemosBuilder(jc *JournalConfig) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		slog.Info("Received request to handle a memos event...")

		if r.Method != http.MethodPost {
			slog.Error("Invalid request method", "method", r.Method)
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		b, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("Failed to read request body", "error", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		wh := &MemosWebhook{}
		if err := json.Unmarshal(b, wh); err != nil {
			slog.Error("Failed to unmarshal request body", "error", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		if wh.ActivityType != "memos.memo.created" {
			slog.Error("Invalid activity type", "activity_type", wh.ActivityType)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := WriteJournal(jc, wh); err != nil {
			slog.Error("Failed to write journal", "error", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.WriteHeader(http.StatusOK)
	}
}
