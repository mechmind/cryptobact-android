package main

import (
    "log"
    "net/http"
)

func runServer() {
    log.Println("goerror: starting server")
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("hello from fucking goandroid\n"))
    })

    err := http.ListenAndServe(":8088", nil)
    if err != nil {
        log.Println("goerror: cannot start server: ", err.Error())
    }
}
