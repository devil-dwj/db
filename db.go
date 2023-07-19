package db

import (
	"errors"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	defaultMaxIdleConns = 8
	defaultMaxOpenConns = 4
	defaultMaxIdleTime  = time.Second * 30
	defaultMaxLifeTime  = time.Minute * 10
)

var (
	ErrNotFindRecord = errors.New("not find record")
)

type Option func(*options)

type options struct {
	maxIdleConns     int
	maxIdleOpenConns int
	maxIdleTime      time.Duration
	maxLifeTime      time.Duration
	trace            bool
}

func WithMaxIdleConns(n int) Option {
	return func(o *options) {
		o.maxIdleConns = n
	}
}

func WithMaxIdleOpenConns(n int) Option {
	return func(o *options) {
		o.maxIdleOpenConns = n
	}
}

func WithMaxIdleTime(t time.Duration) Option {
	return func(o *options) {
		o.maxIdleTime = t
	}
}

func WithMaxLifetime(t time.Duration) Option {
	return func(o *options) {
		o.maxLifeTime = t
	}
}

func WithTrace(b bool) Option {
	return func(o *options) {
		o.trace = b
	}
}

type DB struct {
	*gorm.DB

	dsn string
	opt options
}

func NewDB(dsn string, opts ...Option) *DB {
	o := options{
		maxIdleConns:     defaultMaxIdleConns,
		maxIdleOpenConns: defaultMaxOpenConns,
		maxIdleTime:      defaultMaxIdleTime,
		maxLifeTime:      defaultMaxLifeTime,
	}
	for _, opt := range opts {
		opt(&o)
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: NewLogger()})
	if err != nil {
		panic(err)
	}
	rawdb, err := db.DB()
	if err != nil {
		panic(err)
	}
	rawdb.SetMaxIdleConns(o.maxIdleConns)
	rawdb.SetMaxOpenConns(o.maxIdleOpenConns)
	rawdb.SetConnMaxIdleTime(o.maxIdleTime)
	rawdb.SetConnMaxLifetime(o.maxLifeTime)
	if o.trace {
		if err := db.Use(otelgorm.NewPlugin()); err != nil {
			panic(err)
		}
	}
	return &DB{DB: db, dsn: dsn, opt: o}
}
