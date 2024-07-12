package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"time"
)

type JournalConfig struct {
	FilenameTemplate string `json:"filename_template"`
	JournalDir       string `json:"journal_dir"`
	HeadingName      string `json:"heading_name"`
}

func NewJournalConfig(path string) (*JournalConfig, error) {
	jc := &JournalConfig{}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	by, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(by, jc)
	if err != nil {
		return nil, err
	}

	return jc, nil
}

func WriteJournal(jc *JournalConfig, mc *MemosWebhook) error {
	// create a file path to the journal file
	fn := time.Now().Format(jc.FilenameTemplate)
	p, err := filepath.Abs(path.Join(jc.JournalDir, fn))
	if err != nil {
		slog.Error("Failed to get absolute path", "error", err)
		return err
	}

	slog.Info("Updating journal", "file_name", p)

	// open the journal file
	f, err := os.OpenFile(p, os.O_RDWR|os.O_CREATE, os.FileMode(0644))
	if err != nil {
		slog.Error("Failed to create journal file", "file_path", p, "error", err)
		return err
	}
	defer f.Close()

	byt, err := io.ReadAll(f)
	if err != nil {
		slog.Error("Failed to read journal file", "file_path", p, "error", err)
		return err
	}

	// reset content
	f.Seek(0, 0)

	// write the header in if this is a new file
	if len(byt) == 0 {
		_, err = f.WriteString(fmt.Sprintf("# %s\n\n", jc.HeadingName))
		if err != nil {
			slog.Error("Failed to write heading", "error", err)
			return nil
		}
	}

	// seek to the end and write the file
	f.Seek(0, io.SeekEnd)

	_, err = f.WriteString(fmt.Sprintf("## %s\n\n", mc.CreateTime))
	if err != nil {
		slog.Error("Failed to write display time", "error", err)
		return err
	}

	_, err = f.WriteString(fmt.Sprintf("%s\n\n", mc.Memo.Content))
	if err != nil {
		slog.Error("Failed to write content", "error", err)
		return err
	}

	return nil
}
