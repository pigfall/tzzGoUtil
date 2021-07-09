package http

import(
	stdhttp	"net/http"
)

func WrapHandlerAllowedCros(handler stdhttp.Handler)stdhttp.Handler {
    return stdhttp.HandlerFunc(
        func(w stdhttp.ResponseWriter,req *stdhttp.Request){
            // (w).Header().Set("Access-Control-Allow-Origin", "*")
            (w).Header().Set("Access-Control-Allow-Origin", "*")
            (w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
            (w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization,user")
            if req.Method == "OPTIONS"{
                return
            }
            handler.ServeHTTP(w,req)
        },
    )
}
