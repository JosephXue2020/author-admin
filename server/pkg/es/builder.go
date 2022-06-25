package es

import (
	"goweb/author-admin/server/pkg/setting"
	"goweb/author-admin/server/pkg/util"
	"log"
	"time"
)

type Builder struct {
	IDs []*Index
	// Temp file stores timestamp.
	Path string
	// The timestamp of last scan of DB.
	Timestamp int
	// Scan period.
	// The unit is second.
	// It is loaded from config file.
	Period int
	// The unit is nanosecond.
	Overlap int
	Ticker  *time.Ticker
	// worker number of bulk operation.
	WorkerNum int
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
	for _, idx := range b.IDs {
		scanner := idx.IndexScanner
		// process update
		scs := scanner.ScanUpdate(b.Timestamp-b.Overlap, newTimestamp)
		if scs != nil {
			var docs []*Doc
			for _, sc := range scs {
				doc := NewDocFromScanner(sc)
				docs = append(docs, doc)
			}
			err = IndexDocBulk(docs, b.WorkerNum)
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

func NewDefaultBuilder(ids []*Index) *Builder {
	builder := &Builder{}

	if ids == nil {
		log.Println("No indices need to build.")
		return nil
	}
	builder.IDs = ids

	builder.Path = "./var/timestamp.gob"
	builder.Period = setting.DuraSecond
	builder.Overlap = 1000
	builder.WorkerNum = 10

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

func Build(ids []*Index) {
	builder := NewDefaultBuilder(ids)
	builder.Start()
}
