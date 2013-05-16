package test

import (
    "fmt"
    "harken/base/http"
)

func Test(session string, data string) *http.OutgoingMsg {
    response := http.OutgoingMsg{ Ext: "test", Method: "Test" }
    fmt.Println("Received", data);
    response.Data = "Received " + data
    return &response
}
