./build.sh
goexec 'http.ListenAndServe(`:8081`, http.FileServer(http.Dir(`.`)))'
