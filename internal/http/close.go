package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"todoes/pb/todoes"
)

func (s *HTTPServer) CloseTodos(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	httpArgs := s.agregator(r.Form)
	fmt.Println("Close")

	id, err := strconv.Atoi(httpArgs["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Id invalid"))
		return
	}
	answer := response{}
	res, err := s.grpc.DeleteTodo(context.Background(), &todoes.DeleteTodoRequest{IdTodo: int32(id)})
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
