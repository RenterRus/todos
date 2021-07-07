package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"todoes/pb/todoes"
)

func (s *HTTPServer) SetExpTimeout(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	answer := response{}
	httpArgs := s.agregator(r.Form)
	fmt.Println("SetExpTimeout")

	id, err := strconv.Atoi(httpArgs["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Id invalid"))
		return
	}

	date := strings.Split(httpArgs["date"], ".")
	day, err := strconv.Atoi(date[0])
	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error in parse day: " + err.Error()))
		return
	}
	mount, err := strconv.Atoi(date[1])
	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error in parse mount: " + err.Error()))
		return
	}
	year, err := strconv.Atoi(date[2])
	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error in parse year: " + err.Error()))
		return
	}
	times := strings.Split(httpArgs["time"], ":")
	hours, err := strconv.Atoi(times[0])
	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error in parse hours: " + err.Error()))
		return
	}
	minuts, err := strconv.Atoi(times[1])
	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error in parse minuts: " + err.Error()))
		return
	}
	seconds, err := strconv.Atoi(times[2])
	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error in parse seconds: " + err.Error()))
		return
	}

	l, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		panic(err)
	}
	tTime := time.Date(year, time.Month(mount), day, hours, minuts, seconds, 0, l)

	res, err := s.grpc.SeTodoExpirationTimeoutTodo(context.Background(), &todoes.SeTodoExpirationTimeoutTodoRequest{IdTodo: int32(id), Deadline: tTime.Unix()})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		answer = response{
			Err: 500,
			Messgae: err.Error(),
		}
	}
	if res != nil{
		w.WriteHeader(http.StatusOK)
		answer = response{
			Err: res.Code,
			Messgae: res.Msg,
		}
	} else{
		w.WriteHeader(http.StatusInternalServerError)
		answer = response{
			Err: 500,
			Messgae: "Answer is empty",
		}
	}

	jsonAnswer, err := json.Marshal(&answer)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(jsonAnswer)
}