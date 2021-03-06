package client

import (
	"context"
	"fmt"
	"log"
	"sync"
	"sync/atomic"

	"github.com/overnest/strongdoc-go-sdk/proto"
	"github.com/overnest/strongdoc-go-sdk/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var clientInit uint32 = 0
var clientMutex sync.Mutex

type locationConfig struct {
	HostPort string
	Cert     string
}

// ServiceLocation specifies the location of the StrongDoc service
type ServiceLocation string

const (
	// DEFAULT is the default production service location
	DEFAULT ServiceLocation = "DEFAULT"
	// QA is the QA service used only for testing
	QA ServiceLocation = "QA"
	// LOCAL is the local service location used only for testing
	LOCAL ServiceLocation = "LOCAL"
	// unset specifies that the service location is not set
	unset ServiceLocation = "UNSET"
)

var serviceLocations = map[ServiceLocation]locationConfig{
	DEFAULT: locationConfig{"api.strongsalt.com:9090", "./certs/ssca.cert.pem"},
	QA:      locationConfig{"api.strongsaltqa.com:9090", "./certs/ssca.cert.pem"},
	LOCAL:   locationConfig{"localhost:9090", "./certs/localhost.crt"},
}

type strongDocManagerObj struct {
	location   ServiceLocation
	noAuthConn *grpc.ClientConn
	authConn   *grpc.ClientConn
	authToken  string
}

// Singletons
var serviceLocation ServiceLocation = unset
var manager *strongDocManagerObj = nil

// StrongDocManager encapsulates the client object that allows connection to the remote service
type StrongDocManager interface {
	Login(userID, password, orgID string) (token string, err error)
	GetNoAuthConn() *grpc.ClientConn
	GetAuthConn() *grpc.ClientConn
	GetClient() proto.StrongDocServiceClient
	Close()
}

// InitStrongDocManager initializes a singleton StrongDocManager
func InitStrongDocManager(location ServiceLocation, reset bool) (StrongDocManager, error) {
	_, ok := serviceLocations[location]
	if !ok || location == unset {
		return nil, fmt.Errorf("The ServiceLocation %v is not supported", location)
	}

	if atomic.LoadUint32(&clientInit) == 1 {
		if location == serviceLocation {
			return manager, nil
		} else if !reset {
			return nil, fmt.Errorf("Can not initialize StrongDocManager with service location %v. "+
				"Singleton already initialized with %v", location, serviceLocation)
		}
	}

	clientMutex.Lock()
	defer clientMutex.Unlock()

	if manager != nil {
		manager.Close()
	}
	serviceLocation = location
	config := serviceLocations[location]
	noAuthConn, err := getNoAuthConn(config.HostPort, config.Cert)
	if err != nil {
		return nil, err
	}
	manager = &strongDocManagerObj{location, noAuthConn, nil, ""}

	atomic.StoreUint32(&clientInit, 1)
	return manager, nil
}

// GetStrongDocManager gets a singleton StrongDocManager
func GetStrongDocManager() (StrongDocManager, error) {
	if atomic.LoadUint32(&clientInit) == 1 {
		if manager != nil {
			return manager, nil
		}
	}
	return nil, fmt.Errorf("Can not get StrongDocManager. Please call InitStrongDocManager to initialize")
}

// GetStrongDocClient gets a singleton StrongDocServiceClient
func GetStrongDocClient() (proto.StrongDocServiceClient, error) {
	if atomic.LoadUint32(&clientInit) == 1 {
		if manager != nil {
			return manager.GetClient(), nil
		}
	}
	return nil, fmt.Errorf("Can not get StrongDocClient. Please call InitStrongDocManager to initialize")
}

// Login attempts a log in. If successful, it generates an authenticatecd GRPC connection
func (c *strongDocManagerObj) Login(userID, password, orgID string) (token string, err error) {
	token = ""
	noAuthConn := c.GetNoAuthConn()
	if noAuthConn == nil || err != nil {
		log.Fatalf("Can not obtain none authenticated connection %s", err)
		return
	}

	noAuthClient := proto.NewStrongDocServiceClient(noAuthConn)
	res, err := noAuthClient.Login(context.Background(), &proto.LoginReq{
		UserID: userID, Password: password, OrgID: orgID})
	if err != nil {
		err = fmt.Errorf("Login err: [%v]", err)
		return
	}

	token = res.Token
	config := serviceLocations[c.location]
	authConn, err := getAuthConn(token, config.HostPort, config.Cert)
	if err != nil {
		return
	}

	// Close existing authenticated connection
	if c.authConn != nil {
		c.authConn.Close()
	}

	c.authConn = authConn
	c.authToken = token
	return
}

// GetNoAuthConn get the unauthenticated GRPC connection. This is always available, but will not work in most API calls
func (c *strongDocManagerObj) GetNoAuthConn() *grpc.ClientConn {
	return c.noAuthConn
}

// GetAuthConn gets an authenticated GRPC connection. This is available after a successful login.
func (c *strongDocManagerObj) GetAuthConn() *grpc.ClientConn {
	return c.authConn
}

// GetClient returns a StrongDocServiceClient used to call GRPC functions
func (c *strongDocManagerObj) GetClient() proto.StrongDocServiceClient {
	if c.GetAuthConn() != nil {
		return proto.NewStrongDocServiceClient(c.authConn)
	}
	return proto.NewStrongDocServiceClient(c.GetNoAuthConn())
}

// Close closes all the connections.
func (c *strongDocManagerObj) Close() {
	if c.GetAuthConn() != nil {
		c.GetAuthConn().Close()
	}
	if c.GetNoAuthConn() != nil {
		c.GetNoAuthConn().Close()
	}
}

func getNoAuthConn(hostport, cert string) (conn *grpc.ClientConn, err error) {
	certFilePath, err := utils.FetchFileLoc(cert)

	// Create the client TLS credentials
	creds, err := credentials.NewClientTLSFromFile(certFilePath, "")
	if err != nil {
		err = fmt.Errorf("Can not load TLS cert at %v", cert)
		return
	}

	// Initiate a connection with the server
	return grpc.DialContext(context.Background(), hostport,
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(&grpcNoAuthCred{}),
	)
}

func getAuthConn(token, hostport, cert string) (conn *grpc.ClientConn, err error) {
	certFilePath, err := utils.FetchFileLoc(cert)

	// Create the client TLS credentials
	creds, err := credentials.NewClientTLSFromFile(certFilePath, "")
	if err != nil {
		log.Fatalf("could not load tls cert: %s", err)
		return
	}

	// Initiate a connection with the server
	return grpc.DialContext(context.Background(), hostport,
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(&grpcAuthCred{token}),
	)
}
