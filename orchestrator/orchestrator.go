package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mendium/orchestrator-c/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"
)

var (
	delayPlus     = 1000
	delayMinus    = 1000
	delayMultiply = 1000
	delayDivide   = 1000

	calcWorkers = 3
)

type Server struct {
	proto.OrchestratorServiceServer
}

type ComputingProfile struct {
	DelayPlus     uint `json:"delayPlus"`
	DelayMinus    uint `json:"delayMinus"`
	DelayMultiply uint `json:"delayMultiply"`
	DelayDivide   uint `json:"delayDivide"`
	CalcWorkers   uint `json:"calcWorkers"`
}

const dsn = "docker_test_exo:1111@tcp(db:3306)/docker_test"

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Orchestrate(ctx context.Context, exp *proto.Expression) (*proto.StatusCode, error) {
	log.Println("Новое выражение в очереди вычисления: " + exp.Expression)
	// канал для получения результатов
	resultChan := make(chan int)
	var wg sync.WaitGroup

	// создаем пул горутин с ограничением на 5 вычислителей-воркеров
	sem := make(chan struct{}, 5)

	// горутина для каждой операции
	wg.Add(1)
	sem <- struct{}{} // Резервируем место в пуле горутин
	go calculate(exp.Expression, resultChan, &wg, sem)

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	result := <-resultChan

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "UPDATE Tasks SET status = ?, answer = ? WHERE task_id = ?"

	_, err = db.Exec(query, "ready", result, exp.TaskId)
	if err != nil {
		log.Println("Ошибка при записи результата вычисления в базу данных")
		return nil, err
	}
	return &proto.StatusCode{StatusCode: "ok"}, nil
}

func calculate(expression string, resultChan chan int, wg *sync.WaitGroup, sem chan struct{}) {
	defer func() {
		wg.Done()
		<-sem // освобождение места в пуле горутинок
	}()

	// субвыражения
	operations := regexp.MustCompile(`\d+|\+|\-|\*|\/`).FindAllString(expression, -1)

	eval := func(op1, op2 int, operator rune) int {
		switch operator {
		case '+':
			time.Sleep(time.Duration(delayPlus))
			return op1 + op2
		case '-':
			time.Sleep(time.Duration(delayMinus))
			return op1 - op2
		case '*':
			time.Sleep(time.Duration(delayMultiply))
			return op1 * op2
		case '/':
			time.Sleep(time.Duration(delayDivide))
			return op1 / op2
		}
		return 0
	}

	// вычисляем выражение с учетом приоритета операций
	var stack []int
	var opStack []rune

	for _, op := range operations {
		if num, err := strconv.Atoi(op); err == nil {
			stack = append(stack, num)
		} else {
			currOp := rune(op[0])
			for len(opStack) > 0 && priority(opStack[len(opStack)-1]) >= priority(currOp) {
				stack[len(stack)-2] = eval(stack[len(stack)-2], stack[len(stack)-1], opStack[len(opStack)-1])
				stack = stack[:len(stack)-1]
				opStack = opStack[:len(opStack)-1]
			}
			opStack = append(opStack, currOp)
		}
	}

	for len(opStack) > 0 {
		stack[len(stack)-2] = eval(stack[len(stack)-2], stack[len(stack)-1], opStack[len(opStack)-1])
		stack = stack[:len(stack)-1]
		opStack = opStack[:len(opStack)-1]
	}

	// отправляем результат в канал
	resultChan <- stack[0]
}

func priority(operator rune) int {
	switch operator {
	case '*', '/':
		return 2
	case '+', '-':
		return 1
	}
	return 0
}

func main() {
	lis, err := net.Listen("tcp", "app2:5000") // будем ждать запросы по этому адресу

	if err != nil {
		log.Println("error starting tcp listener: ", err)
		os.Exit(1)
	}

	log.Println("tcp listener started at port: ", "5000")
	//Creating new GRPC Server
	grpcServer := grpc.NewServer()
	orchestrateServiceServer := NewServer()
	proto.RegisterOrchestratorServiceServer(grpcServer, orchestrateServiceServer)
	if err := grpcServer.Serve(lis); err != nil {
		log.Println("error serving grpc: ", err)
		os.Exit(1)
	}
}

/* func changeHandler(w http.ResponseWriter, r *http.Request) {
	var userProfile ComputingProfile
	err := json.NewDecoder(r.Body).Decode(&userProfile)
	if err != nil {
		http.Error(w, "Ошибка при чтении JSON", http.StatusBadRequest)
		return
	}
	delayPlus = int(userProfile.DelayPlus)
	delayMinus = int(userProfile.DelayMinus)
	delayMultiply = int(userProfile.DelayMultiply)
	delayDivide = int(userProfile.DelayDivide)
	calcWorkers = int(userProfile.CalcWorkers)

	fmt.Fprintf(w, "Профиль вычислителя успешно изменён.")
} */
