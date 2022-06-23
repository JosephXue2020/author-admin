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

func (b *Builder) SingleRun() {
	newTimestamp := util.CurrentTimestamp()
	var err error
	for _, sc := range b.Scanners {
		docs := sc.ScanUpdate(b.Timestamp, newTimestamp)
		if docs != nil {
			for _, doc := range docs {
				err = CreateDoc(sc.IndexName(), doc, sc.Depth())
			}

		}
	}

	if err == nil {
		b.Timestamp = newTimestamp
		b.SaveTimestamp()
	} else {
		log.Println("创建doc失败：", err)
	}
}

func (b *Builder) run() {
	for {
		select {
		case <-b.Ticker.C:
			log.Println("index scanner将要运行：")
			b.SingleRun()
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
