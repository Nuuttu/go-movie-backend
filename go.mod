module back-movie

go 1.17

require (
	example.com/endpoints v0.0.0-00010101000000-000000000000
	example.com/excelporter v0.0.0-00010101000000-000000000000
	example.com/mysctructs v0.0.0-00010101000000-000000000000
	github.com/darahayes/go-boom v0.0.0-20200826120415-fa5cb724143a
	github.com/go-playground/validator/v10 v10.10.1
	github.com/gorilla/mux v1.8.0
	github.com/rs/xid v1.4.0
	github.com/xuri/excelize/v2 v2.6.0
	golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4
)

require (
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/richardlehane/mscfb v1.0.4 // indirect
	github.com/richardlehane/msoleps v1.0.1 // indirect
	github.com/xuri/efp v0.0.0-20220407160117-ad0f7a785be8 // indirect
	github.com/xuri/nfp v0.0.0-20220409054826-5e722a1d9e22 // indirect
	golang.org/x/net v0.0.0-20220412020605-290c469a71a5 // indirect
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad // indirect
	golang.org/x/text v0.3.7 // indirect
)

replace (
	example.com/endpoints => ./endpoints
	example.com/excelporter => ./utils/excelporter
	example.com/mysctructs => ./structs
)
