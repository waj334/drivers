module pkg.si-go.dev/drivers/cypress/cyw4343w/example

replace pkg.si-go.dev/chip => /home/waj334/Projects/chip

replace pkg.si-go.dev/drivers => ../../../..

go 1.26.1

require (
	pkg.si-go.dev/chip v0.0.0-20250822234245-1e12b76519dd
	pkg.si-go.dev/drivers v0.0.0-00010101000000-000000000000
)

require golang.org/x/exp v0.0.0-20250620022241-b7579e27df2b // indirect
