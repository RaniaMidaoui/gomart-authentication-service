package test

import (
	"context"
	"log"
	"net"
	"net/http"
	"testing"

	"github.com/RaniaMidaoui/goMart-authentication-service/pkg/config"
	"github.com/RaniaMidaoui/goMart-authentication-service/pkg/db"
	"github.com/RaniaMidaoui/goMart-authentication-service/pkg/pb"
	"github.com/RaniaMidaoui/goMart-authentication-service/pkg/services"
	"github.com/RaniaMidaoui/goMart-authentication-service/pkg/utils"
	"github.com/stretchr/testify/assert"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func server(ctx context.Context) (pb.AuthServiceClient, func(), error) {
	lis := bufconn.Listen(1024 * 1024)

	s := grpc.NewServer()
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatal(err)
	}

	h := db.Mock()

	jwt := utils.JwtWrapper{
		SecretKey:       c.JWTSecretKey,
		Issuer:          "gomart-auth-service",
		ExpirationHours: 24 * 365,
	}

	ss := services.Server{
		H:   h,
		Jwt: jwt,
	}

	pb.RegisterAuthServiceServer(s, &ss)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, nil, err
	}

	closer := func() {
		err := lis.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
		s.Stop()
	}

	return pb.NewAuthServiceClient(conn), closer, nil
}

func TestLogin(t *testing.T) {
	ctx := context.Background()
	client, closer, err := server(ctx)

	if err != nil {
		t.Fatal(err)
	}

	defer closer()

	type expectation struct {
		Status int64
	}

	// register_tests := map[string]struct {
	// 	req  *pb.RegisterRequest
	// 	want expectation
	// }{
	// 	"register": {
	// 		req: &pb.RegisterRequest{
	// 			Email:    "test@admida0ui.tech",
	// 			Password: "P@ssw0rd",
	// 		},
	// 		want: expectation{
	// 			Status: http.StatusCreated,
	// 		},
	// 	},
	// }

	login_tests := map[string]struct {
		req  *pb.LoginRequest
		want expectation
	}{
		"login": {
			req: &pb.LoginRequest{
				Email:    "test@admida0ui.tech",
				Password: "P@ssw0rd",
			},
			want: expectation{
				Status: http.StatusOK,
			},
		},

		"login with wrong email": {
			req: &pb.LoginRequest{
				Email:    "wrong@email.com",
				Password: "123456789",
			},
			want: expectation{
				Status: http.StatusNotFound,
			},
		},

		"login with wrong password": {
			req: &pb.LoginRequest{
				Email:    "test@admida0ui.tech",
				Password: "wrongpassword",
			},
			want: expectation{
				Status: http.StatusNotFound,
			},
		},
	}

	// for name, tc := range register_tests {
	// 	t.Run(name, func(t *testing.T) {
	// 		res, err := client.Register(ctx, tc.req)

	// 		if err != nil {
	// 			t.Fatal(err)
	// 		}

	// 		assert.Equal(t, tc.want.Status, res.Status)

	// 	})
	// }

	for name, tc := range login_tests {
		t.Run(name, func(t *testing.T) {
			res, err := client.Login(ctx, tc.req)

			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tc.want.Status, res.Status)

		})
	}

}
