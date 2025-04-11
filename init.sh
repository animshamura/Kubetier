go mod init deployer
go get k8s.io/client-go@v0.28.1
go get k8s.io/apimachinery@v0.28.1
go get k8s.io/api@v0.28.1
go get k8s.io/client-go/plugin/pkg/client/auth # for auth providers (GCP, Azure, etc.)
