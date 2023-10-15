package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	apiCfg := ApiConfig{}
	apiCfg.initialize()
	defer apiCfg.close()

	apiCfg.ServerPort = os.Getenv("PORT")
	if apiCfg.ServerPort == "" {
		apiCfg.ServerPort = "8080"
	}

	server := http.Server{
		Addr: ":" + apiCfg.ServerPort,
	}

	http.HandleFunc("/server/status", apiCfg.HandleServerStatus)
	http.HandleFunc("/auth/register", apiCfg.HandleOwnerAccountRegistration)
	http.HandleFunc("/auth/login", apiCfg.HandleOwnerAccountLogin)
	http.HandleFunc("/auth/sa/login", apiCfg.HandleServiceAccountLogin)
	http.HandleFunc("/auth/change-password", apiCfg.authMiddleware(apiCfg.HandleOwnerAccountChangePassword))
	http.HandleFunc("/auth/check", apiCfg.authMiddleware(apiCfg.HandleCheckAuth))
	http.HandleFunc("/ws/check-avl", apiCfg.authMiddleware(apiCfg.handleCheckWorkspaceAvailability))
	http.HandleFunc("/ws/op", apiCfg.authMiddleware(apiCfg.handleWorkspaceOperations))
	http.HandleFunc("/ws/account", apiCfg.authMiddleware(apiCfg.handleWorkspaceAccountOperations))
	http.HandleFunc("/fs/dir/query", apiCfg.authMiddleware(apiCfg.HandleDirectoryQuery))
	http.HandleFunc("/fs/file/query", apiCfg.authMiddleware(apiCfg.HandleFileQuery))
	http.HandleFunc("/fs/dir/details", apiCfg.authMiddleware(apiCfg.handleDirDetailsQuery))
	http.HandleFunc("/fs/file/details", apiCfg.authMiddleware(apiCfg.handleFileDetailsQuery))
	http.HandleFunc("/fs/shared/query", apiCfg.authMiddleware(apiCfg.HandleFSShared))
	http.HandleFunc("/fs/upload/", apiCfg.authMiddleware(apiCfg.handleFileUpload))
	http.HandleFunc("/fs/download/", apiCfg.authMiddleware(apiCfg.handleFileDownload))
	http.HandleFunc("/role/op", apiCfg.authMiddleware(apiCfg.HandleRolesOperations))
	http.HandleFunc("/role/assign", apiCfg.authMiddleware(apiCfg.HandleAssignRoleToSA))
	http.HandleFunc("/roles/sa", apiCfg.authMiddleware(apiCfg.HandleGetAllAccountRoles))
	http.HandleFunc("/roles/details", apiCfg.authMiddleware(apiCfg.HandleGetAllRolesInWorkspace))
	http.HandleFunc("/rbac/fs", apiCfg.authMiddleware(apiCfg.HandleGetRoleFSPermissions))

	log.Default().Printf("Server starting at %v \n", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Default().Println("Server error while starting to listen : ", err.Error())
	}
}
