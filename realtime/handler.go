package realtime

import (
	"fmt"
	"net/http"
)

func outputSerialFlow() http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		job1 := makeBatchFileName("job1")
		job2 := makeBatchFileName("job2")
		job3 := makeBatchFileName("job3")
		fmt.Fprintf(w, `{"flow":"%s->%s->%s"}`, job1, job2, job3)

	}
	return http.HandlerFunc(f)
}

func outputParallelFlow() http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		job1 := makeBatchFileName("job1")
		job2 := makeBatchFileName("job2")
		job3 := makeBatchFileName("job3")
		fmt.Fprintf(w, `{"flow":"[%s,%s->%s]"}`, job1, job2, job3)

	}
	return http.HandlerFunc(f)
}

func outputWithJobDetail(name string) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		job := makeBatchFileName(name)
		fmt.Fprintf(w, `{"flow":"testjob","jobs":[{"name":"testjob","path":"%s"}]}`, job)

	}
	return http.HandlerFunc(f)
}

func outputTestJobOnly() http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"flow":"testjob"}`)

	}
	return http.HandlerFunc(f)
}
