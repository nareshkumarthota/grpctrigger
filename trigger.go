package grpctrigger

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/nareshkumarthota/grpctrigger/pb"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

var addr string

//GRPCTriggerFactory gRPC Trigger factory
type GRPCTriggerFactory struct {
	metadata *trigger.Metadata
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &GRPCTriggerFactory{metadata: md}
}

//New Creates a new trigger instance for a given id
func (t *GRPCTriggerFactory) New(config *trigger.Config) trigger.Trigger {
	return &GRPCTrigger{metadata: t.metadata, config: config}
}

// GRPCTrigger is a stub for your Trigger implementation
type GRPCTrigger struct {
	metadata *trigger.Metadata
	runner   action.Runner
	config   *trigger.Config
}

// Init implements trigger.Trigger.Init
func (t *GRPCTrigger) Init(runner action.Runner) {

	if t.config.Settings == nil {
		panic(fmt.Sprintf("No Settings found for trigger '%s'", t.config.Id))
	}

	if _, ok := t.config.Settings["port"]; !ok {
		panic(fmt.Sprintf("No Port found for trigger '%s' in settings", t.config.Id))
	}

	addr = ":" + t.config.GetSetting("port")

	t.runner = runner

}

// Metadata implements trigger.Trigger.Metadata
func (t *GRPCTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Stop implements trigger.Trigger.Start
func (t *GRPCTrigger) Stop() error {
	// stop the trigger
	return nil
}

// Start implements trigger.Trigger.Start
func (t *GRPCTrigger) Start() error {
	// start the trigger
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	creds, err := credentials.NewServerTLSFromFile("cert.pem", "key.pem")
	if err != nil {
		log.Fatal(err)
	}

	opts := []grpc.ServerOption{grpc.Creds(creds)}

	s := grpc.NewServer(opts...)

	pb.RegisterEmployeeServiceServer(s, new(employeeService))

	log.Println("Starting server on port: ", addr)

	go func() {
		s.Serve(lis)
	}()

	log.Println("Server started")
	return nil
}

type employeeService struct{}

func (s *employeeService) GetByBadgeNumber(ctx context.Context, req *pb.GetByBadgeNumberRequest) (*pb.EmployeeResponse, error) {

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		fmt.Printf("metadata  received : %v\n", md)
	}

	for _, e := range employees {
		if e.BadgeNumber == req.BadgeNumber {
			return &pb.EmployeeResponse{Employee: &e}, nil
		}
	}

	return nil, errors.New("Employee not found")
}

func (s *employeeService) GetAll(req *pb.GetAllRequest, stream pb.EmployeeService_GetAllServer) error {
	for _, e := range employees {
		stream.Send(&pb.EmployeeResponse{Employee: &e})
	}
	return nil
}

func (s *employeeService) Save(ctx context.Context, req *pb.EmployeeRequest) (*pb.EmployeeResponse, error) {
	return nil, nil
}

func (s *employeeService) SaveAll(stream pb.EmployeeService_SaveAllServer) error {

	for {
		emp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		employees = append(employees, *emp.Employee)
		stream.Send(&pb.EmployeeResponse{Employee: emp.Employee})
	}

	for _, e := range employees {
		fmt.Println(e)
	}

	return nil
}

func (s *employeeService) AddPhoto(stream pb.EmployeeService_AddPhotoServer) error {

	md, ok := metadata.FromIncomingContext(stream.Context())
	if ok {
		fmt.Printf("recevied photo request for badge number %v \n", md["badgenumber"][0])
	}

	imageData := []byte{}

	for {
		data, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("filereceived with length %v \n", len(imageData))
			return stream.Send(&pb.AddPhotoResponse{IsOk: true})
		}
		if err != nil {
			return err
		}
		fmt.Printf("Received %v bytes \n", len(data.Data))
		imageData = append(imageData, data.Data...)
	}
}
