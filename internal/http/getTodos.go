package http

import (
	"context"
	"encoding/json"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
)

func (s *HTTPServer) GetTodos(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	res, err := s.grpc.GetTodos(context.Background(), &emptypb.Empty{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error: " + err.Error()))
		return
	}
	if res != nil{
		var resp []Todo
		errun := json.Unmarshal(res.Msg, &resp)
		if errun != nil{
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errun.Error()))
			return
		}
		jsonAnswer, errMarshal := json.Marshal(resp)

		if errMarshal != nil{
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errMarshal.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(jsonAnswer)
		return
	} else{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Answer is nil"))
		return
	}
}

