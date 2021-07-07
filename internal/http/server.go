package http

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"todoes/pb/todoes"
)

type HTTPServer struct {
	s                  http.Server
	grpc  todoes.TodoesClient
}

func GetGRPCConnect(grpcAddr string) *grpc.ClientConn {
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure(), grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                60 * time.Second,
		Timeout:             10 * time.Second,
		PermitWithoutStream: true,
	}))
	if err != nil {
		fmt.Println(err.Error())
	}
	return conn
}

func NewServer(addr, grpcAddr string) *HTTPServer {
	s := new(HTTPServer)
	s.s = http.Server{
		Addr: addr,
	}
	s.grpc = todoes.NewTodoesClient(GetGRPCConnect(grpcAddr))
	return s
}

func (s *HTTPServer) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/GetTodos", s.GetTodos)
	mux.HandleFunc("/Create", s.CreateTodos)
	mux.HandleFunc("/Update", s.UpdateTodo)
	mux.HandleFunc("/Close", s.CloseTodos)
	mux.HandleFunc("/SetExpTimeout", s.SetExpTimeout)
	s.s.Handler = mux
	return s.s.ListenAndServe()
}

func (s *HTTPServer) GraceShutdown() {
	s.s.Shutdown(context.Background())
}

type response struct {
	Err int32
	Messgae string
}

func (s *HTTPServer) SetExpTimeout(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	answer := response{}
	httpArgs := s.agregator(r.Form)
	fmt.Println("SetExpTimeout")

	id, err := strconv.Atoi(httpArgs["id"])
	if err != nil {
		w.Write([]byte("Id invalid"))
	}

	date := strings.Split(httpArgs["date"], ".")
	day, err := strconv.Atoi(date[0])
	if err != nil{
		w.Write([]byte("Error in parse day: " + err.Error()))
		return
	}
	mount, err := strconv.Atoi(date[1])
	if err != nil{
		w.Write([]byte("Error in parse mount: " + err.Error()))
		return
	}
	year, err := strconv.Atoi(date[2])
	if err != nil{
		w.Write([]byte("Error in parse year: " + err.Error()))
		return
	}
	times := strings.Split(httpArgs["time"], ":")
	hours, err := strconv.Atoi(times[0])
	if err != nil{
		w.Write([]byte("Error in parse hours: " + err.Error()))
		return
	}
	minuts, err := strconv.Atoi(times[1])
	if err != nil{
		w.Write([]byte("Error in parse minuts: " + err.Error()))
		return
	}
	seconds, err := strconv.Atoi(times[2])
	if err != nil{
		w.Write([]byte("Error in parse seconds: " + err.Error()))
		return
	}

	l, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		panic(err)
	}
	tTime := time.Date(year, time.Month(mount), day, hours, minuts, seconds, 0, l)

	fmt.Println(tTime.String())
	fmt.Println(tTime.Unix())


	t := time.Unix(0, 0)
	epoch := t.Unix()
	fmt.Println(epoch)
	fmt.Println(time.Now())
	fmt.Println(time.Now().Unix())

	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
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
		w.Write([]byte(err.Error()))
	}

	w.Write(jsonAnswer)
}

func (s *HTTPServer) agregator(form url.Values) map[string]string {
	httpArgs := map[string]string{}
	for k, v := range form {
		httpArgs[k] = strings.Join(v, "")
	}
	return httpArgs
}

func (s *HTTPServer) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	httpArgs := s.agregator(r.Form)
	fmt.Println("Update")

	id, err := strconv.Atoi(httpArgs["id"])
	if err != nil {
		w.Write([]byte("Id invalid"))
	}

	answer := response{}
	res, err := s.grpc.UpdateTodo(context.Background(), &todoes.UpdateTodoRequest{IdTodo: int32(id), Message: httpArgs["message"]})
	if err != nil {
		answer = response{
			Err: 500,
			Messgae: err.Error(),
		}
	}
	if res != nil{
		answer = response{
			Err: res.Code,
			Messgae: res.Msg,
		}
	} else{
		answer = response{
			Err: 500,
			Messgae: "Answer is empty",
		}
	}

	jsonAnswer, err := json.Marshal(&answer)
	if err != nil{
		w.Write([]byte(err.Error()))
	}

	w.Write(jsonAnswer)
}


func (s *HTTPServer) CloseTodos(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	httpArgs := s.agregator(r.Form)
	fmt.Println("Close")

	id, err := strconv.Atoi(httpArgs["id"])
	if err != nil {
		w.Write([]byte("Id invalid"))
	}
	answer := response{}
	res, err := s.grpc.DeleteTodo(context.Background(), &todoes.DeleteTodoRequest{IdTodo: int32(id)})
	if err != nil {
		answer = response{
			Err: 500,
			Messgae: err.Error(),
		}
	}
	if res != nil{
		answer = response{
			Err: res.Code,
			Messgae: res.Msg,
		}
	} else{
		answer = response{
			Err: 500,
			Messgae: "Answer is empty",
		}
	}

	jsonAnswer, err := json.Marshal(&answer)
	if err != nil{
		w.Write([]byte(err.Error()))
	}

	w.Write(jsonAnswer)
}

func (s *HTTPServer) CreateTodos(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	httpArgs := s.agregator(r.Form)
	fmt.Println("Create")

	answer := response{}
	res, err := s.grpc.CreateTodo(context.Background(), &todoes.CreateTodoRequest{Message: httpArgs["message"]})
	if err != nil {
		answer = response{
			Err: 500,
			Messgae: err.Error(),
		}
	}
	if res != nil{
		answer = response{
			Err: res.Code,
			Messgae: res.Msg,
		}
	} else{
		answer = response{
			Err: 500,
			Messgae: "Answer is empty",
		}
	}

	jsonAnswer, err := json.Marshal(&answer)
	if err != nil{
		w.Write([]byte(err.Error()))
	}

	w.Write(jsonAnswer)
}

type Todo struct{
	Id int32 `json:"Id"`
	Message string `json:"Message"`
	Deadline int64 `json:"Deadline"`
	Status string `json:"Status"`
}

type ResponseTodo struct{
	Message []Todo
}

func (s *HTTPServer) GetTodos(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	res, err := s.grpc.GetTodos(context.Background(), &emptypb.Empty{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error: " + err.Error()))
	}
	if res != nil{
		var resp []Todo
		errun := json.Unmarshal(res.Msg, &resp)
		if errun != nil{
			w.Write([]byte(errun.Error()))
		}
		jsonAnswer, errMarshal := json.Marshal(resp)

		if errMarshal != nil{
			w.Write([]byte(errMarshal.Error()))
		}
		w.WriteHeader(http.StatusOK)
		w.Write(jsonAnswer)
	} else{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Answer is nil"))
	}
}

