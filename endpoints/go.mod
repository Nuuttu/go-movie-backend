module example.com/endpoints

go 1.17

require (
	back-movie v0.0.0-00010101000000-000000000000
	example.com/mysctructs v0.0.0-00010101000000-000000000000
)

replace (
	back-movie => ../
	example.com/mysctructs => ../structs
)
