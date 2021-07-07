package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"todoes/pb/todoes"
)

func (s *HTTPServer) CreateTodos(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	httpArgs := s.agregator(r.Form)
	fmt.Println("Create")

	answer := response{}
	res, err := s.grpc.CreateTodo(context.Background(), &todoes.CreateTodoRequest{Message: httpArgs["message"]})
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