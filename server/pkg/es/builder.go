package es

import (
	"goweb/author-admin/server/pkg/setting"
	"goweb/author-admin/server/pkg/util"
	"log"
	"time"
)

type Builder struct {
	Scanners []Scanner

	// Temp file stores timestamp.
	Path string

	// The timestamp of last scan of DB.
	Timestamp int

	// Scan period.
	// It is loaded from config file.
	Period int

	Ticker *time.Ticker
}

func (b *Builder) SaveTimestamp() {
	util.SaveGob(b.Path, b.Timestamp)
}

// RunOnce has 2 parts.
// 1st: create if not exist; update if exist. Use Index api.
// 2rd: delete docs.
func (b *Builder) RunOnce() {
	newTimestamp := util.CurrentTimestamp()
	var err error
	for _, sc := range b.Scanners {
		// process update
		docs := sc.ScanUpdate(b.Timestamp, newTimestamp)
		if docs != nil {
			err = IndexDocBulk(sc.IndexName(), docs, sc.Depth())
		}
		// process delete
	}

	if err != nil {
		log.Println("Failed to update indices: ", err)
		return
	}

	b.Timestamp = newTimestamp
	b.SaveTimestamp()
}

func (b *Builder) run() {
	for {
		select {
		case <-b.Ticker.C:
			log.Println("Indices builder runs once at ", time.Now().Format("2006-01-02 15:04:05"))
			b.RunOnce()
		default:
		}
	}
}

func (b *Builder) Start() {
	go b.run()
}

func NewBuilder(scs []Scanner) *Builder {
	builder := &Builder{}

	if scs == nil {
		log.Println("No indices need to build.")
		return nil
	}
	builder.Scanners = scs

	builder.Path = "./var/timestamp.gob"
	builder.Period = setting.DuraSecond

	if util.FileExist(builder.Path) {
		x := new(int)
		err := util.LoadGob(builder.Path, x)
		if err != nil {
			log.Println("Failed to decode timestamp file: ", err)
			return nil
		}
		builder.Timestamp = *x
	} else {
		log.Println("Timestamp file does not exist. It will be created and init with 0.")
		builder.Timestamp = 0
		_, err := util.SaveGob(builder.Path, 0)
		if err != nil {
			log.Println("Failed to create timestamp file: ", err)
		}
	}

	builder.Ticker = time.NewTicker(time.Second * time.Duration(builder.Period))

	return builder
}

func Build(scs []Scanner) {
	builder := NewBuilder(scs)
	builder.Start()
}
