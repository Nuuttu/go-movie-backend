module example.com/excelporter

go 1.17

require (
	example.com/mysctructs v0.0.0-00010101000000-000000000000
)

replace (
	example.com/mysctructs => ../structs
)
