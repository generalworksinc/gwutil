module github.com/generalworksinc/goutil

go 1.20

require (
	github.com/dsnet/compress v0.0.1
	github.com/getsentry/sentry-go v0.9.0
	github.com/google/uuid v1.2.0
	github.com/oklog/ulid/v2 v2.0.2
	github.com/pkg/errors v0.9.1
	github.com/yeka/zip v0.0.0-20180914125537-d046722c6feb
	golang.org/x/crypto v0.0.0-20201217014255-9d1352758620
	golang.org/x/text v0.3.2
)

require golang.org/x/sys v0.0.0-20191026070338-33540a1f6037 // indirect

replace github.com/generalworksinc/goutil => ./
