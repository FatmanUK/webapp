module dreamtrack.net/webapp

go 1.18

replace dreamtrack.net/fatgo/serialisers => ./fatgo/serialisers
require dreamtrack.net/fatgo/serialisers v0.0.0-00010101000000-000000000000 // indirect
toolchain go1.22.0

require (
	github.com/glebarez/sqlite v1.11.0
	gorm.io/gorm v1.25.10
)

require (
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/glebarez/go-sqlite v1.21.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	go.etcd.io/bbolt v1.3.10 // indirect
	golang.org/x/sys v0.7.0 // indirect
	gorm.io/driver/sqlite v1.5.5 // indirect
	modernc.org/libc v1.22.5 // indirect
	modernc.org/mathutil v1.5.0 // indirect
	modernc.org/memory v1.5.0 // indirect
	modernc.org/sqlite v1.23.1 // indirect
)
