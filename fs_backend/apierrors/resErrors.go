package apierrors

const (
	ResErrMethodNotAllowed       = "method-not-allowed"
	ResErrTokenNotFound          = "no-token-found"
	ResErrTokenInvalid           = "invalid-token"
	ResErrWeakPassword           = "weak-password"
	ResErrOAAlreadyExists        = "oa-already-exists"
	ResErrInvalidCredentials     = "invalid-credentials"
	ResErrInvalidData            = "invalid-data"
	ResErrServerError            = "server-error"
	ResErrAccountHasWorkspace    = "account-has-workspace"
	ResErrInvalidLocation        = "invalid-location"
	ResErrPermissionDenied       = "permission-denied"
	ResErrDirAlreadyExists       = "dir-already-exists"
	ResErrDirNotEmpty            = "dir-not-empty"
	ResErrInvalidUploadId        = "invalid-upload-id"
	ResErrInvalidWorkspaceName   = "invalid-workspace"
	ResErrWorkspaceAlreadyExists = "workspace-exists"
	ResErrSAAlreadyExists        = "sa-already-exists"
	ResErrSANotFound             = "sa-not-found"
	ResErrRoleNotFound           = "role-not-found"
	ResErrRoleAlreadyAssigned    = "role-already-assigned"
	ResErrRoleNotAssigned        = "role-not-assigned"
	ResErrResourceNotFound       = "resource-not-found"
)

func GetErrorCodeDescription(errorCode string) string {
	switch errorCode {
	case ResErrMethodNotAllowed:
		return "Method not allowed."
	case ResErrTokenNotFound:
		return "No token found in request header."
	case ResErrTokenInvalid:
		return "Token provided is invalid. Login again."
	case ResErrWeakPassword:
		return "Password is weak. Should be of atleast length 8 with one number and symbols."
	case ResErrOAAlreadyExists:
		return "Account with the email already. Use a different email address."
	case ResErrInvalidCredentials:
		return "Given username or password is wrong."
	case ResErrInvalidData:
		return "Request has invalid data."
	case ResErrServerError:
		return "Internal server error."
	case ResErrAccountHasWorkspace:
		return "Account is associated with one or more workspaces."
	case ResErrInvalidLocation:
		return "The specified location is not found."
	case ResErrPermissionDenied:
		return "The requested action is not allowed for this account."
	case ResErrDirAlreadyExists:
		return "The requested directory already exists."
	case ResErrDirNotEmpty:
		return "The directory not empty. Delete the containing files."
	case ResErrInvalidUploadId:
		return "The provided upload ID is not valid."
	case ResErrInvalidWorkspaceName:
		return "The provided workspace name is invalid."
	case ResErrWorkspaceAlreadyExists:
		return "A workspace already exists with the given name."
	case ResErrSAAlreadyExists:
		return "Service account already exists."
	case ResErrSANotFound:
		return "Requested service account not found."
	case ResErrRoleNotFound:
		return "Requested role not found"
	case ResErrRoleAlreadyAssigned:
		return "Role already assigned to the user."
	case ResErrResourceNotFound:
		return "The requested resource is not found"
	default:
		return ""
	}
}
