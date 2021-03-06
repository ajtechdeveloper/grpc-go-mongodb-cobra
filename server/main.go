package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	employeepb "github.com/ajtechdeveloper/grpc-go-mongodb-cobra/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EmployeeServiceServer struct {
}

func (s *EmployeeServiceServer) GetEmployee(ctx context.Context, req *employeepb.GetEmployeeRequest) (*employeepb.GetEmployeeResponse, error) {
	// Convert the string id (from proto) to mongoDB ObjectId
	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}
	result := employeedb.FindOne(ctx, bson.M{"_id": oid})
	// Create an empty EmployeeItem to write our decode result to
	data := EmployeeItem{}
	// Decode and write to data
	if err := result.Decode(&data); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find employee with Object Id %s: %v", req.GetId(), err))
	}
	// Cast to GetEmployeeResponse type
	response := &employeepb.GetEmployeeResponse{
		Employee: &employeepb.Employee{
			Id:         oid.Hex(),
			Name:       data.Name,
			Department: data.Department,
			Salary:     data.Salary,
		},
	}
	return response, nil
}

func (s *EmployeeServiceServer) CreateEmployee(ctx context.Context, req *employeepb.CreateEmployeeRequest) (*employeepb.CreateEmployeeResponse, error) {
	// Get the protobuf employee type from the protobuf request type
	// Essentially doing req.Employee to access the struct with a nil check
	employee := req.GetEmployee()
	// Now we have to convert this into a EmployeeItem type to convert into BSON
	data := EmployeeItem{
		Name:       employee.GetName(),
		Department: employee.GetDepartment(),
		Salary:     employee.GetSalary(),
	}

	// Insert the data into the database
	result, err := employeedb.InsertOne(mongoCtx, data)
	// check error
	if err != nil {
		// return internal gRPC error to be handled later
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}
	// Add the ID to employee
	oid := result.InsertedID.(primitive.ObjectID)
	employee.Id = oid.Hex()
	// Return the employee in a CreateEmployeeResponse type
	return &employeepb.CreateEmployeeResponse{Employee: employee}, nil
}

func (s *EmployeeServiceServer) UpdateEmployee(ctx context.Context, req *employeepb.UpdateEmployeeRequest) (*employeepb.UpdateEmployeeResponse, error) {
	// Get the employee data from the request
	employee := req.GetEmployee()

	// Convert the Id string to a MongoDB ObjectId
	oid, err := primitive.ObjectIDFromHex(employee.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Could not convert the supplied employee id to a MongoDB ObjectId: %v", err),
		)
	}

	// Convert the data to be updated into an unordered Bson document
	update := bson.M{
		"name":       employee.GetName(),
		"department": employee.GetDepartment(),
		"salary":     employee.GetSalary(),
	}

	// Convert the oid into an unordered bson document to search by ID
	filter := bson.M{"_id": oid}

	// Responseult is the BSON encoded result
	// To return the updated document instead of original we have to add options.
	result := employeedb.FindOneAndUpdate(ctx, filter, bson.M{"$set": update}, options.FindOneAndUpdate().SetReturnDocument(1))

	// Decode result and write it to 'decoded'
	decoded := EmployeeItem{}
	err = result.Decode(&decoded)
	if err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Could not find employee with supplied ID: %v", err),
		)
	}
	return &employeepb.UpdateEmployeeResponse{
		Employee: &employeepb.Employee{
			Id:         decoded.ID.Hex(),
			Name:       decoded.Name,
			Department: decoded.Department,
			Salary:     decoded.Salary,
		},
	}, nil
}

func (s *EmployeeServiceServer) DeleteEmployee(ctx context.Context, req *employeepb.DeleteEmployeeRequest) (*employeepb.DeleteEmployeeResponse, error) {
	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}
	// DeleteOne returns DeleteResponseult which is a struct containing the amount of deleted docs (in this case only 1 always)
	// So we return a boolean instead
	_, err = employeedb.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find/delete employee with id %s: %v", req.GetId(), err))
	}
	return &employeepb.DeleteEmployeeResponse{
		Success: true,
	}, nil
}

func (s *EmployeeServiceServer) GetAllEmployees(req *employeepb.GetAllEmployeesRequest, stream employeepb.EmployeeService_GetAllEmployeesServer) error {
	// Initiate a EmployeeItem type to write decoded data to
	data := &EmployeeItem{}
	// collection.Find returns a cursor for our (empty) query
	cursor, err := employeedb.Find(context.Background(), bson.M{})
	if err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Unknown internal error: %v", err))
	}
	// An expression with defer will be called at the end of the function
	defer cursor.Close(context.Background())
	// cursor.Next() returns a boolean, if false there are no more items and loop will break
	for cursor.Next(context.Background()) {
		// Decode the data at the current pointer and write it to data
		err := cursor.Decode(data)
		if err != nil {
			return status.Errorf(codes.Unavailable, fmt.Sprintf("Could not decode data: %v", err))
		}
		// If no error is found send employee over stream
		stream.Send(&employeepb.GetAllEmployeesResponse{
			Employee: &employeepb.Employee{
				Id:         data.ID.Hex(),
				Name:       data.Name,
				Salary:     data.Salary,
				Department: data.Department,
			},
		})
	}
	if err := cursor.Err(); err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Unkown cursor error: %v", err))
	}
	return nil
}

type EmployeeItem struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"name"`
	Salary     int32              `bson:"salary"`
	Department string             `bson:"department"`
}

var db *mongo.Client
var employeedb *mongo.Collection
var mongoCtx context.Context

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Starting server on port :50051...")

	// 50051 is the default port for gRPC
	// Ideally we'd use 0.0.0.0 instead of localhost as well
	listener, err := net.Listen("tcp", ":50051")

	if err != nil {
		log.Fatalf("Unable to listen on port :50051: %v", err)
	}

	// Slice of gRPC options
	// Here we can configure things like TLS
	opts := []grpc.ServerOption{}
	// var s *grpc.Server
	s := grpc.NewServer(opts...)
	// var srv *EmployeeServiceServer
	srv := &EmployeeServiceServer{}

	employeepb.RegisterEmployeeServiceServer(s, srv)

	// Initialize MongoDB client
	fmt.Println("Connecting to MongoDB...")
	mongoCtx = context.Background()
	db, err = mongo.Connect(mongoCtx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping(mongoCtx, nil)
	if err != nil {
		log.Fatalf("Could not connect to MongoDB: %v\n", err)
	} else {
		fmt.Println("Connected to MongoDB")
	}

	employeedb = db.Database("softwaredevelopercentral").Collection("employee")

	// Start the server in a child routine
	go func() {
		if err := s.Serve(listener); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()
	fmt.Println("Server succesfully started on port :50051")

	// Create a channel to receive OS signals
	c := make(chan os.Signal)

	// Relay os.Interrupt to our channel (os.Interrupt = CTRL+C)
	// Ignore other incoming signals
	signal.Notify(c, os.Interrupt)

	// Block main routine until a signal is received
	// As long as user does not press CTRL+C a message is not passed and our main routine keeps running
	// If the main routine were to shutdown, the child routine that is Serving the server would shutdown
	<-c

	fmt.Println("\nStopping the server...")
	s.Stop()
	listener.Close()
	fmt.Println("Closing MongoDB connection")
	db.Disconnect(mongoCtx)
	fmt.Println("Done.")
}
